package main

import (
	_ "i-rpc/codec"

	"i-rpc/client"
	"i-rpc/server"
)

func main() {
	severAddr := startServer()
	startClient(severAddr)
	for {
		select {}
	}
}

func startServer() string {
	s := server.GetDefault()
	c := make(chan string)

	go s.Start(c)

	addr := <-c
	return addr
}

func startClient(addr string) {
	cli := client.NewClient(addr)
	cli.Connecte()
}
