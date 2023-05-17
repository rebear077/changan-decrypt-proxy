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
	orders    *conf.Config
	Decrypter *Decrypter
}

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
		orders:    config,
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

// 解析入池数据URL
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

// 融资意向申请URL
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

// 回款账户YRL
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

// 借贷合同URL
func (s *SqlCtr) FinancingContractIndex(request *http.Request) *types.FinancingContractSearch {
	query := request.URL.Query()
	pageid := "1"
	searchType := "increase"
	FinanceId := ""
	if len(query["pageid"]) > 0 {
		pageid = query["pageid"][0]
	}
	if len(query["FinanceId"]) > 0 {
		pageid = query["FinanceId"][0]
	}
	if len(query["searchType"]) > 0 {
		searchType = query["searchType"][0]
	}
	index := types.FinancingContractSearch{
		PageId:     pageid,
		SearchType: searchType,
		FinanceId:  FinanceId,
	}
	return &index
}

// 还款信息URL
func (s *SqlCtr) RepaymentRecordIndex(request *http.Request) *types.RepaymentRecordSearch {
	query := request.URL.Query()
	FinanceId := ""
	pageid := "1"
	searchType := "increase"
	if len(query["id"]) > 0 {
		FinanceId = query["FinanceId"][0]
	}
	if len(query["pageid"]) > 0 {
		pageid = query["pageid"][0]
	}
	if len(query["searchType"]) > 0 {
		searchType = query["searchType"][0]
	}
	index := types.RepaymentRecordSearch{
		PageId:     pageid,
		SearchType: searchType,
		FinanceId:  FinanceId,
	}
	return &index
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 将需要解密的数据，解密后打包成结构体
// 输入参数是解密后的发票信息，转换成redis存储所需要的数据结构
func (s *SqlCtr) InvoiceinfoToMap(ret []string) []*types.InvoiceInformation {
	return HandleInvoiceInfo(ret)
}

// 打包融资意向申请
func (s *SqlCtr) IntensioninfoToMap(ret []string) []*types.FinancingIntention {

	return HandleFinancingIntention(ret)
}

// 打包回款账户信息
func (s *SqlCtr) AccountinfoToMap(ret []string) []*types.CollectionAccount {

	return HandleCollectionAccount(ret)
}

// 打包借贷合同
func (s *SqlCtr) FinancingContractToMap(ret []*types.RawFinancingContractData) []*types.FinancingContract {

	return HandleFinancingContract(ret)
}

// 打包还款信息
func (s *SqlCtr) RepaymentRecordToMap(ret []*types.RawRepaymentRecord) []*types.RepaymentRecord {

	return HandleRepaymentRecord(ret)
}

// ////////////////////////////////////////////////////////////////////////////////////////
// 查询mysql数据库中加密后的发票信息，如果id为空，则查找全部的信息
func (s *SqlCtr) QueryInvoiceInformation(id string) []string {
	var ret []string
	if id == "" {
		ret, _ = s.QueryInvoiceByOrder(s.orders.InvoiceSQLQueryAll)
	} else {
		ret, _ = s.QueryInvoiceByOrder(s.orders.InvoiceSQLQueryByID + id + "%")
	}
	return ret
}
func (s *SqlCtr) QueryInvoiceInforsBySQLID(_id_ string) []string {
	var ret []string
	ret, _ = s.QueryInvoiceByOrder(s.orders.InvoiceSQLQueryBy_ID + _id_)
	return ret
}
func (s *SqlCtr) QueryInvoiceInformationLength() int {
	var length int
	err := s.db.QueryRow(s.orders.InvoiceSQLQueryLength).Scan(&length)
	if err != nil {
		logrus.Fatalln(err)
	}
	return length
}

// //////////////////////////////////////////////////////////////////////////////////////////////////////
func (s *SqlCtr) QueryHistoricalTransUsedInfos(id string) []string {
	var ret []string
	if id == "" {
		ret, _ = s.QueryHistoricalTransByOrder(s.orders.HistoricalUsedSQLQueryAll)
	} else {
		ret, _ = s.QueryHistoricalTransByOrder(s.orders.HistoricalUsedSQLQueryByID + id)
	}
	return ret
}
func (s *SqlCtr) QueryHistoricalTransUsedInfosBySQLID(_id_ string) []string {
	var ret []string
	ret, _ = s.QueryHistoricalTransByOrder(s.orders.HistoricalUsedSQLQueryBy_ID + _id_)
	return ret
}

func (s *SqlCtr) QueryHistoricalTransSettleInfos(id string) []string {
	var ret []string
	if id == "" {
		ret, _ = s.QueryHistoricalTransByOrder(s.orders.HistoricalSettleSQLQueryAll)
	} else {
		ret, _ = s.QueryHistoricalTransByOrder(s.orders.HistoricalSettleSQLQueryByID + id)
	}
	return ret
}
func (s *SqlCtr) QueryHistoricalTransSettleInfosBySQLID(_id_ string) []string {
	var ret []string
	ret, _ = s.QueryHistoricalTransByOrder(s.orders.HistoricalSettleSQLQueryBy_ID + _id_)
	return ret
}
func (s *SqlCtr) QueryHistoricalTransOrderInfos(id string) []string {
	var ret []string
	if id == "" {
		ret, _ = s.QueryHistoricalTransByOrder(s.orders.HistoricalOrderSQLQueryAll)
	} else {
		ret, _ = s.QueryHistoricalTransByOrder(s.orders.HistoricalOrderSQLQueryByID + id)
	}
	return ret
}
func (s *SqlCtr) QueryHistoricalTransOrderInfosBySQLID(_id_ string) []string {
	var ret []string
	ret, _ = s.QueryHistoricalTransByOrder(s.orders.HistoricalOrderSQLQueryBy_ID + _id_)
	return ret
}

func (s *SqlCtr) QueryHistoricalTransReceivableInfos(id string) []string {
	var ret []string
	if id == "" {
		ret, _ = s.QueryHistoricalTransByOrder(s.orders.HistoricalReceivableSQLQueryAll)
	} else {
		ret, _ = s.QueryHistoricalTransByOrder(s.orders.HistoricalReceivableSQLQueryByID + id)
	}
	return ret
}
func (s *SqlCtr) QueryHistoricalTransReceivableInfosBySQLID(_id_ string) []string {
	var ret []string
	ret, _ = s.QueryHistoricalTransByOrder(s.orders.HistoricalReceivableSQLQueryBy_ID + _id_)
	return ret
}

// //////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (s *SqlCtr) QueryEnterpoolDataPlanInfos(id string) []string {
	var ret []string
	if id == "" {
		ret, _ = s.QueryPoolDataByOrder(s.orders.EnterPoolPlanSQLQueryAll)
	} else {
		ret, _ = s.QueryPoolDataByOrder(s.orders.EnterPoolPlanSQLQueryByID + id)
	}
	return ret
}
func (s *SqlCtr) QueryEnterpoolDataPlanInfosBySQLID(_id_ string) []string {
	var ret []string
	ret, _ = s.QueryPoolDataByOrder(s.orders.EnterPoolPlanSQLQueryBy_ID + _id_)
	return ret
}
func (s *SqlCtr) QueryEnterpoolDataUsedInfos(id string) []string {
	var ret []string
	if id == "" {
		ret, _ = s.QueryPoolDataByOrder(s.orders.EnterPoolUsedSQLQueryAll)
	} else {
		ret, _ = s.QueryPoolDataByOrder(s.orders.EnterPoolUsedSQLQueryByID + id)
	}
	return ret
}
func (s *SqlCtr) QueryEnterpoolDataUsedInfosBySQLID(_id_ string) []string {
	var ret []string
	ret, _ = s.QueryPoolDataByOrder(s.orders.EnterPoolUsedSQLQueryBy_ID + _id_)
	return ret
}

