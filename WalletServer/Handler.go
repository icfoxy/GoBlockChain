package main

import (
	"net/http"
	"text/template"

	blockchain "github.com/icfoxy/GoBlockChain/BlockChain"
	"github.com/icfoxy/GoTools"
)

// const templatesdir = "WalletServer/templates"

type SignTransfer struct {
	SenderAddr    string `json:"sender_addr"`
	ReceriverAddr string `json:"receiver_addr"`
	Value         int    `json:"value"`
	Info          string `json:"info"`
	PrivateKey    string `json:"private_key"`
}

func OpenIndexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./templates/index.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, "")
}

func TestAliveHandler(w http.ResponseWriter, r *http.Request) {
	GoTools.RespondByJSON(w, 200, "I AM ALIVE")
}

func CreateWalletHandler(w http.ResponseWriter, r *http.Request) {
	err := GoTools.RespondByJSON(w, 200, blockchain.NewWallet().ToTransfer())
	if err != nil {
		panic(err)
	}
}

func SignHandler(w http.ResponseWriter, r *http.Request) {
	sf := new(SignTransfer)
	err := GoTools.GetAnyFromBody(r.Body, sf)
	if err != nil {
		panic(err)
	}
	transaction := blockchain.Transaction{
		SenderAddr:   sf.SenderAddr,
		ReceiverAddr: sf.ReceriverAddr,
		Value:        sf.Value,
		Info:         sf.Info,
	}
	privatekey, err := blockchain.StringToPrivateKey(sf.PrivateKey)
	if err != nil {
		panic(err)
	}
	signature := blockchain.SignTransaction(&transaction, privatekey).ToString()
	GoTools.RespondByJSON(w, 200, signature)
}
