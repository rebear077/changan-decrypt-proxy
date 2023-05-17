package sql

import (
	"strings"

	types "github.com/FISCO-BCOS/go-sdk/type"
	_ "github.com/go-sql-driver/mysql"
)

// 针对发票信息，进入的参数是解密后的明文，转换成结构体
func HandleInvoiceInfo(data []string) []*types.InvoiceInformation {
	//如果其他输入中存在[]怎么办？
	//最后返回的结果，目前是结构体的切片
	var INV []*types.InvoiceInformation
	for i := 0; i < len(data); i++ {
		str := data[i]
		str_split := strings.Split(str, ",")

		ICfo := types.InvoiceInformation{
			Certificateid:   str_split[0],
			Customerid:      str_split[1],
			Corpname:        str_split[2],
			Certificatetype: str_split[3],
			Intercustomerid: str_split[4],
			Invoicenotaxamt: str_split[5],
			Invoiceccy:      str_split[6],
			Sellername:      str_split[7],
			Invoicetype:     str_split[8],
			Buyername:       str_split[9],
			Buyerusccode:    str_split[10],
			Invoicedate:     str_split[11],
			Sellerusccode:   str_split[12],
			Invoicecode:     str_split[13],
			Invoicenum:      str_split[14],
			Checkcode:       str_split[15],
			Invoiceamt:      str_split[16],
			Owner:           str_split[17],
		}
		INV = append(INV, &ICfo)
	}
	return INV
}

// 针对历史交易信息的used infos，将解密后的明文转换成结构体
func HandleHistoricaltransactionUsedinfos(data []string) []*types.TempTransactionHistoryUsedinfos {
	var HUI []*types.TempTransactionHistoryUsedinfos
	for i := 0; i < len(data); i++ {
		str := data[i]
		usedinfo_split := strings.Split(str, ",")
		usedInfos := types.Usedinfos{
			Tradeyearmonth: usedinfo_split[7],
			Usedamount:     usedinfo_split[8],
			Ccy:            usedinfo_split[9],
		}
		trui := types.TempTransactionHistoryUsedinfos{
			Customergrade:   usedinfo_split[0],
			Certificatetype: usedinfo_split[1],
			Intercustomerid: usedinfo_split[2],
			Corpname:        usedinfo_split[3],
			Financeid:       usedinfo_split[4],
			Certificateid:   usedinfo_split[5],
			Customerid:      usedinfo_split[6],
			Usedinfos:       usedInfos,
		}
		HUI = append(HUI, &trui)
	}
	return HUI
}

// 针对历史交易信息的 settle infos，将解密后的明文转换成结构体
func HandleHistoricaltransactionSettleinfos(data []string) []*types.TempTransactionHistorySettleinfos {
	var HSI []*types.TempTransactionHistorySettleinfos
	for i := 0; i < len(data); i++ {
		str := data[i]
		settle_split := strings.Split(str, ",")
		settleInfos := types.Settleinfos{
			Tradeyearmonth: settle_split[7],
			Settleamount:   settle_split[8],
			Ccy:            settle_split[9],
		}
		trsi := types.TempTransactionHistorySettleinfos{
			Customergrade:   settle_split[0],
			Certificatetype: settle_split[1],
			Intercustomerid: settle_split[2],
			Corpname:        settle_split[3],
			Financeid:       settle_split[4],
			Certificateid:   settle_split[5],
			Customerid:      settle_split[6],
			Settleinfos:     settleInfos,
		}
		HSI = append(HSI, &trsi)
	}
	return HSI
}

// 针对历史交易信息的 order infos，将解密后的明文转换成结构体
func HandleHistoricaltransactionOrderinfos(data []string) []*types.TempTransactionHistoryOrderinfos {
	var HOI []*types.TempTransactionHistoryOrderinfos
	for i := 0; i < len(data); i++ {
		str := data[i]
		order_split := strings.Split(str, ",")
		orderInfos := types.Orderinfos{
			Tradeyearmonth: order_split[7],
			Orderamount:    order_split[8],
			Ccy:            order_split[9],
		}
		troi := types.TempTransactionHistoryOrderinfos{
			Customergrade:   order_split[0],
			Certificatetype: order_split[1],
			Intercustomerid: order_split[2],
			Corpname:        order_split[3],
			Financeid:       order_split[4],
			Certificateid:   order_split[5],
			Customerid:      order_split[6],
			Orderinfos:      orderInfos,
		}
		HOI = append(HOI, &troi)
	}
	return HOI
}

// 针对历史交易信息的 receivable infos，将解密后的明文转换成结构体
func HandleHistoricaltransactionReceivableinfos(data []string) []*types.TempTransactionHistoryReceivableinfos {
	var HRI []*types.TempTransactionHistoryReceivableinfos
	for i := 0; i < len(data); i++ {
		str := data[i]
		receivable_split := strings.Split(str, ",")
		RIfo := types.Receivableinfos{
			Tradeyearmonth:   receivable_split[7],
			Receivableamount: receivable_split[8],
			Ccy:              receivable_split[9],
		}

		trri := types.TempTransactionHistoryReceivableinfos{
			Customergrade:   receivable_split[0],
			Certificatetype: receivable_split[1],
			Intercustomerid: receivable_split[2],
			Corpname:        receivable_split[3],
			Financeid:       receivable_split[4],
			Certificateid:   receivable_split[5],
			Customerid:      receivable_split[6],
			Receivableinfos: RIfo,
		}
		HRI = append(HRI, &trri)
	}
	return HRI
}