// //////////////////////////////////////////////////////////////////////////////////////////////////////
// 查询mysql数据库中融资意向信息，如果id为空，则查找全部的信息
func (s *SqlCtr) QueryFinancingIntention(id string) []string {
	var ret []string
	if id == "" {
		ret, _ = s.QueryFinancingByOrder(s.orders.FinancingSQLQueryAll)
	} else {
		ret, _ = s.QueryFinancingByOrder(s.orders.FinancingSQLQueryByID + id)
	}
	return ret
}
func (s *SqlCtr) QueryFinancingIntentionBySQLID(_id_ string) []string {
	var ret []string
	ret, _ = s.QueryFinancingByOrder(s.orders.FinancingSQLQueryBy_ID + _id_)
	return ret
}

// //////////////////////////////////////////////////////////////////////////////
// 查询mysql数据库中回款账户信息，如果id为空，则查找全部的信息
func (s *SqlCtr) QueryCollectionAccount(id string) []string {
	var ret []string
	if id == "" {
		ret, _ = s.QueryAccountsByOrder(s.orders.AccountsSQLQueryAll)
	} else {
		ret, _ = s.QueryAccountsByOrder(s.orders.AccountsSQLQueryByID + id)
	}
	return ret
}
func (s *SqlCtr) QueryCollectionAccountBySQLID(_id_ string) []string {
	var ret []string
	ret, _ = s.QueryAccountsByOrder(s.orders.AccountsSQLQueryBy_ID + _id_)
	return ret
}

