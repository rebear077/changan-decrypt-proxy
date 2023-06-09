package querytable

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strconv"
	"strings"
	"time"

	server "github.com/FISCO-BCOS/go-sdk/backend"
	"github.com/FISCO-BCOS/go-sdk/conf"
	"github.com/FISCO-BCOS/go-sdk/proxy"
	sql "github.com/FISCO-BCOS/go-sdk/sqlController"
	types "github.com/FISCO-BCOS/go-sdk/type"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type FrontEnd struct {
	Server *server.Server
	url    *conf.Config
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// TokenResponse represents the response structure for token
type TokenResponse struct {
	Token string `json:"token"`
}

// JWTClaims represents the JWT claims
type JWTClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var jwtKey = []byte("your_secret_key")

var users = make(map[string]User)
var passwordResetTokens = make(map[string]time.Time)

func NewFrontEnd() *FrontEnd {
	server := server.NewServer()
	configs, err := conf.ParseConfigFile("./configs/config.toml")
	if err != nil {
		logrus.Fatalln(err)
	}
	config := &configs[0]
	return &FrontEnd{
		url:    config,
		Server: server,
	}
}

func (front *FrontEnd) DecryptSelectToApplicationInformation(writer http.ResponseWriter, request *http.Request) {
	var message types.SelectedInfoToApplication
	if json.NewDecoder(request.Body).Decode(&message) != nil {
		jsonData := wrongJsonType()
		fmt.Fprint(writer, jsonData)
	} else {
		if !verifyConsistency(message) {
			jsonData := unconsistencyCode()
			fmt.Fprint(writer, jsonData)
		} else {
			targetURL := front.url.TargetUrl
			targetJSON, err := json.Marshal(message)
			if err != nil {
				logrus.Errorln(err)
				jsonData := failedCode()
				fmt.Fprint(writer, jsonData)
				return
			}
			fmt.Println(string(targetJSON))
			response := handle(targetURL, string(targetJSON))
			jsonData, err := json.Marshal(response)
			if err != nil {
				logrus.Errorln(err)
				jsonData := failedCode()
				fmt.Fprint(writer, jsonData)
				return
			}
			fmt.Println("收到回执", string(jsonData))
			fmt.Fprint(writer, jsonData)
		}

	}
}

// 发票信息查询URL
func (front *FrontEnd) DecryptInvoiceInformation(writer http.ResponseWriter, request *http.Request) {
	order := make(map[string]string)
	Sql := sql.NewSqlCtr()
	slice := Sql.InvoiceInformationIndex(request)
	order["id"] = slice.Id
	order["financeID"] = slice.FinanceID
	order["invoiceType"] = slice.InvoiceType
	order["num"] = slice.InvoiceNum
	order["time"] = slice.Time
	order["pageid"] = slice.PageId
	currentPage, _ := strconv.Atoi(order["pageid"])
	order["searchType"] = slice.SearchType
	invoices, totalcount := front.Server.SearchInvoiceFromRedis(order)
	jsonData := front.Server.PackToInvoiceJson(invoices, totalcount, currentPage)
	fmt.Fprint(writer, jsonData)
}

// 历史交易信息URL
func (front *FrontEnd) DecryptHistoricaltransaction(writer http.ResponseWriter, request *http.Request) {
	order := make(map[string]string)
	Sql := sql.NewSqlCtr()
	slice := Sql.HistoryTransactionIndex(request)
	fmt.Println(slice)
	order["id"] = slice.Id
	order["tradeyearmonth"] = slice.Tradeyearmonth
	order["financeId"] = slice.FinanceId
	order["pageid"] = slice.PageId
	currentPage, _ := strconv.Atoi(order["pageid"])
	order["searchType"] = slice.SearchType
	fmt.Println(order)
	txs, totalcount := front.Server.SearchHistoryTXFromRedis(order)
	jsonData := front.Server.PackToHistoryTXJson(txs, totalcount, currentPage)
	fmt.Fprint(writer, jsonData)

}

// 入池数据
func (front *FrontEnd) DecryptEnterPoolData(writer http.ResponseWriter, request *http.Request) {
	order := make(map[string]string)
	Sql := sql.NewSqlCtr()
	slice := Sql.PoolDataIndex(request)
	order["id"] = slice.Id
	order["tradeyearmonth"] = slice.Tradeyearmonth
	order["pageid"] = slice.PageId
	currentPage, _ := strconv.Atoi(order["pageid"])
	order["searchType"] = slice.SearchType
	enterpools, totalcount := front.Server.SearchEnterPoolFromRedis(order)
	// fmt.Println(txs[0])
	jsonData := front.Server.PackToEnterPoolJson(enterpools, totalcount, currentPage)
	fmt.Fprint(writer, jsonData)
}

func (front *FrontEnd) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		jsonResponse(w, http.StatusMethodNotAllowed, "方法不是Post")
		return
	}
	DB := sql.NewSqlCtr()
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		jsonResponse(w, http.StatusBadRequest, "数据格式不对")
		return
	}

	// Check if the user exists
	// storedUser, ok := users[user.Username]
	// if !ok {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	jsonResponse(w, http.StatusUnauthorized, "用户名不存在")
	// 	return
	// }
	exists, err := DB.CheckUsernameExists(user.Username)
	if err != nil {
		logrus.Errorln(err)
	}
	if !exists {
		w.WriteHeader(http.StatusConflict)
		jsonResponse(w, http.StatusConflict, "用户名不存在")
		return
	}

	// Compare hashed password with the provided password
	password, err := DB.GetPasswordByUsername(user.Username)
	if err != nil {
		logrus.Errorln(err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(user.Password)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		jsonResponse(w, http.StatusUnauthorized, "密码错误")
		return
	}

	// Generate a JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(20 * time.Second).Unix(),
		},
	})

	// Sign the token with the secret key
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		jsonResponse(w, http.StatusInternalServerError, "")
		return
	}

	// Send the token as response
	response := TokenResponse{Token: tokenString}
	jsonResponse(w, http.StatusOK, response)
}

