package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/FISCO-BCOS/go-sdk/querytable"
)

func main() {
	front := querytable.NewFrontEnd()
	front.Server.ForceSynchronous()
	fmt.Println("xxx")
	time.Sleep(1 * time.Second)
	go front.Server.DumpFromCanal()

	http.HandleFunc("/asl/universal/decryptInvoiceInformation", front.DecryptInvoiceInformation)
	http.HandleFunc("/asl/universal/decryptHistoricaltransaction", front.DecryptHistoricaltransaction)
	http.HandleFunc("/asl/universal/decryptEnterpoolDataInfos", front.DecryptEnterPoolData)
	http.HandleFunc("/asl/universal/decryptFinancingIntention", front.DecryptIntensionInformation)
	http.HandleFunc("/asl/universal/decryptCollectionAccount", front.DecryptAccountInformation)
	http.HandleFunc("/asl/universal/handle/", front.ParesTXInfo)

	// err := http.ListenAndServeTLS(":8440", "connApi/confs/server.pem", "connApi/confs/server.key", nil)
	err := http.ListenAndServe(":8440", nil)
	if err != nil {
		log.Fatalf("启动 HTTPS 服务器失败: %v", err)
	}
	// }
}
