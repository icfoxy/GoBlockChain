package main

import "fmt"

func main() {
	bc := NewBlockChain()
	for i := 0; i < 10000; i++ {
		bc.PushTransaction("000aa", "000bb", 2.0, "nothing")
		bc.PushTransaction("000bb", "007aa", 2.02, "nothing either")
		bc.PushTransaction("000cc", "000bb", 2.03, "nothing")
	}
	// bc.Print()
	fmt.Println("mining...")
	nonce := bc.ProofOfWork()
	fmt.Println("nonce:", nonce)
}
