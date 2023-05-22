package server

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/FISCO-BCOS/go-sdk/redis"
	types "github.com/FISCO-BCOS/go-sdk/type"
	"github.com/sirupsen/logrus"
)

// 存储发票信息到redis数据库中
func (s *Server) StoreInvoicesToRedis(data []*types.InvoiceInformation) {
	ctx := context.Background()
	for _, invoice := range data {
		values := make(map[string]interface{})
		key := invoice.Customerid + ":" + invoice.Checkcode + ":" + invoice.Invoicedate
		values["Certificateid"] = invoice.Certificateid
		values["Customerid"] = invoice.Customerid
		values["Corpname"] = invoice.Corpname
		values["Certificatetype"] = invoice.Certificatetype
		values["Intercustomerid"] = invoice.Intercustomerid
		values["Invoicenotaxamt"] = invoice.Invoicenotaxamt
		values["Invoiceccy"] = invoice.Invoiceccy
		values["Sellername"] = invoice.Sellername
		values["Invoicetype"] = invoice.Invoicetype
		values["Buyername"] = invoice.Buyername
		values["Buyerusccode"] = invoice.Buyerusccode
		values["Invoicedate"] = invoice.Invoicedate
		values["Sellerusccode"] = invoice.Sellerusccode
		values["Invoicecode"] = invoice.Invoicecode
		values["Invoicenum"] = invoice.Invoicenum
		values["Checkcode"] = invoice.Checkcode
		values["Invoiceamt"] = invoice.Invoiceamt
		values["Owner"] = invoice.Owner
		err := s.redisInvoice.MultipleSet(ctx, key, values)
		if err != nil {
			logrus.Errorln(err)
		}
	}
}

// 根据指令从redis中查询发票信息
func (s *Server) SearchInvoiceFromRedis(order map[string]string) ([]*types.InvoiceInformation, int) {
	pageid, err := strconv.ParseInt(order["pageid"], 10, 64)
	if err != nil {
		logrus.Errorln(err)
		return nil, 0
	}
	invoices := s.searchInvoiceByIDFromRedis(order["id"], order["searchType"])
	fmt.Println(len(invoices))
	//redis未命中
	if len(invoices) == 0 {
		//同步mysql到redis
		s.DumpInvoiceFromMysqlToRedis(order["id"])
		time.Sleep(500 * time.Millisecond)
		//二次查询
		invoices = s.searchInvoiceByIDFromRedis(order["id"], order["searchType"])
		if len(invoices) == 0 {
			return nil, 0
		}
	}
	fmt.Println(order)
	filterBytype := s.filterByInvoiceType(invoices, order["invoiceType"])
	filterByFinanceID := s.filterByFinanceID(filterBytype, order["financeID"])
	filterByTime := s.filterByInvoiceTimeStamp(filterByFinanceID, order["time"])
	filterByNum := s.filterByInvoiceNum(filterByTime, order["num"])
	filterByPageId := s.filterByInvoicePageId(filterByNum, pageid)
	totalcount := len(filterByNum)
	return filterByPageId, totalcount
}

// 根据id的信息从redis中查询数据，如果结构体是空的，那么说明redis未命中，需要去mysql数据库中查询
func (s *Server) searchInvoiceByIDFromRedis(id string, order string) []*types.InvoiceInformation {
	ctx := context.Background()
	invoices := make([]*types.InvoiceInformation, 0)
	keys := s.GetMutipleInvoiceKeys(id)
	s.redisInvoice.Del(ctx, "invoice")
	s.StoreInvoiceKeyAndScoreToZset(keys)
	//如果redis未命中,返回空的结构体
	if len(keys) == 0 {
		return nil
	}
	for _, key := range keys {
		resmap, err := s.redisInvoice.GetAll(ctx, key)
		if err != nil {
			logrus.Errorln("err", err)
			continue
		}
		invoice := packToInvoiceStruct(resmap)
		invoices = append(invoices, invoice)
	}
	return invoices
}

// 根据发票的时间戳进行过滤，调用此函数前，需要先通过id进行第一次检索
func (s *Server) filterByInvoiceType(messages []*types.InvoiceInformation, invocietype string) []*types.InvoiceInformation {
	if invocietype == "" {
		return messages
	}
	result := make([]*types.InvoiceInformation, 0)
	for _, message := range messages {
		if message.Invoicetype == invocietype {
			result = append(result, message)
		}
	}
	return result
}
func (s *Server) filterByFinanceID(messages []*types.InvoiceInformation, financeID string) []*types.InvoiceInformation {
	if financeID == "" {
		return messages
	}
	result := make([]*types.InvoiceInformation, 0)
	for _, message := range messages {
		if message.Owner == financeID {
			result = append(result, message)
		}
	}
	return result
}

