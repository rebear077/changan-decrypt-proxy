package server

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/FISCO-BCOS/go-sdk/redis"
	sql "github.com/FISCO-BCOS/go-sdk/sqlController"
	types "github.com/FISCO-BCOS/go-sdk/type"
	"github.com/sirupsen/logrus"
)

func (s *Server) SearchEnterPoolByID(id string) []types.EnterpoolData {
	//对四个子表的mysql数据库进行检索，将检索的结果以[]string的形式返回，[]string中的每一个元素对应mysql数据表中的每一行数据
	plan_ret := s.sql.QueryEnterpoolDataPlanInfos(id)
	used_ret := s.sql.QueryEnterpoolDataUsedInfos(id)
	//记录所有出现过的Customerid与Tradeyearmonth的组合
	var CustomeridwithTradeyearmonth []string
	CustomeridwithTradeyearmonthMap := make(map[string]map[string]int)
	CustomeridwithTradeyearmonthSonMap := make(map[string]int)
	//属于同一笔交易的四个字表单头部信息相同，pooldataheader_map用于先记录头部信息
	//string Customerid+"|"+Tradeyearmonth
	pooldataheader_map := make(map[string]types.EnterpoolDataHeader)
	//将[]string转换成[]struct形式，结构体数组中的每一个元素对应一个plan表单
	plan_struct := sql.HandleEnterpoolDataPlaninfos(plan_ret)
	//构造双map，key值：Customerid Tradeyearmonth
	plan_map := make(map[string]map[string]types.EnterpoolDataPlaninfos)
	plan_sonmap := make(map[string]types.EnterpoolDataPlaninfos)
	//
	for _, planinfo := range plan_struct {
		plan_sonmap[planinfo.Planinfos[0].Tradeyearmonth] = *planinfo
		plan_map[planinfo.Customerid] = plan_sonmap
		if CustomeridwithTradeyearmonthMap[planinfo.Customerid][planinfo.Planinfos[0].Tradeyearmonth] == 0 {
			//如果该Customerid与Tradeyearmonth的组合还没有出现过，则记录
			CustomeridwithTradeyearmonth = append(CustomeridwithTradeyearmonth, planinfo.Customerid+"|"+planinfo.Planinfos[0].Tradeyearmonth)
			CustomeridwithTradeyearmonthSonMap[planinfo.Planinfos[0].Tradeyearmonth] = 1
			CustomeridwithTradeyearmonthMap[planinfo.Customerid] = CustomeridwithTradeyearmonthSonMap
			//记录该Customerid与Tradeyearmonth的组合所对应的头部信息
			var temp types.EnterpoolDataHeader
			temp.Datetimepoint = planinfo.Datetimepoint
			temp.Ccy = planinfo.Ccy
			temp.Customerid = planinfo.Customerid
			temp.Intercustomerid = planinfo.Intercustomerid
			temp.Receivablebalance = planinfo.Receivablebalance
			pooldataheader_map[planinfo.Customerid+"|"+planinfo.Planinfos[0].Tradeyearmonth] = temp
		}
	}
	//**********************************************************************************************************
	used_struct := sql.HandleEnterpoolDataProviderusedinfos(used_ret)
	used_map := make(map[string]map[string]types.EnterpoolDataProviderusedinfos)
	used_sonmap := make(map[string]types.EnterpoolDataProviderusedinfos)
	for _, usedinfo := range used_struct {
		used_sonmap[usedinfo.Providerusedinfos[0].Tradeyearmonth] = *usedinfo
		used_map[usedinfo.Customerid] = used_sonmap
		if CustomeridwithTradeyearmonthMap[usedinfo.Customerid][usedinfo.Providerusedinfos[0].Tradeyearmonth] == 0 {
			CustomeridwithTradeyearmonth = append(CustomeridwithTradeyearmonth, usedinfo.Customerid+"|"+usedinfo.Providerusedinfos[0].Tradeyearmonth)
			CustomeridwithTradeyearmonthSonMap[usedinfo.Providerusedinfos[0].Tradeyearmonth] = 1
			CustomeridwithTradeyearmonthMap[usedinfo.Customerid] = CustomeridwithTradeyearmonthSonMap

			var temp types.EnterpoolDataHeader
			temp.Datetimepoint = usedinfo.Datetimepoint
			temp.Ccy = usedinfo.Ccy
			temp.Customerid = usedinfo.Customerid
			temp.Intercustomerid = usedinfo.Intercustomerid
			temp.Receivablebalance = usedinfo.Receivablebalance
			pooldataheader_map[usedinfo.Customerid+"|"+usedinfo.Providerusedinfos[0].Tradeyearmonth] = temp
		}
	}

	//********************************************************************************************************************
	//拼接成EnterpoolData类型
	var ans []types.EnterpoolData
	for _, v := range CustomeridwithTradeyearmonth {
		var t types.EnterpoolData
		fields := strings.Split(v, "|")
		customerid := fields[0]
		tradeyearmonth := fields[1]

		t.Datetimepoint = pooldataheader_map[v].Datetimepoint
		t.Ccy = pooldataheader_map[v].Ccy
		t.Customerid = pooldataheader_map[v].Customerid
		t.Intercustomerid = pooldataheader_map[v].Intercustomerid
		t.Receivablebalance = pooldataheader_map[v].Receivablebalance

		t.Planinfos = plan_map[customerid][tradeyearmonth].Planinfos
		t.Providerusedinfos = used_map[customerid][tradeyearmonth].Providerusedinfos

		ans = append(ans, t)
	}
	return ans
}
func (s *Server) SearchEnterPoolBySQLID(id string) []types.EnterpoolData {
	//对四个子表的mysql数据库进行检索，将检索的结果以[]string的形式返回，[]string中的每一个元素对应mysql数据表中的每一行数据
	plan_ret := s.sql.QueryEnterpoolDataPlanInfosBySQLID(id)
	used_ret := s.sql.QueryEnterpoolDataUsedInfosBySQLID(id)
	//记录所有出现过的Customerid与Tradeyearmonth的组合
	var CustomeridwithTradeyearmonth []string
	CustomeridwithTradeyearmonthMap := make(map[string]map[string]int)
	CustomeridwithTradeyearmonthSonMap := make(map[string]int)
	//属于同一笔交易的四个字表单头部信息相同，pooldataheader_map用于先记录头部信息
	//string Customerid+"|"+Tradeyearmonth
	pooldataheader_map := make(map[string]types.EnterpoolDataHeader)
	//将[]string转换成[]struct形式，结构体数组中的每一个元素对应一个plan表单
	plan_struct := sql.HandleEnterpoolDataPlaninfos(plan_ret)
	//构造双map，key值：Customerid Tradeyearmonth
	plan_map := make(map[string]map[string]types.EnterpoolDataPlaninfos)
	plan_sonmap := make(map[string]types.EnterpoolDataPlaninfos)
	//
	for _, planinfo := range plan_struct {
		plan_sonmap[planinfo.Planinfos[0].Tradeyearmonth] = *planinfo
		plan_map[planinfo.Customerid] = plan_sonmap
		if CustomeridwithTradeyearmonthMap[planinfo.Customerid][planinfo.Planinfos[0].Tradeyearmonth] == 0 {
			//如果该Customerid与Tradeyearmonth的组合还没有出现过，则记录
			CustomeridwithTradeyearmonth = append(CustomeridwithTradeyearmonth, planinfo.Customerid+"|"+planinfo.Planinfos[0].Tradeyearmonth)
			CustomeridwithTradeyearmonthSonMap[planinfo.Planinfos[0].Tradeyearmonth] = 1
			CustomeridwithTradeyearmonthMap[planinfo.Customerid] = CustomeridwithTradeyearmonthSonMap
			//记录该Customerid与Tradeyearmonth的组合所对应的头部信息
			var temp types.EnterpoolDataHeader
			temp.Datetimepoint = planinfo.Datetimepoint
			temp.Ccy = planinfo.Ccy
			temp.Customerid = planinfo.Customerid
			temp.Intercustomerid = planinfo.Intercustomerid
			temp.Receivablebalance = planinfo.Receivablebalance
			pooldataheader_map[planinfo.Customerid+"|"+planinfo.Planinfos[0].Tradeyearmonth] = temp
		}
	}
	//**********************************************************************************************************
	used_struct := sql.HandleEnterpoolDataProviderusedinfos(used_ret)
	used_map := make(map[string]map[string]types.EnterpoolDataProviderusedinfos)
	used_sonmap := make(map[string]types.EnterpoolDataProviderusedinfos)
	for _, usedinfo := range used_struct {
		used_sonmap[usedinfo.Providerusedinfos[0].Tradeyearmonth] = *usedinfo
		used_map[usedinfo.Customerid] = used_sonmap
		if CustomeridwithTradeyearmonthMap[usedinfo.Customerid][usedinfo.Providerusedinfos[0].Tradeyearmonth] == 0 {
			CustomeridwithTradeyearmonth = append(CustomeridwithTradeyearmonth, usedinfo.Customerid+"|"+usedinfo.Providerusedinfos[0].Tradeyearmonth)
			CustomeridwithTradeyearmonthSonMap[usedinfo.Providerusedinfos[0].Tradeyearmonth] = 1
			CustomeridwithTradeyearmonthMap[usedinfo.Customerid] = CustomeridwithTradeyearmonthSonMap

			var temp types.EnterpoolDataHeader
			temp.Datetimepoint = usedinfo.Datetimepoint
			temp.Ccy = usedinfo.Ccy
			temp.Customerid = usedinfo.Customerid
			temp.Intercustomerid = usedinfo.Intercustomerid
			temp.Receivablebalance = usedinfo.Receivablebalance
			pooldataheader_map[usedinfo.Customerid+"|"+usedinfo.Providerusedinfos[0].Tradeyearmonth] = temp
		}
	}

	//********************************************************************************************************************
	//拼接成EnterpoolData类型
	var ans []types.EnterpoolData
	for _, v := range CustomeridwithTradeyearmonth {
		var t types.EnterpoolData
		fields := strings.Split(v, "|")
		customerid := fields[0]
		tradeyearmonth := fields[1]

		t.Datetimepoint = pooldataheader_map[v].Datetimepoint
		t.Ccy = pooldataheader_map[v].Ccy
		t.Customerid = pooldataheader_map[v].Customerid
		t.Intercustomerid = pooldataheader_map[v].Intercustomerid
		t.Receivablebalance = pooldataheader_map[v].Receivablebalance

		t.Planinfos = plan_map[customerid][tradeyearmonth].Planinfos
		t.Providerusedinfos = used_map[customerid][tradeyearmonth].Providerusedinfos

		ans = append(ans, t)
	}
	return ans
}

