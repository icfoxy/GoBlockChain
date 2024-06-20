package blockchain

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"
)

const difficulty int = 4
const Reward int = 2
const MainNetAddr = "BlockChainMainNet"

type BlockChain struct {
	TransactionPool []Transaction `json:"transaction_pool"`
	Chain           []*Block      `json:"chain"`
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
	fmt.Printf("\n%sPool:%s\n", strings.Repeat("=", 25), strings.Repeat("=", 25))
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

func (bc *BlockChain) PushTransaction(senderAddr, receiverAddr string, value int, info string) {
	t := Transaction{
		SenderAddr:   senderAddr,
		ReceiverAddr: receiverAddr,
		Value:        value,
		Info:         info,
	}
	bc.TransactionPool = append(bc.TransactionPool, t)
}

func (bc *BlockChain) ValidTransaction(
	senderPublicKey *ecdsa.PublicKey, signature Signature, transaction Transaction) bool {
	data, _ := json.Marshal(transaction)
	hash := sha256.Sum256([]byte(data))
	return ecdsa.Verify(senderPublicKey, hash[:], signature.R, signature.S)
}

func (bc *BlockChain) ValidAndPushTransaction(
	senderPublicKey *ecdsa.PublicKey, signature Signature, transaction Transaction) bool {
	if !bc.ValidTransaction(senderPublicKey, signature, transaction) {
		return false
	}
	// 测试时不开
	// total := bc.GetTotalValue(transaction.SenderAddr)
	// if total-transaction.Value < 0 {
	// 	return false
	// }
	bc.PushTransaction(transaction.SenderAddr, transaction.ReceiverAddr, transaction.Value, transaction.Info)
	return true
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

func (bc *BlockChain) ValidProof(guessBlock *Block) bool {
	zeros := strings.Repeat("0", difficulty)
	guessHashstr := fmt.Sprintf("%x", guessBlock.Hash())
	return guessHashstr[:difficulty] == zeros
}

func (bc *BlockChain) ProofOfWork() *Block {
	transactions := bc.CopyTransactionsFromPool()
	guessBlock := NewBlock(len(bc.Chain), 0, bc.Chain[len(bc.Chain)-1].Hash(), transactions)
	for !bc.ValidProof(guessBlock) {
		guessBlock.Nonce++
	}
	return guessBlock
}

func (bc *BlockChain) Mine(minerAddr string) *Block {
	if bc.TransactionPool == nil {
		return nil
	}
	bc.PushTransaction(MainNetAddr, minerAddr, Reward, "Get Reward")
	block := bc.ProofOfWork()
	bc.Chain = append(bc.Chain, block)
	bc.TransactionPool = nil
	return block
}

func (bc *BlockChain) GetTotalValue(addr string) int {
	var total int = 0.0
	for _, b := range bc.Chain {
		for _, t := range b.Transactions {
			if t.ReceiverAddr == addr {
				total += t.Value
			} else if t.SenderAddr == addr {
				total -= t.Value
			}
		}
	}
	return total
}
