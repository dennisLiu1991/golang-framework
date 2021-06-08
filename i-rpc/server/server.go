package server

import (
	"fmt"
	"i-rpc/codec"
	"i-rpc/model"
	"net"
)

// Server is rpc server
type Server struct {
	Addr string
}

var defaultServer = &Server{}

// GetDefault returns the default server
func GetDefault() *Server {
	return defaultServer
}

// Start the server
func (s *Server) Start(c chan<- string) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		fmt.Printf("【server】start server fail! %v", err)
		return
	}
	s.Addr = l.Addr().String()
	s.serve(l)

	fmt.Printf("【server】start server on %s\n", s.Addr)
	c <- s.Addr
}

func (s *Server) serve(l net.Listener) {
	go func() {
		for {
			fmt.Printf("【server】server waite for connextion ...\n")
			conn, err := l.Accept()
			if err != nil {
				fmt.Printf("【server】shutdown server ! serve err %v\n", err)
				continue
			}
			if conn == nil {
				continue
			}
			fmt.Println("【server】 recevie connection!")

			newCodec := codec.GetCodec(codec.Gob)
			go s.serveForConn(newCodec(conn))
		}
	}()
}

func (s *Server) serveForConn(cc codec.CodeC) {
	for {
		req, err := readRequest(cc)
		if err != nil {
			continue
		}
		go s.handleRequest(req, cc)
	}
}

func readRequest(cc codec.CodeC) (model.Request, error) {
	var req model.Request
	if err := cc.Read(&req); err != nil {
		fmt.Printf("【server】read fail %v\n", err)
		return req, err
	}
	fmt.Printf("【server】got request : %s\n", req)
	return req, nil
}

func (s *Server) handleRequest(req model.Request, cc codec.CodeC) {
	s.sendResponse(cc)
}

func (s *Server) sendResponse(cc codec.CodeC) {
	response := model.Response{
		ID: "1",
	}
	if err := cc.Write(&response); err != nil {
		fmt.Printf("【server】write response  %v", err)
		return
	}
}
