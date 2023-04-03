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
	case "a5daec12":
		ret, err := ParseIssueResult.UnpackInput("issueAPChannelInfo", common.FromHex(inputStr)[4:]) /////
		if err != nil {
			fmt.Printf("parse return value failed, err: %v\n", err)
			return nil, "false"
		}
		parseRet, ok := ret.([]interface{})
		if !ok {
			fmt.Println("issueAPChannelInfo解析错误.......")
		} else {
			addr, _ := parseRet[0].(string)
			issueAPInfo, _ := parseRet[1].(string)
			fmt.Printf("节点%s 上传APInfo:%s\n", addr, issueAPInfo)
			resStrArray["addr"] = addr
			resStrArray["issueAPInfo"] = issueAPInfo
			return resStrArray, "issueAPChannelInfo"
		}
	case "b4e2c79b": //
		ret, err := ParseIssueResult.UnpackInput("updateAPChannelInfo", common.FromHex(inputStr)[4:])
		if err != nil {
			fmt.Printf("parse return value failed, err: %v\n", err)
			return nil, "false"
		}
		parseRet, ok := ret.([]interface{})
		if !ok {
			fmt.Println("updateAPChannelInfo解析错误.......")
		} else {
			addr, _ := parseRet[0].(string)
			updateAPInfo, _ := parseRet[1].(string)
			fmt.Printf("节点%s 上传APInfo:%s\n", addr, updateAPInfo)
			resStrArray["addr"] = addr
			resStrArray["updateAPInfo"] = updateAPInfo
			return resStrArray, "updateAPChannelInfo"
		}
	case "9a898c74":
		ret, err := ParseIssueResult.UnpackInput("issueBidingPriceInfo", common.FromHex(inputStr)[4:])
		if err != nil {
			fmt.Printf("parse return value failed, err: %v\n", err)
			return nil, "false"
		}
		parseRet, ok := ret.([]interface{})
		fmt.Println(parseRet)
		if !ok {
			fmt.Println("issueBidingPriceInfo解析错误.......")
		} else {
			issueBid, _ := parseRet[1].(string)
			fmt.Println(issueBid)
			resStrArray["issueBid"] = issueBid
			return resStrArray, "issueBidingPriceInfo"
		}
	case "ebe3efe5":
		ret, err := ParseIssueResult.UnpackInput("issueChannelSwitchInfo", common.FromHex(inputStr)[4:])
		if err != nil {
			fmt.Printf("parse return value failed, err: %v\n", err)
			return nil, "false"
		}
		parseRet, ok := ret.([]interface{})
		fmt.Println(parseRet)
		if !ok {
			fmt.Println("issueChannelSwitchInfo解析错误.......")
		} else {
			issueChSwitch, _ := parseRet[1].(string)
			fmt.Println(issueChSwitch)
			resStrArray["issueChSwitch"] = issueChSwitch
			return resStrArray, "issueChannelSwitchInfo"
		}
	case "a351d088":
		ret, err := ParseIssueResult.UnpackInput("issueChannelDealInfo", common.FromHex(inputStr)[4:])
		if err != nil {
			fmt.Printf("parse return value failed, err: %v\n", err)
			return nil, "false"
		}
		parseRet, ok := ret.([]interface{})
		fmt.Println(parseRet)
		if !ok {
			fmt.Println("issueChannelDealInfo解析错误.......")
		} else {
			issueChDeal, _ := parseRet[1].(string)
			fmt.Println(issueChDeal)
			resStrArray["issueChDeal"] = issueChDeal
			return resStrArray, "issueChannelDealInfo"
		}
	case "f2ba888f":
		ret, err := ParseIssueResult.UnpackInput("queryAPChannelInfo", common.FromHex(inputStr)[4:])
		if err != nil {
			fmt.Printf("parse return value failed, err: %v\n", err)
			return nil, "false"
		}
		parseRet, ok := ret.([]interface{})
		if !ok {
			fmt.Println("queryAPChannelInfo解析错误.......")
		} else {
			addr, _ := parseRet[0].(string)
			fmt.Println("节点：", addr, "发起queryAPChannelInfo请求.....")
			resStrArray["addr"] = addr
			return resStrArray, "queryAPChannelInfo"
		}
	}
	return nil, "false"
}

const DBControllerABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_time\",\"type\":\"string\"}],\"name\":\"queryBidingPriceInfo\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"time\",\"type\":\"string\"},{\"name\":\"_bidingPrice\",\"type\":\"string\"}],\"name\":\"issueBidingPriceInfo\",\"outputs\":[{\"name\":\"\",\"type\":\"int256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"time\",\"type\":\"string\"},{\"name\":\"_channeldeal\",\"type\":\"string\"},{\"name\":\"_buyeraddr\",\"type\":\"string\"},{\"name\":\"_selleraddr\",\"type\":\"string\"},{\"name\":\"_key\",\"type\":\"string\"}],\"name\":\"issueChannelDealInfo\",\"outputs\":[{\"name\":\"\",\"type\":\"int256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"addr\",\"type\":\"string\"},{\"name\":\"_apchannelInfo\",\"type\":\"string\"}],\"name\":\"issueAPChannelInfo\",\"outputs\":[{\"name\":\"\",\"type\":\"int256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_time\",\"type\":\"string\"}],\"name\":\"queryChannelSwitchInfo\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"addr\",\"type\":\"string\"},{\"name\":\"_apchannelInfo\",\"type\":\"string\"}],\"name\":\"updateAPChannelInfo\",\"outputs\":[{\"name\":\"\",\"type\":\"int256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_addr\",\"type\":\"string\"}],\"name\":\"queryAPChannelInfo\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_time\",\"type\":\"string\"}],\"name\":\"queryChannelDealInfo\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"time\",\"type\":\"string\"},{\"name\":\"_channelswitch\",\"type\":\"string\"}],\"name\":\"issueChannelSwitchInfo\",\"outputs\":[{\"name\":\"\",\"type\":\"int256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_addr\",\"type\":\"string\"}],\"name\":\"queryAPChannelInfoInString\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"apchannelInfo\",\"type\":\"string\"},{\"indexed\":false,\"name\":\"addr\",\"type\":\"string\"}],\"name\":\"IssueAPChannelInfo\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"apchannelInfo\",\"type\":\"string\"},{\"indexed\":false,\"name\":\"addr\",\"type\":\"string\"}],\"name\":\"UpdateAPChannelInfo\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"bidingPrice\",\"type\":\"string\"},{\"indexed\":false,\"name\":\"time\",\"type\":\"string\"}],\"name\":\"IssueBidingPriceInfo\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"channeldeal\",\"type\":\"string\"},{\"indexed\":false,\"name\":\"time\",\"type\":\"string\"}],\"name\":\"IssueChannelDealInfo\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"channelswitch\",\"type\":\"string\"},{\"indexed\":false,\"name\":\"time\",\"type\":\"string\"}],\"name\":\"IssueChannelSwitchInfo\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"buyeraddr\",\"type\":\"string\"},{\"indexed\":false,\"name\":\"selleraddr\",\"type\":\"string\"}],\"name\":\"AutoUpdateAPChannelInfo\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"}]"
