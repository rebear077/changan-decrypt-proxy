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
	"github.com/FISCO-BCOS/go-sdk/structure"
	types "github.com/FISCO-BCOS/go-sdk/type"
	"github.com/sirupsen/logrus"
)

// 根据传入的id从mysql数据库中查找出四个历史交易信息表单数据，解密后打包成types.TransactionHistory的数据结构
func (s *Server) SearchHistoryTXByID(id string) []types.TransactionHistory {
	// 对四个子表的mysql数据库进行检索，将检索的结果以[]string的形式返回，[]string中的每一个元素对应mysql数据表中的每一行数据
	used_ret := s.sql.QueryHistoricalTransUsedInfos(id)
	settle_ret := s.sql.QueryHistoricalTransSettleInfos(id)
	order_ret := s.sql.QueryHistoricalTransOrderInfos(id)
	receivable_ret := s.sql.QueryHistoricalTransReceivableInfos(id)
	used_struct := sql.HandleHistoricaltransactionUsedinfos(used_ret)
	settle_struct := sql.HandleHistoricaltransactionSettleinfos(settle_ret)
	order_struct := sql.HandleHistoricaltransactionOrderinfos(order_ret)
	receivable_struct := sql.HandleHistoricaltransactionReceivableinfos(receivable_ret)

	hisSettleMap := make(map[types.TransactionHistoryHeader][]types.Settleinfos)
	hisUsedMap := make(map[types.TransactionHistoryHeader][]types.Usedinfos)
	hisOrderMap := make(map[types.TransactionHistoryHeader][]types.Orderinfos)
	hisReceivableMap := make(map[types.TransactionHistoryHeader][]types.Receivableinfos)

	headerSet := structure.NewSet()
	dateSet := structure.NewSet()
	for _, tempUsedInfo := range used_struct {
		header := types.TransactionHistoryHeader{
			Customergrade:   tempUsedInfo.Customergrade,
			Certificatetype: tempUsedInfo.Certificatetype,
			Intercustomerid: tempUsedInfo.Intercustomerid,
			Corpname:        tempUsedInfo.Corpname,
			Financeid:       tempUsedInfo.Financeid,
			Certificateid:   tempUsedInfo.Certificateid,
			Customerid:      tempUsedInfo.Customerid,
		}
		headerSet.Add(header)
		hisUsedMap[header] = append(hisUsedMap[header], tempUsedInfo.Usedinfos)
		dateSet.Add(tempUsedInfo.Usedinfos.Tradeyearmonth)
	}
	for _, tempSettleInfo := range settle_struct {
		header := types.TransactionHistoryHeader{
			Customergrade:   tempSettleInfo.Customergrade,
			Certificatetype: tempSettleInfo.Certificatetype,
			Intercustomerid: tempSettleInfo.Intercustomerid,
			Corpname:        tempSettleInfo.Corpname,
			Financeid:       tempSettleInfo.Financeid,
			Certificateid:   tempSettleInfo.Certificateid,
			Customerid:      tempSettleInfo.Customerid,
		}
		headerSet.Add(header)
		hisSettleMap[header] = append(hisSettleMap[header], tempSettleInfo.Settleinfos)
		dateSet.Add(tempSettleInfo.Settleinfos.Tradeyearmonth)
	}
	for _, tempOrderInfo := range order_struct {
		header := types.TransactionHistoryHeader{
			Customergrade:   tempOrderInfo.Customergrade,
			Certificatetype: tempOrderInfo.Certificatetype,
			Intercustomerid: tempOrderInfo.Intercustomerid,
			Corpname:        tempOrderInfo.Corpname,
			Financeid:       tempOrderInfo.Financeid,
			Certificateid:   tempOrderInfo.Certificateid,
			Customerid:      tempOrderInfo.Customerid,
		}
		headerSet.Add(header)
		hisOrderMap[header] = append(hisOrderMap[header], tempOrderInfo.Orderinfos)
		dateSet.Add(tempOrderInfo.Orderinfos.Tradeyearmonth)
	}
	for _, tempReceivableInfo := range receivable_struct {
		header := types.TransactionHistoryHeader{
			Customergrade:   tempReceivableInfo.Customergrade,
			Certificatetype: tempReceivableInfo.Certificatetype,
			Intercustomerid: tempReceivableInfo.Intercustomerid,
			Corpname:        tempReceivableInfo.Corpname,
			Financeid:       tempReceivableInfo.Financeid,
			Certificateid:   tempReceivableInfo.Certificateid,
			Customerid:      tempReceivableInfo.Customerid,
		}
		headerSet.Add(header)
		hisReceivableMap[header] = append(hisReceivableMap[header], tempReceivableInfo.Receivableinfos)
		dateSet.Add(tempReceivableInfo.Receivableinfos.Tradeyearmonth)
	}
	tempHisInfos := make([]types.TransactionHistory, 0)
	headerList := headerSet.List()
	dateList := dateSet.List()
	for _, rawHeader := range headerList.([]interface{}) {
		for _, rawDate := range dateList.([]interface{}) {
			header := rawHeader.(types.TransactionHistoryHeader)
			date := rawDate.(string)
			historicalInfo := types.TransactionHistory{
				Customergrade:   header.Customergrade,
				Certificatetype: header.Certificatetype,
				Intercustomerid: header.Intercustomerid,
				Corpname:        header.Corpname,
				Financeid:       header.Financeid,
				Certificateid:   header.Certificateid,
				Customerid:      header.Customerid,
			}
			if _, ok := hisUsedMap[header]; ok {
				for _, usedInfo := range hisUsedMap[header] {
					if usedInfo.Tradeyearmonth == date {
						historicalInfo.Usedinfos = append(historicalInfo.Usedinfos, usedInfo)
					}
				}
			}
			if _, ok := hisSettleMap[header]; ok {
				for _, settleInfo := range hisSettleMap[header] {
					if settleInfo.Tradeyearmonth == date {
						historicalInfo.Settleinfos = append(historicalInfo.Settleinfos, settleInfo)
					}
				}
			}
			if _, ok := hisOrderMap[header]; ok {
				for _, orderInfo := range hisOrderMap[header] {
					if orderInfo.Tradeyearmonth == date {
						historicalInfo.Orderinfos = append(historicalInfo.Orderinfos, orderInfo)
					}
				}
			}
			if _, ok := hisReceivableMap[header]; ok {
				for _, receivableInfo := range hisReceivableMap[header] {
					if receivableInfo.Tradeyearmonth == date {
						historicalInfo.Receivableinfos = append(historicalInfo.Receivableinfos, receivableInfo)
					}
				}
			}
			tempHisInfos = append(tempHisInfos, historicalInfo)
		}
	}
	return tempHisInfos

}
func (s *Server) SearchHistoryTXBySQLID(id string) []types.TransactionHistory {
	// 对四个子表的mysql数据库进行检索，将检索的结果以[]string的形式返回，[]string中的每一个元素对应mysql数据表中的每一行数据
	used_ret := s.sql.QueryHistoricalTransUsedInfosBySQLID(id)
	settle_ret := s.sql.QueryHistoricalTransSettleInfosBySQLID(id)
	order_ret := s.sql.QueryHistoricalTransOrderInfosBySQLID(id)
	receivable_ret := s.sql.QueryHistoricalTransReceivableInfosBySQLID(id)
	used_struct := sql.HandleHistoricaltransactionUsedinfos(used_ret)
	settle_struct := sql.HandleHistoricaltransactionSettleinfos(settle_ret)
	order_struct := sql.HandleHistoricaltransactionOrderinfos(order_ret)
	receivable_struct := sql.HandleHistoricaltransactionReceivableinfos(receivable_ret)

	hisSettleMap := make(map[types.TransactionHistoryHeader][]types.Settleinfos)
	hisUsedMap := make(map[types.TransactionHistoryHeader][]types.Usedinfos)
	hisOrderMap := make(map[types.TransactionHistoryHeader][]types.Orderinfos)
	hisReceivableMap := make(map[types.TransactionHistoryHeader][]types.Receivableinfos)

	headerSet := structure.NewSet()
	dateSet := structure.NewSet()
	for _, tempUsedInfo := range used_struct {
		header := types.TransactionHistoryHeader{
			Customergrade:   tempUsedInfo.Customergrade,
			Certificatetype: tempUsedInfo.Certificatetype,
			Intercustomerid: tempUsedInfo.Intercustomerid,
			Corpname:        tempUsedInfo.Corpname,
			Financeid:       tempUsedInfo.Financeid,
			Certificateid:   tempUsedInfo.Certificateid,
			Customerid:      tempUsedInfo.Customerid,
		}
		headerSet.Add(header)
		hisUsedMap[header] = append(hisUsedMap[header], tempUsedInfo.Usedinfos)
		dateSet.Add(tempUsedInfo.Usedinfos.Tradeyearmonth)
	}
	for _, tempSettleInfo := range settle_struct {
		header := types.TransactionHistoryHeader{
			Customergrade:   tempSettleInfo.Customergrade,
			Certificatetype: tempSettleInfo.Certificatetype,
			Intercustomerid: tempSettleInfo.Intercustomerid,
			Corpname:        tempSettleInfo.Corpname,
			Financeid:       tempSettleInfo.Financeid,
			Certificateid:   tempSettleInfo.Certificateid,
			Customerid:      tempSettleInfo.Customerid,
		}
		headerSet.Add(header)
		hisSettleMap[header] = append(hisSettleMap[header], tempSettleInfo.Settleinfos)
		dateSet.Add(tempSettleInfo.Settleinfos.Tradeyearmonth)
	}
	for _, tempOrderInfo := range order_struct {
		header := types.TransactionHistoryHeader{
			Customergrade:   tempOrderInfo.Customergrade,
			Certificatetype: tempOrderInfo.Certificatetype,
			Intercustomerid: tempOrderInfo.Intercustomerid,
			Corpname:        tempOrderInfo.Corpname,
			Financeid:       tempOrderInfo.Financeid,
			Certificateid:   tempOrderInfo.Certificateid,
			Customerid:      tempOrderInfo.Customerid,
		}
		headerSet.Add(header)
		hisOrderMap[header] = append(hisOrderMap[header], tempOrderInfo.Orderinfos)
		dateSet.Add(tempOrderInfo.Orderinfos.Tradeyearmonth)
	}
	for _, tempReceivableInfo := range receivable_struct {
		header := types.TransactionHistoryHeader{
			Customergrade:   tempReceivableInfo.Customergrade,
			Certificatetype: tempReceivableInfo.Certificatetype,
			Intercustomerid: tempReceivableInfo.Intercustomerid,
			Corpname:        tempReceivableInfo.Corpname,
			Financeid:       tempReceivableInfo.Financeid,
			Certificateid:   tempReceivableInfo.Certificateid,
			Customerid:      tempReceivableInfo.Customerid,
		}
		headerSet.Add(header)
		hisReceivableMap[header] = append(hisReceivableMap[header], tempReceivableInfo.Receivableinfos)
		dateSet.Add(tempReceivableInfo.Receivableinfos.Tradeyearmonth)
	}
	tempHisInfos := make([]types.TransactionHistory, 0)
	headerList := headerSet.List()
	dateList := dateSet.List()
	for _, rawHeader := range headerList.([]interface{}) {
		for _, rawDate := range dateList.([]interface{}) {
			header := rawHeader.(types.TransactionHistoryHeader)
			date := rawDate.(string)
			historicalInfo := types.TransactionHistory{
				Customergrade:   header.Customergrade,
				Certificatetype: header.Certificatetype,
				Intercustomerid: header.Intercustomerid,
				Corpname:        header.Corpname,
				Financeid:       header.Financeid,
				Certificateid:   header.Certificateid,
				Customerid:      header.Customerid,
			}
			if _, ok := hisUsedMap[header]; ok {
				for _, usedInfo := range hisUsedMap[header] {
					if usedInfo.Tradeyearmonth == date {
						historicalInfo.Usedinfos = append(historicalInfo.Usedinfos, usedInfo)
					}
				}
			}
			if _, ok := hisSettleMap[header]; ok {
				for _, settleInfo := range hisSettleMap[header] {
					if settleInfo.Tradeyearmonth == date {
						historicalInfo.Settleinfos = append(historicalInfo.Settleinfos, settleInfo)
					}
				}
			}
			if _, ok := hisOrderMap[header]; ok {
				for _, orderInfo := range hisOrderMap[header] {
					if orderInfo.Tradeyearmonth == date {
						historicalInfo.Orderinfos = append(historicalInfo.Orderinfos, orderInfo)
					}
				}
			}
			if _, ok := hisReceivableMap[header]; ok {
				for _, receivableInfo := range hisReceivableMap[header] {
					if receivableInfo.Tradeyearmonth == date {
						historicalInfo.Receivableinfos = append(historicalInfo.Receivableinfos, receivableInfo)
					}
				}
			}
			tempHisInfos = append(tempHisInfos, historicalInfo)
		}
	}
	return tempHisInfos
}

