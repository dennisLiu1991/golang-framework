package client

import (
	"fmt"
	"net"

	"i-rpc/codec"
	"i-rpc/model"
)

type Client struct {
	ServerAddr string
}

func NewClient(addr string) *Client {
	return &Client{
		ServerAddr: addr,
	}
}

func (c *Client) Connecte() {
	conn, err := net.Dial("tcp", c.ServerAddr)
	if err != nil {
		fmt.Printf("【client】connect to server fail! %v \n", err)
		return
	}
	f := codec.GetCodec("gob")
	cc := f(conn)

	fmt.Printf("【client】connect to server!\n")
	req := model.Request{
		Name: "abc123123123",
	}
	if err := cc.Write(&req); err != nil {
		fmt.Printf("【client】client write fail %v\n", err)
		return
	}

	fmt.Println("【client】send request done")
	c.waitResponse(cc)
}

func (c *Client) waitResponse(cc codec.CodeC) {
	var rsp model.Response
	if err := cc.Read(&rsp); err != nil {
		fmt.Printf("【client】read response fail! %v", err)
		return
	}
	fmt.Printf("【client】read response : %v\n", rsp)
}