// 根据发票的时间戳进行过滤，调用此函数前，需要先通过id、发票类型检索
func (s *Server) filterByInvoiceTimeStamp(messages []*types.InvoiceInformation, invoiceTimeStamp string) []*types.InvoiceInformation {
	// 根据前端界面，时间只精确到年月，时间选择框选定的应该是一个时间范围，即“XXX年X月至XXX年X月”
	// 最后在url中，对于时间的搜索应该是“?time=xxx-xxtoxxx-xx”
	// 也有可能是?time=toxxx-xx,即查看某时间之前的所有信息
	// 也有可能是?time=xxx-xxto,即查看某时间之后的所有信息
	// 先根据to分割时间，得到起始时间和结束时间，再根据时间是否位于该时间段内将结果加入返回的结果
	// 对于不含?time的url，默认全部返回

	if invoiceTimeStamp == "" {
		return messages
	}
	fields := strings.Split(invoiceTimeStamp, "to")
	if fields[0] == "" {
		fields[0] = "1970-01"
	}
	fields[0] = fields[0] + "-01"
	if fields[1] == "" {
		fields[1] = "2100-01"
	}
	fields[1] = fields[1] + "-01"
	beginData, _ := time.Parse("2006-01-02", fields[0])
	endData, _ := time.Parse("2006-01-02", fields[1])
	beginData = beginData.AddDate(0, 0, -1)
	endData = endData.AddDate(0, 1, 0)
	result := make([]*types.InvoiceInformation, 0)
	for _, message := range messages {
		checkInvoiceData, _ := time.Parse("2006-01-02", message.Invoicedate)
		if checkInvoiceData.After(beginData) && checkInvoiceData.Before(endData) {
			result = append(result, message)
		}
	}
	return result
}

func (s *Server) filterByInvoiceNum(messages []*types.InvoiceInformation, invoicenum string) []*types.InvoiceInformation {
	if invoicenum == "" {
		return messages
	}
	result := make([]*types.InvoiceInformation, 0)
	for _, message := range messages {
		if message.Invoicenum == invoicenum {
			result = append(result, message)
		}
	}
	return result
}

func (s *Server) filterByInvoicePageId(messages []*types.InvoiceInformation, invoicepageid int64) []*types.InvoiceInformation {
	start := (invoicepageid - 1) * 10
	end := invoicepageid * 10
	result := make([]*types.InvoiceInformation, 0)
	for i := start; i < end; i++ {
		if i >= int64(len(messages)) {
			break
		}
		result = append(result, messages[i])
	}
	return result
}

// 将从redis查询出来的数据转换成结构体
func packToInvoiceStruct(message map[string]string) *types.InvoiceInformation {
	invoice := new(types.InvoiceInformation)
	invoice.Certificateid = message["Certificateid"]
	invoice.Customerid = message["Customerid"]
	invoice.Corpname = message["Corpname"]
	invoice.Certificatetype = message["Certificatetype"]
	invoice.Intercustomerid = message["Intercustomerid"]
	invoice.Invoicenotaxamt = message["Invoicenotaxamt"]
	invoice.Invoiceccy = message["Invoiceccy"]
	invoice.Sellername = message["Sellername"]
	invoice.Invoicetype = message["Invoicetype"]
	invoice.Buyername = message["Buyername"]
	invoice.Buyerusccode = message["Buyerusccode"]
	invoice.Invoicedate = message["Invoicedate"]
	invoice.Sellerusccode = message["Sellerusccode"]
	invoice.Invoicecode = message["Invoicecode"]
	invoice.Invoicenum = message["Invoicenum"]
	invoice.Checkcode = message["Checkcode"]
	invoice.Invoiceamt = message["Invoiceamt"]
	invoice.Owner = message["Owner"]
	return invoice

}

// redis未命中的情况下，去查询数据库中的数据，这种情况只适用于指定了id的情况，如果id未指定，则直接从redis数据库中返回信息
// 将mysql查询的数据首先存入redis，然后进行二次过滤
func (s *Server) DumpInvoiceFromMysqlToRedis(id string) {
	plaintext := s.sql.QueryInvoiceInformation(id)
	invoices := s.sql.InvoiceinfoToMap(plaintext)
	s.StoreInvoicesToRedis(invoices)
}

func (s *Server) PackToInvoiceJson(messages []*types.InvoiceInformation, totalcount, currentPage int) string {
	returnresult := types.InvoiceInformationReturn{
		InvoiceInformationList: messages,
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
func (s *Server) GetMutipleInvoiceKeys(id string) []string {
	ctx := context.Background()
	var order string
	if id == "" {
		order = id + "*"
	} else {
		order = id + ":" + "*"
	}

	_, keys := s.redisInvoice.Scan(ctx, order)
	return keys
}

// ********************************************************************************************************
// 将key值和分数存入zset集合中
func (s *Server) StoreInvoiceKeyAndScoreToZset(keys []string) {
	keyscore := make(map[string]float64)
	for _, key := range keys {
		if !strings.Contains(key, ":") {
			continue
		}
		res := strings.Split(key, ":")

		high, err := strconv.ParseFloat(res[0], 64)
		if err != nil {
			// logrus.Errorln(err)
			continue
		}
		low, err := strconv.ParseFloat(res[1], 64)
		if err != nil {
			// logrus.Errorln(err)
			continue
		}
		score := high*1000000 + low
		keyscore[key] = score
	}
	ctx := context.Background()
	s.redisInvoice.ZAdd(ctx, "invoice", keyscore)
}

// 根据下标去查询key的数据
func (s *Server) SearchInvoiceKeysFromZset(ctx context.Context, start, end int64, order string) []string {

	if order == redis.Increase {
		keys := s.redisInvoice.ZrangeIncrease(ctx, "invoice", start, end)
		return keys
	}
	if order == redis.Decrease {
		keys := s.redisInvoice.ZrangeDecrease(ctx, "invoice", start, end)
		return keys
	}
	return nil
}
