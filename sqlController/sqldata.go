package sql

import (
	"strings"

	types "github.com/FISCO-BCOS/go-sdk/type"
)

func handleInvoiceInfo(data []string) []*types.InvoiceInformation {
	//如果其他输入中存在[]怎么办？
	//最后返回的结果，目前是结构体的切片
	var INV []*types.InvoiceInformation
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
		INV = append(INV, &ICfo)
	}
	// fmt.Println(INV)
	return INV
}

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
		// fmt.Println(trsh)
		HUI = append(HUI, &trui)
	}
	// fmt.Println(HUI)
	return HUI
}

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
	// fmt.Println(HSI)
	return HSI
}

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
	// fmt.Println(HOI)
	return HOI
}

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
	// fmt.Println(HRI)
	return HRI
}

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
	// fmt.Println(EPD)
	return EPD
}

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
	// fmt.Println(EPD)
	return EPD
}
