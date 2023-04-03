package main

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// 	"strconv"
// 	"strings"

// 	"github.com/FISCO-BCOS/go-sdk/proxy"
// )

// func relay(w http.ResponseWriter, r *http.Request) {

// 	r.ParseForm()
// 	relaybody, err := ioutil.ReadAll(r.Body)

// 	var jsonSring string //存储获取的jsonrpc命令字符串
// 	for k, _ := range r.Form {
// 		fmt.Printf("接收到的消息:%v\n", k)
// 		jsonSring = k
// 	}
// 	if jsonSring != "" {
// 		relaybody = []byte(jsonSring)
// 	}
// 	url := "http://127.0.0.1:8545"
// 	body := string(relaybody)
// 	response, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(body))) //将获取的jsonrpc命令直接转发给url指定的fisco地址
// 	if err != nil {
// 		fmt.Println("err:", err)
// 	}
// 	b, err := ioutil.ReadAll(response.Body) //等待fisco jsonrpc端口的回复

// 	w.Write(b) //将获取的回复回复给当前的http客户端

// }

// func paresTXInfo(w http.ResponseWriter, r *http.Request) {
// 	r.ParseForm()
// 	relaybody, _ := ioutil.ReadAll(r.Body)

// 	var jsonSring string //存储获取的jsonrpc命令字符串
// 	for k, _ := range r.Form {
// 		jsonSring = k
// 	}
// 	if jsonSring != "" {
// 		relaybody = []byte(jsonSring)
// 	}

// 	jsonCommand := new(proxy.JsonCommand)
// 	json.Unmarshal(relaybody, jsonCommand)

// 	fmt.Printf("jsonCommand.Method:%s\n", jsonCommand.Method)

// 	if jsonCommand.Method != "getBlockByNumber_sp" && jsonCommand.Method != "getBlockByNumber_all" {

// 		client := &http.Client{}
// 		url := "http://127.0.0.1:8545"
// 		body := string(relaybody)

// 		buffer := []byte(body)
// 		reader := bytes.NewReader(buffer)

// 		request, err := http.NewRequest("POST", url, reader)
// 		if err != nil {
// 			fmt.Println("GetHttpSkip Request Error:", err)
// 		}
// 		request.Header.Set("Access-Control-Allow-Origin", "*")
// 		request.Header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
// 		request.Header.Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
// 		request.Header.Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
// 		request.Header.Set("Access-Control-Allow-Credentials", "true")

// 		request.Header.Set("content-type", "application/json") //返回数据格式是json

// 		response, err := client.Do(request)
// 		if err != nil {
// 			fmt.Println()
// 		}
// 		b, err := ioutil.ReadAll(response.Body) //等待fisco jsonrpc端口的回复

// 		w.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
// 		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
// 		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
// 		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
// 		w.Header().Set("Access-Control-Allow-Credentials", "true")

// 		w.Header().Set("content-type", "application/json") //返回数据格式是json
// 		// fmt.Println("request header: ", request.Header)
// 		// fmt.Println("request: ", request)
// 		// fmt.Println("response: ", response)
// 		// fmt.Println("response body: ", response.Body)
// 		fmt.Println("ioutil.ReadAll(response.Body: ", string(b))
// 		w.Write(b) //将获取的回复回复给当前的http客户端
// 		return
// 	} else if jsonCommand.Method == "getBlockByNumber_all" {

// 		jsonCommand.Method = "getBlockNumber" //申请获取当前区块高度
// 		jsonCommand.Params = []interface{}{1}
// 		var relaybody_count []byte = make([]byte, 0)
// 		relaybody_count, _ = json.Marshal(jsonCommand)
// 		fmt.Println(string(relaybody_count))

// 		client := &http.Client{}
// 		url := "http://127.0.0.1:8545"
// 		//body := string(relaybody)

// 		buffer := []byte(relaybody_count)
// 		reader := bytes.NewReader(buffer)

// 		request, err := http.NewRequest("POST", url, reader)
// 		if err != nil {
// 			fmt.Println("GetHttpSkip Request Error:", err)
// 		}
// 		request.Header.Set("Access-Control-Allow-Origin", "*")
// 		request.Header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
// 		request.Header.Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
// 		request.Header.Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
// 		request.Header.Set("Access-Control-Allow-Credentials", "true")

// 		request.Header.Set("content-type", "application/json") //返回数据格式是json

