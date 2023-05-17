package types

type RawCanalData struct {
	SQLId []byte
}

// 数据库中原生发票数据
type RawInvoiceData struct {
	SQLId      string
	Num        string
	Status     string
	ID         string
	CustomerID string
	Time       string
	Type       string
	InvoiceNum string
	Data       string
	Key        string
	Hash       string
	Owner      string
}

// 数据库中原生历史交易信息
// 四种公用一种数据结构
type RawHistoryTransData struct {
	SqlId          string
	Num            string
	Status         string
	Id             string
	CustomerID     string
	TradeYearMonth string
	FinanceId      string
	Data           string
	Key            string
	Hash           string
	Owner          string
}

// 数据库中原生入池信息
type RawEnterPoolData struct {
	SqlId          string
	Num            string
	Status         string
	Id             string
	CustomerID     string
	TradeYearMonth string
	Data           string
	Key            string
	Hash           string
	Owner          string
}

// 数据库中原生的融资意向申请
type RawFinancingData struct {
	SQLId     string
	Num       string
	Status    string
	ID        string
	FinanceId string
	Data      string
	Key       string
	Hash      string
	State     string
}

// 数据库中原生回款账户
type RawAccountsData struct {
	SQLId     string
	Num       string
	Status    string
	ID        string
	FinanceID string
	Data      string
	Key       string
	Hash      string
	State     string
}

// 数据库中原生借贷合同信息
type RawFinancingContractData struct {
	SQLId       string
	Num         string
	Status      string
	ID          string
	FinancingID string
	CustomerID  string
	CorpName    string
	DebtMoney   string
	SupplyDate  string
	ExpireDate  string
	Balance     string
}

// 数据库中原生的还款记录
type RawRepaymentRecord struct {
	SQLId       string
	Num         string
	Status      string
	ID          string
	FinancingID string
	CustomerID  string
	Repay       string
	Time        string
}

// 发票信息查询
type InvoiceInformation struct {
	Certificateid   string `json:"certificateId"`
	Customerid      string `json:"customerId"`
	Corpname        string `json:"corpName"`
	Certificatetype string `json:"certificateType"`
	Intercustomerid string `json:"interCustomerId"`
	Invoicenotaxamt string `json:"invoiceNotaxAmt"`
	Invoiceccy      string `json:"invoiceCcy"`
	Sellername      string `json:"sellerName"`
	Invoicetype     string `json:"invoiceType"`
	Buyername       string `json:"buyerName"`
	Buyerusccode    string `json:"buyerUsccode"`
	Invoicedate     string `json:"invoiceDate"`
	Sellerusccode   string `json:"sellerUsccode"`
	Invoicecode     string `json:"invoiceCode"`
	Invoicenum      string `json:"invoiceNum"`
	Checkcode       string `json:"checkCode"`
	Invoiceamt      string `json:"invoiceAmt"`
	Owner           string `json:"owner"`
}

// 返回给前端的发票信息
type InvoiceInformationReturn struct {
	InvoiceInformationList []*InvoiceInformation `json:"invoiceInformationList"`
	TotalCount             int                   `json:"totalcount"`
	CurrentPage            int                   `json:"currentPage"`
}

// 推送历史交易信息接口
type TransactionHistory struct {
	Customergrade   string            `json:"customerGrade"`
	Certificatetype string            `json:"certificateType"`
	Intercustomerid string            `json:"interCustomerId"`
	Corpname        string            `json:"corpName"`
	Financeid       string            `json:"financeId"`
	Certificateid   string            `json:"certificateId"`
	Customerid      string            `json:"customerId"`
	Usedinfos       []Usedinfos       `json:"usedInfos"`
	Settleinfos     []Settleinfos     `json:"settleInfos"`
	Orderinfos      []Orderinfos      `json:"orderInfos"`
	Receivableinfos []Receivableinfos `json:"receivableInfos"`
}

// 返回给前端的历史交易信息
type TransactionHistoryReturn struct {
	TransactionHistoryList []*TransactionHistory `json:"transactionHistoryList"`
	TotalCount             int                   `json:"totalcount"`
	CurrentPage            int                   `json:"currentPage"`
}