// 存储历史交易信息到redis数据库中
func (s *Server) StoreHistoryTXToRedis(data []types.TransactionHistory) {
	ctx := context.Background()

	for _, tx := range data {
		var key string

		if tx.Usedinfos != nil {
			key = tx.Customerid + ":" + tx.Usedinfos[0].Tradeyearmonth
		} else if tx.Settleinfos != nil {
			key = tx.Customerid + ":" + tx.Settleinfos[0].Tradeyearmonth
		} else if tx.Orderinfos != nil {
			key = tx.Customerid + ":" + tx.Orderinfos[0].Tradeyearmonth
		} else if tx.Receivableinfos != nil {
			key = tx.Customerid + ":" + tx.Receivableinfos[0].Tradeyearmonth
		}

		values := make(map[string]interface{})
		values["Customergrade"] = tx.Customergrade
		values["Certificatetype"] = tx.Certificatetype
		values["Intercustomerid"] = tx.Intercustomerid
		values["Corpname"] = tx.Corpname
		values["Financeid"] = tx.Financeid
		values["Certificateid"] = tx.Certificateid
		values["Customerid"] = tx.Customerid
		values["UsedinfosTradeyearmonth"] = ""
		values["UsedinfosUsedamount"] = ""
		values["UsedinfosCcy"] = ""
		values["SettleinfosTradeyearmonth"] = ""
		values["SettleinfosSettleamount"] = ""
		values["SettleinfosCcy"] = ""
		values["OrderinfosTradeyearmonth"] = ""
		values["OrderinfosOrderamount"] = ""
		values["OrderinfosCcy"] = ""
		values["ReceivableinfosTradeyearmonth"] = ""
		values["ReceivableinfosReceivableamount"] = ""
		values["ReceivableinfosCcy"] = ""

		if tx.Usedinfos != nil {
			values["UsedinfosTradeyearmonth"] = tx.Usedinfos[0].Tradeyearmonth
			values["UsedinfosUsedamount"] = tx.Usedinfos[0].Usedamount
			values["UsedinfosCcy"] = tx.Usedinfos[0].Ccy
		}
		if tx.Settleinfos != nil {
			values["SettleinfosTradeyearmonth"] = tx.Settleinfos[0].Tradeyearmonth
			values["SettleinfosSettleamount"] = tx.Settleinfos[0].Settleamount
			values["SettleinfosCcy"] = tx.Settleinfos[0].Ccy
		}
		if tx.Orderinfos != nil {
			values["OrderinfosTradeyearmonth"] = tx.Orderinfos[0].Tradeyearmonth
			values["OrderinfosOrderamount"] = tx.Orderinfos[0].Orderamount
			values["OrderinfosCcy"] = tx.Orderinfos[0].Ccy
		}
		if tx.Receivableinfos != nil {
			values["ReceivableinfosTradeyearmonth"] = tx.Receivableinfos[0].Tradeyearmonth
			values["ReceivableinfosReceivableamount"] = tx.Receivableinfos[0].Receivableamount
			values["ReceivableinfosCcy"] = tx.Receivableinfos[0].Ccy
		}
		// fmt.Println(values)
		err := s.redisHistoryTX.MultipleSet(ctx, key, values)
		if err != nil {
			logrus.Errorln(err)
		}

	}
}

