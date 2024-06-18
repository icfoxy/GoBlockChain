package main

func main() {
	bc := NewBlockChain()
	bc.PushTransaction("000aa", "000bb", 2.0, "nothing")
	bc.PushTransaction("000bb", "000aa", 2.02, "nothing either")
	bc.PushTransaction("000cc", "000bb", 2.03, "nothing")
	bc.Print()
}
