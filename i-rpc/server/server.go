package server

import (
	"log"
	"net"

	"i-rpc/codec"
	"i-rpc/model"
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
		log.Printf("【server】start server fail! %v", err)
		return
	}
	s.Addr = l.Addr().String()
	s.serve(l)

	log.Printf("【server】start server on %s", s.Addr)
	c <- s.Addr
}

func (s *Server) serve(l net.Listener) {
	go func() {
		for {
			log.Printf("【server】server waite for connextion ...")
			conn, err := l.Accept()
			if err != nil {
				log.Printf("【server】shutdown server ! serve err %v", err)
				continue
			}
			if conn == nil {
				continue
			}
			log.Println("【server】 recevie connection!")

			newCodec := codec.GetCodec(codec.Gob)
			go s.serveForConn(newCodec(conn))
		}
	}()
}

func (s *Server) serveForConn(cc codec.CodeC) {
	for {
		req, h, err := readRequest(cc)
		if err != nil {
			continue
		}
		go s.handleRequest(req, h, cc)
	}
}

func readRequest(cc codec.CodeC) (model.Request, model.Header, error) {
	// 先读header
	var h model.Header
	if err := cc.Read(&h); err != nil {
		log.Printf("【server】read header %v", err)
		return model.Request{}, model.Header{}, err
	}
	log.Printf("【server】got header : %v", h)

	var req model.Request
	if err := cc.Read(&req); err != nil {
		log.Printf("【server】read fail %v", err)
		return req, h, err
	}
	log.Printf("【server】got request : %s", req)
	return req, h, nil
}

func (s *Server) handleRequest(req model.Request, h model.Header, cc codec.CodeC) {
	s.sendResponse(cc, h)
}

func (s *Server) sendResponse(cc codec.CodeC, h model.Header) {
	response := model.Response{
		ID: "123456789",
	}
	if err := cc.Write(&h, &response); err != nil {
		log.Printf("【server】write response  %v", err)
		return
	}
}
