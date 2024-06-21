package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var stop = make(chan bool, 1)
var start = make(chan bool, 1)

func init() {
	transport := &http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			localAddr, err := net.ResolveTCPAddr(network, "localhost:"+os.Getenv("SendPort"))
			if err != nil {
				return nil, err
			}
			remoteAddr, err := net.ResolveTCPAddr(network, addr)
			if err != nil {
				return nil, err
			}
			return net.DialTCP(network, localAddr, remoteAddr)
		},
	}
	client = &http.Client{
		Transport: transport,
	}
}

func main() {
	nodeName := flag.String("nodeName", "node_0", "name of node")
	flag.Parse()

	// 读取.env文件
	envVars, err := godotenv.Read()
	if err != nil {
		log.Println(".env load wrong")
		return
	}
	// 获取nodeName对应的值
	UserAddr = envVars[*nodeName]
	// 创建一个空的字符串切片来存储其他节点的地址
	var otherNodes []string
	// 遍历.env文件中的变量
	for key, value := range envVars {
		// 如果这个变量不是nodeName对应的值，就将它添加到切片中
		if key != *nodeName {
			otherNodes = append(otherNodes, value)
		}
	}
	OtherAddr = otherNodes

	server := NewServer(UserAddr[:10], GetPortFromAddr(UserAddr))
	server.AddHandlerFunc("/GetChain", GetChainHandler)
	server.AddHandlerFunc("/PushTransaction", PushTransactionHandler)
	server.AddHandlerFunc("/UpdatePool", UpdatePoolHandler)
	server.AddHandlerFunc("/DelayMining", DelayMiningHandler)
	server.AddHandlerFunc("/AskForChainUpdate", AskForChainUpdate)
	server.AddHandlerFunc("/ReceiveChain", ReceiveChain)
	go func() {
		for {
			result := MainChain.Mine("test", stop, start)
			if result == nil {
				time.Sleep(2 * time.Second)
			} else {
				AskForChainUpdate(*new(http.ResponseWriter), new(http.Request))
			}
		}
	}()
	server.Run()
}