func (front *FrontEnd) HandleProtected(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		fmt.Println("00")
		w.WriteHeader(http.StatusMethodNotAllowed)
		jsonResponse(w, http.StatusMethodNotAllowed, "方法不是Post")
		return
	}

	// Read the token from the Authorization header
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		w.WriteHeader(http.StatusUnauthorized)
		jsonResponse(w, http.StatusUnauthorized, "Authorization字段不含token")
		return
	}
	// fmt.Println(1)
	claims := &JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			jsonResponse(w, http.StatusUnauthorized, "jwt.ErrSignatureInvalid")
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		jsonResponse(w, http.StatusBadRequest, "token过期")
		return
	}
	// fmt.Println(2)
	if !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		jsonResponse(w, http.StatusUnauthorized, "token失效")
		return
	}

	// Token is valid, continue with protected logic
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("token验证成功，登录界面"))
}

func jsonResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// 融资意向
// func (front *FrontEnd) DecryptIntensionInformation(writer http.ResponseWriter, request *http.Request) {
// 	order := make(map[string]string)
// 	Sql := sql.NewSqlCtr()
// 	slice := Sql.FinancingIntentionIndex(request)
// 	order["id"] = slice.Id
// 	order["financingId"] = slice.FinanceId
// 	order["pageid"] = slice.PageId
// 	currentPage, _ := strconv.Atoi(order["pageid"])
// 	order["searchType"] = slice.SearchType
// 	fmt.Println(order)
// 	intensions, totalcount := front.Server.SearchIntensionFromRedis(order)
// 	jsonData := front.Server.PackToIntensionJson(intensions, totalcount, currentPage)
// 	fmt.Fprint(writer, jsonData)
// }

// 回款账户
// func (front *FrontEnd) DecryptAccountInformation(writer http.ResponseWriter, request *http.Request) {
// 	order := make(map[string]string)
// 	Sql := sql.NewSqlCtr()
// 	slice := Sql.CollectionAccountIndex(request)
// 	order["id"] = slice.Id
// 	order["financeId"] = slice.FinanceId
// 	order["pageid"] = slice.PageId
// 	currentPage, _ := strconv.Atoi(order["pageid"])
// 	order["searchType"] = slice.SearchType
// 	accounts, totalcount := front.Server.SearchAccountFromRedis(order)
// 	jsonData := front.Server.PackToAccountJson(accounts, totalcount, currentPage)
// 	fmt.Fprint(writer, jsonData)
// }

// 借贷合同
// func (front *FrontEnd) DecryptFinancingContractInformation(writer http.ResponseWriter, request *http.Request) {
// 	order := make(map[string]string)
// 	Sql := sql.NewSqlCtr()
// 	slice := Sql.FinancingContractIndex(request)
// 	order["id"] = slice.Id
// 	order["pageid"] = slice.PageId
// 	currentPage, _ := strconv.Atoi(order["pageid"])
// 	order["searchType"] = slice.SearchType
// 	order["financeId"] = slice.FinanceId
// 	contracts, totalcount := front.Server.SearchFinancingContractFromRedis(order)
// 	jsonData := front.Server.PackToFinancingContractJson(contracts, totalcount, currentPage)
// 	fmt.Fprint(writer, jsonData)
// }

