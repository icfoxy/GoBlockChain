package main

import "fmt"

func main() {
	ws := NewWalletServer("localhost:", 8089, "")
	ws.AddHandlerFunc("/TestAlive", TestAliveHandler)
	ws.AddHandlerFunc("/index", OpenIndexHandler)
	ws.AddHandlerFunc("/CreateNewWallet", CreateWalletHandler)
	ws.AddHandlerFunc("/Sign", SignHandler)
	fmt.Println("server starts at:8089")
	ws.Run()
}
