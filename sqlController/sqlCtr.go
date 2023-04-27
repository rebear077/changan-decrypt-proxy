package sql

import (
	"database/sql"
	"net/http"

	"github.com/FISCO-BCOS/go-sdk/conf"
	types "github.com/FISCO-BCOS/go-sdk/type"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

type SqlCtr struct {
	db        *sql.DB
	Decrypter *Decrypter
}

const (
	//发票查询指令
	InvoiceSQLQueryAll    = "select * from u_t_invoice_information1"
	InvoiceSQLQueryByID   = "select * from u_t_invoice_information1 where id ="
	InvoiceSQLQueryBy_ID  = "select * from u_t_invoice_information1 where _id_= "
	InvoiceSQLQueryLength = "select count(*) from u_t_invoice_information1"
	// 历史交易信息-used信息
	HistoricalUsedSQLQueryAll   = "select * from u_t_history_used_information2"
	HistoricalUsedSQLQueryByID  = "select * from u_t_history_used_information2 where id = "
	HistoricalUsedSQLQueryBy_ID = "select * from u_t_history_used_information2 where _id_= "
	// 历史交易信息-settle信息
	HistoricalSettleSQLQueryAll   = "select * from u_t_history_settle_information2"
	HistoricalSettleSQLQueryByID  = "select * from u_t_history_settle_information2 where id = "
	HistoricalSettleSQLQueryBy_ID = "select * from u_t_history_settle_information2 where _id_= "
	//历史交易信息-order信息
	HistoricalOrderSQLQueryAll   = "select * from u_t_history_order_information2"
	HistoricalOrderSQLQueryByID  = "select * from u_t_history_order_information2 where id = "
	HistoricalOrderSQLQueryBy_ID = "select * from u_t_history_order_information2 where _id_= "
	//历史交易信息-receivalble信息
	HistoricalReceivableSQLQueryAll   = "select * from u_t_history_receivable_information2"
	HistoricalReceivableSQLQueryByID  = "select * from u_t_history_receivable_information2 where id = "
	HistoricalReceivableSQLQueryBy_ID = "select * from u_t_history_receivable_information2 where _id_= "
	// 入池数据-plan信息
	EnterPoolPlanSQLQueryAll   = "select * from u_t_pool_plan_information1"
	EnterPoolPlanSQLQueryByID  = "select * from u_t_pool_plan_information1 where id = "
	EnterPoolPlanSQLQueryBy_ID = "select * from u_t_pool_plan_information1 where _id_= "
	//入池数据-used信息
	EnterPoolUsedSQLQueryAll   = "select * from u_t_pool_used_information1"
	EnterPoolUsedSQLQueryByID  = "select * from u_t_pool_used_information1 where id = "
	EnterPoolUsedSQLQueryBy_ID = "select * from u_t_pool_used_information1 where _id_= "
	//融资意向信息
	FinancingSQLQueryAll   = "select * from u_t_supplier_financing_application1"
	FinancingSQLQueryByID  = "select * from u_t_supplier_financing_application1 where id = "
	FinancingSQLQueryBy_ID = "select * from u_t_supplier_financing_application1 where _id_= "
	//回款账户信息
	AccountsSQLQueryAll   = "select * from u_t_push_payment_accounts1"
	AccountsSQLQueryByID  = "select * from u_t_push_payment_accounts1 where id = "
	AccountsSQLQueryBy_ID = "select * from u_t_push_payment_accounts1 where _id_= "
	// 借贷合同信息
	FinancingContractSQLAll   = "select * from u_t_finance_contract1"
	FinancingContractSQLByID  = "select * from u_t_finance_contract1 where id="
	FinancingContractSQLBy_ID = "select * from u_t_finance_contract1 where _id_="
)

func NewSqlCtr() *SqlCtr {
	configs, err := conf.ParseConfigFile("./configs/config.toml")
	if err != nil {
		logrus.Fatalln(err)
	}
	config := &configs[0]
	// db, err := sql.Open("mysql", "root:123456@/db_node0")
	str := config.MslUsername + ":" + config.MslPasswd + "@/" + config.MslName
	db, err := sql.Open("mysql", str)
	if err != nil {
		logrus.Fatalln(err)
	}
	de := NewDecrypter()
	return &SqlCtr{
		db:        db,
		Decrypter: de,
	}
}

// //////////////////////////////////////////////////////////////////////////////////////////////////
// 解析前端发来的URL请求，获取检索条件，结构体形式返回
// 发票信息
func (s *SqlCtr) InvoiceInformationIndex(request *http.Request) *types.InvoiceInformationSearch {
	query := request.URL.Query()
	id := ""
	time := ""
	invoiceType := ""
	invoiceNum := ""
	searchType := "increase"
	pageid := "1"
	if len(query["id"]) > 0 {
		id = query["id"][0]
	}
	if len(query["time"]) > 0 {
		time = query["time"][0]
	}
	if len(query["invoiceType"]) > 0 {
		invoiceType = query["invoiceType"][0]
	}
	if len(query["invoiceNum"]) > 0 {
		invoiceNum = query["invoiceNum"][0]
	}
	if len(query["searchType"]) > 0 {
		searchType = query["searchType"][0]
	}
	if len(query["pageid"]) > 0 {
		pageid = query["pageid"][0]
	}
	index := types.InvoiceInformationSearch{
		Id:          id,
		Time:        time,
		InvoiceType: invoiceType,
		InvoiceNum:  invoiceNum,
		PageId:      pageid,
		SearchType:  searchType,
	}
	return &index
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 解析历史交易信息URL
func (s *SqlCtr) HistoryTransactionIndex(request *http.Request) *types.HistoryTransactionSearch {
	query := request.URL.Query()
	id := ""
	tradeyearmonth := ""
	financeid := ""
	pageid := "1"
	searchType := "increase"
	if len(query["id"]) > 0 {
		id = query["id"][0]
	}
	if len(query["tradeyearmonth"]) > 0 {
		tradeyearmonth = query["tradeyearmonth"][0]
	}
	if len(query["financeid"]) > 0 {
		financeid = query["financeid"][0]
	}
	if len(query["pageid"]) > 0 {
		pageid = query["pageid"][0]
	}
	if len(query["searchType"]) > 0 {
		searchType = query["searchType"][0]
	}
	index := types.HistoryTransactionSearch{
		Id:             id,
		Tradeyearmonth: tradeyearmonth,
		FinanceId:      financeid,
		PageId:         pageid,
		SearchType:     searchType,
	}
	return &index
}
func (s *SqlCtr) PoolDataIndex(request *http.Request) *types.PoolDataSearch {
	query := request.URL.Query()
	id := ""
	pageid := "1"
	tradeyearmonth := ""
	searchType := "increase"
	if len(query["id"]) > 0 {
		id = query["id"][0]
	}
	if len(query["tradeyearmonth"]) > 0 {
		tradeyearmonth = query["tradeyearmonth"][0]
	}
	if len(query["pageid"]) > 0 {
		pageid = query["pageid"][0]
	}
	if len(query["searchType"]) > 0 {
		searchType = query["searchType"][0]
	}
	index := types.PoolDataSearch{
		Id:             id,
		PageId:         pageid,
		SearchType:     searchType,
		Tradeyearmonth: tradeyearmonth,
	}
	return &index
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 输入参数是解密后的发票信息，转换成redis存储所需要的数据结构
func (s *SqlCtr) InvoiceinfoToMap(ret []string) []*types.InvoiceInformation {

	return handleInvoiceInfo(ret)
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (s *SqlCtr) FinancingIntentionIndex(request *http.Request) *types.FinancingIntentionSearch {
	query := request.URL.Query()
	id := ""
	financeid := ""
	pageid := "1"
	searchType := "increase"
	if len(query["id"]) > 0 {
		id = query["id"][0]
	}
	if len(query["financeid"]) > 0 {
		financeid = query["financeid"][0]
	}
	if len(query["pageid"]) > 0 {
		pageid = query["pageid"][0]
	}
	if len(query["searchType"]) > 0 {
		searchType = query["searchType"][0]
	}
	index := types.FinancingIntentionSearch{
		Id:         id,
		FinanceId:  financeid,
		PageId:     pageid,
		SearchType: searchType,
	}
	return &index
}
func (s *SqlCtr) IntensioninfoToMap(ret []string) []*types.FinancingIntention {

	return handleFinancingIntention(ret)
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (s *SqlCtr) CollectionAccountIndex(request *http.Request) *types.CollectionAccountSearch {
	query := request.URL.Query()
	id := ""
	pageid := "1"
	searchType := "increase"
	if len(query["id"]) > 0 {
		id = query["id"][0]
	}
	if len(query["pageid"]) > 0 {
		pageid = query["pageid"][0]
	}
	if len(query["searchType"]) > 0 {
		searchType = query["searchType"][0]
	}
	index := types.CollectionAccountSearch{
		Id:         id,
		PageId:     pageid,
		SearchType: searchType,
	}
	return &index
}
func (s *SqlCtr) AccountinfoToMap(ret []string) []*types.CollectionAccount {

	return handleCollectionAccount(ret)
}

// 查询mysql数据库中加密后的发票信息，如果id为空，则查找全部的信息
func (s *SqlCtr) QueryInvoiceInformation(id string) []string {
	var ret []string
	if id == "" {
		ret, _ = s.QueryInvoiceByOrder(InvoiceSQLQueryAll)
	} else {
		ret, _ = s.QueryInvoiceByOrder(InvoiceSQLQueryByID + id)
	}
	return ret
}
func (s *SqlCtr) QueryInvoiceInforsBySQLID(_id_ string) []string {
	var ret []string
	ret, _ = s.QueryInvoiceByOrder(InvoiceSQLQueryBy_ID + _id_)
	return ret
}
func (s *SqlCtr) QueryInvoiceInformationLength() int {
	var length int
	err := s.db.QueryRow(InvoiceSQLQueryLength).Scan(&length)
	if err != nil {
		logrus.Fatalln(err)
	}
	return length
}

// //////////////////////////////////////////////////////////////////////////////////////////////////////
func (s *SqlCtr) QueryHistoricalTransUsedInfos(id string) []string {
	var ret []string
	if id == "" {
		ret, _ = s.QueryHistoricalTransByOrder(HistoricalUsedSQLQueryAll)
	} else {
		ret, _ = s.QueryHistoricalTransByOrder(HistoricalUsedSQLQueryByID + id)
	}
	return ret
}
func (s *SqlCtr) QueryHistoricalTransUsedInfosBySQLID(_id_ string) []string {
	var ret []string
	ret, _ = s.QueryHistoricalTransByOrder(HistoricalUsedSQLQueryBy_ID + _id_)
	return ret
}

func (s *SqlCtr) QueryHistoricalTransSettleInfos(id string) []string {
	var ret []string
	if id == "" {
		ret, _ = s.QueryHistoricalTransByOrder(HistoricalSettleSQLQueryAll)
	} else {
		ret, _ = s.QueryHistoricalTransByOrder(HistoricalSettleSQLQueryByID + id)
	}
	return ret
}
func (s *SqlCtr) QueryHistoricalTransSettleInfosBySQLID(_id_ string) []string {
	var ret []string
	ret, _ = s.QueryHistoricalTransByOrder(HistoricalSettleSQLQueryBy_ID + _id_)
	return ret
}
func (s *SqlCtr) QueryHistoricalTransOrderInfos(id string) []string {
	var ret []string
	if id == "" {
		ret, _ = s.QueryHistoricalTransByOrder(HistoricalOrderSQLQueryAll)
	} else {
		ret, _ = s.QueryHistoricalTransByOrder(HistoricalOrderSQLQueryByID + id)
	}
	return ret
}
func (s *SqlCtr) QueryHistoricalTransOrderInfosBySQLID(_id_ string) []string {
	var ret []string
	ret, _ = s.QueryHistoricalTransByOrder(HistoricalOrderSQLQueryBy_ID + _id_)
	return ret
}

func (s *SqlCtr) QueryHistoricalTransReceivableInfos(id string) []string {
	var ret []string
	if id == "" {
		ret, _ = s.QueryHistoricalTransByOrder(HistoricalReceivableSQLQueryAll)
	} else {
		ret, _ = s.QueryHistoricalTransByOrder(HistoricalReceivableSQLQueryByID + id)
	}
	return ret
}
func (s *SqlCtr) QueryHistoricalTransReceivableInfosBySQLID(_id_ string) []string {
	var ret []string
	ret, _ = s.QueryHistoricalTransByOrder(HistoricalReceivableSQLQueryBy_ID + _id_)
	return ret
}

// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (s *SqlCtr) QueryEnterpoolDataPlanInfos(id string) []string {
	var ret []string
	if id == "" {
		ret, _ = s.QueryPoolDataByOrder(EnterPoolPlanSQLQueryAll)
	} else {
		ret, _ = s.QueryPoolDataByOrder(EnterPoolPlanSQLQueryByID + id)
	}
	return ret
}
func (s *SqlCtr) QueryEnterpoolDataPlanInfosBySQLID(_id_ string) []string {
	var ret []string
	ret, _ = s.QueryPoolDataByOrder(EnterPoolPlanSQLQueryBy_ID + _id_)
	return ret
}
func (s *SqlCtr) QueryEnterpoolDataUsedInfos(id string) []string {
	var ret []string
	if id == "" {
		ret, _ = s.QueryPoolDataByOrder(EnterPoolUsedSQLQueryAll)
	} else {
		ret, _ = s.QueryPoolDataByOrder(EnterPoolUsedSQLQueryByID + id)
	}
	return ret
}
func (s *SqlCtr) QueryEnterpoolDataUsedInfosBySQLID(_id_ string) []string {
	var ret []string
	ret, _ = s.QueryPoolDataByOrder(EnterPoolUsedSQLQueryBy_ID + _id_)
	return ret
}

// //////////////////////////////////////////////////////////////////////////////////////////////////////
// 查询mysql数据库中融资意向信息，如果id为空，则查找全部的信息
func (s *SqlCtr) QueryFinancingIntention(id string) []string {
	var ret []string
	if id == "" {
		ret, _ = s.QueryFinancingByOrder(FinancingSQLQueryAll)
	} else {
		ret, _ = s.QueryFinancingByOrder(FinancingSQLQueryByID + id)
	}
	return ret
}
func (s *SqlCtr) QueryFinancingIntentionBySQLID(_id_ string) []string {
	var ret []string
	ret, _ = s.QueryFinancingByOrder(FinancingSQLQueryBy_ID + _id_)
	return ret
}

// //////////////////////////////////////////////////////////////////////////////
// 查询mysql数据库中回款账户信息，如果id为空，则查找全部的信息
func (s *SqlCtr) QueryCollectionAccount(id string) []string {
	var ret []string
	if id == "" {
		ret, _ = s.QueryAccountsByOrder(AccountsSQLQueryAll)
	} else {
		ret, _ = s.QueryAccountsByOrder(AccountsSQLQueryByID + id)
	}
	return ret
}
func (s *SqlCtr) QueryCollectionAccountBySQLID(_id_ string) []string {
	var ret []string
	ret, _ = s.QueryAccountsByOrder(AccountsSQLQueryBy_ID + _id_)
	return ret
}

// //////////////////////////////////////////////////////////////////////////////////
func (s *SqlCtr) QueryFinancingContract(id string) []types.RawFinancingContractData {
	var ret []types.RawFinancingContractData
	if id == "" {
		ret, _ = s.QueryFinancingContractByOrder(FinancingContractSQLAll)
	} else {
		ret, _ = s.QueryFinancingContractByOrder(AccountsSQLQueryByID + id)
	}
	return ret
}
func (s *SqlCtr) QueryFinancingContractBySQLID(_id_ string) []types.RawFinancingContractData {
	var ret []types.RawFinancingContractData
	ret, _ = s.QueryFinancingContractByOrder(AccountsSQLQueryBy_ID + _id_)
	return ret
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////
func (s *SqlCtr) QueryInvoiceByOrder(order string) ([]string, error) {
	in_stmt, err := s.db.Prepare(order)
	if err != nil {
		return nil, err
	}
	rows, err := in_stmt.Query()
	if err != nil {
		return nil, err
	}
	ret := make([]string, 0)
	count := 0
	i := 0
	for rows.Next() {
		record := &types.RawInvoiceData{}
		err = rows.Scan(&record.SQLId, &record.Num, &record.Status, &record.ID, &record.CustomerID, &record.Time, &record.Type, &record.InvoiceNum, &record.Data, &record.Key, &record.Hash, &record.Owner)
		if err != nil {
			logrus.Errorln(err)
			count++
			continue
		}
		symkey, err := s.Decrypter.DecryptSymkey([]byte(record.Key))
		if err != nil {
			logrus.Errorln("利用私钥解密对称密钥失败")
		}
		data, err := s.Decrypter.DecryptData(record.Data, symkey)
		if err != nil {
			logrus.Errorln("利用对称密钥解密数据失败")
		}
		if s.Decrypter.ValidateHash([]byte(record.Hash), data) {
			ret = append(ret, string(data))
		} else {
			logrus.Errorln("哈希值与数据对应错误")
		}
		i = i + 1
	}
	return ret, nil
}
func (s *SqlCtr) QueryHistoricalTransByOrder(order string) ([]string, error) {
	in_stmt, err := s.db.Prepare(order)
	if err != nil {
		return nil, err
	}
	rows, err := in_stmt.Query()
	if err != nil {
		return nil, err
	}
	ret := make([]string, 0)
	count := 0
	i := 0
	for rows.Next() {
		record := &types.RawHistoryTransData{}
		err = rows.Scan(&record.SqlId, &record.Num, &record.Status, &record.Id, &record.CustomerID, &record.TradeYearMonth, &record.FinanceId, &record.Data, &record.Key, &record.Hash, &record.Owner)
		if err != nil {
			logrus.Errorln(err)
			count++
			continue
		}
		symkey, err := s.Decrypter.DecryptSymkey([]byte(record.Key))
		if err != nil {
			logrus.Infof("利用私钥解密对称密钥失败")
		}
		data, err := s.Decrypter.DecryptData(record.Data, symkey)
		if err != nil {
			logrus.Infof("利用对称密钥解密数据失败")
		}
		if s.Decrypter.ValidateHash([]byte(record.Hash), data) {
			ret = append(ret, string(data))
		} else {
			logrus.Infof("哈希值与数据对应错误")
		}
		i = i + 1
	}
	return ret, nil
}

func (s *SqlCtr) QueryPoolDataByOrder(order string) ([]string, error) {
	in_stmt, err := s.db.Prepare(order)
	if err != nil {
		return nil, err
	}
	rows, err := in_stmt.Query()
	if err != nil {
		return nil, err
	}
	ret := make([]string, 0)
	count := 0
	i := 0
	for rows.Next() {
		record := &types.RawEnterPoolData{}
		err = rows.Scan(&record.SqlId, &record.Num, &record.Status, &record.Id, &record.CustomerID, &record.TradeYearMonth, &record.Data, &record.Key, &record.Hash, &record.Owner)
		if err != nil {
			logrus.Errorln(err)
			count++
			continue
		}
		symkey, err := s.Decrypter.DecryptSymkey([]byte(record.Key))
		if err != nil {
			logrus.Infof("利用私钥解密对称密钥失败")
		}
		data, err := s.Decrypter.DecryptData(record.Data, symkey)
		if err != nil {
			logrus.Infof("利用对称密钥解密数据失败")
		}
		if s.Decrypter.ValidateHash([]byte(record.Hash), data) {
			ret = append(ret, string(data))
		} else {
			logrus.Infof("哈希值与数据对应错误")
		}
		i = i + 1
	}
	return ret, nil
}

func (s *SqlCtr) QueryFinancingByOrder(order string) ([]string, error) {
	in_stmt, err := s.db.Prepare(order)
	if err != nil {
		return nil, err
	}
	rows, err := in_stmt.Query()
	if err != nil {
		return nil, err
	}
	ret := make([]string, 0)
	count := 0
	i := 0
	for rows.Next() {
		record := &types.RawFinancingData{}
		err = rows.Scan(&record.SQLId, &record.Num, &record.Status, &record.ID, &record.FinanceId, &record.Data, &record.Key, &record.Hash)
		if err != nil {
			logrus.Errorln(err)
			count++
			continue
		}
		symkey, err := s.Decrypter.DecryptSymkey([]byte(record.Key))
		if err != nil {
			logrus.Errorln("利用私钥解密对称密钥失败")
		}
		data, err := s.Decrypter.DecryptData(record.Data, symkey)
		if err != nil {
			logrus.Errorln("利用对称密钥解密数据失败")
		}
		if s.Decrypter.ValidateHash([]byte(record.Hash), data) {
			ret = append(ret, string(data))
		} else {
			logrus.Errorln("哈希值与数据对应错误")
		}
		i = i + 1
	}
	return ret, nil
}

// 输入命令，比如“select * from u_t_push_payment_accounts”,查询出加密后的密文然后自动解密，返回明文[]string
func (s *SqlCtr) QueryAccountsByOrder(order string) ([]string, error) {
	in_stmt, err := s.db.Prepare(order)
	if err != nil {
		return nil, err
	}
	rows, err := in_stmt.Query()
	if err != nil {
		return nil, err
	}
	ret := make([]string, 0)
	count := 0
	i := 0
	for rows.Next() {
		record := &types.RawAccountsData{}
		err = rows.Scan(&record.SQLId, &record.Num, &record.Status, &record.ID, &record.Data, &record.Key, &record.Hash)
		if err != nil {
			logrus.Errorln(err)
			count++
			continue
		}
		symkey, err := s.Decrypter.DecryptSymkey([]byte(record.Key))
		if err != nil {
			logrus.Errorln("利用私钥解密对称密钥失败")
		}
		data, err := s.Decrypter.DecryptData(record.Data, symkey)
		if err != nil {
			logrus.Errorln("利用对称密钥解密数据失败")
		}
		if s.Decrypter.ValidateHash([]byte(record.Hash), data) {
			ret = append(ret, string(data))
		} else {
			logrus.Errorln("哈希值与数据对应错误")
		}
		i = i + 1
	}
	return ret, nil
}

// 输入命令，比如“select * from u_t_push_payment_accounts”,
func (s *SqlCtr) QueryFinancingContractByOrder(order string) ([]types.RawFinancingContractData, error) {
	in_stmt, err := s.db.Prepare(order)
	if err != nil {
		return nil, err
	}
	rows, err := in_stmt.Query()
	if err != nil {
		return nil, err
	}
	ret := make([]types.RawFinancingContractData, 0)
	count := 0
	for rows.Next() {
		record := &types.RawFinancingContractData{}
		err = rows.Scan(&record.SQLId, &record.Num, &record.Status, &record.ID, &record.FinancingID, &record.CustomerID, &record.CorpName, &record.DebtMoney, &record.SupplyDate, &record.ExpireDate, &record.Balance)
		if err != nil {
			logrus.Errorln(err)
			count++
			continue
		}
		ret = append(ret, *record)
	}
	return ret, nil
}