// 还款记录
// func (front *FrontEnd) DecryptRepaymentRecordInformation(writer http.ResponseWriter, request *http.Request) {
// 	order := make(map[string]string)
// 	Sql := sql.NewSqlCtr()
// 	slice := Sql.RepaymentRecordIndex(request)
// 	order["pageid"] = slice.PageId
// 	currentPage, _ := strconv.Atoi(order["pageid"])
// 	order["searchType"] = slice.SearchType
// 	order["financeId"] = slice.FinanceId
// 	order["id"] = slice.Id
// 	fmt.Println(order)
// 	records, totalcount := front.Server.SearchRepaymentRecordFromRedis(order)
// 	jsonData := front.Server.PackToRepaymentRecordJson(records, totalcount, currentPage)
// 	fmt.Fprint(writer, jsonData)
// }

// 接口，负责发送勾选数据至其他服务端
func handle(targetUrl string, targetJson string) SucessCode {
	targetData := []byte(targetJson)
	reader := bytes.NewReader(targetData)
	var resp *http.Response
	var data []byte
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	// 获取 request请求
	request, err := http.NewRequest("POST", targetUrl, reader)
	if err != nil {

		fmt.Println("GetHttpSkip Request Error:", err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Connection", "close")
	//启用cookie
	client.Jar, _ = cookiejar.New(nil)
	resp, err = client.Do(request)
	check(err)
	if data, err = ioutil.ReadAll(resp.Body); err == nil {
		fmt.Printf("%s\n", data)
	}
	err = resp.Body.Close()
	if err != nil {
		logrus.Errorln(err)
	}
	// 解析返回的JSON数据
	var message SucessCode
	err = json.Unmarshal(data, &message)

	check(err)
	return message
}

func (front *FrontEnd) Relay(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	relaybody, _ := ioutil.ReadAll(r.Body)

	var jsonSring string //存储获取的jsonrpc命令字符串
	for k := range r.Form {
		fmt.Printf("接收到的消息:%v\n", k)
		jsonSring = k
	}
	if jsonSring != "" {
		relaybody = []byte(jsonSring)
	}
	url := front.url.FiscoUrl
	body := string(relaybody)
	response, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(body))) //将获取的jsonrpc命令直接转发给url指定的fisco地址
	if err != nil {
		fmt.Println("err:", err)
	}
	b, _ := ioutil.ReadAll(response.Body) //等待fisco jsonrpc端口的回复

	w.Write(b) //将获取的回复回复给当前的http客户端

}

