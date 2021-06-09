package client

import (
	"log"
	"net"

	"i-rpc/codec"
	"i-rpc/model"
)

// Client is rpc client
type Client struct {
	ServerAddr string
	cc         codec.CodeC
	calls      map[int]call
	seq        int
}

type call struct {
	reply interface{}
	req   interface{}
	seq   int
	done  chan bool
}

// NewClient returns a rpc client
func NewClient(addr string) *Client {
	return &Client{
		ServerAddr: addr,
		calls:      make(map[int]call),
	}
}

func (c *Client) initClient() {
	conn, err := net.Dial("tcp", c.ServerAddr)
	if err != nil {
		log.Printf("【client】connect to server fail! %v ", err)
		return
	}
	log.Printf("【client】connect to server!")

	f := codec.GetCodec(codec.Gob)
	c.cc = f(conn)

	// 监听响应
	go c.receive()
}

func (c *Client) receive() {
	for {
		var h model.Header
		if err := c.cc.Read(&h); err != nil {
			log.Printf("【client】read header fail! %v", err)
			continue
		}
		call := c.calls[h.CallSeq]
		delete(c.calls, h.CallSeq)

		if err := c.cc.Read(call.reply); err != nil {
			log.Printf("【client】read rsp fail! %v", err)
			continue
		}
		call.done <- true
	}
}

// SendRequest send a rpc request
func (c *Client) SendRequest(req model.Request, rsp *model.Response) error {
	// 如果是第一次发起请求需要初始化
	if c.cc == nil {
		c.initClient()
	}

	seq := c.seq
	h := model.Header{
		Service: "srv1",
		Method:  "mothod1",
		CallSeq: seq,
	}

	call := call{
		req:   req,
		reply: rsp,
		seq:   seq,
		done:  make(chan bool, 1),
	}
	c.calls[seq] = call

	c.seq++

	if err := c.cc.Write(&h, &req); err != nil {
		return err
	}
	log.Println("【client】send request done")
	<-call.done

	rsp = call.reply.(*model.Response)
	log.Printf("【client】got rsp :%v", rsp)
	return nil
}