// 针对入池信息的 plan infos，将解密后的明文转换成结构体
func HandleEnterpoolDataPlaninfos(data []string) []*types.TempEnterpoolDataPlaninfos {
	//如果其他输入中存在[]怎么办？

	var EPD []*types.TempEnterpoolDataPlaninfos
	for _, slice := range data {
		strs := strings.Split(slice, ",")
		planinfo := types.Planinfos{
			Tradeyearmonth: strs[5],
			Planamount:     strs[6],
			Currency:       strs[7],
		}
		poolPlan := types.TempEnterpoolDataPlaninfos{
			Datetimepoint:     strs[0],
			Ccy:               strs[1],
			Customerid:        strs[2],
			Intercustomerid:   strs[3],
			Receivablebalance: strs[4],
			Planinfos:         planinfo,
		}
		EPD = append(EPD, &poolPlan)
	}
	return EPD
}

// 针对入池信息的 provider used infos，将解密后的明文转换成结构体
func HandleEnterpoolDataProviderusedinfos(data []string) []*types.TempEnterpoolDataProviderusedinfos {
	var EDP []*types.TempEnterpoolDataProviderusedinfos
	for _, slice := range data {
		strs := strings.Split(slice, ",")
		usedinfo := types.Providerusedinfos{
			Tradeyearmonth: strs[5],
			Usedamount:     strs[6],
			Currency:       strs[7],
		}
		poolUsed := types.TempEnterpoolDataProviderusedinfos{
			Datetimepoint:     strs[0],
			Ccy:               strs[1],
			Customerid:        strs[2],
			Intercustomerid:   strs[3],
			Receivablebalance: strs[4],
			Providerusedinfos: usedinfo,
		}
		EDP = append(EDP, &poolUsed)
	}
	return EDP
}

// 处理融资意向信息，转换成结构体
func HandleFinancingIntention(data []string) []*types.FinancingIntention {
	var FCI []*types.FinancingIntention
	for i := 0; i < len(data); i++ {
		str := data[i]
		//fmt.Println(str)
		flag := 0
		header := ""
		for index, val := range str {
			if index+1 >= len(str) {
				break
			}
			if flag == 0 {
				if str[index] == ',' && str[index+1] == '[' {
					flag = 1
				} else {
					header = header + string(val)
				}
			}
		}
		header_split := strings.Split(header, ",")
		fcin := types.FinancingIntention{
			Custcdlinkposition: header_split[0],
			Custcdlinkname:     header_split[1],
			Certificateid:      header_split[2],
			Corpname:           header_split[3],
			Remark:             header_split[4],
			Bankcontact:        header_split[5],
			Banklinkname:       header_split[6],
			Custcdcontact:      header_split[7],
			Customerid:         header_split[8],
			Financeid:          header_split[9],
			Cooperationyears:   header_split[10],
			Certificatetype:    header_split[11],
			Intercustomerid:    header_split[12],
			State:              header_split[13],
		}
		FCI = append(FCI, &fcin)
	}
	// fmt.Println(FCI)
	return FCI
}

// 处理回款账户信息，转换成结构体
func HandleCollectionAccount(data []string) []*types.CollectionAccount {
	var COLA []*types.CollectionAccount
	for i := 0; i < len(data); i++ {
		str := data[i]
		//fmt.Println(str)
		flag := 0
		header := ""
		for index, val := range str {
			if index+1 >= len(str) {
				break
			}
			if flag == 0 {
				if str[index] == ',' && str[index+1] == '[' {
					flag = 1
				} else {
					header = header + string(val)
				}
			}
		}
		header_split := strings.Split(header, ",")
		cola := types.CollectionAccount{
			Backaccount:     header_split[0],
			Certificateid:   header_split[1],
			Customerid:      header_split[2],
			Corpname:        header_split[3],
			Lockremark:      header_split[4],
			Certificatetype: header_split[5],
			Intercustomerid: header_split[6],
		}
		COLA = append(COLA, &cola)
	}
	// fmt.Println(COLA)
	return COLA
}

// 处理借贷合同信息，转换成结构体
func HandleFinancingContract(data []*types.RawFinancingContractData) []*types.FinancingContract {
	var FC []*types.FinancingContract
	for _, v := range data {
		fc := types.FinancingContract{
			FinancingID: v.FinancingID,
			CustomerID:  v.CustomerID,
			CorpName:    v.CorpName,
			DebtMoney:   v.DebtMoney,
			SupplyDate:  v.SupplyDate,
			ExpireDate:  v.ExpireDate,
			Balance:     v.Balance,
		}
		FC = append(FC, &fc)

	}

	return FC
}

// 处理还款记录信息，转换成结构体
func HandleRepaymentRecord(data []*types.RawRepaymentRecord) []*types.RepaymentRecord {
	var records []*types.RepaymentRecord
	for _, v := range data {
		rr := types.RepaymentRecord{
			FinancingID:     v.FinancingID,
			CustomerID:      v.CustomerID,
			Time:            v.Time,
			RepaymentAmount: v.Repay,
		}
		records = append(records, &rr)
	}
	return records
}
