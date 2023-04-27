package server

import (
	"context"
	"fmt"
	"time"

	"github.com/FISCO-BCOS/go-sdk/canal"
	"github.com/FISCO-BCOS/go-sdk/conf"
	"github.com/FISCO-BCOS/go-sdk/redis"
	sql "github.com/FISCO-BCOS/go-sdk/sqlController"
	types "github.com/FISCO-BCOS/go-sdk/type"
	"github.com/sirupsen/logrus"
)

type Server struct {
	sql                     *sql.SqlCtr
	redisInvoice            *redis.RedisOperator
	redisHistoryTX          *redis.RedisOperator
	redisEnterPool          *redis.RedisOperator
	redisFinancingIntention *redis.RedisOperator
	redisCollectionAccount  *redis.RedisOperator
	redisFinancingContract  *redis.RedisOperator
	canal                   *canal.Connector
}

// 初始化
func NewServer() *Server {
	configs, err := conf.ParseConfigFile("./configs/config.toml")
	if err != nil {
		logrus.Fatalln(err)
	}
	config := &configs[0]
	sqlCtr := sql.NewSqlCtr()
	invoice := redis.NewRedisOperator(0)
	historyTX := redis.NewRedisOperator(1)
	enterpool := redis.NewRedisOperator(2)
	financingIntention := redis.NewRedisOperator(3)
	collectAccount := redis.NewRedisOperator(4)
	finangcingContract := redis.NewRedisOperator(5)
	// canal := canal.NewConnector("db0\\.u_t.*")
	canal := canal.NewConnector(config.CanalConnectedDB + "\\.u_t.*")
	return &Server{
		sql:                     sqlCtr,
		redisInvoice:            invoice,
		redisHistoryTX:          historyTX,
		redisEnterPool:          enterpool,
		redisFinancingIntention: financingIntention,
		redisCollectionAccount:  collectAccount,
		redisFinancingContract:  finangcingContract,
		canal:                   canal,
	}
}

// 将redis的数据全部删除后，与mysql中的数据进行同步
func (s *Server) ForceSynchronous() {
	//首先清除redis数据
	ctx := context.Background()
	s.redisInvoice.FlushData(ctx)
	//获取mysql中的数据,获取的是解密后的明文
	//同步发票信息
	plaintext := s.sql.QueryInvoiceInformation("")
	invoices := s.sql.InvoiceinfoToMap(plaintext)
	// for _, invoice := range invoices {
	// 	fmt.Println(invoice)
	// }
	s.StoreInvoicesToRedis(invoices)
	//同步历史交易信息
	txs := s.SearchHistoryTXByID("")
	s.StoreHistoryTXToRedis(txs)
	//同步入池数据信息
	enterpools := s.SearchEnterPoolByID("")
	s.StoreEnterPoolToRedis(enterpools)
	//同步融资意向信息
	plaintextIntension := s.sql.QueryFinancingIntention("")
	intensions := s.sql.IntensioninfoToMap(plaintextIntension)
	s.StoreFinanacingIntensionToRedis(intensions)
	//同步回款账户信息
	plaintextAccounts := s.sql.QueryCollectionAccount("")
	accounts := s.sql.AccountinfoToMap(plaintextAccounts)
	s.StoreAccountToRedis(accounts)

}

// ********************************************************************************************
// cannal操作
// 从canal中进行同步消息到redis中
func (s *Server) DumpFromCanal() {
	go s.canal.Start()
	for {
		time.Sleep(3 * time.Second)
		s.canal.Lock.Lock()
		fmt.Println(s.canal.RawData)
		if len(s.canal.RawData[canal.Invoice]) != 0 {
			messages := make([]*types.RawCanalData, 0)
			messages = append(messages, s.canal.RawData[canal.Invoice]...)
			delete(s.canal.RawData, canal.Invoice)
			s.CannalStoreInvoiceToredis(messages)
		}
		if len(s.canal.RawData[canal.Intension]) != 0 {
			messages := make([]*types.RawCanalData, 0)
			messages = append(messages, s.canal.RawData[canal.Intension]...)
			delete(s.canal.RawData, canal.Intension)
			s.CannalStoreIntensionToredis(messages)
		}
		if len(s.canal.RawData[canal.Accounts]) != 0 {
			messages := make([]*types.RawCanalData, 0)
			messages = append(messages, s.canal.RawData[canal.Accounts]...)
			delete(s.canal.RawData, canal.Accounts)
			s.CannalStoreAccountsToredis(messages)
		}
		if len(s.canal.RawData[canal.HisOrder]) != 0 {
			messages := make([]*types.RawCanalData, 0)
			messages = append(messages, s.canal.RawData[canal.HisOrder]...)
			delete(s.canal.RawData, canal.HisOrder)
			s.CannalStoreHisOrderToredis(messages)
		}
		if len(s.canal.RawData[canal.HisReceivable]) != 0 {
			messages := make([]*types.RawCanalData, 0)
			messages = append(messages, s.canal.RawData[canal.HisReceivable]...)
			delete(s.canal.RawData, canal.HisReceivable)
			s.CannalStoreHisReceivableToredis(messages)
		}
		if len(s.canal.RawData[canal.HisSettle]) != 0 {
			messages := make([]*types.RawCanalData, 0)
			messages = append(messages, s.canal.RawData[canal.HisSettle]...)
			delete(s.canal.RawData, canal.HisSettle)
			s.CannalStoreHisSettleToredis(messages)
		}
		if len(s.canal.RawData[canal.HisUsed]) != 0 {
			messages := make([]*types.RawCanalData, 0)
			messages = append(messages, s.canal.RawData[canal.HisUsed]...)
			delete(s.canal.RawData, canal.HisUsed)
			s.CannalStoreHisUsedToredis(messages)
		}
		if len(s.canal.RawData[canal.PoolPlan]) != 0 {
			messages := make([]*types.RawCanalData, 0)
			messages = append(messages, s.canal.RawData[canal.PoolPlan]...)
			delete(s.canal.RawData, canal.PoolPlan)
			s.CannalStoreEnterPoolPlanToredis(messages)
		}
		if len(s.canal.RawData[canal.PoolUsed]) != 0 {
			messages := make([]*types.RawCanalData, 0)
			messages = append(messages, s.canal.RawData[canal.PoolUsed]...)
			delete(s.canal.RawData, canal.PoolUsed)
			s.CannalStoreEnterPoolUsedToredis(messages)
		}
		s.canal.Lock.Unlock()
	}
}

