package canal

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/FISCO-BCOS/go-sdk/conf"
	queue "github.com/FISCO-BCOS/go-sdk/structure"
	types "github.com/FISCO-BCOS/go-sdk/type"
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
	"github.com/withlin/canal-go/client"
	pbe "github.com/withlin/canal-go/protocol/entry"
)

const (
	Invoice           = "u_t_invoice_information3"
	Accounts          = "u_t_push_payment_accounts"
	Intension         = "u_t_supplier_financing_application1"
	HisOrder          = "u_t_history_order_information3"
	HisUsed           = "u_t_history_used_information3"
	HisSettle         = "u_t_history_settle_information4"
	HisReceivable     = "u_t_history_receivable_information3"
	PoolPlan          = "u_t_pool_plan_information2"
	PoolUsed          = "u_t_pool_used_information2"
	FinancingContract = "u_t_finance_contract1"
)

type Connector struct {
	conn    *client.SimpleCanalConnector
	queue   *queue.CircleQueue
	RawData map[string][]*types.RawCanalData
	Lock    sync.RWMutex
}

func NewConnector(table string) *Connector {
	configs, err := conf.ParseConfigFile("./configs/config.toml")
	if err != nil {
		logrus.Fatalln(err)
	}
	config := &configs[0]
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
	queue := queue.NewCircleQueue(20)
	raw := make(map[string][]*types.RawCanalData)
	return &Connector{
		conn:    connector,
		queue:   queue,
		RawData: raw,
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
		batchId := message.Id
		if batchId == -1 || len(message.Entries) <= 0 {
			time.Sleep(1 * time.Second)
			logrus.Println("===没有数据了===")
			continue
		}
		c.dealMessage(message.Entries)
	}
}
func (c *Connector) dealMessage(entrys []pbe.Entry) {
	fmt.Println("1223456334")
	for _, entry := range entrys {
		fmt.Println(entry.GetEntryType())
		if entry.GetEntryType() == pbe.EntryType_TRANSACTIONBEGIN || entry.GetEntryType() == pbe.EntryType_TRANSACTIONEND {
			continue
		}
		rowChange := new(pbe.RowChange)
		err := proto.Unmarshal(entry.GetStoreValue(), rowChange)
		checkError(err)
		if rowChange != nil {
			eventType := rowChange.GetEventType()
			header := entry.GetHeader()
			fmt.Println(fmt.Sprintf("================> binlog[%s : %d],name[%s,%s], eventType: %s", header.GetLogfileName(), header.GetLogfileOffset(), header.GetSchemaName(), header.GetTableName(), header.GetEventType()))
			switch header.GetTableName() {
			case Invoice:
				for _, rowData := range rowChange.GetRowDatas() {
					if eventType == pbe.EventType_INSERT {
						c.dealInsertInvoiceMessage(rowData.GetAfterColumns())

					}
				}
			case Intension:
				for _, rowData := range rowChange.GetRowDatas() {
					if eventType == pbe.EventType_INSERT {
						c.dealInsertIntensionMessage(rowData.GetAfterColumns())
					}
				}
			case Accounts:
				for _, rowData := range rowChange.GetRowDatas() {
					if eventType == pbe.EventType_INSERT {
						c.dealInsertAccountMessage(rowData.GetAfterColumns())
					}
				}
			case HisOrder:
				for _, rowData := range rowChange.GetRowDatas() {
					if eventType == pbe.EventType_INSERT {
						c.dealInsertHistoryOrderMessage(rowData.GetAfterColumns())
					}
				}
			case HisReceivable:
				for _, rowData := range rowChange.GetRowDatas() {
					if eventType == pbe.EventType_INSERT {
						c.dealInsertHistoryReceivableMessage(rowData.GetAfterColumns())
					}
				}
			case HisSettle:
				for _, rowData := range rowChange.GetRowDatas() {
					if eventType == pbe.EventType_INSERT {
						c.dealInsertHistorySettleMessage(rowData.GetAfterColumns())
					}
				}
			case HisUsed:
				for _, rowData := range rowChange.GetRowDatas() {
					if eventType == pbe.EventType_INSERT {
						c.dealInsertHistoryUsedMessage(rowData.GetAfterColumns())
					}
				}
			case PoolPlan:
				for _, rowData := range rowChange.GetRowDatas() {
					if eventType == pbe.EventType_INSERT {
						c.dealInsertPoolPlanMessage(rowData.GetAfterColumns())
					}
				}
			case PoolUsed:
				for _, rowData := range rowChange.GetRowDatas() {
					if eventType == pbe.EventType_INSERT {
						c.dealInsertPoolUsedMessage(rowData.GetAfterColumns())
					}
				}
			case FinancingContract:
				for _, rowData := range rowChange.GetRowDatas() {
					if eventType == pbe.EventType_INSERT {
						c.dealInsertFinancingContractMessage(rowData.GetAfterColumns())
					}
				}
			default:
				logrus.Warnln("未知的数据库")
			}
		}
	}
}
func (c *Connector) dealInsertInvoiceMessage(columns []*pbe.Column) {
	rawdata := new(types.RawCanalData)
	rawdata.SQLId = []byte(columns[0].GetValue())
	fmt.Println(rawdata)
	c.Lock.Lock()
	c.RawData[Invoice] = append(c.RawData[Invoice], rawdata)
	c.Lock.Unlock()

}
func (c *Connector) dealInsertHistoryOrderMessage(columns []*pbe.Column) {
	rawdata := new(types.RawCanalData)
	rawdata.SQLId = []byte(columns[0].GetValue())
	fmt.Println(rawdata)
	c.Lock.Lock()
	c.RawData[HisOrder] = append(c.RawData[HisOrder], rawdata)
	c.Lock.Unlock()

}
func (c *Connector) dealInsertHistoryUsedMessage(columns []*pbe.Column) {
	rawdata := new(types.RawCanalData)
	rawdata.SQLId = []byte(columns[0].GetValue())
	fmt.Println(rawdata)
	c.Lock.Lock()
	c.RawData[HisUsed] = append(c.RawData[HisUsed], rawdata)
	c.Lock.Unlock()

}
func (c *Connector) dealInsertHistorySettleMessage(columns []*pbe.Column) {
	rawdata := new(types.RawCanalData)
	rawdata.SQLId = []byte(columns[0].GetValue())
	fmt.Println(rawdata)
	c.Lock.Lock()
	c.RawData[HisSettle] = append(c.RawData[HisSettle], rawdata)
	c.Lock.Unlock()

}
func (c *Connector) dealInsertHistoryReceivableMessage(columns []*pbe.Column) {
	rawdata := new(types.RawCanalData)
	rawdata.SQLId = []byte(columns[0].GetValue())
	fmt.Println(rawdata)
	c.Lock.Lock()
	c.RawData[HisReceivable] = append(c.RawData[HisReceivable], rawdata)
	c.Lock.Unlock()

}
func (c *Connector) dealInsertIntensionMessage(columns []*pbe.Column) {
	rawdata := new(types.RawCanalData)
	rawdata.SQLId = []byte(columns[0].GetValue())
	fmt.Println(rawdata)
	c.Lock.Lock()
	c.RawData[Intension] = append(c.RawData[Intension], rawdata)
	c.Lock.Unlock()
}
func (c *Connector) dealInsertAccountMessage(columns []*pbe.Column) {
	rawdata := new(types.RawCanalData)
	rawdata.SQLId = []byte(columns[0].GetValue())
	fmt.Println(rawdata)
	c.Lock.Lock()
	c.RawData[Accounts] = append(c.RawData[Accounts], rawdata)
	c.Lock.Unlock()

}
func (c *Connector) dealInsertPoolPlanMessage(columns []*pbe.Column) {
	rawdata := new(types.RawCanalData)
	rawdata.SQLId = []byte(columns[0].GetValue())
	fmt.Println(rawdata)
	c.Lock.Lock()
	c.RawData[PoolPlan] = append(c.RawData[PoolPlan], rawdata)
	c.Lock.Unlock()

}
func (c *Connector) dealInsertPoolUsedMessage(columns []*pbe.Column) {
	rawdata := new(types.RawCanalData)
	rawdata.SQLId = []byte(columns[0].GetValue())
	fmt.Println(rawdata)
	c.Lock.Lock()
	c.RawData[PoolUsed] = append(c.RawData[PoolUsed], rawdata)
	c.Lock.Unlock()

}
func (c *Connector) dealInsertFinancingContractMessage(columns []*pbe.Column) {
	rawdata := new(types.RawCanalData)
	rawdata.SQLId = []byte(columns[0].GetValue())
	fmt.Println(rawdata)
	c.Lock.Lock()
	c.RawData[FinancingContract] = append(c.RawData[FinancingContract], rawdata)
	c.Lock.Unlock()

}
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
