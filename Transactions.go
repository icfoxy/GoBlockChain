package main

import "fmt"

type Transaction struct {
	SenderAddr   string  `json:"sender_addr"`
	ReceiverAddr string  `json:"receiver_addr"`
	Value        float32 `json:"value"`
	Info         string  `json:"info"`
}

func (t *Transaction) Print() {
	fmt.Printf("      SenderAddr:%v\n", t.SenderAddr)
	fmt.Printf("      ReceiverAddr:%v\n", t.ReceiverAddr)
	fmt.Printf("      Value:%v\n", t.Value)
	fmt.Printf("      Info:%v\n", t.Info)
}
