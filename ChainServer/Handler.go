package main

import (
	"bytes"
	"crypto/ecdsa"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	blockchain "github.com/icfoxy/GoBlockChain/BlockChain"
	"github.com/icfoxy/GoTools"
)

type TransactionTansfer struct {
	SenderAddr    string `json:"sender_addr"`
	ReceriverAddr string `json:"receiver_addr"`
	Value         int    `json:"value"`
	Info          string `json:"info"`
	PublicKey     string `json:"public_key"`
	Signature     string `json:"signature"`
}

func GetChainHandler(w http.ResponseWriter, r *http.Request) {
	GoTools.RespondByJSON(w, 200, *MainChain)
}

func TestAliveHandler(w http.ResponseWriter, r *http.Request) {
	GoTools.RespondByJSON(w, 200, "I AM ALIVE")
}

// TODO:可能有bug，要优化
func PushTransactionHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("receive push request")
	// 读取请求体并存储以便后续使用
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request body:", err)
		return
	}
	// 使用存储的请求体获取交易数据
	transaction, publicKey, signature := GetTansactionDataFromBody(bytes.NewReader(bodyBytes))
	result := MainChain.ValidAndPushTransaction(publicKey, *signature, transaction)
	MainChain.Print()
	if result {
		GoTools.RespondByJSON(w, 200, "TransAction Pushed")
		// 使用并发发送请求
		var wg sync.WaitGroup
		for _, v := range OtherAddr {
			wg.Add(1)
			go func(addr string) {
				defer wg.Done()
				resp, err := client.Post("http://"+addr+"/UpdatePool", "json", bytes.NewReader(bodyBytes))
				if err != nil {
					log.Println(addr, "not reachable")
				}
				log.Println(resp)
			}(v)
		}
		wg.Wait() // 等待所有请求完成
	} else {
		GoTools.RespondByErr(w, 801, "TransAction not added", "high")
	}
}

func UpdatePoolHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("receive push request")
	// 读取请求体并存储以便后续使用
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request body:", err)
		return
	}
	// 使用存储的请求体获取交易数据
	transaction, publicKey, signature := GetTansactionDataFromBody(bytes.NewReader(bodyBytes))
	result := MainChain.ValidAndPushTransaction(publicKey, *signature, transaction)
	MainChain.Print()
	if result {
		GoTools.RespondByJSON(w, 200, "TransAction Pushed")
	} else {
		GoTools.RespondByErr(w, 801, "TransAction not added", "high")
	}
}

func GetTansactionDataFromBody(body io.Reader) (blockchain.Transaction, *ecdsa.PublicKey, *blockchain.Signature) {
	tf := new(TransactionTansfer)
	err := GoTools.GetAnyFromBody(body, tf)
	if err != nil {
		panic(err)
	}
	transaction := blockchain.Transaction{
		SenderAddr:   tf.SenderAddr,
		ReceiverAddr: tf.ReceriverAddr,
		Value:        tf.Value,
		Info:         tf.Info,
	}
	publicKey, err := blockchain.StringToPublicKey(tf.PublicKey)
	if err != nil {
		panic(err)
	}
	signature, err := blockchain.StringToSignature(tf.Signature)
	if err != nil {
		panic(err)
	}
	return transaction, publicKey, signature
}
