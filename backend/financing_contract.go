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

// 存储借贷合同信息到redis数据库中
func (s *Server) StoreFinancingContractToRedis(contracts []*types.FinancingContract) {
	ctx := context.Background()
	for _, contract := range contracts {
		values := make(map[string]interface{})
		key := contract.FinancingID
		fmt.Println(key)
		values["FinancingID"] = contract.FinancingID
		values["CustomerID"] = contract.CustomerID
		values["CorpName"] = contract.CorpName
		values["DebtMoney"] = contract.DebtMoney
		values["SupplyDate"] = contract.SupplyDate
		values["ExpireDate"] = contract.ExpireDate
		values["Balance"] = contract.Balance
		err := s.redisFinancingContract.MultipleSet(ctx, key, values)
		if err != nil {
			logrus.Errorln(err)
		}
	}
}

// 根据指令从redis中查询发票信息
func (s *Server) SearchFinancingContractFromRedis(order map[string]string) ([]*types.FinancingContract, int) {
	pageid, err := strconv.ParseInt(order["pageid"], 10, 64)
	if err != nil {
		logrus.Errorln(err)
		return nil, 0
	}
	contracts := s.searchFinancingContractByIDFromRedis(order["FinanceId"], order["searchType"])
	//redis未命中
	if len(contracts) == 0 {
		//同步mysql到redis
		s.DumpFinancingContractFromMysqlToRedis(order["FinanceId"])
		time.Sleep(500 * time.Millisecond)
		//二次查询
		contracts := s.searchFinancingContractByIDFromRedis(order["FinanceId"], order["searchType"])
		if len(contracts) == 0 {
			return nil, 0
		}
	}
	filterByPageId := s.filterByFinancingContractPageId(contracts, pageid)
	totalcount := len(filterByPageId)
	return filterByPageId, totalcount
}

// 根据id的信息从redis中查询数据，如果结构体是空的，那么说明redis未命中，需要去mysql数据库中查询
func (s *Server) searchFinancingContractByIDFromRedis(id string, order string) []*types.FinancingContract {
	ctx := context.Background()
	contracts := make([]*types.FinancingContract, 0)
	keys := s.GetMutipleFinangcingContractKeys(id)
	s.redisFinancingContract.Del(ctx, "financingContract")
	s.StoreFinancingContractKeyAndScoreToZset(keys)
	fmt.Println(keys)
	//如果redis未命中,返回空的结构体
	if len(keys) == 0 {
		return nil
	}
	for _, key := range keys {
		fmt.Println(key, "-----------------")
		resmap, err := s.redisFinancingContract.GetAll(ctx, key)
		if err != nil {
			logrus.Errorln("err", err)
			continue
		}
		contract := packToFinancingContractStruct(resmap)
		contracts = append(contracts, contract)
	}
	return contracts
}

func (s *Server) filterByFinancingContractPageId(messages []*types.FinancingContract, financingContractpageid int64) []*types.FinancingContract {
	start := (financingContractpageid - 1) * 10
	end := financingContractpageid * 10
	result := make([]*types.FinancingContract, 0)
	for i := start; i < end; i++ {
		if i >= int64(len(messages)) {
			break
		}
		result = append(result, messages[i])
	}
	return result
}

// 将从redis查询出来的数据转换成结构体
func packToFinancingContractStruct(message map[string]string) *types.FinancingContract {
	contract := new(types.FinancingContract)
	contract.FinancingID = message["FinancingID"]
	contract.CustomerID = message["CustomerID"]
	contract.CorpName = message["CorpName"]
	contract.DebtMoney = message["DebtMoney"]
	contract.SupplyDate = message["SupplyDate"]
	contract.ExpireDate = message["ExpireDate"]
	contract.Balance = message["Balance"]
	return contract

}

// redis未命中的情况下，去查询数据库中的数据，这种情况只适用于指定了id的情况，如果id未指定，则直接从redis数据库中返回信息
// 将mysql查询的数据首先存入redis，然后进行二次过滤
func (s *Server) DumpFinancingContractFromMysqlToRedis(id string) {
	rawContracts := s.sql.QueryFinancingContract(id)
	contracts := s.sql.FinancingContractToMap(rawContracts)
	s.StoreFinancingContractToRedis(contracts)
}

func (s *Server) PackToFinancingContractJson(messages []*types.FinancingContract, totalcount, currentPage int) string {
	returnresult := types.FinancingContractReturn{
		FinancingContractList: messages,
		TotalCount:            totalcount,
		CurrentPage:           currentPage,
	}
	// ans, err := json.Marshal(messages)
	ans, err := json.Marshal(returnresult)
	if err != nil {
		logrus.Errorln(err)
	}
	return string(ans)

}
func (s *Server) GetMutipleFinangcingContractKeys(id string) []string {
	ctx := context.Background()
	var order string
	if id == "" {
		order = id + "*"
	} else {
		order = id + ":" + "*"
	}

	_, keys := s.redisFinancingContract.Scan(ctx, order)
	return keys
}

// ********************************************************************************************************
// 将key值和分数存入zset集合中
func (s *Server) StoreFinancingContractKeyAndScoreToZset(keys []string) {
	keyscore := make(map[string]float64)
	snow := new(structure.Snowflake)
	for _, key := range keys {
		score := float64(snow.NextVal())
		keyscore[key] = score
	}
	ctx := context.Background()
	s.redisFinancingContract.ZAdd(ctx, "financingContract", keyscore)
}

// 根据下标去查询key的数据
func (s *Server) SearchFinancingContractKeysFromZset(ctx context.Context, start, end int64, order string) []string {

	if order == redis.Increase {
		keys := s.redisFinancingContract.ZrangeIncrease(ctx, "financingContract", start, end)
		return keys
	}
	if order == redis.Decrease {
		keys := s.redisFinancingContract.ZrangeDecrease(ctx, "financingContract", start, end)
		return keys
	}
	return nil
}