// 存储入池信息到redis数据库中
func (s *Server) StoreEnterPoolToRedis(data []types.EnterpoolData) {
	ctx := context.Background()

	for _, enterPool := range data {
		var key string
		if enterPool.Planinfos != nil {
			key = enterPool.Customerid + ":" + enterPool.Planinfos[0].Tradeyearmonth
		} else if enterPool.Providerusedinfos != nil {
			key = enterPool.Customerid + ":" + enterPool.Providerusedinfos[0].Tradeyearmonth
		}
		values := make(map[string]interface{})
		values["Datetimepoint"] = enterPool.Datetimepoint
		values["Ccy"] = enterPool.Ccy
		values["Customerid"] = enterPool.Customerid
		values["Intercustomerid"] = enterPool.Intercustomerid
		values["Receivablebalance"] = enterPool.Receivablebalance
		values["PlaninfosTradeyearmonth"] = ""
		values["PlaninfosPlanamount"] = ""
		values["PlaninfosCurrency"] = ""
		values["ProviderusedinfosTradeyearmonth"] = ""
		values["ProviderusedinfosUsedamount"] = ""
		values["ProviderusedinfosCurrency"] = ""
		if enterPool.Planinfos != nil {
			values["PlaninfosTradeyearmonth"] = enterPool.Planinfos[0].Tradeyearmonth
			values["PlaninfosPlanamount"] = enterPool.Planinfos[0].Planamount
			values["PlaninfosCurrency"] = enterPool.Planinfos[0].Currency
		}
		if enterPool.Providerusedinfos != nil {
			values["ProviderusedinfosTradeyearmonth"] = enterPool.Providerusedinfos[0].Tradeyearmonth
			values["ProviderusedinfosUsedamount"] = enterPool.Providerusedinfos[0].Usedamount
			values["ProviderusedinfosCurrency"] = enterPool.Providerusedinfos[0].Currency
		}
		err := s.redisEnterPool.MultipleSet(ctx, key, values)
		if err != nil {
			logrus.Errorln(err)
		}
	}

}

