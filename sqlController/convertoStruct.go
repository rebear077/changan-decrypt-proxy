package sql

import (
	"encoding/json"
	"fmt"
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
func handleInvoiceInformation(data []string) []types.InvoiceInformation {
	//如果其他输入中存在[]怎么办？
	//最后返回的结果，目前是结构体的切片
	var INV []types.InvoiceInformation
	for i := 0; i < len(data); i++ {
		str := data[i]
		//fmt.Println(str)
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
		}
		INV = append(INV, ICfo)
	}
	// fmt.Println(INV)
	return INV
}

// 针对历史交易信息的used infos，将解密后的明文转换成结构体
func handleHistoricaltransactionUsedinfos(data []string) []types.TransactionHistoryUsedinfos {
	var HUI []types.TransactionHistoryUsedinfos
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
		// fmt.Println(trsh)
		HUI = append(HUI, trui)
	}
	// fmt.Println(HUI)
	return HUI
}

// 针对历史交易信息的 settle infos，将解密后的明文转换成结构体
func handleHistoricaltransactionSettleinfos(data []string) []types.TransactionHistorySettleinfos {
	var HSI []types.TransactionHistorySettleinfos
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
		HSI = append(HSI, trsi)
	}
	// fmt.Println(HSI)
	return HSI
}

// 针对历史交易信息的 order infos，将解密后的明文转换成结构体
func handleHistoricaltransactionOrderinfos(data []string) []types.TransactionHistoryOrderinfos {
	var HOI []types.TransactionHistoryOrderinfos
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
		HOI = append(HOI, troi)
	}
	// fmt.Println(HOI)
	return HOI
}

// 针对历史交易信息的 receivable infos，将解密后的明文转换成结构体
func handleHistoricaltransactionReceivableinfos(data []string) []types.TransactionHistoryReceivableinfos {
	var HRI []types.TransactionHistoryReceivableinfos
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
		HRI = append(HRI, trri)
	}
	// fmt.Println(HRI)
	return HRI
}

// 针对入池信息的 plan infos，将解密后的明文转换成结构体
func handleEnterpoolDataPlaninfos(data []string) []types.EnterpoolDataPlaninfos {
	//如果其他输入中存在[]怎么办？
	//最后返回的结果，目前是结构体的切片
	var EPD []types.EnterpoolDataPlaninfos
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
		EPD = append(EPD, epdt)
	}
	// fmt.Println(EPD)
	return EPD
}

// 针对入池信息的 provider used infos，将解密后的明文转换成结构体
func handleEnterpoolDataProviderusedinfos(data []string) []types.EnterpoolDataProviderusedinfos {
	var EPD []types.EnterpoolDataProviderusedinfos
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
		EPD = append(EPD, epdt)
	}
	// fmt.Println(EPD)
	return EPD
}

// 处理融资意向信息，转换成结构体
func handleFinancingIntention(data []string) []*types.FinancingIntention {
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
		}
		FCI = append(FCI, &fcin)
	}
	// fmt.Println(FCI)
	return FCI
}

