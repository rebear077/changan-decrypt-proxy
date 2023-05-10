package sql

import (
	"strings"

	types "github.com/FISCO-BCOS/go-sdk/type"
	_ "github.com/go-sql-driver/mysql"
)

// 针对历史交易和入池信息这两个表单，由于包含多个list，需要将header和list的内容分别提取成header和info，header包含多个字段。字段之间通过逗号分割，list形式为[**,**,**|**,**,**]
func sliceinfohandler(str string) (string, string) {
	flag := 0
	header := ""
	infos := ""
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
		} else if flag == 1 {
			//应该是防止有[,]的情况，即子表单中无内容
			if str[index] == '[' && str[index+1] == ',' {
				flag = 2
			} else if str[index] == ']' {
				flag = 2
			} else if str[index] != '[' && str[index] != ']' {
				infos = infos + string(val)
			}
		} else if flag == 2 {
			break
		}
	}
	return header, infos
}

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
func HandleHistoricaltransactionUsedinfos(data []string) []*types.TransactionHistoryUsedinfos {
	var HUI []*types.TransactionHistoryUsedinfos
	for i := 0; i < len(data); i++ {
		str := data[i]
		header, usedinfos := sliceinfohandler(str)
		header_split := strings.Split(header, ",")
		var UsedInfos []types.Usedinfos
		usedinfos_split := strings.Split(usedinfos, "|")
		if usedinfos_split[0] != "" {
			for i := 0; i < len(usedinfos_split); i++ {
				us := strings.Split(usedinfos_split[i], ",")
				UIfo := types.Usedinfos{
					Tradeyearmonth: us[0],
					Usedamount:     us[1],
					Ccy:            us[2],
				}
				UsedInfos = append(UsedInfos, UIfo)
			}
		}
		trui := types.TransactionHistoryUsedinfos{

			Customergrade:   header_split[0],
			Certificatetype: header_split[1],
			Intercustomerid: header_split[2],
			Corpname:        header_split[3],
			Financeid:       header_split[4],
			Certificateid:   header_split[5],
			Customerid:      header_split[6],
			Usedinfos:       UsedInfos,
		}
		HUI = append(HUI, &trui)
	}
	return HUI
}

// 针对历史交易信息的 settle infos，将解密后的明文转换成结构体
func HandleHistoricaltransactionSettleinfos(data []string) []*types.TransactionHistorySettleinfos {
	var HSI []*types.TransactionHistorySettleinfos
	for i := 0; i < len(data); i++ {
		str := data[i]
		header, settleinfos := sliceinfohandler(str)
		header_split := strings.Split(header, ",")
		var SettleInfos []types.Settleinfos
		settleinfos_split := strings.Split(settleinfos, "|")
		if settleinfos_split[0] != "" {
			for i := 0; i < len(settleinfos_split); i++ {
				st := strings.Split(settleinfos_split[i], ",")
				SIfo := types.Settleinfos{
					Tradeyearmonth: st[0],
					Settleamount:   st[1],
					Ccy:            st[2],
				}
				SettleInfos = append(SettleInfos, SIfo)
			}
		}
		trsi := types.TransactionHistorySettleinfos{
			Customergrade:   header_split[0],
			Certificatetype: header_split[1],
			Intercustomerid: header_split[2],
			Corpname:        header_split[3],
			Financeid:       header_split[4],
			Certificateid:   header_split[5],
			Customerid:      header_split[6],
			Settleinfos:     SettleInfos,
		}
		HSI = append(HSI, &trsi)
	}
	return HSI
}

// 针对历史交易信息的 order infos，将解密后的明文转换成结构体
func HandleHistoricaltransactionOrderinfos(data []string) []*types.TransactionHistoryOrderinfos {
	var HOI []*types.TransactionHistoryOrderinfos
	for i := 0; i < len(data); i++ {
		str := data[i]
		header, orderinfos := sliceinfohandler(str)
		header_split := strings.Split(header, ",")
		var OrderInfos []types.Orderinfos
		orderinfos_split := strings.Split(orderinfos, "|")
		if orderinfos_split[0] != "" {
			for i := 0; i < len(orderinfos_split); i++ {
				od := strings.Split(orderinfos_split[i], ",")
				OIfo := types.Orderinfos{
					Tradeyearmonth: od[0],
					Orderamount:    od[1],
					Ccy:            od[2],
				}
				OrderInfos = append(OrderInfos, OIfo)
			}
		}
		troi := types.TransactionHistoryOrderinfos{
			Customergrade:   header_split[0],
			Certificatetype: header_split[1],
			Intercustomerid: header_split[2],
			Corpname:        header_split[3],
			Financeid:       header_split[4],
			Certificateid:   header_split[5],
			Customerid:      header_split[6],
			Orderinfos:      OrderInfos,
		}
		HOI = append(HOI, &troi)
	}
	return HOI
}

// 针对历史交易信息的 receivable infos，将解密后的明文转换成结构体
func HandleHistoricaltransactionReceivableinfos(data []string) []*types.TransactionHistoryReceivableinfos {
	var HRI []*types.TransactionHistoryReceivableinfos
	for i := 0; i < len(data); i++ {
		str := data[i]
		header, receivableinfos := sliceinfohandler(str)
		header_split := strings.Split(header, ",")
		var ReceivableInfos []types.Receivableinfos
		receivableinfos_split := strings.Split(receivableinfos, "|")
		if receivableinfos_split[0] != "" {
			for i := 0; i < len(receivableinfos_split); i++ {
				rc := strings.Split(receivableinfos_split[i], ",")
				RIfo := types.Receivableinfos{
					Tradeyearmonth:   rc[0],
					Receivableamount: rc[1],
					Ccy:              rc[2],
				}
				ReceivableInfos = append(ReceivableInfos, RIfo)
			}
		}
		trri := types.TransactionHistoryReceivableinfos{
			Customergrade:   header_split[0],
			Certificatetype: header_split[1],
			Intercustomerid: header_split[2],
			Corpname:        header_split[3],
			Financeid:       header_split[4],
			Certificateid:   header_split[5],
			Customerid:      header_split[6],
			Receivableinfos: ReceivableInfos,
		}
		HRI = append(HRI, &trri)
	}
	return HRI
}

