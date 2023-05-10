// Package conf parse config to configuration
package conf

import (
	"bytes"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config contains configuration items for sdk
type Config struct {
	//mysql
	MslUrl      string
	MslUsername string
	MslPasswd   string
	MslName     string
	MslProtocol string
	//redis
	RedisUrl      string
	RedisPassword string
	//fisco
	FiscoUrl string
	//logDB
	LogDBUrl      string
	LogDBUsername string
	LogDBPasswd   string
	LogDBName     string
	LogDBProtocol string
	//canal
	CanalIP          string
	CanalPort        int
	CanalUsername    string
	CanalPassword    string
	CanalDestination string
	CanalConnectedDB string
	//tables
	FinanceContract      string
	HistoricalOrder      string
	HistoricalReceivable string
	HistoricalSettle     string
	HistoricalUsed       string
	InvoiceInfos         string
	PoolPlanInfos        string
	PoolUsedInfos        string
	Accounts             string
	RepaymentRecord      string
	FinanceApplication   string
	//SQLOrders
	//发票
	InvoiceSQLQueryAll    string
	InvoiceSQLQueryByID   string
	InvoiceSQLQueryBy_ID  string
	InvoiceSQLQueryLength string
	//历史交易-used
	HistoricalUsedSQLQueryAll   string
	HistoricalUsedSQLQueryByID  string
	HistoricalUsedSQLQueryBy_ID string
	//历史交易-settle
	HistoricalSettleSQLQueryAll   string
	HistoricalSettleSQLQueryByID  string
	HistoricalSettleSQLQueryBy_ID string
	//历史交易-order
	HistoricalOrderSQLQueryAll   string
	HistoricalOrderSQLQueryByID  string
	HistoricalOrderSQLQueryBy_ID string
	//历史交易-receivable
	HistoricalReceivableSQLQueryAll   string
	HistoricalReceivableSQLQueryByID  string
	HistoricalReceivableSQLQueryBy_ID string
	//入池数据-plan
	EnterPoolPlanSQLQueryAll   string
	EnterPoolPlanSQLQueryByID  string
	EnterPoolPlanSQLQueryBy_ID string
	//入池数据-used
	EnterPoolUsedSQLQueryAll   string
	EnterPoolUsedSQLQueryByID  string
	EnterPoolUsedSQLQueryBy_ID string
	//融资意向
	FinancingSQLQueryAll   string
	FinancingSQLQueryByID  string
	FinancingSQLQueryBy_ID string
	//回款账户
	AccountsSQLQueryAll   string
	AccountsSQLQueryByID  string
	AccountsSQLQueryBy_ID string
	//借贷合同
	FinancingContractSQLAll   string
	FinancingContractSQLByID  string
	FinancingContractSQLBy_ID string
	//还款记录
	RepaymentRecordSQLAll   string
	RepaymentRecordSQLByID  string
	RepaymentRecordSQLBy_ID string
}

// ParseConfigFile parses the configuration from toml config file
func ParseConfigFile(cfgFile string) ([]Config, error) {
	file, err := os.Open(cfgFile)
	if err != nil {
		return nil, fmt.Errorf("open file failed, err: %v", err)
	}

	defer func() {
		err = file.Close()
		if err != nil {
			logrus.Fatalf("close file failed, err: %v", err)
		}
	}()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("file is not found, err: %v", err)
	}

	fileSize := fileInfo.Size()
	buffer := make([]byte, fileSize)

	_, err = file.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("read file failed, err: %v", err)
	}
	return ParseConfig(buffer)
}

