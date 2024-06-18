package main

import (
	"fmt"
	"strings"
)

type BlockChain struct {
	TransactionPool []Transaction
	Chain           []*Block
}

func NewBlockChain() *BlockChain {
	bc := new(BlockChain)
	bc.Chain = append(bc.Chain, NewBlock(0, 0, [32]byte{}, nil))
	return bc
}

func (bc *BlockChain) AddBlock(nonce int, transactions []Transaction) *Block {
	b := NewBlock(len(bc.Chain), nonce, bc.getLastBlock().Hash(), transactions)
	bc.Chain = append(bc.Chain, b)
	return b
}

func (bc *BlockChain) getLastBlock() *Block {
	return bc.Chain[len(bc.Chain)-1]
}

func (bc *BlockChain) Print() {
	fmt.Printf("%sPool:%s\n", strings.Repeat("=", 25), strings.Repeat("=", 25))
	for i, v := range bc.TransactionPool {
		fmt.Printf("    Num:%d\n", i)
		v.Print()
		fmt.Printf("    %s\n", strings.Repeat("-", 35))
	}
	fmt.Printf("%sChain:%s\n", strings.Repeat("=", 25), strings.Repeat("=", 25))
	for _, v := range bc.Chain {
		v.Print()
	}
	fmt.Println()
}

func (bc *BlockChain) PushTransaction(senderAddr, receiverAddr string, value float32, info string) {
	t := Transaction{
		SenderAddr:   senderAddr,
		ReceiverAddr: receiverAddr,
		Value:        value,
		Info:         info,
	}
	bc.TransactionPool = append(bc.TransactionPool, t)
}

func (bc *BlockChain) CopyTransactionsFromPool() []Transaction {
	t := make([]Transaction, 0, len(bc.TransactionPool))
	for _, v := range bc.TransactionPool {
		t = append(t, Transaction{
			SenderAddr:   v.SenderAddr,
			ReceiverAddr: v.ReceiverAddr,
			Value:        v.Value,
			Info:         v.Info,
		})
	}
	return t
}
