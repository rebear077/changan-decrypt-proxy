package main

import (
	"log"
	"net/http"
	"time"

	"github.com/FISCO-BCOS/go-sdk/querytable"
)

func main() {
	front := querytable.NewFrontEnd()
	front.Server.ForceSynchronous()
	time.Sleep(1 * time.Second)
	go front.Server.DumpFromCanal()

	http.HandleFunc("/asl/universal/decryptInvoiceInformation", front.DecryptInvoiceInformation)
	http.HandleFunc("/asl/universal/decryptHistoricaltransaction", front.DecryptHistoricaltransaction)
	http.HandleFunc("/asl/universal/decryptEnterpoolDataInfos", front.DecryptEnterPoolData)
	// http.HandleFunc("/asl/universal/decryptFinancingIntention", front.DecryptIntensionInformation)
	// http.HandleFunc("/asl/universal/decryptCollectionAccount", front.DecryptAccountInformation)
	http.HandleFunc("/asl/universal/selectedToApplication", front.DecryptSelectToApplicationInformation)
	// http.HandleFunc("/asl/universal/decryptFinancingContract", front.DecryptFinancingContractInformation)
	// http.HandleFunc("/asl/universal/decryptRepaymentRecord", front.DecryptRepaymentRecordInformation)
	http.HandleFunc("/asl/universal/handle/", front.ParesTXInfo)

	http.HandleFunc("/login", front.HandleLogin)         //登录
	http.HandleFunc("/protected", front.HandleProtected) //token登录
	// http.HandleFunc("/register", handleRegister)

	log.Fatal(http.ListenAndServe(":8080", nil))

	// err := http.ListenAndServeTLS(":8440", "connApi/confs/server.pem", "connApi/confs/server.key", nil)
	err := http.ListenAndServe(":8440", nil)
	if err != nil {
		log.Fatalf("启动 HTTPS 服务器失败: %v", err)
	}
}
