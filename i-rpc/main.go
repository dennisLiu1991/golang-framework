package main

import (
	_ "i-rpc/codec"
	"i-rpc/model"
	"log"

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

	req := model.Request{
		Name: "abc",
	}
	var rsp model.Response
	err := cli.SendRequest(req, &rsp)
	log.Printf("client rsp = %v,err = %v", rsp, err)
}
