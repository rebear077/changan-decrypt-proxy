package server

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/FISCO-BCOS/go-sdk/redis"
	types "github.com/FISCO-BCOS/go-sdk/type"
	"github.com/sirupsen/logrus"
)

// 存储融资意向信息到redis数据库中
func (s *Server) StoreFinanacingIntensionToRedis(data []*types.FinancingIntention) {
	ctx := context.Background()
	for _, intension := range data {
		values := make(map[string]interface{})
		key := intension.Customerid + ":" + intension.Financeid
		values["Custcdlinkposition"] = intension.Custcdlinkposition
		values["Custcdlinkname"] = intension.Custcdlinkname
		values["Certificateid"] = intension.Certificateid
		values["Corpname"] = intension.Corpname
		values["Remark"] = intension.Remark
		values["Bankcontact"] = intension.Bankcontact
		values["Banklinkname"] = intension.Banklinkname
		values["Custcdcontact"] = intension.Custcdcontact
		values["Customerid"] = intension.Customerid
		values["Financeid"] = intension.Financeid
		values["Cooperationyears"] = intension.Cooperationyears
		values["Certificatetype"] = intension.Certificatetype
		values["Intercustomerid"] = intension.Intercustomerid
		values["State"] = intension.State
		err := s.redisFinancingIntention.MultipleSet(ctx, key, values)
		if err != nil {
			logrus.Errorln(err)
		}
	}
}

// 根据指令从redis中查询融资意向信息
// func (s *Server) SearchIntensionFromRedis(order map[string]string) ([]*types.FinancingIntention, int) {
// 	pageid, err := strconv.ParseInt(order["pageid"], 10, 64)
// 	if err != nil {
// 		logrus.Errorln(err)
// 		return nil, 0
// 	}
// 	intensions := s.searchIntensionByIDFromRedis(order["id"], order["searchType"])
// 	//redis未命中
// 	if len(intensions) == 0 {
// 		//同步mysql到redis
// 		s.DumpIntensionFromMysqlToRedis(order["id"])
// 		time.Sleep(500 * time.Millisecond)
// 		//二次查询
// 		intensions = s.searchIntensionByIDFromRedis(order["id"], order["searchType"])
// 		if len(intensions) == 0 {
// 			return nil, 0
// 		}
// 	}
// 	fliterByFinancingId := s.fliterByIntensionID(intensions, order["financingId"])
// 	filterByPageId := s.filterByIntensionPageId(fliterByFinancingId, pageid)
// 	totalcount := len(fliterByFinancingId)
// 	return filterByPageId, totalcount
// }

// 根据id的信息从redis中查询数据，如果结构体是空的，那么说明redis未命中，需要去mysql数据库中查询
func (s *Server) searchIntensionByIDFromRedis(id string, order string) []*types.FinancingIntention {
	ctx := context.Background()
	intensions := make([]*types.FinancingIntention, 0)
	keys := s.GetMutipleIntensionKeys(id)
	s.redisFinancingIntention.Del(ctx, "intension")
	s.StoreIntensionKeyAndScoreToZset(keys)
	fmt.Println(keys)
	//如果redis未命中,返回空的结构体
	if len(keys) == 0 {
		return nil
	}
	for _, key := range keys {
		resmap, err := s.redisFinancingIntention.GetAll(ctx, key)
		if err != nil {
			logrus.Errorln(err)
			continue
		}
		intension := packToIntensionStruct(resmap)
		intensions = append(intensions, intension)
	}
	return intensions
}

// 根据融资请求编号进行过滤，调用此函数前，需要先通过id进行第一次检索
func (s *Server) fliterByIntensionID(messages []*types.FinancingIntention, financeId string) []*types.FinancingIntention {
	if financeId == "" {
		return messages
	}
	result := make([]*types.FinancingIntention, 0)
	for _, message := range messages {
		if message.Financeid == financeId {
			result = append(result, message)
		}
	}
	return result
}

