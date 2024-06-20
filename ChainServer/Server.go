package main

import (
	"fmt"
	"log"
	"net/http"

	blockchain "github.com/icfoxy/GoBlockChain/BlockChain"
)

var MainChain *blockchain.BlockChain

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
	MainChain = initBlockChain()
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

func initBlockChain() *blockchain.BlockChain {
	//todo: 从主网络拉取MainChain
	return blockchain.NewBlockChain()
}