// 针对入池信息的 plan infos，将解密后的明文转换成结构体
func HandleEnterpoolDataPlaninfos(data []string) []*types.EnterpoolDataPlaninfos {
	//如果其他输入中存在[]怎么办？
	//最后返回的结果，目前是结构体的切片
	var EPD []*types.EnterpoolDataPlaninfos
	for i := 0; i < len(data); i++ {
		str := data[i]

		header, planinfos := sliceinfohandler(str)

		header_split := strings.Split(header, ",")
		var PlanInfos []types.Planinfos
		planinfos_split := strings.Split(planinfos, "|")
		if planinfos_split[0] != "" {
			for i := 0; i < len(planinfos_split); i++ {
				pl := strings.Split(planinfos_split[i], ",")
				PLfo := types.Planinfos{
					Tradeyearmonth: pl[0],
					Planamount:     pl[1],
					Currency:       pl[2],
				}
				PlanInfos = append(PlanInfos, PLfo)
			}
		}

		epdt := types.EnterpoolDataPlaninfos{
			Datetimepoint:     header_split[0],
			Ccy:               header_split[1],
			Customerid:        header_split[2],
			Intercustomerid:   header_split[3],
			Receivablebalance: header_split[4],
			Planinfos:         PlanInfos,
		}
		EPD = append(EPD, &epdt)
	}
	return EPD
}

// 针对入池信息的 provider used infos，将解密后的明文转换成结构体
func HandleEnterpoolDataProviderusedinfos(data []string) []*types.EnterpoolDataProviderusedinfos {
	var EPD []*types.EnterpoolDataProviderusedinfos
	for i := 0; i < len(data); i++ {
		str := data[i]
		header, providerusedinfos := sliceinfohandler(str)
		header_split := strings.Split(header, ",")
		var ProviderusedInfos []types.Providerusedinfos
		providerusedinfos_split := strings.Split(providerusedinfos, "|")
		if providerusedinfos_split[0] != "" {
			for i := 0; i < len(providerusedinfos_split); i++ {
				pr := strings.Split(providerusedinfos_split[i], ",")
				PRfo := types.Providerusedinfos{
					Tradeyearmonth: pr[0],
					Usedamount:     pr[1],
					Currency:       pr[2],
				}
				ProviderusedInfos = append(ProviderusedInfos, PRfo)
			}
		}

		epdt := types.EnterpoolDataProviderusedinfos{
			Datetimepoint:     header_split[0],
			Ccy:               header_split[1],
			Customerid:        header_split[2],
			Intercustomerid:   header_split[3],
			Receivablebalance: header_split[4],
			Providerusedinfos: ProviderusedInfos,
		}
		EPD = append(EPD, &epdt)
	}
	return EPD
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

// 将从数据库解密出来的数据从[]string先转换成结构体数组，然后转换成json
// func (s *SqlCtr) ConvertoStruct(method string, data []string) string {
// 	switch method {
// 	case "HistoricaltransactionUsedinfos":
// 		result := HandleHistoricaltransactionUsedinfos(data)
// 		ans, err := json.Marshal(result)
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 		// fmt.Println(string(ans))
// 		return string(ans)

// 	case "HistoricaltransactionSettleinfos":
// 		result := HandleHistoricaltransactionSettleinfos(data)
// 		ans, err := json.Marshal(result)
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 		// fmt.Println(string(ans))
// 		return string(ans)

// 	case "HistoricaltransactionOrderinfos":
// 		result := HandleHistoricaltransactionOrderinfos(data)
// 		ans, err := json.Marshal(result)
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 		// fmt.Println(string(ans))
// 		return string(ans)

// 	case "HistoricaltransactionReceivableinfos":
// 		result := HandleHistoricaltransactionReceivableinfos(data)
// 		ans, err := json.Marshal(result)
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 		// fmt.Println(string(ans))
// 		return string(ans)

// 	case "InvoiceInformation":
// 		result := HandleInvoiceInfo(data)
// 		ans, err := json.Marshal(result)
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 		fmt.Println(string(ans))
// 		return string(ans)

// 	case "EnterpoolDataPlaninfos":
// 		result := HandleEnterpoolDataPlaninfos(data)
// 		ans, err := json.Marshal(result)
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 		// fmt.Println(string(ans))
// 		return string(ans)

// 	case "EnterpoolDataUsedinfos":
// 		result := HandleEnterpoolDataProviderusedinfos(data)
// 		ans, err := json.Marshal(result)
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 		// fmt.Println(string(ans))
// 		return string(ans)

// 	case "FinancingIntention":
// 		result := HandleFinancingIntention(data)
// 		ans, err := json.Marshal(result)
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 		// fmt.Println(string(ans))
// 		return string(ans)

// 	case "CollectionAccount":
// 		result := HandleCollectionAccount(data)
// 		ans, err := json.Marshal(result)
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 		// fmt.Println(string(ans))
// 		return string(ans)
// 	}
// 	return ""
// }
