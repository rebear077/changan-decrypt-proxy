package proxy

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/FISCO-BCOS/go-sdk/abi"
)

func Decoder(input []byte) (map[string]string, string) {
	ParseIssueResult, err := abi.JSON(strings.NewReader(DBControllerABI))
	if err != nil {
		fmt.Printf("parse abi failed, err: %v\n", err)
	}

	inputStr := string(input)
	methodID := inputStr[2:10] //前八位是方法id
	resStrArray := make(map[string]string, 0)

	switch methodID {
	case "784103d2":
		ret, err := ParseIssueResult.UnpackInput("issueInvoiceInformationStorage", common.FromHex(inputStr)[4:]) /////
		if err != nil {
			fmt.Printf("parse return value failed, err: %v\n", err)
			return nil, "false"
		}
		parseRet, ok := ret.([]interface{})
		if !ok {
			fmt.Println("issueInvoiceInformationStorage解析错误.......")
		} else {
			// addr, _ := parseRet[0].(string)
			issueInvoiceInfo, _ := parseRet[1].(string)
			// fmt.Printf("节点%s 上传APInfo:%s\n", addr, issueInvoiceInfo)
			// resStrArray["addr"] = addr
			resStrArray["issueInvoiceInfo"] = issueInvoiceInfo
			return resStrArray, "issueInvoiceInformationStorage"
		}
	// case "b4e2c79b":
	// 	ret, err := ParseIssueResult.UnpackInput("updateAPChannelInfo", common.FromHex(inputStr)[4:])
	// 	if err != nil {
	// 		fmt.Printf("parse return value failed, err: %v\n", err)
	// 		return nil, "false"
	// 	}
	// 	parseRet, ok := ret.([]interface{})
	// 	if !ok {
	// 		fmt.Println("updateAPChannelInfo解析错误.......")
	// 	} else {
	// 		addr, _ := parseRet[0].(string)
	// 		updateAPInfo, _ := parseRet[1].(string)
	// 		fmt.Printf("节点%s 上传APInfo:%s\n", addr, updateAPInfo)
	// 		resStrArray["addr"] = addr
	// 		resStrArray["updateAPInfo"] = updateAPInfo
	// 		return resStrArray, "updateAPChannelInfo"
	// 	}
	case "737ad475":
		ret, err := ParseIssueResult.UnpackInput("issueHistoricalSettleInformation", common.FromHex(inputStr)[4:])
		if err != nil {
			fmt.Printf("parse return value failed, err: %v\n", err)
			return nil, "false"
		}
		parseRet, ok := ret.([]interface{})
		fmt.Println(parseRet)
		if !ok {
			fmt.Println("issueHistoricalSettleInformation解析错误.......")
		} else {
			issueHistoricalSettleInfo, _ := parseRet[1].(string)
			fmt.Println(issueHistoricalSettleInfo)
			resStrArray["issueHistoricalSettleInfo"] = issueHistoricalSettleInfo
			return resStrArray, "issueHistoricalSettleInformation"
		}
	case "05f2a57a":
		ret, err := ParseIssueResult.UnpackInput("issueHistoricalReceivableInformation", common.FromHex(inputStr)[4:])
		if err != nil {
			fmt.Printf("parse return value failed, err: %v\n", err)
			return nil, "false"
		}
		parseRet, ok := ret.([]interface{})
		fmt.Println(parseRet)
		if !ok {
			fmt.Println("issueHistoricalReceivableInformation解析错误.......")
		} else {
			issueHistoricalReceivableInfo, _ := parseRet[1].(string)
			fmt.Println(issueHistoricalReceivableInfo)
			resStrArray["issueHistoricalReceivableInfo"] = issueHistoricalReceivableInfo
			return resStrArray, "issueHistoricalReceivableInformation"
		}
	case "6e382033":
		ret, err := ParseIssueResult.UnpackInput("issueHistoricalUsedInformation", common.FromHex(inputStr)[4:])
		if err != nil {
			fmt.Printf("parse return value failed, err: %v\n", err)
			return nil, "false"
		}
		parseRet, ok := ret.([]interface{})
		fmt.Println(parseRet)
		if !ok {
			fmt.Println("issueHistoricalUsedInformation解析错误.......")
		} else {
			issueHistoricalUsedInfo, _ := parseRet[1].(string)
			fmt.Println(issueHistoricalUsedInfo)
			resStrArray["issueHistoricalOrderInfo"] = issueHistoricalUsedInfo
			return resStrArray, " issueHistoricalUsedInformation"
		}
	case "66fffcbc":
		ret, err := ParseIssueResult.UnpackInput("issueHistoricalOrderInformation", common.FromHex(inputStr)[4:])
		if err != nil {
			fmt.Printf("parse return value failed, err: %v\n", err)
			return nil, "false"
		}
		parseRet, ok := ret.([]interface{})
		fmt.Println(parseRet)
		if !ok {
			fmt.Println("issueHistoricalOrderInformation解析错误.......")
		} else {
			issueHistoricalOrderInfo, _ := parseRet[1].(string)
			fmt.Println(issueHistoricalOrderInfo)
			resStrArray["issueHistoricalOrderInfo"] = issueHistoricalOrderInfo
			return resStrArray, "issueHistoricalOrderInformation"
		}
	case "ebf1db00":
		ret, err := ParseIssueResult.UnpackInput("issuePoolPlanInformation", common.FromHex(inputStr)[4:])
		if err != nil {
			fmt.Printf("parse return value failed, err: %v\n", err)
			return nil, "false"
		}
		parseRet, ok := ret.([]interface{})
		fmt.Println(parseRet)
		if !ok {
			fmt.Println("issuePoolPlanInformation解析错误.......")
		} else {
			issuePoolPlanInfo, _ := parseRet[1].(string)
			fmt.Println(issuePoolPlanInfo)
			resStrArray["issuePoolPlanInfo"] = issuePoolPlanInfo
			return resStrArray, "issuePoolPlanInformation"
		}
	case "89dd75fd":
		ret, err := ParseIssueResult.UnpackInput("issuePoolUsedInformation", common.FromHex(inputStr)[4:])
		if err != nil {
			fmt.Printf("parse return value failed, err: %v\n", err)
			return nil, "false"
		}
		parseRet, ok := ret.([]interface{})
		fmt.Println(parseRet)
		if !ok {
			fmt.Println("issuePoolUsedInformation解析错误.......")
		} else {
			issuePoolUsedInfo, _ := parseRet[1].(string)
			fmt.Println(issuePoolUsedInfo)
			resStrArray["issuePoolUsedInfo"] = issuePoolUsedInfo
			return resStrArray, "issuePoolUsedInformation"
		}
	}
	return nil, "false"
}

const DBControllerABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_time\",\"type\":\"string\"}],\"name\":\"queryBidingPriceInfo\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"time\",\"type\":\"string\"},{\"name\":\"_bidingPrice\",\"type\":\"string\"}],\"name\":\"issueBidingPriceInfo\",\"outputs\":[{\"name\":\"\",\"type\":\"int256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"time\",\"type\":\"string\"},{\"name\":\"_channeldeal\",\"type\":\"string\"},{\"name\":\"_buyeraddr\",\"type\":\"string\"},{\"name\":\"_selleraddr\",\"type\":\"string\"},{\"name\":\"_key\",\"type\":\"string\"}],\"name\":\"issueChannelDealInfo\",\"outputs\":[{\"name\":\"\",\"type\":\"int256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"addr\",\"type\":\"string\"},{\"name\":\"_apchannelInfo\",\"type\":\"string\"}],\"name\":\"issueAPChannelInfo\",\"outputs\":[{\"name\":\"\",\"type\":\"int256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_time\",\"type\":\"string\"}],\"name\":\"queryChannelSwitchInfo\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"addr\",\"type\":\"string\"},{\"name\":\"_apchannelInfo\",\"type\":\"string\"}],\"name\":\"updateAPChannelInfo\",\"outputs\":[{\"name\":\"\",\"type\":\"int256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_addr\",\"type\":\"string\"}],\"name\":\"queryAPChannelInfo\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_time\",\"type\":\"string\"}],\"name\":\"queryChannelDealInfo\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"time\",\"type\":\"string\"},{\"name\":\"_channelswitch\",\"type\":\"string\"}],\"name\":\"issueChannelSwitchInfo\",\"outputs\":[{\"name\":\"\",\"type\":\"int256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_addr\",\"type\":\"string\"}],\"name\":\"queryAPChannelInfoInString\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"apchannelInfo\",\"type\":\"string\"},{\"indexed\":false,\"name\":\"addr\",\"type\":\"string\"}],\"name\":\"IssueAPChannelInfo\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"apchannelInfo\",\"type\":\"string\"},{\"indexed\":false,\"name\":\"addr\",\"type\":\"string\"}],\"name\":\"UpdateAPChannelInfo\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"bidingPrice\",\"type\":\"string\"},{\"indexed\":false,\"name\":\"time\",\"type\":\"string\"}],\"name\":\"IssueBidingPriceInfo\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"channeldeal\",\"type\":\"string\"},{\"indexed\":false,\"name\":\"time\",\"type\":\"string\"}],\"name\":\"IssueChannelDealInfo\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"channelswitch\",\"type\":\"string\"},{\"indexed\":false,\"name\":\"time\",\"type\":\"string\"}],\"name\":\"IssueChannelSwitchInfo\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"buyeraddr\",\"type\":\"string\"},{\"indexed\":false,\"name\":\"selleraddr\",\"type\":\"string\"}],\"name\":\"AutoUpdateAPChannelInfo\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"}]"