// 根据指令从redis中查询历史交易信息
func (s *Server) SearchHistoryTXFromRedis(order map[string]string) ([]*types.TransactionHistory, int) {
	pageid, err := strconv.ParseInt(order["pageid"], 10, 64)
	if err != nil {
		logrus.Errorln(err)
		return nil, 0
	}
	txs := s.searchHistoryTXByIDFromRedis(order["id"], order["searchType"])

	if len(txs) == 0 {
		s.DumpHistoryTXFromMysqlToRedis(order["id"])
		time.Sleep(500 * time.Millisecond)
		//二次查询
		txs = s.searchHistoryTXByIDFromRedis(order["id"], order["searchType"])
		if len(txs) == 0 {
			return nil, 0
		}
	}
	for _, v := range txs {
		fmt.Println(v)
	}
	fliterByTime := s.fliterByHistoryTXTimeStamp(txs, order["tradeyearmonth"])
	fliterByFinanceId := s.fliterByHistoryTXFinanceId(fliterByTime, order["financeId"])
	filterByPageId := s.filterByHistoryTXPageId(fliterByFinanceId, pageid)
	totalcount := len(fliterByFinanceId)
	return filterByPageId, totalcount
}

// 根据发票的时间戳进行过滤，调用此函数前，需要先通过id进行第一次检索
// 根据时间戳进行过滤，调用此函数前，需要先通过id进行第一次检索
func (s *Server) fliterByHistoryTXTimeStamp(messages []*types.TransactionHistory, txTimpeStamp string) []*types.TransactionHistory {
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
	result := make([]*types.TransactionHistory, 0)
	for _, message := range messages {
		var checkInvoiceData time.Time
		if message.Usedinfos != nil {
			checkInvoiceData, _ = time.Parse("2006-01", message.Usedinfos[0].Tradeyearmonth)
		} else if message.Settleinfos != nil {
			checkInvoiceData, _ = time.Parse("2006-01", message.Settleinfos[0].Tradeyearmonth)
		} else if message.Orderinfos != nil {
			checkInvoiceData, _ = time.Parse("2006-01", message.Orderinfos[0].Tradeyearmonth)
		} else if message.Receivableinfos != nil {
			checkInvoiceData, _ = time.Parse("2006-01", message.Receivableinfos[0].Tradeyearmonth)
		}
		if checkInvoiceData.After(beginData) && checkInvoiceData.Before(endData) {
			result = append(result, message)
		}
	}
	return result
}

