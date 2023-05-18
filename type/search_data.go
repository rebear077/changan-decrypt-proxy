package types

type InvoiceInformationSearch struct {
	Id          string
	Time        string
	InvoiceType string
	FinanceID   string
	InvoiceNum  string
	PageId      string
	SearchType  string
}

type HistoryTransactionSearch struct {
	Id             string
	Tradeyearmonth string
	FinanceId      string
	PageId         string
	SearchType     string
}

type PoolDataSearch struct {
	Id             string
	Tradeyearmonth string
	PageId         string
	SearchType     string
}

type FinancingIntentionSearch struct {
	Id         string
	FinanceId  string
	PageId     string
	SearchType string
}

type CollectionAccountSearch struct {
	Id         string
	FinanceId  string
	PageId     string
	SearchType string
}
type FinancingContractSearch struct {
	Id         string
	FinanceId  string
	PageId     string
	SearchType string
}
type RepaymentRecordSearch struct {
	Id         string
	FinanceId  string
	PageId     string
	SearchType string
}
type SelectedInfoToApplication struct {
	FinanceIntention FinancingIntention   `json:"financingIntension"`
	Invoice          []InvoiceInformation `json:"invoice"`
}
