package main

import (
	"fmt"
	"log"
	"net/http"
)

type WalletServer struct {
	Addr    string
	port    uint16
	gateway string
	Mux     *http.ServeMux
}

func NewWalletServer(addr string, port uint16, gateway string) *WalletServer {
	ws := new(WalletServer)
	ws.Addr = addr
	ws.port = port
	ws.gateway = gateway
	ws.Mux = http.NewServeMux()
	return ws
}

func (ws *WalletServer) GetPort() uint16 {
	return ws.port
}

func (ws *WalletServer) GetGateWay() string {
	return ws.gateway
}

func (ws *WalletServer) AddHandlerFunc(url string, f func(http.ResponseWriter, *http.Request)) {
	ws.Mux.HandleFunc(url, f)
}

func (ws *WalletServer) Run() {
	httpServer := &http.Server{
		Addr:    ws.Addr + fmt.Sprint(ws.port),
		Handler: ws.Mux,
	}
	log.Println("server starts at:", ws.port)
	err := httpServer.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
