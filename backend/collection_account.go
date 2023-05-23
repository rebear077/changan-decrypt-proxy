package server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/FISCO-BCOS/go-sdk/redis"
	"github.com/FISCO-BCOS/go-sdk/structure"
	types "github.com/FISCO-BCOS/go-sdk/type"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// 存储回款账户信息到redis数据库中
func (s *Server) StoreAccountToRedis(data []*types.CollectionAccount) {
	ctx := context.Background()

	for _, account := range data {
		id, err := uuid.NewUUID()
		if err != nil {
			logrus.Errorln(err)
		}
		values := make(map[string]interface{})
		key := account.Customerid + ":" + id.String()
		values["Backaccount"] = account.Backaccount
		values["Certificateid"] = account.Certificateid
		values["Customerid"] = account.Customerid
		values["Corpname"] = account.Corpname
		values["Lockremark"] = account.Lockremark
		values["Certificatetype"] = account.Certificatetype
		values["Intercustomerid"] = account.Intercustomerid
		err = s.redisCollectionAccount.MultipleSet(ctx, key, values)
		if err != nil {
			logrus.Errorln(err)
		}
	}
}

// 根据指令从redis中查询信息
// func (s *Server) SearchAccountFromRedis(order map[string]string) ([]*types.CollectionAccount, int) {
// 	pageid, err := strconv.ParseInt(order["pageid"], 10, 64)
// 	if err != nil {
// 		logrus.Errorln(err)
// 		return nil, 0
// 	}
// 	accounts := s.searchAccountByIDFromRedis(order["id"], order["searchType"])
// 	//redis未命中
// 	if len(accounts) == 0 {
// 		//同步mysql到redis
// 		s.DumpAccountFromMysqlToRedis(order["id"])
// 		time.Sleep(500 * time.Millisecond)
// 		//二次查询
// 		accounts = s.searchAccountByIDFromRedis(order["id"], order["searchType"])
// 		if len(accounts) == 0 {
// 			return nil, 0
// 		}
// 	}
// 	filterByPageId := s.filterByAccountPageId(accounts, pageid)
// 	totalcount := len(accounts)
// 	return filterByPageId, totalcount
// }

// 根据id的信息从redis中查询数据，如果结构体是空的，那么说明redis未命中，需要去mysql数据库中查询
func (s *Server) searchAccountByIDFromRedis(id string, order string) []*types.CollectionAccount {
	ctx := context.Background()
	accounts := make([]*types.CollectionAccount, 0)
	keys := s.GetMutipleAccountKeys(id)
	fmt.Println(keys)
	s.redisCollectionAccount.Del(ctx, "account")
	s.StoreAccountKeyAndScoreToZset(keys)
	//如果redis未命中,返回空的结构体
	if len(keys) == 0 {
		return nil
	}
	// if start == -1 || end == -1 {
	for _, key := range keys {
		resmap, err := s.redisCollectionAccount.GetAll(ctx, key)
		if err != nil {
			logrus.Errorln(err)
			continue
		}
		account := packToAccountStruct(resmap)
		accounts = append(accounts, account)
	}
	return accounts
	// } else {
	// 	keys := s.SearchAccountKeysFromZset(ctx, start, end, order)
	// 	for _, key := range keys {
	// 		resmap, err := s.redisCollectionAccount.GetAll(ctx, key)
	// 		if err != nil {
	// 			logrus.Errorln(err)
	// 			continue
	// 		}
	// 		account := packToAccountStruct(resmap)
	// 		accounts = append(accounts, account)
	// 	}
	// 	return accounts
	// }
}

func (s *Server) filterByAccountPageId(messages []*types.CollectionAccount, pageid int64) []*types.CollectionAccount {
	start := (pageid - 1) * 10
	end := pageid * 10
	result := make([]*types.CollectionAccount, 0)
	for i := start; i < end; i++ {
		if i >= int64(len(messages)) {
			break
		}
		result = append(result, messages[i])
	}
	return result
}

// 将从redis查询出来的数据转换成结构体
func packToAccountStruct(message map[string]string) *types.CollectionAccount {
	account := new(types.CollectionAccount)
	account.Backaccount = message["Backaccount"]
	account.Certificateid = message["Certificateid"]
	account.Customerid = message["Customerid"]
	account.Corpname = message["Corpname"]
	account.Lockremark = message["Lockremark"]
	account.Certificatetype = message["Certificatetype"]
	account.Intercustomerid = message["Intercustomerid"]
	return account

}

// redis未命中的情况下，去查询数据库中的数据，这种情况只适用于指定了id的情况，如果id未指定，则直接从redis数据库中返回信息
// 将mysql查询的数据首先存入redis，然后进行二次过滤
//
//	func (s *Server) DumpAccountFromMysqlToRedis(id string) {
//		plaintext := s.sql.QueryCollectionAccount(id)
//		accounts := s.sql.AccountinfoToMap(plaintext)
//		s.StoreAccountToRedis(accounts)
//	}
func (s *Server) PackToAccountJson(messages []*types.CollectionAccount, totalcount, currentPage int) string {
	returnresult := types.CollectionAccountReturn{
		CollectionAccountList: messages,
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
func (s *Server) GetMutipleAccountKeys(id string) []string {
	ctx := context.Background()
	var order string
	if id == "" {
		order = id + "*"
	} else {
		order = id + ":" + "*"
	}
	_, keys := s.redisCollectionAccount.Scan(ctx, order)
	return keys
}

// ********************************************************************************************************
// 将key值和分数存入zset集合中
func (s *Server) StoreAccountKeyAndScoreToZset(keys []string) {
	keyscore := make(map[string]float64)
	snow := new(structure.Snowflake)
	for _, key := range keys {
		//雪花算法生成分数
		score := float64(snow.NextVal())
		keyscore[key] = score
	}
	ctx := context.Background()
	s.redisCollectionAccount.ZAdd(ctx, "account", keyscore)
}

// 根据下标去查询key的数据
func (s *Server) SearchAccountKeysFromZset(ctx context.Context, start, end int64, order string) []string {

	if order == redis.Increase {
		keys := s.redisCollectionAccount.ZrangeIncrease(ctx, "account", start, end)
		return keys
	}
	if order == redis.Decrease {
		keys := s.redisCollectionAccount.ZrangeDecrease(ctx, "account", start, end)
		return keys
	}
	return nil
}