// 处理回款账户信息，转换成结构体
func handleCollectionAccount(data []string) []*types.CollectionAccount {
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
func handleFinancingContract(data []types.RawFinancingContractData) []*types.FinancingContract {
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

// 将从数据库解密出来的数据从[]string先转换成结构体数组，然后转换成json
func (s *SqlCtr) ConvertoStruct(method string, data []string) string {
	switch method {
	case "HistoricaltransactionUsedinfos":
		result := handleHistoricaltransactionUsedinfos(data)
		ans, err := json.Marshal(result)
		if err != nil {
			fmt.Println(err)
		}
		// fmt.Println(string(ans))
		return string(ans)

	case "HistoricaltransactionSettleinfos":
		result := handleHistoricaltransactionSettleinfos(data)
		ans, err := json.Marshal(result)
		if err != nil {
			fmt.Println(err)
		}
		// fmt.Println(string(ans))
		return string(ans)

	case "HistoricaltransactionOrderinfos":
		result := handleHistoricaltransactionOrderinfos(data)
		ans, err := json.Marshal(result)
		if err != nil {
			fmt.Println(err)
		}
		// fmt.Println(string(ans))
		return string(ans)

	case "HistoricaltransactionReceivableinfos":
		result := handleHistoricaltransactionReceivableinfos(data)
		ans, err := json.Marshal(result)
		if err != nil {
			fmt.Println(err)
		}
		// fmt.Println(string(ans))
		return string(ans)

	case "InvoiceInformation":
		result := handleInvoiceInformation(data)
		ans, err := json.Marshal(result)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(ans))
		return string(ans)

	case "EnterpoolDataPlaninfos":
		result := handleEnterpoolDataPlaninfos(data)
		ans, err := json.Marshal(result)
		if err != nil {
			fmt.Println(err)
		}
		// fmt.Println(string(ans))
		return string(ans)

	case "EnterpoolDataUsedinfos":
		result := handleEnterpoolDataProviderusedinfos(data)
		ans, err := json.Marshal(result)
		if err != nil {
			fmt.Println(err)
		}
		// fmt.Println(string(ans))
		return string(ans)

	case "FinancingIntention":
		result := handleFinancingIntention(data)
		ans, err := json.Marshal(result)
		if err != nil {
			fmt.Println(err)
		}
		// fmt.Println(string(ans))
		return string(ans)

	case "CollectionAccount":
		result := handleCollectionAccount(data)
		ans, err := json.Marshal(result)
		if err != nil {
			fmt.Println(err)
		}
		// fmt.Println(string(ans))
		return string(ans)
	}
	return ""
}

// func handleHistoricaltransactionInformation(data []string) ([]TransactionHistoryUsedinfos, []TransactionHistorySettleinfos, []TransactionHistoryOrderinfos, []TransactionHistoryReceivableinfos) {
// 	var HUI []TransactionHistoryUsedinfos
// 	var HSI []TransactionHistorySettleinfos
// 	var HOI []TransactionHistoryOrderinfos
// 	var HRI []TransactionHistoryReceivableinfos
// 	for i := 0; i < len(data); i++ {
// 		str := data[i]
// 		flag := 0
// 		header := ""
// 		usedinfos := ""
// 		settleinfos := ""
// 		orderinfos := ""
// 		receivableinfos := ""
// 		for index, val := range str {
// 			if index+1 >= len(str) {
// 				break
// 			}
// 			if flag == 0 {
// 				if str[index] == ',' && str[index+1] == '[' {
// 					flag = 1
// 				} else {
// 					header = header + string(val)
// 				}
// 			} else if flag == 1 {
// 				if str[index] == '[' && str[index+1] == ',' {
// 					flag = 2
// 				} else if str[index] == ']' {
// 					flag = 2
// 				} else if str[index] != '[' && str[index] != ']' {
// 					usedinfos = usedinfos + string(val)
// 				}
// 			} else if flag == 2 {
// 				if str[index] == '[' && str[index+1] == ',' {
// 					flag = 3
// 				} else if str[index] == ']' {
// 					flag = 3
// 				} else if str[index] != '[' && str[index] != ']' {
// 					if len(settleinfos) == 0 && str[index] == ',' {
// 						continue
// 					} else {
// 						settleinfos = settleinfos + string(val)
// 					}
// 				}
// 			} else if flag == 3 {
// 				if str[index] == '[' && str[index+1] == ',' {
// 					flag = 4
// 				} else if str[index] == ']' {
// 					flag = 4
// 				} else if str[index] != '[' && str[index] != ']' {
// 					if len(orderinfos) == 0 && str[index] == ',' {
// 						continue
// 					} else {
// 						orderinfos = orderinfos + string(val)
// 					}
// 				}
// 			} else if flag == 4 {
// 				if str[index] == '[' && str[index+1] == ',' {
// 					flag = 5
// 				} else if str[index] == ']' {
// 					flag = 5
// 				} else if str[index] != '[' && str[index] != ']' {
// 					if len(receivableinfos) == 0 && str[index] == ',' {
// 						continue
// 					} else {
// 						receivableinfos = receivableinfos + string(val)
// 					}
// 				}
// 			}
// 		}
// 		header_split := strings.Split(header, ",")
// 		var UsedInfos []Usedinfos
// 		var SettleInfos []Settleinfos
// 		var OrderInfos []Orderinfos
// 		var ReceivableInfos []Receivableinfos

// 		usedinfos_split := strings.Split(usedinfos, "|")
// 		if usedinfos_split[0] != "" {
// 			for i := 0; i < len(usedinfos_split); i++ {
// 				us := strings.Split(usedinfos_split[i], ",")
// 				UIfo := Usedinfos{
// 					us[0],
// 					us[1],
// 					us[2],
// 				}
// 				UsedInfos = append(UsedInfos, UIfo)
// 			}
// 		}

// 		settleinfos_split := strings.Split(settleinfos, "|")
// 		if settleinfos_split[0] != "" {
// 			for i := 0; i < len(settleinfos_split); i++ {
// 				st := strings.Split(settleinfos_split[i], ",")
// 				SIfo := Settleinfos{
// 					st[0],
// 					st[1],
// 					st[2],
// 				}
// 				SettleInfos = append(SettleInfos, SIfo)
// 			}
// 		}

// 		orderinfos_split := strings.Split(orderinfos, "|")
// 		if orderinfos_split[0] != "" {
// 			for i := 0; i < len(orderinfos_split); i++ {
// 				od := strings.Split(orderinfos_split[i], ",")
// 				OIfo := Orderinfos{
// 					od[0],
// 					od[1],
// 					od[2],
// 				}
// 				OrderInfos = append(OrderInfos, OIfo)
// 			}
// 		}

// 		receivableinfos_split := strings.Split(receivableinfos, "|")
// 		if receivableinfos_split[0] != "" {
// 			for i := 0; i < len(receivableinfos_split); i++ {
// 				rc := strings.Split(receivableinfos_split[i], ",")
// 				RIfo := Receivableinfos{
// 					rc[0],
// 					rc[1],
// 					rc[2],
// 				}
// 				ReceivableInfos = append(ReceivableInfos, RIfo)
// 			}
// 		}

// 		trui := TransactionHistoryUsedinfos{
// 			header_split[0],
// 			header_split[1],
// 			header_split[2],
// 			header_split[3],
// 			header_split[4],
// 			header_split[5],
// 			header_split[6],
// 			UsedInfos,
// 		}
// 		// fmt.Println(trsh)
// 		HUI = append(HUI, trui)
// 		trsi := TransactionHistorySettleinfos{
// 			header_split[0],
// 			header_split[1],
// 			header_split[2],
// 			header_split[3],
// 			header_split[4],
// 			header_split[5],
// 			header_split[6],
// 			SettleInfos,
// 		}
// 		HSI = append(HSI, trsi)
// 		troi := TransactionHistoryOrderinfos{
// 			header_split[0],
// 			header_split[1],
// 			header_split[2],
// 			header_split[3],
// 			header_split[4],
// 			header_split[5],
// 			header_split[6],
// 			OrderInfos,
// 		}
// 		HOI = append(HOI, troi)
// 		trri := TransactionHistoryReceivableinfos{
// 			header_split[0],
// 			header_split[1],
// 			header_split[2],
// 			header_split[3],
// 			header_split[4],
// 			header_split[5],
// 			header_split[6],
// 			ReceivableInfos,
// 		}
// 		HRI = append(HRI, trri)
// 	}
// 	return HUI, HSI, HOI, HRI
// }

// func handleEnterpoolData(data []string) []EnterpoolData {
// 	//如果其他输入中存在[]怎么办？
// 	//最后返回的结果，目前是结构体的切片
// 	var EPD []EnterpoolData
// 	for i := 0; i < len(data); i++ {
// 		str := data[i]
// 		flag := 0
// 		header := ""
// 		planinfos := ""
// 		providerusedinfos := ""
// 		for index, val := range str {
// 			if index+1 >= len(str) {
// 				break
// 			}
// 			if flag == 0 {
// 				if str[index] == ',' && str[index+1] == '[' {
// 					flag = 1
// 				} else {
// 					header = header + string(val)
// 				}
// 			} else if flag == 1 {
// 				if str[index] == '[' && str[index+1] == ',' {
// 					flag = 2
// 				} else if str[index] == ']' {
// 					flag = 2
// 				} else if str[index] != '[' && str[index] != ']' {
// 					planinfos = planinfos + string(val)
// 				}
// 			} else if flag == 2 {
// 				if str[index] == '[' && str[index+1] == ',' {
// 					flag = 3
// 				} else if str[index] == ']' {
// 					flag = 3
// 				} else if str[index] != '[' && str[index] != ']' {
// 					if len(providerusedinfos) == 0 && str[index] == ',' {
// 						continue
// 					} else {
// 						providerusedinfos = providerusedinfos + string(val)
// 					}
// 				}
// 			}
// 		}
// 		header_split := strings.Split(header, ",")
// 		var PlanInfos []Planinfos
// 		planinfos_split := strings.Split(planinfos, "|")
// 		if planinfos_split[0] != "" {
// 			for i := 0; i < len(planinfos_split); i++ {
// 				pl := strings.Split(planinfos_split[i], ",")
// 				PLfo := Planinfos{
// 					pl[0],
// 					pl[1],
// 					pl[2],
// 				}
// 				PlanInfos = append(PlanInfos, PLfo)
// 			}
// 		}
// 		var ProviderusedInfos []Providerusedinfos
// 		providerusedinfos_split := strings.Split(providerusedinfos, "|")
// 		if providerusedinfos_split[0] != "" {
// 			for i := 0; i < len(providerusedinfos_split); i++ {
// 				pr := strings.Split(providerusedinfos_split[i], ",")
// 				PRfo := Providerusedinfos{
// 					pr[0],
// 					pr[1],
// 					pr[2],
// 				}
// 				ProviderusedInfos = append(ProviderusedInfos, PRfo)
// 			}
// 		}

// 		epdt := EnterpoolData{
// 			header_split[0],
// 			header_split[1],
// 			header_split[2],
// 			header_split[3],
// 			header_split[4],
// 			PlanInfos,
// 			ProviderusedInfos,
// 		}
// 		EPD = append(EPD, epdt)
// 	}
// 	fmt.Println(EPD)
// 	return EPD
// }