func (s *Server) fliterByHistoryTXFinanceId(messages []*types.TransactionHistory, financeId string) []*types.TransactionHistory {
	if financeId == "" {
		return messages
	}
	result := make([]*types.TransactionHistory, 0)
	for _, message := range messages {
		if message.Financeid == financeId {
			result = append(result, message)
		}
	}
	return result
}

func (s *Server) filterByHistoryTXPageId(messages []*types.TransactionHistory, pageid int64) []*types.TransactionHistory {
	start := (pageid - 1) * 10
	end := pageid * 10
	result := make([]*types.TransactionHistory, 0)
	for i := start; i < end; i++ {
		if i >= int64(len(messages)) {
			break
		}
		result = append(result, messages[i])
	}
	return result
}

// 根据id信息从redis中查询历史交易信息，如果结构体是空的，那么说明redis未命中，需要去mysql数据库中查询
func (s *Server) searchHistoryTXByIDFromRedis(id string, order string) []*types.TransactionHistory {
	ctx := context.Background()
	txs := make([]*types.TransactionHistory, 0)
	keys := s.GetMutipleHistoryTXKeys(id)
	s.redisHistoryTX.Del(ctx, "hisTX")
	s.StoreHisTXKeyAndScoreToZset(keys)
	fmt.Println("historykey", keys)
	if len(keys) == 0 {
		return nil
	}
	// if start == -1 || end == -1 {
	for _, key := range keys {
		resmap, err := s.redisHistoryTX.GetAll(ctx, key)
		if err != nil {
			logrus.Errorln(err)
			continue
		}
		tx := packToHistoryTXStruct(resmap)
		txs = append(txs, tx)
	}
	return txs
}