// 根据指令从redis中查询入池信息
func (s *Server) SearchEnterPoolFromRedis(order map[string]string) ([]*types.EnterpoolData, int) {
	pageid, err := strconv.ParseInt(order["pageid"], 10, 64)
	if err != nil {
		logrus.Errorln(err)
		return nil, 0
	}
	enterPool := s.searchEnterPoolByIDFromRedis(order["id"], order["searchType"])
	if len(enterPool) == 0 {
		s.DumpEnterPoolFromMysqlToRedis(order["id"])
		time.Sleep(500 * time.Millisecond)
		//二次查询
		enterPool = s.searchEnterPoolByIDFromRedis(order["id"], order["searchType"])
		if len(enterPool) == 0 {
			return nil, 0
		}
	}
	fliterByTime := s.fliterByEnterPoolTimeStamp(enterPool, order["Tradeyearmonth"])
	filterByPageId := s.filterByEnterPoolPageId(fliterByTime, pageid)
	totalcount := len(fliterByTime)
	return filterByPageId, totalcount
}

// 根据时间戳进行过滤，调用此函数前，需要先通过id进行第一次检索
func (s *Server) fliterByEnterPoolTimeStamp(messages []*types.EnterpoolData, txTimpeStamp string) []*types.EnterpoolData {
	if txTimpeStamp == "" {
		return messages
	}
	fields := strings.Split(txTimpeStamp, "to")
	if fields[0] == "" {
		fields[0] = "1970-01"
	}
	if fields[1] == "" {
		fields[1] = "2100-01"
	}
	beginData, _ := time.Parse("2006-01", fields[0])
	endData, _ := time.Parse("2006-01", fields[1])
	beginData = beginData.AddDate(0, -1, 0)
	endData = endData.AddDate(0, 1, 0)
	result := make([]*types.EnterpoolData, 0)
	for _, message := range messages {
		var checkInvoiceData time.Time
		if message.Planinfos != nil {
			checkInvoiceData, _ = time.Parse("2006-01", message.Planinfos[0].Tradeyearmonth)
		} else if message.Providerusedinfos != nil {
			checkInvoiceData, _ = time.Parse("2006-01", message.Providerusedinfos[0].Tradeyearmonth)
		}
		if checkInvoiceData.After(beginData) && checkInvoiceData.Before(endData) {
			result = append(result, message)
		}
	}
	return result

}

