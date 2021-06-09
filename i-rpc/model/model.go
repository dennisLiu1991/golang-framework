package model

type Request struct {
	Name string
}

type Response struct {
	ID string
}

type Header struct {
	Service string
	Method  string
	CallSeq int
}