// 将从redis查询出来的数据转换成结构体
func packToHistoryTXStruct(messages map[string]string) *types.TransactionHistory {
	tx := new(types.TransactionHistory)
	usedinfo := new(types.Usedinfos)
	settleinfo := new(types.Settleinfos)
	orderinfo := new(types.Orderinfos)
	receivableinfo := new(types.Receivableinfos)
	tx.Customergrade = messages["Customergrade"]
	tx.Certificatetype = messages["Certificatetype"]
	tx.Intercustomerid = messages["Intercustomerid"]
	tx.Corpname = messages["Corpname"]
	tx.Financeid = messages["Financeid"]
	tx.Certificateid = messages["Certificateid"]
	tx.Customerid = messages["Customerid"]
	_, ok := messages["UsedinfosTradeyearmonth"]
	if ok {

		usedinfo.Usedamount = messages["UsedinfosUsedamount"]
		usedinfo.Tradeyearmonth = messages["UsedinfosTradeyearmonth"]
		usedinfo.Ccy = messages["UsedinfosCcy"]
		tx.Usedinfos = append(tx.Usedinfos, *usedinfo)
	}
	_, ok = messages["SettleinfosTradeyearmonth"]

	if ok {
		settleinfo.Ccy = messages["SettleinfosCcy"]
		settleinfo.Tradeyearmonth = messages["SettleinfosTradeyearmonth"]
		settleinfo.Settleamount = messages["SettleinfosSettleamount"]
		tx.Settleinfos = append(tx.Settleinfos, *settleinfo)
	}
	_, ok = messages["OrderinfosTradeyearmonth"]
	if ok {

		orderinfo.Tradeyearmonth = messages["OrderinfosTradeyearmonth"]
		orderinfo.Orderamount = messages["OrderinfosOrderamount"]
		orderinfo.Ccy = messages["OrderinfosCcy"]
		tx.Orderinfos = append(tx.Orderinfos, *orderinfo)
	}
	_, ok = messages["ReceivableinfosTradeyearmonth"]
	if ok {

		receivableinfo.Tradeyearmonth = messages["ReceivableinfosTradeyearmonth"]
		receivableinfo.Receivableamount = messages["ReceivableinfosReceivableamount"]
		receivableinfo.Ccy = messages["ReceivableinfosCcy"]
		tx.Receivableinfos = append(tx.Receivableinfos, *receivableinfo)
	}

	return tx
}