// 历史交易信息
type TransactionHistoryHeader struct {
	Customergrade   string `json:"customerGrade"`
	Certificatetype string `json:"certificateType"`
	Intercustomerid string `json:"interCustomerId"`
	Corpname        string `json:"corpName"`
	Financeid       string `json:"financeId"`
	Certificateid   string `json:"certificateId"`
	Customerid      string `json:"customerId"`
}
type TempTransactionHistoryUsedinfos struct {
	Customergrade   string
	Certificatetype string
	Intercustomerid string
	Corpname        string
	Financeid       string
	Certificateid   string
	Customerid      string
	Usedinfos       Usedinfos
}
type TransactionHistoryUsedinfos struct {
	Customergrade   string      `json:"customerGrade"`
	Certificatetype string      `json:"certificateType"`
	Intercustomerid string      `json:"interCustomerId"`
	Corpname        string      `json:"corpName"`
	Financeid       string      `json:"financeId"`
	Certificateid   string      `json:"certificateId"`
	Customerid      string      `json:"customerId"`
	Usedinfos       []Usedinfos `json:"usedInfos"`
}
type TempTransactionHistorySettleinfos struct {
	Customergrade   string
	Certificatetype string
	Intercustomerid string
	Corpname        string
	Financeid       string
	Certificateid   string
	Customerid      string
	Settleinfos     Settleinfos
}
type TransactionHistorySettleinfos struct {
	Customergrade   string        `json:"customerGrade"`
	Certificatetype string        `json:"certificateType"`
	Intercustomerid string        `json:"interCustomerId"`
	Corpname        string        `json:"corpName"`
	Financeid       string        `json:"financeId"`
	Certificateid   string        `json:"certificateId"`
	Customerid      string        `json:"customerId"`
	Settleinfos     []Settleinfos `json:"settleInfos"`
}
type TempTransactionHistoryOrderinfos struct {
	Customergrade   string
	Certificatetype string
	Intercustomerid string
	Corpname        string
	Financeid       string
	Certificateid   string
	Customerid      string
	Orderinfos      Orderinfos
}
type TransactionHistoryOrderinfos struct {
	Customergrade   string       `json:"customerGrade"`
	Certificatetype string       `json:"certificateType"`
	Intercustomerid string       `json:"interCustomerId"`
	Corpname        string       `json:"corpName"`
	Financeid       string       `json:"financeId"`
	Certificateid   string       `json:"certificateId"`
	Customerid      string       `json:"customerId"`
	Orderinfos      []Orderinfos `json:"orderInfos"`
}
type TempTransactionHistoryReceivableinfos struct {
	Customergrade   string
	Certificatetype string
	Intercustomerid string
	Corpname        string
	Financeid       string
	Certificateid   string
	Customerid      string
	Receivableinfos Receivableinfos
}
type TransactionHistoryReceivableinfos struct {
	Customergrade   string            `json:"customerGrade"`
	Certificatetype string            `json:"certificateType"`
	Intercustomerid string            `json:"interCustomerId"`
	Corpname        string            `json:"corpName"`
	Financeid       string            `json:"financeId"`
	Certificateid   string            `json:"certificateId"`
	Customerid      string            `json:"customerId"`
	Receivableinfos []Receivableinfos `json:"receivableInfos"`
}
type Usedinfos struct {
	Tradeyearmonth string `json:"tradeYearMonth"`
	Usedamount     string `json:"usedAmount"`
	Ccy            string `json:"ccy"`
}
type Settleinfos struct {
	Tradeyearmonth string `json:"tradeYearMonth"`
	Settleamount   string `json:"settleAmount"`
	Ccy            string `json:"ccy"`
}
type Orderinfos struct {
	Tradeyearmonth string `json:"tradeYearMonth"`
	Orderamount    string `json:"orderAmount"`
	Ccy            string `json:"ccy"`
}
type Receivableinfos struct {
	Tradeyearmonth   string `json:"tradeYearMonth"`
	Receivableamount string `json:"receivableAmount"`
	Ccy              string `json:"ccy"`
}

// 入池信息
type EnterpoolDataHeader struct {
	Datetimepoint     string `json:"dateTimePoint"`
	Ccy               string `json:"ccy"`
	Customerid        string `json:"customerId"`
	Intercustomerid   string `json:"interCustomerId"`
	Receivablebalance string `json:"receivableBalance"`
}
type TempEnterpoolDataPlaninfos struct {
	Datetimepoint     string
	Ccy               string
	Customerid        string
	Intercustomerid   string
	Receivablebalance string
	Planinfos         Planinfos
}
type EnterpoolDataPlaninfos struct {
	Datetimepoint     string      `json:"dateTimePoint"`
	Ccy               string      `json:"ccy"`
	Customerid        string      `json:"customerId"`
	Intercustomerid   string      `json:"interCustomerId"`
	Receivablebalance string      `json:"receivableBalance"`
	Planinfos         []Planinfos `json:"planInfos"`
}
type TempEnterpoolDataProviderusedinfos struct {
	Datetimepoint     string
	Ccy               string
	Customerid        string
	Intercustomerid   string
	Receivablebalance string
	Providerusedinfos Providerusedinfos
}
type EnterpoolDataProviderusedinfos struct {
	Datetimepoint     string              `json:"dateTimePoint"`
	Ccy               string              `json:"ccy"`
	Customerid        string              `json:"customerId"`
	Intercustomerid   string              `json:"interCustomerId"`
	Receivablebalance string              `json:"receivableBalance"`
	Providerusedinfos []Providerusedinfos `json:"ProviderUsedInfos"`
}

