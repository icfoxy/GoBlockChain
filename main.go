package main

import (
	chainserver "github.com/icfoxy/GoBlockChain/ChainServer"
)

func main() {
	server := chainserver.NewServer("localhost:", 8088)
	server.AddHandlerFunc("/TestAlive", chainserver.TestAliveHandler)
	server.AddHandlerFunc("/GetChain", chainserver.GetChainHandler)
	chainserver.MainChain.PushTransaction("a", "b", 2, "")
	chainserver.MainChain.Mine("mxh")
	chainserver.MainChain.PushTransaction("a", "c", 2, "")
	chainserver.MainChain.Mine("mxh")
	server.Run()
	// wA := blockchain.NewWallet()
	// wB := blockchain.NewWallet()
	// wM := blockchain.NewWallet()
	// transaction, sign := wA.NewSignedTransaction(wA.GetAddr(), wB.GetAddr(), 1, "test")
	// bc := blockchain.NewBlockChain()
	// bc.PushTransaction(blockchain.MainNetAddr, wA.GetAddr(), 2, "")
	// bc.Mine(wA.GetAddr())
	// bc.ValidAndPushTransaction(wA.GetPublicKey(), *sign, *transaction)
	// fmt.Println("mining...")
	// bc.Mine(wM.GetAddr())
	// bc.Print()
	// fmt.Println("A:", bc.GetTotalValue(wA.GetAddr()))
	// fmt.Println("B:", bc.GetTotalValue(wB.GetAddr()))
	// fmt.Println("M:", bc.GetTotalValue(wM.GetAddr()))
}