// 从数据库原始的数据，先解密，然后转换格式后存入redis中
func (s *Server) CannalStoreInvoiceToredis(datas []*types.RawCanalData) {
	for _, data := range datas {
		plaintext := s.sql.QueryInvoiceInforsBySQLID(string(data.SQLId))
		invoices := s.sql.InvoiceinfoToMap(plaintext)
		s.StoreInvoicesToRedis(invoices)
	}
}
func (s *Server) CannalStoreAccountsToredis(datas []*types.RawCanalData) {
	for _, data := range datas {
		plaintext := s.sql.QueryCollectionAccountBySQLID(string(data.SQLId))
		accounts := s.sql.AccountinfoToMap(plaintext)
		s.StoreAccountToRedis(accounts)
	}
}
func (s *Server) CannalStoreIntensionToredis(datas []*types.RawCanalData) {
	fmt.Println("CannalStoreIntensionToredis")
	for _, data := range datas {
		plaintext := s.sql.QueryFinancingIntentionBySQLID(string(data.SQLId))
		fmt.Println(plaintext)
		intensions := s.sql.IntensioninfoToMap(plaintext)
		fmt.Println(intensions)
		s.StoreFinanacingIntensionToRedis(intensions)
	}
}
func (s *Server) CannalStoreHisUsedToredis(datas []*types.RawCanalData) {
	for _, data := range datas {
		use := s.sql.QueryHistoricalTransUsedInfosBySQLID(string(data.SQLId))
		txUsed := sql.HandleHistoricaltransactionUsedinfos(use)
		for _, used := range txUsed {
			values := make(map[string]interface{})
			ctx := context.Background()
			key := used.Customerid + ":" + used.Usedinfos[0].Tradeyearmonth
			resmap, _ := s.redisHistoryTX.GetAll(ctx, key)
			fmt.Println(len(resmap))
			if len(resmap) == 0 {
				hisTX := s.SearchHistoryTXBySQLID(string(data.SQLId))
				s.StoreHistoryTXToRedis(hisTX)
				continue
			}
			values["UsedinfosTradeyearmonth"] = used.Usedinfos[0].Tradeyearmonth
			values["UsedinfosUsedamount"] = used.Usedinfos[0].Usedamount
			values["UsedinfosCcy"] = used.Usedinfos[0].Ccy
			s.redisHistoryTX.MultipleSet(ctx, key, values)
		}

	}
}
func (s *Server) CannalStoreHisSettleToredis(datas []*types.RawCanalData) {
	for _, data := range datas {
		settle := s.sql.QueryHistoricalTransSettleInfosBySQLID(string(data.SQLId))
		txSettled := sql.HandleHistoricaltransactionSettleinfos(settle)
		for _, settled := range txSettled {
			fmt.Println(settled, "--------------------")
			values := make(map[string]interface{})
			ctx := context.Background()
			key := settled.Customerid + ":" + settled.Settleinfos[0].Tradeyearmonth
			resmap, _ := s.redisHistoryTX.GetAll(ctx, key)
			fmt.Println(len(resmap))
			if len(resmap) == 0 {
				fmt.Println("resmap nil")
				hisTX := s.SearchHistoryTXBySQLID(string(data.SQLId))
				s.StoreHistoryTXToRedis(hisTX)
				continue
			}
			fmt.Println("resmap not nul")
			values["SettleinfosTradeyearmonth"] = settled.Settleinfos[0].Tradeyearmonth
			values["SettleinfosSettleamount"] = settled.Settleinfos[0].Settleamount
			values["SettleinfosCcy"] = settled.Settleinfos[0].Ccy
			s.redisHistoryTX.MultipleSet(ctx, key, values)
		}

	}
}