// ParseConfig parses the configuration from []byte
func ParseConfig(buffer []byte) ([]Config, error) {
	viper.SetConfigType("toml")
	viper.SetDefault("SMCrypto", false)
	err := viper.ReadConfig(bytes.NewBuffer(buffer))
	if err != nil {
		return nil, fmt.Errorf("viper .ReadConfig failed, err: %v", err)
	}
	config := new(Config)
	var configs []Config
	if viper.IsSet("Mysql") {
		config.MslUrl = viper.GetString("Mysql.MslUrl")
		config.MslUsername = viper.GetString("Mysql.MslUsername")
		config.MslPasswd = viper.GetString("Mysql.MslPasswd")
		config.MslName = viper.GetString("Mysql.MslName")
		config.MslProtocol = viper.GetString("Mysql.MslProtocol")
	}
	if viper.IsSet("LogDB") {
		config.LogDBUrl = viper.GetString("LogDB.LogDBUrl")
		config.LogDBUsername = viper.GetString("LogDB.LogDBUsername")
		config.LogDBPasswd = viper.GetString("LogDB.LogDBPasswd")
		config.LogDBName = viper.GetString("LogDB.LogDBName")
		config.LogDBProtocol = viper.GetString("LogDB.LogDBProtocol")
	}
	if viper.IsSet("Fisco") {
		config.FiscoUrl = viper.GetString("Fisco.FiscoUrl")

	}
	if viper.IsSet("Redis") {
		config.RedisUrl = viper.GetString("Redis.RedisUrl")
		config.RedisPassword = viper.GetString("Redis.RedisPassword")
	}
	if viper.IsSet("Canal") {
		config.CanalIP = viper.GetString("Canal.CanalIP")
		config.CanalPort = viper.GetInt("Canal.CanalPort")
		config.CanalUsername = viper.GetString("Canal.CanalUsername")
		config.CanalPassword = viper.GetString("Canal.CanalPassword")
		config.CanalDestination = viper.GetString("Canal.CanalDestination")
		config.CanalConnectedDB = viper.GetString("Canal.CanalConnectedDB")
	}
	if viper.IsSet("SubscribeTable") {
		config.FinanceContract = viper.GetString("SubscribeTable.FinanceContract")
		config.HistoricalOrder = viper.GetString("SubscribeTable.HistoricalOrder")
		config.HistoricalReceivable = viper.GetString("SubscribeTable.HistoricalReceivable")
		config.HistoricalSettle = viper.GetString("SubscribeTable.HistoricalSettle")
		config.HistoricalUsed = viper.GetString("SubscribeTable.HistoricalUsed")
		config.InvoiceInfos = viper.GetString("SubscribeTable.InvoiceInfos")
		config.PoolPlanInfos = viper.GetString("SubscribeTable.PoolPlanInfos")
		config.PoolUsedInfos = viper.GetString("SubscribeTable.PoolUsedInfos")
		config.Accounts = viper.GetString("SubscribeTable.Accounts")
		config.RepaymentRecord = viper.GetString("SubscribeTable.RepaymentRecord")
		config.FinanceApplication = viper.GetString("SubscribeTable.FinanceApplication")
	}
	if viper.IsSet("SQLOrder") {
		config.InvoiceSQLQueryAll = viper.GetString("SQLOrder.InvoiceSQLQueryAll")
		config.InvoiceSQLQueryByID = viper.GetString("SQLOrder.InvoiceSQLQueryByID")
		config.InvoiceSQLQueryBy_ID = viper.GetString("SQLOrder.InvoiceSQLQueryBy_ID")
		config.InvoiceSQLQueryLength = viper.GetString("SQLOrder.InvoiceSQLQueryLength")

		config.HistoricalUsedSQLQueryAll = viper.GetString("SQLOrder.HistoricalUsedSQLQueryAll")
		config.HistoricalUsedSQLQueryByID = viper.GetString("SQLOrder.HistoricalUsedSQLQueryByID")
		config.HistoricalUsedSQLQueryBy_ID = viper.GetString("SQLOrder.HistoricalUsedSQLQueryBy_ID")

		config.HistoricalSettleSQLQueryAll = viper.GetString("SQLOrder.HistoricalSettleSQLQueryAll")
		config.HistoricalSettleSQLQueryByID = viper.GetString("SQLOrder.HistoricalSettleSQLQueryByID")
		config.HistoricalSettleSQLQueryBy_ID = viper.GetString("SQLOrder.HistoricalSettleSQLQueryBy_ID")

		config.HistoricalOrderSQLQueryAll = viper.GetString("SQLOrder.HistoricalOrderSQLQueryAll")
		config.HistoricalOrderSQLQueryByID = viper.GetString("SQLOrder.HistoricalOrderSQLQueryByID")
		config.HistoricalOrderSQLQueryBy_ID = viper.GetString("SQLOrder.HistoricalOrderSQLQueryBy_ID")

		config.HistoricalReceivableSQLQueryAll = viper.GetString("SQLOrder.HistoricalReceivableSQLQueryAll")
		config.HistoricalReceivableSQLQueryByID = viper.GetString("SQLOrder.HistoricalReceivableSQLQueryByID")
		config.HistoricalReceivableSQLQueryBy_ID = viper.GetString("SQLOrder.HistoricalReceivableSQLQueryBy_ID")

		config.EnterPoolPlanSQLQueryAll = viper.GetString("SQLOrder.EnterPoolPlanSQLQueryAll")
		config.EnterPoolPlanSQLQueryByID = viper.GetString("SQLOrder.EnterPoolPlanSQLQueryByID")
		config.EnterPoolPlanSQLQueryBy_ID = viper.GetString("SQLOrder.EnterPoolPlanSQLQueryBy_ID")

		config.EnterPoolUsedSQLQueryAll = viper.GetString("SQLOrder.EnterPoolUsedSQLQueryAll")
		config.EnterPoolUsedSQLQueryByID = viper.GetString("SQLOrder.EnterPoolUsedSQLQueryByID")
		config.EnterPoolUsedSQLQueryBy_ID = viper.GetString("SQLOrder.EnterPoolUsedSQLQueryBy_ID")

		config.FinancingSQLQueryAll = viper.GetString("SQLOrder.FinancingSQLQueryAll")
		config.FinancingSQLQueryByID = viper.GetString("SQLOrder.FinancingSQLQueryByID")
		config.FinancingSQLQueryBy_ID = viper.GetString("SQLOrder.FinancingSQLQueryBy_ID")

		config.AccountsSQLQueryAll = viper.GetString("SQLOrder.AccountsSQLQueryAll")
		config.AccountsSQLQueryByID = viper.GetString("SQLOrder.AccountsSQLQueryByID")
		config.AccountsSQLQueryBy_ID = viper.GetString("SQLOrder.AccountsSQLQueryBy_ID")

		config.FinancingContractSQLAll = viper.GetString("SQLOrder.FinancingContractSQLAll")
		config.FinancingContractSQLByID = viper.GetString("SQLOrder.FinancingContractSQLByID")
		config.FinancingContractSQLBy_ID = viper.GetString("SQLOrder.FinancingContractSQLBy_ID")

		config.RepaymentRecordSQLAll = viper.GetString("SQLOrder.RepaymentRecordSQLAll")
		config.RepaymentRecordSQLByID = viper.GetString("SQLOrder.RepaymentRecordSQLByID")
		config.RepaymentRecordSQLBy_ID = viper.GetString("SQLOrder.RepaymentRecordSQLBy_ID")
	}
	configs = append(configs, *config)
	return configs, nil
}