func (s *Server) filterByEnterPoolPageId(messages []*types.EnterpoolData, pageid int64) []*types.EnterpoolData {
	start := (pageid - 1) * 10
	end := pageid * 10
	result := make([]*types.EnterpoolData, 0)
	for i := start; i < end; i++ {
		if i >= int64(len(messages)) {
			break
		}
		result = append(result, messages[i])
	}
	return result
}

// 根据id信息从redis中查询历史交易信息，如果结构体是空的，那么说明redis未命中，需要去mysql数据库中查询
func (s *Server) searchEnterPoolByIDFromRedis(id string, order string) []*types.EnterpoolData {
	ctx := context.Background()
	enterpool := make([]*types.EnterpoolData, 0)

	keys := s.GetMutipleEnterPoolKeys(id)
	fmt.Println(len(keys), "+++++++++++++++++++++++++++")
	s.redisEnterPool.Del(ctx, "enterpool")
	s.StoreEnterPoolKeyAndScoreToZset(keys)
	fmt.Println("enterpool", keys)
	if len(keys) == 0 {

		return nil
	}
	// if start == -1 || end == -1 {

	for _, key := range keys {
		resmap, err := s.redisEnterPool.GetAll(ctx, key)
		if err != nil {
			logrus.Errorln(err)
			continue
		}
		tx := packToEnterPoolStruct(resmap)
		enterpool = append(enterpool, tx)
	}

	return enterpool
	// } else {
	// 	keys := s.SearchHisTXKeysFromZset(ctx, start, end, order)
	// 	for _, key := range keys {
	// 		resmap, err := s.redisEnterPool.GetAll(ctx, key)
	// 		if err != nil {
	// 			logrus.Errorln(err)
	// 			continue
	// 		}
	// 		tx := packToEnterPoolStruct(resmap)
	// 		enterpool = append(enterpool, tx)
	// 	}
	// 	return enterpool
	// }
}

