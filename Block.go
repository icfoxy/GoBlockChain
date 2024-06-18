package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Block struct {
	BlockNum     int           `json:"blockNum"`
	Nonce        int           `json:"nonce"`
	PreHash      [32]byte      `json:"preHash"`
	Transactions []Transaction `json:"transaction"`
	Timestamp    int64         `json:"timestamp"`
}

func NewBlock(blockNum, nonce int, preHash [32]byte, transactions []Transaction) *Block {
	return &Block{
		BlockNum:     blockNum,
		Nonce:        nonce,
		PreHash:      preHash,
		Transactions: transactions,
		Timestamp:    time.Now().UnixNano(),
	}
}

func (b *Block) Hash() [32]byte {
	data, _ := json.Marshal(b)
	hash := sha256.Sum256(data)
	return hash
}

func (b *Block) Print() {
	fmt.Printf("%sNum:%v%s\n", strings.Repeat("-", 25), b.BlockNum, strings.Repeat("-", 25))
	fmt.Printf("Nonce---------%v\n", b.Nonce)
	fmt.Printf("PreHash-------%x\n", b.PreHash)
	fmt.Printf("TransActions:\n")
	for i, v := range b.Transactions {
		fmt.Printf("    Num:%d\n", i)
		v.Print()
	}
	fmt.Printf("Timstamp------%v\n", b.Timestamp)
}
