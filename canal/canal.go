package canal

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/FISCO-BCOS/go-sdk/conf"
	types "github.com/FISCO-BCOS/go-sdk/type"
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
	"github.com/withlin/canal-go/client"
	pbe "github.com/withlin/canal-go/protocol/entry"
)

type Connector struct {
	conn    *client.SimpleCanalConnector
	RawData map[string][]*types.RawCanalData
	tables  *conf.Config
	Lock    sync.RWMutex
}

func NewConnector(table string) *Connector {
	configs, err := conf.ParseConfigFile("./configs/config.toml")
	if err != nil {
		logrus.Fatalln(err)
	}
	config := &configs[0]
	fmt.Println(config.CanalIP, config.CanalPort, config.CanalUsername, config.CanalPassword, config.CanalDestination)
	connector := client.NewSimpleCanalConnector(config.CanalIP, config.CanalPort, config.CanalUsername, config.CanalPassword, config.CanalDestination, 60000, 60*60*1000)
	err = connector.Connect()
	if err != nil {
		logrus.Errorln(err)
		os.Exit(1)
	}
	// https://github.com/alibaba/canal/wiki/AdminGuide
	//mysql 数据解析关注的表，Perl正则表达式.
	//
	//多个正则之间以逗号(,)分隔，转义符需要双斜杠(\\)
	//
	//常见例子：
	//
	//  1.  所有表：.*   or  .*\\..*
	//	2.  canal schema下所有表： canal\\..*
	//	3.  canal下的以canal打头的表：canal\\.canal.*
	//	4.  canal schema下的一张表：canal\\.test1
	//  5.  多个规则组合使用：canal\\..*,mysql.test1,mysql.test2 (逗号分隔)

	// err = connector.Subscribe("db_node1\\.u_t_history_settle_information")
	logrus.Infoln(table)
	err = connector.Subscribe(table)
	if err != nil {
		logrus.Errorln(err)
		os.Exit(1)
	}
	raw := make(map[string][]*types.RawCanalData)
	return &Connector{
		conn:    connector,
		RawData: raw,
		tables:  config,
	}
}

