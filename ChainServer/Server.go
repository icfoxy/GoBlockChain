package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	blockchain "github.com/icfoxy/GoBlockChain/BlockChain"
	"github.com/icfoxy/GoTools"
)

var UserAddr = "User"
var MainChain *blockchain.BlockChain
var OtherAddr []string
var client *http.Client

type Server struct {
	Addr string
	Port uint16
	Mux  *http.ServeMux
}

func NewServer(addr string, port uint16) *Server {
	server := new(Server)
	server.Addr = addr
	server.Port = port
	server.Mux = http.NewServeMux()
	MainChain = blockchain.NewBlockChain()
	UpdateChain()
	return server
}

func (s *Server) AddHandlerFunc(url string, f func(http.ResponseWriter, *http.Request)) {
	s.Mux.HandleFunc(url, f)
}

func (s *Server) Run() {
	httpServer := &http.Server{
		Addr:    s.Addr + fmt.Sprint(s.Port),
		Handler: s.Mux,
	}
	log.Println("server starts at:", s.Port)
	err := httpServer.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

// todo
func ValidChain(chain *blockchain.BlockChain) bool {
	return true
}

func UpdateChain() {
	for _, v := range OtherAddr {
		resp, err := client.Get("http://" + v + "/GetChain")
		if err != nil {
			log.Println("node:", v, "not found")
			continue
		}
		defer func() {
			io.ReadAll(resp.Body)
			resp.Body.Close()
		}()
		newChain := blockchain.NewBlockChain()
		err = GoTools.GetAnyFromBody(resp.Body, newChain)
		if err != nil {
			panic(err)
		}
		newChain.Print()
		if len(newChain.Chain) > len(MainChain.Chain) && ValidChain(newChain) {
			MainChain = newChain
		}
	}
}

// 定义一个函数，它接受一个字符串参数，并返回一个uint16和一个error
func GetPortFromAddr(addr string) uint16 {
	// 分割字符串，获取最后四个数字
	parts := strings.Split(addr, ":")
	lastFourDigits := parts[len(parts)-1]

	// 将字符串转换为uint64
	u64, err := strconv.ParseUint(lastFourDigits, 10, 16)
	if err != nil {
		panic(err)
	}

	// 将uint64转换为uint16
	u16 := uint16(u64)

	return u16
}