// redis未命中的情况下，去查询数据库中的数据，这种情况只适用于指定了id的情况，如果id未指定，则直接从redis数据库中返回信息
// 将mysql查询的数据首先存入redis，然后进行二次过滤
func (s *Server) DumpHistoryTXFromMysqlToRedis(id string) {
	txs := s.SearchHistoryTXByID(id)
	s.StoreHistoryTXToRedis(txs)
}
func (s *Server) PackToHistoryTXJson(messages []*types.TransactionHistory, totalcount, currentPage int) string {
	returnresult := types.TransactionHistoryReturn{
		TransactionHistoryList: messages,
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
func (s *Server) GetMutipleHistoryTXKeys(id string) []string {
	ctx := context.Background()
	var order string
	if id == "" {
		order = "*"
	} else {
		order = id + ":" + "*"
	}

	_, keys := s.redisHistoryTX.Scan(ctx, order)
	return keys
}

// ********************************************************************************************************
// 将key值和分数存入zset集合中
func (s *Server) StoreHisTXKeyAndScoreToZset(keys []string) {
	keyscore := make(map[string]float64)
	for _, key := range keys {
		if !strings.Contains(key, ":") {
			continue
		}
		res := strings.Split(key, ":")
		t, _ := time.Parse("2006-01", res[1])
		time, _ := strconv.Atoi(t.Format("200601"))
		score := float64(time)
		keyscore[key] = score
	}
	ctx := context.Background()
	s.redisHistoryTX.ZAdd(ctx, "hisTX", keyscore)
}

// 根据下标去查询key的数据
func (s *Server) SearchHisTXKeysFromZset(ctx context.Context, start, end int64, order string) []string {

	if order == redis.Increase {
		keys := s.redisHistoryTX.ZrangeIncrease(ctx, "hisTX", start, end)
		return keys
	}
	if order == redis.Decrease {
		keys := s.redisHistoryTX.ZrangeDecrease(ctx, "hisTX", start, end)
		return keys
	}
	return nil
}