func (s *Server) CannalStoreHisOrderToredis(datas []*types.RawCanalData) {
	for _, data := range datas {
		order := s.sql.QueryHistoricalTransOrderInfosBySQLID(string(data.SQLId))
		txOrderd := sql.HandleHistoricaltransactionOrderinfos(order)
		for _, orderd := range txOrderd {
			fmt.Println(orderd, "--------------------")
			values := make(map[string]interface{})
			ctx := context.Background()
			key := orderd.Customerid + ":" + orderd.Orderinfos[0].Tradeyearmonth
			resmap, _ := s.redisHistoryTX.GetAll(ctx, key)
			fmt.Println(resmap)
			if len(resmap) == 0 {
				fmt.Println("resmap nul")
				hisTX := s.SearchHistoryTXBySQLID(string(data.SQLId))
				s.StoreHistoryTXToRedis(hisTX)
				continue
			}
			fmt.Println("resmap not nul")
			values["OrderinfosTradeyearmonth"] = orderd.Orderinfos[0].Tradeyearmonth
			values["OrderinfosOrderamount"] = orderd.Orderinfos[0].Orderamount
			values["OrderinfosCcy"] = orderd.Orderinfos[0].Ccy

			s.redisHistoryTX.MultipleSet(ctx, key, values)
		}

	}
}
func (s *Server) CannalStoreHisReceivableToredis(datas []*types.RawCanalData) {
	for _, data := range datas {
		receivable := s.sql.QueryHistoricalTransReceivableInfosBySQLID(string(data.SQLId))
		txReceivabled := sql.HandleHistoricaltransactionReceivableinfos(receivable)
		for _, receivabled := range txReceivabled {
			values := make(map[string]interface{})
			ctx := context.Background()
			key := receivabled.Customerid + ":" + receivabled.Receivableinfos[0].Tradeyearmonth
			resmap, _ := s.redisHistoryTX.GetAll(ctx, key)
			if len(resmap) == 0 {
				hisTX := s.SearchHistoryTXBySQLID(string(data.SQLId))
				s.StoreHistoryTXToRedis(hisTX)
				continue
			}
			values["ReceivableinfosTradeyearmonth"] = receivabled.Receivableinfos[0].Tradeyearmonth
			values["ReceivableinfosReceivableamount"] = receivabled.Receivableinfos[0].Receivableamount
			values["ReceivableinfosCcy"] = receivabled.Receivableinfos[0].Ccy
			fmt.Println(values)
			err := s.redisHistoryTX.MultipleSet(ctx, key, values)
			if err != nil {
				fmt.Println(err)
			}
		}

	}
}

// ************************************************************************************
func (s *Server) CannalStoreEnterPoolPlanToredis(datas []*types.RawCanalData) {
	for _, data := range datas {
		plan := s.sql.QueryEnterpoolDataPlanInfosBySQLID(string(data.SQLId))
		txPlan := sql.HandleEnterpoolDataPlaninfos(plan)
		for _, pland := range txPlan {
			values := make(map[string]interface{})
			ctx := context.Background()

			key := pland.Customerid + ":" + pland.Planinfos[0].Tradeyearmonth
			resmap, _ := s.redisEnterPool.GetAll(ctx, key)
			if len(resmap) == 0 {
				poolplan := s.SearchEnterPoolBySQLID(string(data.SQLId))
				s.StoreEnterPoolToRedis(poolplan)
				continue
			}
			values["PlaninfosTradeyearmonth"] = pland.Planinfos[0].Tradeyearmonth
			values["PlaninfosPlanamount"] = pland.Planinfos[0].Planamount
			values["PlaninfosCurrency"] = pland.Planinfos[0].Currency
			fmt.Println(values)
			s.redisEnterPool.MultipleSet(ctx, key, values)
		}

	}
}

func (s *Server) CannalStoreEnterPoolUsedToredis(datas []*types.RawCanalData) {
	for _, data := range datas {
		use := s.sql.QueryEnterpoolDataUsedInfosBySQLID(string(data.SQLId))
		txUsed := sql.HandleEnterpoolDataProviderusedinfos(use)
		for _, used := range txUsed {
			values := make(map[string]interface{})
			ctx := context.Background()
			key := used.Customerid + ":" + used.Providerusedinfos[0].Tradeyearmonth
			resmap, _ := s.redisEnterPool.GetAll(ctx, key)
			if len(resmap) == 0 {
				poolused := s.SearchEnterPoolBySQLID(string(data.SQLId))
				s.StoreEnterPoolToRedis(poolused)
				continue
			}
			values["ProviderusedinfosTradeyearmonth"] = used.Providerusedinfos[0].Tradeyearmonth
			values["ProviderusedinfosUsedamount"] = used.Providerusedinfos[0].Usedamount
			values["ProviderusedinfosCurrency"] = used.Providerusedinfos[0].Currency
			fmt.Println(values)
			err := s.redisEnterPool.MultipleSet(ctx, key, values)
			if err != nil {
				logrus.Errorln(err)
			}
		}

	}
}