// //////////////////////////////////////////////////////////////////////////////////
func (s *SqlCtr) QueryFinancingContract(id string) []*types.RawFinancingContractData {
	var ret []*types.RawFinancingContractData
	if id == "" {
		ret, _ = s.QueryFinancingContractByOrder(s.orders.FinancingContractSQLAll)
	} else {
		ret, _ = s.QueryFinancingContractByOrder(s.orders.FinancingContractSQLByID + id)
	}
	return ret
}
func (s *SqlCtr) QueryFinancingContractBySQLID(_id_ string) []*types.RawFinancingContractData {
	var ret []*types.RawFinancingContractData
	ret, _ = s.QueryFinancingContractByOrder(s.orders.FinancingContractSQLBy_ID + _id_)
	return ret
}

// //////////////////////////////////////////////////////////////////////////////////////////
func (s *SqlCtr) QueryRepaymentRecord(id string) []*types.RawRepaymentRecord {
	var ret []*types.RawRepaymentRecord
	if id == "" {
		ret, _ = s.QueryRepaymentRecordByOrder(s.orders.RepaymentRecordSQLAll)
	} else {
		ret, _ = s.QueryRepaymentRecordByOrder(s.orders.RepaymentRecordSQLByID + id)
	}
	return ret
}
func (s *SqlCtr) QueryRepaymentRecordBySQLID(_id_ string) []*types.RawRepaymentRecord {
	var ret []*types.RawRepaymentRecord
	ret, _ = s.QueryRepaymentRecordByOrder(s.orders.RepaymentRecordSQLBy_ID + _id_)
	return ret
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////
// 查询解密发票信息
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
			StrData := string(data) + "," + record.Owner
			ret = append(ret, StrData)
		} else {
			logrus.Errorln("哈希值与数据对应错误")
		}
		i = i + 1
	}
	return ret, nil
}

// 查询解密历史交易信息
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

// 查询解密入池信息
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

// 查询解密融资意向申请
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
		err = rows.Scan(&record.SQLId, &record.Num, &record.Status, &record.ID, &record.FinanceId, &record.CustomerId, &record.Data, &record.Key, &record.Hash, &record.State)
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
			data = []byte(string(data) + "," + record.State)
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
		err = rows.Scan(&record.SQLId, &record.Num, &record.Status, &record.ID, &record.FinanceID, &record.Data, &record.Key, &record.Hash, &record.State)
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
func (s *SqlCtr) QueryFinancingContractByOrder(order string) ([]*types.RawFinancingContractData, error) {
	in_stmt, err := s.db.Prepare(order)
	if err != nil {
		return nil, err
	}
	rows, err := in_stmt.Query()
	if err != nil {
		return nil, err
	}
	ret := make([]*types.RawFinancingContractData, 0)
	count := 0
	for rows.Next() {
		record := &types.RawFinancingContractData{}
		err = rows.Scan(&record.SQLId, &record.Num, &record.Status, &record.ID, &record.FinancingID, &record.CustomerID, &record.CorpName, &record.DebtMoney, &record.SupplyDate, &record.ExpireDate, &record.Balance)
		if err != nil {
			logrus.Errorln(err)
			count++
			continue
		}
		ret = append(ret, record)
	}
	return ret, nil
}

func (s *SqlCtr) QueryRepaymentRecordByOrder(order string) ([]*types.RawRepaymentRecord, error) {
	in_stmt, err := s.db.Prepare(order)
	if err != nil {
		return nil, err
	}
	rows, err := in_stmt.Query()
	if err != nil {
		return nil, err
	}
	ret := make([]*types.RawRepaymentRecord, 0)
	count := 0
	for rows.Next() {
		record := &types.RawRepaymentRecord{}
		err = rows.Scan(&record.SQLId, &record.Num, &record.Status, &record.ID, &record.FinancingID, &record.CustomerID, &record.Repay, &record.Time)
		if err != nil {
			logrus.Errorln(err)
			count++
			continue
		}
		ret = append(ret, record)
	}
	return ret, nil
}