// 		//response, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(body))) //将获取的jsonrpc命令直接转发给url指定的fisco地址
// 		//if err != nil {
// 		//	fmt.Println("err:", err)
// 		//}
// 		response, err := client.Do(request)
// 		if err != nil {
// 			fmt.Println()
// 		}

// 		b, err := ioutil.ReadAll(response.Body) //等待fisco jsonrpc端口的回复
// 		blockCount := new(proxy.JsonBlockNumber)
// 		json.Unmarshal(b, blockCount)
// 		// fmt.Println("request header: ", request.Header)
// 		// fmt.Println("request: ", request)
// 		// fmt.Println("response: ", response)
// 		// fmt.Println("response body: ", response.Body)
// 		// fmt.Println("ioutil.ReadAll(response.Body: ", string(b))
// 		fmt.Printf("当前区块数目为: %s\n", blockCount.Result)

// 		countInt, err := strconv.ParseInt(strings.TrimLeft(blockCount.Result, "0x"), 16, 64)
// 		if err != nil {
// 			fmt.Printf("解析十六进制字符串失败，err:%v\n", err)
// 		}

// 		var BlockArray proxy.BlockArray //包装所有从jsonrpc处获取的区块
// 		BlockArray.BlockCount = countInt
// 		fmt.Printf("countInt = %d\n", countInt)
// 		var i int64
// 		for i = 0; i <= countInt; i++ { //申请获取从0开始的到countInt的所有区块
// 			blockNumber := fmt.Sprintf("0x%x", i)
// 			jsonCommand.Method = "getBlockByNumber"
// 			jsonCommand.Params = []interface{}{1, blockNumber, true}
// 			var block []byte = make([]byte, 0)
// 			block, _ = json.Marshal(jsonCommand)

// 			client := &http.Client{}
// 			url := "http://127.0.0.1:8545"
// 			//body := string(relaybody)

// 			buffer := []byte(block)
// 			reader := bytes.NewReader(buffer)

// 			request, err := http.NewRequest("POST", url, reader)
// 			if err != nil {
// 				fmt.Println("GetHttpSkip Request Error:", err)
// 			}
// 			request.Header.Set("Access-Control-Allow-Origin", "*")
// 			request.Header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
// 			request.Header.Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
// 			request.Header.Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
// 			request.Header.Set("Access-Control-Allow-Credentials", "true")

// 			request.Header.Set("content-type", "application/json") //返回数据格式是json

// 			//body := string(block)
// 			fmt.Printf("本次申请获取区块的请求:%s\n", string(block))

// 			response, err := client.Do(request)
// 			if err != nil {
// 				fmt.Println()
// 			}
// 			b, err := ioutil.ReadAll(response.Body) //等待fisco jsonrpc端口的回复
// 			// fmt.Println("request header: ", request.Header)
// 			// fmt.Println("request: ", request)
// 			// fmt.Println("response: ", response)
// 			// fmt.Println("response body: ", response.Body)
// 			// fmt.Println("ioutil.ReadAll(response.Body: ", string(b))
// 			blockInfo := new(proxy.JsonBlockByNumber) //获取到对应编号的区块
// 			json.Unmarshal(b, blockInfo)
// 			BlockArray.Blocks = append(BlockArray.Blocks, blockInfo.Result) //将该区块填充
// 		}

// 		BlockArrayMarshal, err := json.Marshal(BlockArray)
// 		if err != nil {
// 			fmt.Println("BlockArray marshal is failed,err:", err)
// 		}

// 		w.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
// 		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
// 		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
// 		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
// 		w.Header().Set("Access-Control-Allow-Credentials", "true")

// 		w.Header().Set("content-type", "application/json") //返回数据格式是json

// 		w.Write(BlockArrayMarshal)

// 	} else if jsonCommand.Method == "getBlockByNumber_sp" {
// 		jsonCommand.Method = "getBlockByNumber"
// 		var relaybody_sp []byte = make([]byte, 0)
// 		relaybody_sp, _ = json.Marshal(jsonCommand)
// 		fmt.Println(string(relaybody_sp))

// 		client := &http.Client{}
// 		url := "http://127.0.0.1:8545"
// 		//body := string(relaybody)

// 		buffer := []byte(relaybody_sp)
// 		reader := bytes.NewReader(buffer)

