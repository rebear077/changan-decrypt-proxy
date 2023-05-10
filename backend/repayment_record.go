package server

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/FISCO-BCOS/go-sdk/redis"
	"github.com/FISCO-BCOS/go-sdk/structure"
	types "github.com/FISCO-BCOS/go-sdk/type"
	"github.com/sirupsen/logrus"
)

// 存储还款信息到redis数据库中
func (s *Server) StoreRepaymentRecordToRedis(records []*types.RepaymentRecord) {
	ctx := context.Background()
	for _, record := range records {
		values := make(map[string]interface{})
		key := record.FinancingID
		fmt.Println(key)
		values["FinancingID"] = record.FinancingID
		values["CustomerID"] = record.CustomerID
		values["Time"] = record.Time
		values["RepaymentAmount"] = record.RepaymentAmount
		err := s.redisFinancingContract.MultipleSet(ctx, key, values)
		if err != nil {
			logrus.Errorln(err)
		}
	}
}

// 根据指令从redis中查询发票信息
func (s *Server) SearchRepaymentRecordFromRedis(order map[string]string) ([]*types.RepaymentRecord, int) {
	pageid, err := strconv.ParseInt(order["pageid"], 10, 64)
	if err != nil {
		logrus.Errorln(err)
		return nil, 0
	}
	records := s.searchRepaymentRecordByIDFromRedis(order["id"], order["searchType"])
	//redis未命中
	if len(records) == 0 {
		//同步mysql到redis
		s.DumpRepaymentRecordFromMysqlToRedis(order["id"])
		time.Sleep(500 * time.Millisecond)
		//二次查询
		records := s.searchRepaymentRecordByIDFromRedis(order["id"], order["searchType"])
		if len(records) == 0 {
			return nil, 0
		}
	}
	filterByFinancingID := s.fliterByFinanceID(records, order["FinanceId"])
	filterByPageId := s.filterByRepaymentRecordPageId(filterByFinancingID, pageid)
	totalcount := len(filterByPageId)
	return filterByPageId, totalcount
}

// 根据id的信息从redis中查询数据，如果结构体是空的，那么说明redis未命中，需要去mysql数据库中查询
func (s *Server) searchRepaymentRecordByIDFromRedis(id string, order string) []*types.RepaymentRecord {
	ctx := context.Background()
	records := make([]*types.RepaymentRecord, 0)
	keys := s.GetMutipleRepaymentRecordKeys(id)
	s.redisRepaymentRecord.Del(ctx, "repaymentRecord")
	s.StoreRepaymentRecordKeyAndScoreToZset(keys)
	fmt.Println(keys)
	//如果redis未命中,返回空的结构体
	if len(keys) == 0 {
		return nil
	}
	for _, key := range keys {
		fmt.Println(key, "-----------------")
		resmap, err := s.redisRepaymentRecord.GetAll(ctx, key)
		if err != nil {
			logrus.Errorln("err", err)
			continue
		}
		record := packToRepaymentRecordStruct(resmap)
		records = append(records, record)
	}
	return records
}
func (s *Server) fliterByFinanceID(messages []*types.RepaymentRecord, financeId string) []*types.RepaymentRecord {
	if financeId == "" {
		return messages
	}
	result := make([]*types.RepaymentRecord, 0)
	for _, message := range messages {
		if message.FinancingID == financeId {
			result = append(result, message)
		}
	}
	return result
}
func (s *Server) filterByRepaymentRecordPageId(messages []*types.RepaymentRecord, repaymentRecordpageid int64) []*types.RepaymentRecord {
	start := (repaymentRecordpageid - 1) * 10
	end := repaymentRecordpageid * 10
	result := make([]*types.RepaymentRecord, 0)
	for i := start; i < end; i++ {
		if i >= int64(len(messages)) {
			break
		}
		result = append(result, messages[i])
	}
	return result
}

// 将从redis查询出来的数据转换成结构体
func packToRepaymentRecordStruct(message map[string]string) *types.RepaymentRecord {
	record := new(types.RepaymentRecord)
	record.FinancingID = message["FinancingID"]
	record.CustomerID = message["CustomerID"]
	record.Time = message["Time"]
	record.RepaymentAmount = message["RepaymentAmount"]
	return record

}

// redis未命中的情况下，去查询数据库中的数据，这种情况只适用于指定了id的情况，如果id未指定，则直接从redis数据库中返回信息
// 将mysql查询的数据首先存入redis，然后进行二次过滤
func (s *Server) DumpRepaymentRecordFromMysqlToRedis(id string) {
	rawRecords := s.sql.QueryRepaymentRecord(id)
	records := s.sql.RepaymentRecordToMap(rawRecords)
	s.StoreRepaymentRecordToRedis(records)
}

func (s *Server) PackToRepaymentRecordJson(messages []*types.RepaymentRecord, totalcount, currentPage int) string {
	returnresult := types.RepaymentRecordReturn{
		RepaymentList: messages,
		TotalCount:    totalcount,
		CurrentPage:   currentPage,
	}
	ans, err := json.Marshal(returnresult)
	if err != nil {
		logrus.Errorln(err)
	}
	return string(ans)

}
func (s *Server) GetMutipleRepaymentRecordKeys(id string) []string {
	ctx := context.Background()
	var order string
	if id == "" {
		order = id + "*"
	} else {
		order = id + ":" + "*"
	}

	_, keys := s.redisRepaymentRecord.Scan(ctx, order)
	return keys
}

// ********************************************************************************************************
// 将key值和分数存入zset集合中
func (s *Server) StoreRepaymentRecordKeyAndScoreToZset(keys []string) {
	keyscore := make(map[string]float64)
	snow := new(structure.Snowflake)
	for _, key := range keys {
		score := float64(snow.NextVal())
		keyscore[key] = score
	}
	ctx := context.Background()
	s.redisRepaymentRecord.ZAdd(ctx, "repaymentRecord", keyscore)
}

// 根据下标去查询key的数据
func (s *Server) SearchRepaymentRecordKeysFromZset(ctx context.Context, start, end int64, order string) []string {

	if order == redis.Increase {
		keys := s.redisRepaymentRecord.ZrangeIncrease(ctx, "repaymentRecord", start, end)
		return keys
	}
	if order == redis.Decrease {
		keys := s.redisRepaymentRecord.ZrangeDecrease(ctx, "repaymentRecord", start, end)
		return keys
	}
	return nil
}
