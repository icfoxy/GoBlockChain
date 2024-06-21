package main

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

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

func PushTransactionHandler(w http.ResponseWriter, r *http.Request) {
	stop <- true
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
	start <- true
	if result {
		GoTools.RespondByJSON(w, 200, "TransAction Pushed")
		// 使用并发发送请求
		var wg sync.WaitGroup
		for _, v := range OtherAddr {
			wg.Add(1)
			go func(addr string) {
				defer wg.Done()
				_, err := client.Post("http://"+addr+"/UpdatePool", "json", bytes.NewReader(bodyBytes))
				if err != nil {
					log.Println(addr, "not reachable")
				}
			}(v)
		}
		wg.Wait() // 等待所有请求完成
	} else {
		GoTools.RespondByErr(w, 801, "TransAction not added", "high")
	}
}

func UpdatePoolHandler(w http.ResponseWriter, r *http.Request) {
	stop <- true
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
	// 设置随机数种子
	rand.Seed(time.Now().UnixNano())

	// 生成一个0-2秒的随机暂停时间
	pause := time.Duration(rand.Intn(3)) * time.Second

	// 暂停指定的时间
	time.Sleep(pause)
	start <- true
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

func DelayMiningHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("stopping mining...")
	stop <- true
	time.Sleep(3 * time.Second)
	start <- true
	GoTools.RespondByJSON(w, 200, "mine restarted")
}

// TODO:验证可行性
func AskForChainUpdate(w http.ResponseWriter, r *http.Request) {
	jsonData, err := json.Marshal(MainChain)
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	for _, v := range OtherAddr {
		wg.Add(1)
		go func(addr string) {
			defer wg.Done()
			req, err := http.NewRequest("POST", "http://"+addr+"/ReceiveChain", bytes.NewBuffer(jsonData))
			if err != nil {
				panic(err)
			}
			req.Header.Set("Content-Type", "application/json")
			_, err = client.Do(req)
			if err != nil {
				log.Println(addr, "not reachable")
			}
		}(v)
	}
	wg.Wait()
}

func ReceiveChain(w http.ResponseWriter, r *http.Request) {
	stop <- true
	newChain := new(blockchain.BlockChain)
	err := GoTools.GetAnyFromBody(r.Body, newChain)
	if err != nil {
		panic(err)
	}
	fmt.Println("receive chain:")
	newChain.Print()
	if ValidChain(newChain) && len(newChain.Chain) >= len(MainChain.Chain) {
		fmt.Println("MainChain updated")
		MainChain = newChain
	}
	start <- true
}