// 开始运行canal
func (c *Connector) Start() {
	for {
		message, err := c.conn.Get(100, nil, nil)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		fmt.Println(message)
		batchId := message.Id
		if batchId == -1 || len(message.Entries) <= 0 {
			time.Sleep(3 * time.Second)
			logrus.Println("===没有数据了===")
			continue
		}
		c.dealMessage(message.Entries)
	}
}
func (c *Connector) dealMessage(entrys []pbe.Entry) {
	for _, entry := range entrys {
		fmt.Println(entry.GetEntryType())
		if entry.GetEntryType() == pbe.EntryType_TRANSACTIONBEGIN || entry.GetEntryType() == pbe.EntryType_TRANSACTIONEND {
			continue
		}
		rowChange := pbe.RowChange{}
		err := proto.Unmarshal(entry.GetStoreValue(), &rowChange)
		checkError(err)
		eventType := rowChange.GetEventType()
		header := entry.GetHeader()
		// fmt.Println(fmt.Sprintf("================> binlog[%s : %d],name[%s,%s], eventType: %s", header.GetLogfileName(), header.GetLogfileOffset(), header.GetSchemaName(), header.GetTableName(), header.GetEventType()))
		switch header.GetTableName() {
		case c.tables.InvoiceInfos:
			for _, rowData := range rowChange.GetRowDatas() {
				if eventType == pbe.EventType_INSERT {
					c.dealInsertInvoiceMessage(rowData.GetAfterColumns())

				}
			}
		case c.tables.FinanceApplication:
			for _, rowData := range rowChange.GetRowDatas() {
				if eventType == pbe.EventType_INSERT {
					c.dealInsertIntensionMessage(rowData.GetAfterColumns())
				}
			}
		case c.tables.Accounts:
			for _, rowData := range rowChange.GetRowDatas() {
				if eventType == pbe.EventType_INSERT {
					c.dealInsertAccountMessage(rowData.GetAfterColumns())
				}
			}
		case c.tables.HistoricalOrder:
			for _, rowData := range rowChange.GetRowDatas() {
				if eventType == pbe.EventType_INSERT {
					c.dealInsertHistoryOrderMessage(rowData.GetAfterColumns())
				}
			}
		case c.tables.HistoricalReceivable:
			for _, rowData := range rowChange.GetRowDatas() {
				if eventType == pbe.EventType_INSERT {
					c.dealInsertHistoryReceivableMessage(rowData.GetAfterColumns())
				}
			}
		case c.tables.HistoricalSettle:
			for _, rowData := range rowChange.GetRowDatas() {
				if eventType == pbe.EventType_INSERT {
					c.dealInsertHistorySettleMessage(rowData.GetAfterColumns())
				}
			}
		case c.tables.HistoricalUsed:
			for _, rowData := range rowChange.GetRowDatas() {
				if eventType == pbe.EventType_INSERT {
					c.dealInsertHistoryUsedMessage(rowData.GetAfterColumns())
				}
			}
		case c.tables.PoolPlanInfos:
			for _, rowData := range rowChange.GetRowDatas() {
				if eventType == pbe.EventType_INSERT {
					c.dealInsertPoolPlanMessage(rowData.GetAfterColumns())
				}
			}
		case c.tables.PoolUsedInfos:
			for _, rowData := range rowChange.GetRowDatas() {
				if eventType == pbe.EventType_INSERT {
					c.dealInsertPoolUsedMessage(rowData.GetAfterColumns())
				}
			}
		case c.tables.FinanceContract:
			for _, rowData := range rowChange.GetRowDatas() {
				if eventType == pbe.EventType_INSERT {
					c.dealInsertFinancingContractMessage(rowData.GetAfterColumns())
				}
			}
		case c.tables.RepaymentRecord:
			for _, rowData := range rowChange.GetRowDatas() {
				if eventType == pbe.EventType_INSERT {
					c.dealInsertRepaymentRecordtMessage(rowData.GetAfterColumns())
				}
			}
		default:
			logrus.Warnln("未知的数据库")
		}
	}
}
func (c *Connector) dealInsertInvoiceMessage(columns []*pbe.Column) {
	rawdata := new(types.RawCanalData)
	rawdata.SQLId = []byte(columns[0].GetValue())
	fmt.Println(rawdata)
	c.Lock.Lock()
	c.RawData[c.tables.InvoiceInfos] = append(c.RawData[c.tables.InvoiceInfos], rawdata)
	c.Lock.Unlock()

}
func (c *Connector) dealInsertHistoryOrderMessage(columns []*pbe.Column) {
	rawdata := new(types.RawCanalData)
	rawdata.SQLId = []byte(columns[0].GetValue())
	fmt.Println(rawdata)
	c.Lock.Lock()
	c.RawData[c.tables.HistoricalOrder] = append(c.RawData[c.tables.HistoricalOrder], rawdata)
	c.Lock.Unlock()

}
func (c *Connector) dealInsertHistoryUsedMessage(columns []*pbe.Column) {
	rawdata := new(types.RawCanalData)
	rawdata.SQLId = []byte(columns[0].GetValue())
	fmt.Println(rawdata)
	c.Lock.Lock()
	c.RawData[c.tables.HistoricalUsed] = append(c.RawData[c.tables.HistoricalUsed], rawdata)
	c.Lock.Unlock()

}
func (c *Connector) dealInsertHistorySettleMessage(columns []*pbe.Column) {
	rawdata := new(types.RawCanalData)
	rawdata.SQLId = []byte(columns[0].GetValue())
	fmt.Println(rawdata)
	c.Lock.Lock()
	c.RawData[c.tables.HistoricalSettle] = append(c.RawData[c.tables.HistoricalSettle], rawdata)
	c.Lock.Unlock()

}
func (c *Connector) dealInsertHistoryReceivableMessage(columns []*pbe.Column) {
	rawdata := new(types.RawCanalData)
	rawdata.SQLId = []byte(columns[0].GetValue())
	fmt.Println(rawdata)
	c.Lock.Lock()
	c.RawData[c.tables.HistoricalReceivable] = append(c.RawData[c.tables.HistoricalReceivable], rawdata)
	c.Lock.Unlock()

}
func (c *Connector) dealInsertIntensionMessage(columns []*pbe.Column) {
	rawdata := new(types.RawCanalData)
	rawdata.SQLId = []byte(columns[0].GetValue())
	fmt.Println(rawdata)
	c.Lock.Lock()
	c.RawData[c.tables.FinanceApplication] = append(c.RawData[c.tables.FinanceApplication], rawdata)
	c.Lock.Unlock()
}
func (c *Connector) dealInsertAccountMessage(columns []*pbe.Column) {
	rawdata := new(types.RawCanalData)
	rawdata.SQLId = []byte(columns[0].GetValue())
	fmt.Println(rawdata)
	c.Lock.Lock()
	c.RawData[c.tables.Accounts] = append(c.RawData[c.tables.Accounts], rawdata)
	c.Lock.Unlock()

}
func (c *Connector) dealInsertPoolPlanMessage(columns []*pbe.Column) {
	rawdata := new(types.RawCanalData)
	rawdata.SQLId = []byte(columns[0].GetValue())
	fmt.Println(rawdata)
	c.Lock.Lock()
	c.RawData[c.tables.PoolPlanInfos] = append(c.RawData[c.tables.PoolPlanInfos], rawdata)
	c.Lock.Unlock()

}
func (c *Connector) dealInsertPoolUsedMessage(columns []*pbe.Column) {
	rawdata := new(types.RawCanalData)
	rawdata.SQLId = []byte(columns[0].GetValue())
	fmt.Println(rawdata)
	c.Lock.Lock()
	c.RawData[c.tables.PoolUsedInfos] = append(c.RawData[c.tables.PoolUsedInfos], rawdata)
	c.Lock.Unlock()

}
func (c *Connector) dealInsertFinancingContractMessage(columns []*pbe.Column) {
	rawdata := new(types.RawCanalData)
	rawdata.SQLId = []byte(columns[0].GetValue())
	fmt.Println(rawdata)
	c.Lock.Lock()
	c.RawData[c.tables.FinanceContract] = append(c.RawData[c.tables.FinanceContract], rawdata)
	c.Lock.Unlock()

}
func (c *Connector) dealInsertRepaymentRecordtMessage(columns []*pbe.Column) {
	rawdata := new(types.RawCanalData)
	rawdata.SQLId = []byte(columns[0].GetValue())
	fmt.Println(rawdata)
	c.Lock.Lock()
	c.RawData[c.tables.RepaymentRecord] = append(c.RawData[c.tables.RepaymentRecord], rawdata)
	c.Lock.Unlock()

}
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
