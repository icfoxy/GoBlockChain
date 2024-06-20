package main

import (
	"fmt"

	blockchain "github.com/icfoxy/GoBlockChain/BlockChain"
)

func main() {
	// server := NewServer("localhost:", 8088)
	// server.AddHandlerFunc("/TestAlive", TestAliveHandler)
	// server.AddHandlerFunc("/GetChain", GetChainHandler)
	// MainChain.PushTransaction("a", "b", 2, "")
	// MainChain.Mine("mxh")
	// MainChain.PushTransaction("a", "c", 2, "")
	// MainChain.Mine("mxh")
	// MainChain.Print()
	// server.Run()
	wA := blockchain.NewWallet()
	wB := blockchain.NewWallet()
	wM := blockchain.NewWallet()
	transaction, sign := wA.NewSignedTransaction(wA.GetAddr(), wB.GetAddr(), 1, "test")
	bc := blockchain.NewBlockChain()
	bc.PushTransaction(blockchain.MainNetAddr, wA.GetAddr(), 2, "")
	bc.Mine(wA.GetAddr())
	bc.ValidAndPushTransaction(wA.GetPublicKey(), *sign, *transaction)
	fmt.Println("mining...")
	bc.Mine(wM.GetAddr())
	bc.Print()
	fmt.Println("A:", bc.GetTotalValue(wA.GetAddr()))
	fmt.Println("B:", bc.GetTotalValue(wB.GetAddr()))
	fmt.Println("M:", bc.GetTotalValue(wM.GetAddr()))
}