func (s *Server) filterByIntensionPageId(messages []*types.FinancingIntention, pageid int64) []*types.FinancingIntention {
	start := (pageid - 1) * 10
	end := pageid * 10
	result := make([]*types.FinancingIntention, 0)
	for i := start; i < end; i++ {
		if i >= int64(len(messages)) {
			break
		}
		result = append(result, messages[i])
	}
	return result
}

// 将从redis查询出来的数据转换成结构体
func packToIntensionStruct(message map[string]string) *types.FinancingIntention {
	intension := new(types.FinancingIntention)
	intension.Custcdlinkposition = message["Custcdlinkposition"]
	intension.Custcdlinkname = message["Custcdlinkname"]
	intension.Certificateid = message["Certificateid"]
	intension.Corpname = message["Corpname"]
	intension.Remark = message["Remark"]
	intension.Bankcontact = message["Bankcontact"]
	intension.Banklinkname = message["Banklinkname"]
	intension.Custcdcontact = message["Custcdcontact"]
	intension.Customerid = message["Customerid"]
	intension.Financeid = message["Financeid"]
	intension.Cooperationyears = message["Cooperationyears"]
	intension.Certificatetype = message["Certificatetype"]
	intension.Intercustomerid = message["Intercustomerid"]
	intension.State = message["State"]

	return intension

}

// redis未命中的情况下，去查询数据库中的数据，这种情况只适用于指定了id的情况，如果id未指定，则直接从redis数据库中返回信息
// 将mysql查询的数据首先存入redis，然后进行二次过滤
//
//	func (s *Server) DumpIntensionFromMysqlToRedis(id string) {
//		plaintext := s.sql.QueryFinancingIntention(id)
//		intensions := s.sql.IntensioninfoToMap(plaintext)
//		s.StoreFinanacingIntensionToRedis(intensions)
//	}
func (s *Server) PackToIntensionJson(messages []*types.FinancingIntention, totalcount, currentPage int) string {
	returnresult := types.FinancingIntentionReturn{
		FinancingIntentionList: messages,
		TotalCount:             totalcount,
		CurrentPage:            currentPage,
	}
	// ans, err := json.Marshal(messages)
	ans, err := json.Marshal(returnresult)
	if err != nil {
		logrus.Errorln(err)
	}
	return string(ans)

}
func (s *Server) GetMutipleIntensionKeys(id string) []string {
	ctx := context.Background()
	var order string
	if id == "" {
		order = id + "*"
	} else {
		order = id + ":" + "*"
	}
	_, keys := s.redisFinancingIntention.Scan(ctx, order)
	return keys
}

// ********************************************************************************************************
// 将key值和分数存入zset集合中
func (s *Server) StoreIntensionKeyAndScoreToZset(keys []string) {
	keyscore := make(map[string]float64)
	for _, key := range keys {
		if !strings.Contains(key, ":") {
			continue
		}
		res := strings.Split(key, ":")
		high, err := strconv.ParseFloat(res[0], 64)
		if err != nil {
			logrus.Errorln(err)
			continue
		}
		low, err := strconv.ParseFloat(res[1][len(res[1])-3:], 64)
		if err != nil {
			logrus.Errorln(err)
			continue
		}
		score := high*1000 + low
		keyscore[key] = score
	}
	ctx := context.Background()
	s.redisFinancingIntention.ZAdd(ctx, "intension", keyscore)
}

// 根据下标去查询key的数据
func (s *Server) SearchIntensionKeysFromZset(ctx context.Context, start, end int64, order string) []string {

	if order == redis.Increase {
		keys := s.redisFinancingIntention.ZrangeIncrease(ctx, "intension", start, end)
		return keys
	}
	if order == redis.Decrease {
		keys := s.redisFinancingIntention.ZrangeDecrease(ctx, "intension", start, end)
		return keys
	}
	return nil
}