// 推送入池数据接口
type EnterpoolData struct {
	Datetimepoint     string              `json:"dateTimePoint"`
	Ccy               string              `json:"ccy"`
	Customerid        string              `json:"customerId"`
	Intercustomerid   string              `json:"interCustomerId"`
	Receivablebalance string              `json:"receivableBalance"`
	Planinfos         []Planinfos         `json:"planInfos"`
	Providerusedinfos []Providerusedinfos `json:"ProviderUsedInfos"`
}

type EnterpoolDataReturn struct {
	EnterpoolDataList []*EnterpoolData `json:"enterpoolDataList"`
	TotalCount        int              `json:"totalcount"`
	CurrentPage       int              `json:"currentPage"`
}

type Planinfos struct {
	Tradeyearmonth string `json:"tradeYearMonth"`
	Planamount     string `json:"planAmount"`
	Currency       string `json:"currency"`
}
type Providerusedinfos struct {
	Tradeyearmonth string `json:"tradeYearMonth"`
	Usedamount     string `json:"usedAmount"`
	Currency       string `json:"currency"`
}

// 提交融资意向接口
type FinancingIntention struct {
	Custcdlinkposition string `json:"custcdLinkPosition"`
	Custcdlinkname     string `json:"custcdLinkName"`
	Certificateid      string `json:"certificateId"`
	Corpname           string `json:"corpName"`
	Remark             string `json:"remark"`
	Bankcontact        string `json:"bankContact"`
	Banklinkname       string `json:"bankLinkName"`
	Custcdcontact      string `json:"custcdContact"`
	Customerid         string `json:"customerId"`
	Financeid          string `json:"financeId"`
	Cooperationyears   string `json:"cooperationYears"`
	Certificatetype    string `json:"certificateType"`
	Intercustomerid    string `json:"interCustomerId"`
	State              string `json:"state"`
}

// 返回给前端的融资意向数据结构
type FinancingIntentionReturn struct {
	FinancingIntentionList []*FinancingIntention `json:"financingIntentionList"`
	TotalCount             int                   `json:"totalcount"`
	CurrentPage            int                   `json:"currentPage"`
}

// 推送回款账户接口
type CollectionAccount struct {
	Backaccount     string `json:"backAccount"`
	Certificateid   string `json:"certificateId"`
	Customerid      string `json:"customerId"`
	Corpname        string `json:"corpName"`
	Lockremark      string `json:"lockRemark"`
	Certificatetype string `json:"certificateType"`
	Intercustomerid string `json:"interCustomerId"`
}

type CollectionAccountReturn struct {
	CollectionAccountList []*CollectionAccount `json:"collectionAccountList"`
	TotalCount            int                  `json:"totalcount"`
	CurrentPage           int                  `json:"currentPage"`
}

type RawInvoiceInformation struct {
	Certificateid   string         `json:"certificateId"`
	Customerid      string         `json:"customerId"`
	Corpname        string         `json:"corpName"`
	Certificatetype string         `json:"certificateType"`
	Intercustomerid string         `json:"interCustomerId"`
	Invoiceinfos    []Invoiceinfos `json:"invoiceInfos"`
}

type Invoiceinfos struct {
	Invoicenotaxamt string `json:"invoiceNotaxAmt"`
	Invoiceccy      string `json:"invoiceCcy"`
	Sellername      string `json:"sellerName"`
	Invoicetype     string `json:"invoiceType"`
	Buyername       string `json:"buyerName"`
	Buyerusccode    string `json:"buyerUsccode"`
	Invoicedate     string `json:"invoiceDate"`
	Sellerusccode   string `json:"sellerUsccode"`
	Invoicecode     string `json:"invoiceCode"`
	Invoicenum      string `json:"invoiceNum"`
	Checkcode       string `json:"checkCode"`
	Invoiceamt      string `json:"invoiceAmt"`
}

type FinancingContract struct {
	FinancingID string `json:"financingID"`
	CustomerID  string `json:"customerID"`
	CorpName    string `json:"corpName"`
	DebtMoney   string `json:"debtMoney"`
	SupplyDate  string `json:"supplyDate"`
	ExpireDate  string `json:"expireDate"`
	Balance     string `json:"balance"`
}
type FinancingContractReturn struct {
	FinancingContractList []*FinancingContract `json:"financingContractInformationList"`
	TotalCount            int                  `json:"totalcount"`
	CurrentPage           int                  `json:"currentPage"`
}
type RepaymentRecord struct {
	FinancingID     string `json:"financingID"`
	CustomerID      string `json:"customerID"`
	Time            string `json:"time"`
	RepaymentAmount string `json:"repaymentAmount"`
}
type RepaymentRecordReturn struct {
	RepaymentList []*RepaymentRecord `json:"repaymentRecordList"`
	TotalCount    int                `json:"totalcount"`
	CurrentPage   int                `json:"currentPage"`
}
