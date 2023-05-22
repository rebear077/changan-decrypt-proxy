package proxy

import (
	"fmt"
	"strconv"
	"strings"
)

func DecodeInvoiceInfotoString(invoiceinfo string, txHash string) InvoiceInfo {
	// var err error
	var revData InvoiceInfo = InvoiceInfo{}
	array := strings.Split(invoiceinfo, ",")
	revData.Method = array[0]
	return revData
}

func DecodeHistoricalSettleInfotoString(HistoricalSettleinfo string, txHash string) HistoricalSettleInfo {
	// var err error
	var revData HistoricalSettleInfo = HistoricalSettleInfo{}
	array := strings.Split(HistoricalSettleinfo, ",")
	revData.Method = array[0]
	return revData
}

func DecodeHistoricalOrderInfotoString(HistoricalOrderinfo string, txHash string) HistoricalOrderInfo {
	// var err error
	var revData HistoricalOrderInfo = HistoricalOrderInfo{}
	array := strings.Split(HistoricalOrderinfo, ",")
	revData.Method = array[0]
	return revData
}
func DecodeHistoricalUsedInfotoString(HistoricalUsedinfo string, txHash string) HistoricalUsedInfo {
	// var err error
	var revData HistoricalUsedInfo = HistoricalUsedInfo{}
	array := strings.Split(HistoricalUsedinfo, ",")
	revData.Method = array[0]
	return revData
}
func DecodePoolUsedInfotoString(PoolUsedinfo string, txHash string) PoolUsedInfo {
	// var err error
	var revData PoolUsedInfo = PoolUsedInfo{}
	array := strings.Split(PoolUsedinfo, ",")
	revData.Method = array[0]
	return revData
}
func DecodePoolPlanInfotoString(PoolPlaninfo string, txHash string) PoolPlanInfo {
	// var err error
	var revData PoolPlanInfo = PoolPlanInfo{}
	array := strings.Split(PoolPlaninfo, ",")
	revData.Method = array[0]
	return revData
}
func DecodeHistoricalReceivableInfotoString(HistoricalReceivableinfo string, txHash string) HistoricalReceivableInfo {
	// var err error
	var revData HistoricalReceivableInfo = HistoricalReceivableInfo{}
	array := strings.Split(HistoricalReceivableinfo, ",")
	revData.Method = array[0]
	return revData
}

// 解析update/issueAPInfo字符串
func DecodeAPtoString(apChInfo string, txHash string) APChannelInfo {
	var err error
	var revData APChannelInfo = APChannelInfo{
		APUsingChannelList: make([]int, 0),
	}

	array := strings.Split(apChInfo, ",")

	revData.APaddr = ParseAddrFromString(array[0])
	revData.APaddrStr = fmt.Sprintf("%x", revData.APaddr)
	revData.APid = array[1]
	revData.APCoverageRadius, _ = strconv.ParseFloat(array[2], 64)
	revData.APLocation.Lat, _ = strconv.ParseFloat(array[3], 64)
	revData.APLocation.Lng, _ = strconv.ParseFloat(array[4], 64)
	revData.APLocationStr = "Lat:" + array[3] + ", " + "Lng:" + array[4]

	revData.APUsingChannelListStr = array[5]
	if revData.APUsingChannelList, err = FromStringToIntSlice(array[5]); err != nil {
		fmt.Println(err)
	}
	revData.APChannelOccupancy, _ = strconv.ParseFloat(array[6], 64)
	revData.APPower, _ = strconv.ParseFloat(array[7], 64)
	revData.SpectrumCoinNumber, _ = strconv.ParseFloat(array[8], 64)
	revData.SpectrumCoinNumber /= BalanceFactor

	revData.ApIpinfo = array[9]

	revData.Method = "updateAPChannelInfo"
	revData.TxHash = txHash
	return revData
}

// 解析BidPriceInfo字符串
func DecodeBPtoString(bidPriceInfo string, txHash string) BidingPriceInfo {
	var revData BidingPriceInfo = BidingPriceInfo{
		PricesList: make([]Seller, 0),
	}
	array := strings.Split(bidPriceInfo, ",")
	revData.BuyerAPaddr = ParseAddrFromString(array[1])
	revData.BuyerAPaddrStr = fmt.Sprintf("%x", revData.BuyerAPaddr)
	revData.BuyerAPid = array[2]
	temp := strings.TrimSpace(array[3])
	list := strings.Split(temp, ";")

	for i, v := range list {
		if v == " " {
			continue
		}
		elem := strings.Split(v, "|")
		temp := Seller{}

		temp.SellerAPaddr = ParseAddrFromString(elem[0])
		temp.SellerAPaddrStr = fmt.Sprintf("%x", temp.SellerAPaddr)
		temp.SellerAPid = elem[1]
		temp.ChannelPrice, _ = strconv.ParseFloat(elem[2], 64)
		temp.ChannelPrice /= BalanceFactor
		revData.PricesList = append(revData.PricesList, temp)

		revData.PricesListStr += "SellerAPaddr:" + elem[0] + fmt.Sprintf(" || ") + "SellerAPid:" + elem[1]
		if i < len(list)-1 {
			revData.PricesListStr += fmt.Sprintf(" || ")
		}
	}
	revData.Method = "issueBidingPriceInfo"
	revData.TxHash = txHash
	return revData
}

// 解析channelSwitch字符串
func DecodeCStoString(chSwitchInfo string, txHash string) ChannelSwitchInfo {
	var revData ChannelSwitchInfo = ChannelSwitchInfo{
		SwitchChannels: make([]int, 0),
	}
	array := strings.Split(chSwitchInfo, ",")

	revData.TimeStamp = array[0]
	revData.BuyerAPaddr = ParseAddrFromString(array[1])
	revData.BuyerAPaddrStr = fmt.Sprintf("%x", revData.BuyerAPaddr)
	revData.BuyerAPid = array[2]
	revData.SellerAPaddr = ParseAddrFromString(array[3])
	revData.SellerAPaddrStr = fmt.Sprintf("%x", revData.SellerAPaddr)
	revData.SellerAPid = array[4]
	temp, _ := FromStringToIntSlice(array[5])
	revData.SwitchChannels = append(revData.SwitchChannels, temp...)
	revData.TotalPrice, _ = strconv.ParseFloat(array[6], 64)
	revData.TotalPrice /= BalanceFactor

	revData.Method = "issueChannelSwitchInfo"
	revData.TxHash = txHash
	return revData

}

// 解析channel Deal字符串
func DecodeCDtoString(chDealInfo string, txHash string) ChannelDealInfo {
	var revData ChannelDealInfo = ChannelDealInfo{
		SwitchChannels: make([]int, 0),
	}
	array := strings.Split(chDealInfo, ",")
	//revData.ReferredTradeTxID = src.TransactionHash
	revData.BuyerAPaddr = ParseAddrFromString(array[1])
	revData.BuyerAPaddrStr = fmt.Sprintf("%x", revData.BuyerAPaddr)
	revData.BuyerAPid = array[2]
	revData.SellerAPaddr = ParseAddrFromString(array[3])
	revData.SellerAPaddrStr = fmt.Sprintf("%x", revData.SellerAPaddr)
	revData.SellerAPid = array[4]
	temp, _ := FromStringToIntSlice(array[5])
	revData.SwitchChannels = append(revData.SwitchChannels, temp...)
	revData.SendAmount, _ = strconv.ParseFloat(array[6], 64)
	revData.SendAmount /= BalanceFactor

	revData.Method = "issueChannelDealInfo"
	revData.TxHash = txHash
	return revData
}