func (front *FrontEnd) ParesTXInfo(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	relaybody, _ := ioutil.ReadAll(r.Body)

	var jsonSring string //存储获取的jsonrpc命令字符串
	for k := range r.Form {
		jsonSring = k
	}
	if jsonSring != "" {
		relaybody = []byte(jsonSring)
	}

	jsonCommand := new(proxy.JsonCommand)
	json.Unmarshal(relaybody, jsonCommand)

	fmt.Printf("jsonCommand.Method:%s\n", jsonCommand.Method)

	if jsonCommand.Method != "getBlockByNumber_sp" && jsonCommand.Method != "getBlockByNumber_all" {

		client := &http.Client{}
		url := front.url.FiscoUrl
		body := string(relaybody)

		buffer := []byte(body)
		reader := bytes.NewReader(buffer)

		request, err := http.NewRequest("POST", url, reader)
		if err != nil {
			fmt.Println("GetHttpSkip Request Error:", err)
		}
		request.Header.Set("Access-Control-Allow-Origin", "*")
		request.Header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		request.Header.Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		request.Header.Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		request.Header.Set("Access-Control-Allow-Credentials", "true")

		request.Header.Set("content-type", "application/json") //返回数据格式是json

		response, err := client.Do(request)
		if err != nil {
			fmt.Println()
		}
		b, _ := ioutil.ReadAll(response.Body) //等待fisco jsonrpc端口的回复

		w.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		w.Header().Set("content-type", "application/json") //返回数据格式是json
		// fmt.Println("request header: ", request.Header)
		// fmt.Println("request: ", request)
		// fmt.Println("response: ", response)
		// fmt.Println("response body: ", response.Body)
		fmt.Println("ioutil.ReadAll(response.Body: ", string(b))
		w.Write(b) //将获取的回复回复给当前的http客户端
		return
	} else if jsonCommand.Method == "getBlockByNumber_all" {

		jsonCommand.Method = "getBlockNumber" //申请获取当前区块高度
		jsonCommand.Params = []interface{}{1}
		var relaybody_count []byte = make([]byte, 0)
		relaybody_count, _ = json.Marshal(jsonCommand)
		fmt.Println(string(relaybody_count))

		client := &http.Client{}
		url := front.url.FiscoUrl
		//body := string(relaybody)

		buffer := []byte(relaybody_count)
		reader := bytes.NewReader(buffer)

		request, err := http.NewRequest("POST", url, reader)
		if err != nil {
			fmt.Println("GetHttpSkip Request Error:", err)
		}
		request.Header.Set("Access-Control-Allow-Origin", "*")
		request.Header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		request.Header.Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		request.Header.Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		request.Header.Set("Access-Control-Allow-Credentials", "true")

		request.Header.Set("content-type", "application/json") //返回数据格式是json

		//response, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(body))) //将获取的jsonrpc命令直接转发给url指定的fisco地址
		//if err != nil {
		//	fmt.Println("err:", err)
		//}
		response, err := client.Do(request)
		if err != nil {
			fmt.Println()
		}

		b, _ := ioutil.ReadAll(response.Body) //等待fisco jsonrpc端口的回复
		blockCount := new(proxy.JsonBlockNumber)
		json.Unmarshal(b, blockCount)
		// fmt.Println("request header: ", request.Header)
		// fmt.Println("request: ", request)
		// fmt.Println("response: ", response)
		// fmt.Println("response body: ", response.Body)
		// fmt.Println("ioutil.ReadAll(response.Body: ", string(b))
		fmt.Printf("当前区块数目为: %s\n", blockCount.Result)

		countInt, err := strconv.ParseInt(strings.TrimLeft(blockCount.Result, "0x"), 16, 64)
		if err != nil {
			fmt.Printf("解析十六进制字符串失败，err:%v\n", err)
		}

		var BlockArray proxy.BlockArray //包装所有从jsonrpc处获取的区块
		BlockArray.BlockCount = countInt
		fmt.Printf("countInt = %d\n", countInt)
		var i int64
		for i = 0; i <= countInt; i++ { //申请获取从0开始的到countInt的所有区块
			blockNumber := fmt.Sprintf("0x%x", i)
			jsonCommand.Method = "getBlockByNumber"
			jsonCommand.Params = []interface{}{1, blockNumber, true}
			var block []byte = make([]byte, 0)
			block, _ = json.Marshal(jsonCommand)

			client := &http.Client{}
			url := front.url.FiscoUrl
			//body := string(relaybody)

			buffer := []byte(block)
			reader := bytes.NewReader(buffer)

			request, err := http.NewRequest("POST", url, reader)
			if err != nil {
				fmt.Println("GetHttpSkip Request Error:", err)
			}
			request.Header.Set("Access-Control-Allow-Origin", "*")
			request.Header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			request.Header.Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			request.Header.Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			request.Header.Set("Access-Control-Allow-Credentials", "true")

			request.Header.Set("content-type", "application/json") //返回数据格式是json

			//body := string(block)
			fmt.Printf("本次申请获取区块的请求:%s\n", string(block))

			response, err := client.Do(request)
			if err != nil {
				fmt.Println()
			}
			b, _ := ioutil.ReadAll(response.Body) //等待fisco jsonrpc端口的回复
			// fmt.Println("request header: ", request.Header)
			// fmt.Println("request: ", request)
			// fmt.Println("response: ", response)
			// fmt.Println("response body: ", response.Body)
			// fmt.Println("ioutil.ReadAll(response.Body: ", string(b))
			blockInfo := new(proxy.JsonBlockByNumber) //获取到对应编号的区块
			json.Unmarshal(b, blockInfo)
			BlockArray.Blocks = append(BlockArray.Blocks, blockInfo.Result) //将该区块填充
		}

		BlockArrayMarshal, err := json.Marshal(BlockArray)
		if err != nil {
			fmt.Println("BlockArray marshal is failed,err:", err)
		}

		w.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		w.Header().Set("content-type", "application/json") //返回数据格式是json

		w.Write(BlockArrayMarshal)

	} else if jsonCommand.Method == "getBlockByNumber_sp" {
		jsonCommand.Method = "getBlockByNumber"
		var relaybody_sp []byte = make([]byte, 0)
		relaybody_sp, _ = json.Marshal(jsonCommand)
		fmt.Println(string(relaybody_sp))

		client := &http.Client{}
		url := front.url.FiscoUrl
		//body := string(relaybody)

		buffer := []byte(relaybody_sp)
		reader := bytes.NewReader(buffer)

		request, err := http.NewRequest("POST", url, reader)
		if err != nil {
			fmt.Println("GetHttpSkip Request Error:", err)
		}
		request.Header.Set("Access-Control-Allow-Origin", "*")
		request.Header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		request.Header.Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		request.Header.Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		request.Header.Set("Access-Control-Allow-Credentials", "true")

		request.Header.Set("content-type", "application/json") //返回数据格式是json

		response, err := client.Do(request)
		if err != nil {
			fmt.Println()
		}
		b, _ := ioutil.ReadAll(response.Body) //等待fisco jsonrpc端口的回复
		// fmt.Println("request header: ", request.Header)
		// fmt.Println("request: ", request)
		// fmt.Println("response: ", response)
		// fmt.Println("response body: ", response.Body)
		// fmt.Println("ioutil.ReadAll(response.Body: ", string(b))
		blockInfo := new(proxy.JsonBlockByNumber)
		json.Unmarshal(b, blockInfo)

		txList := new(proxy.TransactionArray)
		txList.Transactions = make([]interface{}, 0)
		txList.TxCount = len(blockInfo.Result.Transactions)

		for index, tx := range blockInfo.Result.Transactions {
			fmt.Printf("这是第%d笔交易的输入input:%v\n", index, tx.Input)
			resStrArray, msgType := proxy.Decoder([]byte(tx.Input))
			switch msgType {
			case "issueInvoiceInformationStorage":
				invoiceInfo := proxy.DecodeInvoiceInfotoString(resStrArray["issueInvoiceInfo"], tx.Hash)
				// apInfo.Method = "Issue Invoice Information"
				txList.Transactions = append(txList.Transactions, invoiceInfo)

			case "issueHistoricalSettleInformation":
				historicalsettleInfo := proxy.DecodeHistoricalSettleInfotoString(resStrArray["issueHistoricalSettleInfo"], tx.Hash)
				// apInfo.Method = "Issue HistoricalSettle Information"
				txList.Transactions = append(txList.Transactions, historicalsettleInfo)

			case "issueHistoricalOrderInformation":
				historicalorderInfo := proxy.DecodeHistoricalOrderInfotoString(resStrArray["issueHistoricalOrderInfo"], tx.Hash)
				// bidPrice.Method = "Buyer:Purchase Request"
				txList.Transactions = append(txList.Transactions, historicalorderInfo)

			case "issueHistoricalUsedInformation":
				historicalUsedinfo := proxy.DecodeHistoricalUsedInfotoString(resStrArray["issueHistoricalUsedInfo"], tx.Hash)
				// chSwitch.Method = "Seller:Willingness to sell"
				txList.Transactions = append(txList.Transactions, historicalUsedinfo)

			case "issuePoolPlanInformation":
				poolplanInfo := proxy.DecodePoolPlanInfotoString(resStrArray["issuePoolPlanInfo"], tx.Hash)
				// chDeal.Method = "Transaction Confirm"
				txList.Transactions = append(txList.Transactions, poolplanInfo)
			case "issuePoolUsedInformation":
				poolusedInfo := proxy.DecodePoolUsedInfotoString(resStrArray["issuePoolUsedInfo"], tx.Hash)
				// chDeal.Method = "Transaction Confirm"
				txList.Transactions = append(txList.Transactions, poolusedInfo)

			// case "queryAPChannelInfo":
			// 	fmt.Printf("%s : %s\n ", msgType, resStrArray["addr"])
			// 	query := new(proxy.QueryAPChannelInfo)
			// 	query.Method = "Update database"
			// 	query.TxHash = tx.Hash
			// 	query.APaddrStr = resStrArray["addr"]

			// 	txList.Transactions = append(txList.Transactions, query)

			case "false":
				fmt.Println("未识别的消息类型.......")
			}
		}

		txListMarshal, err := json.Marshal(txList)
		if err != nil {
			fmt.Println("txList marshal is failed,err:", err)
		}

		// txListStr := string(txListMarshal)
		// txListStrRlp := strings.Replace(txListStr, "\\", "", -1)
		// txListStrRlp += "\n"
		// fmt.Println(txListStrRlp)

		w.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		w.Header().Set("content-type", "application/json") //返回数据格式是json

		w.Write(txListMarshal)
	}
	fmt.Println(w.Header())
}
