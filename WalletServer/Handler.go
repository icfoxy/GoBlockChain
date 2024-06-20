package main

import (
	"net/http"
	"text/template"

	blockchain "github.com/icfoxy/GoBlockChain/BlockChain"
	"github.com/icfoxy/GoTools"
)

// const templatesdir = "WalletServer/templates"

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