// 将从redis查询出来的数据转换成结构体
func packToEnterPoolStruct(messages map[string]string) *types.EnterpoolData {
	enterPool := new(types.EnterpoolData)
	planinfo := new(types.Planinfos)
	providerUsedInfos := new(types.Providerusedinfos)
	enterPool.Datetimepoint = messages["Datetimepoint"]
	enterPool.Ccy = messages["Ccy"]
	enterPool.Customerid = messages["Customerid"]
	enterPool.Intercustomerid = messages["Intercustomerid"]
	enterPool.Receivablebalance = messages["Receivablebalance"]
	_, ok := messages["PlaninfosTradeyearmonth"]
	if ok {
		planinfo.Planamount = messages["PlaninfosPlanamount"]
		planinfo.Tradeyearmonth = messages["PlaninfosTradeyearmonth"]
		planinfo.Currency = messages["PlaninfosCurrency"]
		enterPool.Planinfos = append(enterPool.Planinfos, *planinfo)
	}
	_, ok = messages["ProviderusedinfosTradeyearmonth"]
	if ok {
		providerUsedInfos.Currency = messages["ProviderusedinfosCurrency"]
		providerUsedInfos.Tradeyearmonth = messages["ProviderusedinfosTradeyearmonth"]
		providerUsedInfos.Usedamount = messages["ProviderusedinfosUsedamount"]
		enterPool.Providerusedinfos = append(enterPool.Providerusedinfos, *providerUsedInfos)
	}

	return enterPool
}

// redis未命中的情况下，去查询数据库中的数据，这种情况只适用于指定了id的情况，如果id未指定，则直接从redis数据库中返回信息
// 将mysql查询的数据首先存入redis，然后进行二次过滤
func (s *Server) DumpEnterPoolFromMysqlToRedis(id string) {
	txs := s.SearchEnterPoolByID(id)
	s.StoreEnterPoolToRedis(txs)
}
func (s *Server) PackToEnterPoolJson(messages []*types.EnterpoolData, totalcount, currentPage int) string {
	returnresult := types.EnterpoolDataReturn{
		EnterpoolDataList: messages,
		TotalCount:        totalcount,
		CurrentPage:       currentPage,
	}
	// ans, err := json.Marshal(messages)
	ans, err := json.Marshal(returnresult)
	if err != nil {
		logrus.Errorln(err)
	}
	return string(ans)
}
func (s *Server) GetMutipleEnterPoolKeys(id string) []string {
	ctx := context.Background()
	var order string
	if id == "" {
		order = "*"
	} else {
		order = id + ":" + "*"
	}

	_, keys := s.redisEnterPool.Scan(ctx, order)
	return keys
}

// ********************************************************************************************************
// 将key值和分数存入zset集合中
func (s *Server) StoreEnterPoolKeyAndScoreToZset(keys []string) {
	keyscore := make(map[string]float64)
	fmt.Println(keys, "--------------------------")
	for _, key := range keys {
		if !strings.Contains(key, ":") {
			continue
		}
		res := strings.Split(key, ":")
		t, _ := time.Parse("2006-01", res[1])
		time, _ := strconv.Atoi(t.Format("200601"))
		// score := high*1000000 + float64(time)
		score := float64(time)
		keyscore[key] = score
	}
	ctx := context.Background()
	s.redisEnterPool.ZAdd(ctx, "enterpool", keyscore)
}

// 根据下标去查询key的数据
func (s *Server) SearchEnterPoolKeysFromZset(ctx context.Context, start, end int64, order string) []string {

	if order == redis.Increase {
		keys := s.redisEnterPool.ZrangeIncrease(ctx, "enterpool", start, end)
		return keys
	}
	if order == redis.Decrease {
		keys := s.redisEnterPool.ZrangeDecrease(ctx, "enterpool", start, end)
		return keys
	}
	return nil
}