// 		request, err := http.NewRequest("POST", url, reader)
// 		if err != nil {
// 			fmt.Println("GetHttpSkip Request Error:", err)
// 		}
// 		request.Header.Set("Access-Control-Allow-Origin", "*")
// 		request.Header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
// 		request.Header.Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
// 		request.Header.Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
// 		request.Header.Set("Access-Control-Allow-Credentials", "true")

// 		request.Header.Set("content-type", "application/json") //返回数据格式是json

// 		response, err := client.Do(request)
// 		if err != nil {
// 			fmt.Println()
// 		}
// 		b, err := ioutil.ReadAll(response.Body) //等待fisco jsonrpc端口的回复
// 		// fmt.Println("request header: ", request.Header)
// 		// fmt.Println("request: ", request)
// 		// fmt.Println("response: ", response)
// 		// fmt.Println("response body: ", response.Body)
// 		// fmt.Println("ioutil.ReadAll(response.Body: ", string(b))
// 		blockInfo := new(proxy.JsonBlockByNumber)
// 		json.Unmarshal(b, blockInfo)

// 		txList := new(proxy.TransactionArray)
// 		txList.Transactions = make([]interface{}, 0)
// 		txList.TxCount = len(blockInfo.Result.Transactions)

// 		for index, tx := range blockInfo.Result.Transactions {
// 			fmt.Printf("这是第%d笔交易的输入input:%v\n", index, tx.Input)
// 			resStrArray, msgType := proxy.Decoder([]byte(tx.Input))
// 			switch msgType {
// 			case "issueAPChannelInfo":
// 				apInfo := proxy.DecodeAPtoString(resStrArray["issueAPInfo"], tx.Hash)
// 				apInfo.Method = "Update APInfo FromDB"
// 				txList.Transactions = append(txList.Transactions, apInfo)

// 			case "updateAPChannelInfo":
// 				apInfo := proxy.DecodeAPtoString(resStrArray["updateAPInfo"], tx.Hash)
// 				apInfo.Method = "Update APInfo FromDB"
// 				txList.Transactions = append(txList.Transactions, apInfo)

// 			case "issueBidingPriceInfo":
// 				bidPrice := proxy.DecodeBPtoString(resStrArray["issueBid"], tx.Hash)
// 				bidPrice.Method = "Buyer:Purchase Request"
// 				txList.Transactions = append(txList.Transactions, bidPrice)

// 			case "issueChannelSwitchInfo":
// 				chSwitch := proxy.DecodeCStoString(resStrArray["issueChSwitch"], tx.Hash)
// 				chSwitch.Method = "Seller:Willingness to sell"
// 				txList.Transactions = append(txList.Transactions, chSwitch)

// 			case "issueChannelDealInfo":
// 				chDeal := proxy.DecodeCDtoString(resStrArray["issueChDeal"], tx.Hash)
// 				chDeal.Method = "Transaction Confirm"
// 				txList.Transactions = append(txList.Transactions, chDeal)

// 			case "queryAPChannelInfo":
// 				fmt.Printf("%s : %s\n ", msgType, resStrArray["addr"])
// 				query := new(proxy.QueryAPChannelInfo)
// 				query.Method = "Update database"
// 				query.TxHash = tx.Hash
// 				query.APaddrStr = resStrArray["addr"]

// 				txList.Transactions = append(txList.Transactions, query)

// 			case "false":
// 				fmt.Println("未识别的消息类型.......")
// 			}
// 		}

// 		txListMarshal, err := json.Marshal(txList)
// 		if err != nil {
// 			fmt.Println("txList marshal is failed,err:", err)
// 		}

// 		// txListStr := string(txListMarshal)
// 		// txListStrRlp := strings.Replace(txListStr, "\\", "", -1)
// 		// txListStrRlp += "\n"
// 		// fmt.Println(txListStrRlp)

// 		w.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
// 		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
// 		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
// 		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
// 		w.Header().Set("Access-Control-Allow-Credentials", "true")

// 		w.Header().Set("content-type", "application/json") //返回数据格式是json

// 		w.Write(txListMarshal)
// 	}
// 	fmt.Println(w.Header())
// }

// func main() {
// 	//http.HandleFunc("/", relay) //设置访问的路由
// 	http.HandleFunc("/handle", paresTXInfo) //设置访问的路由
// 	//http.HandleFunc("/txInfo", paresTXInfo)
// 	err := http.ListenAndServe(":8999", nil) //设置监听的端口
// 	if err != nil {
// 		log.Fatal("ListenAndServe: ", err)
// 	}
// }
