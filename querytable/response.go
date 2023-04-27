package querytable

import (
	"fmt"

	types "github.com/FISCO-BCOS/go-sdk/type"
)

func wrongJsonType() string {
	jsonData := []byte(`{
		"msg":"错误的Json格式",
		"result": "{}",
		"code":"ERR04"
	}`)
	return string(jsonData)
}

//	func sucessCode() string {
//		jsonData := []byte(`{
//			"msg":"",
//			"result": "{}",
//			"code":"SUC000000"
//		}`)
//		return string(jsonData)
//	}
func failedCode() string {
	jsonData := []byte(`{
		"msg":"",
		"result": "{}",
		"code":"failed"
	}`)
	return string(jsonData)
}
func unconsistencyCode() string {
	jsonData := []byte(`{
		"msg":"",
		"result": "{}",
		"code":"unconsistency"
	}`)
	return string(jsonData)
}
func verifyConsistency(data types.SelectedInfoToApplication) bool {
	var checkQueue string = ""
	for _, invoice := range data.Invoice {
		if checkQueue == "" {
			checkQueue = invoice.Customerid
		} else {
			if checkQueue != invoice.Customerid {
				return false
			}
		}
	}
	checkQueue = ""
	for _, historyInfo := range data.HistoryInfo {
		if checkQueue == "" {
			checkQueue = historyInfo.Customerid
		} else {
			if checkQueue != historyInfo.Customerid {
				return false
			}
		}
	}
	checkQueue = ""
	for _, poolInfo := range data.PoolInfo {
		if checkQueue == "" {
			checkQueue = poolInfo.Customerid
		} else {
			if checkQueue != poolInfo.Customerid {
				return false
			}
		}
	}
	return true
}
func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

type SucessCode struct {
	Msg    string
	Result string
	Code   string
}
