package codec

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"

	"i-rpc/model"
)

type CodeC interface {
	Write(h *model.Header, v interface{}) error
	ReadHeader(h *model.Header) error
	Read(body interface{}) error
}

type NewCodecFunc func(io.ReadWriteCloser) CodeC

var codecMap map[string]NewCodecFunc

const Gob string = "gob"

func init() {
	codecMap = make(map[string]NewCodecFunc)
	codecMap[Gob] = NewCodec
	fmt.Println("init codec success")
}

func GetCodec(name string) NewCodecFunc {
	return codecMap[name]
}

type defaultCodeC struct {
	conn io.ReadWriteCloser
	buf  *bufio.Writer
	ec   *gob.Encoder
	dc   *gob.Decoder
}

func NewCodec(conn io.ReadWriteCloser) CodeC {
	buf := bufio.NewWriter(conn)
	return &defaultCodeC{
		conn: conn,
		buf:  buf,
		dc:   gob.NewDecoder(conn),
		ec:   gob.NewEncoder(buf),
	}
}

func (d *defaultCodeC) ReadHeader(h *model.Header) error {
	return d.dc.Decode(h)
}

func (d *defaultCodeC) Write(h *model.Header, v interface{}) (err error) {
	defer func() {
		// flush之后才会将buf里的数据写入conn
		_ = d.buf.Flush()
		if err != nil {
			_ = d.Close()
		}
	}()
	if err = d.ec.Encode(h); err != nil {
		return
	}
	if err = d.ec.Encode(v); err != nil {
		return
	}

	return
}

func (d *defaultCodeC) Read(v interface{}) error {
	return d.dc.Decode(v)
}

func (d *defaultCodeC) Close() error {
	return d.conn.Close()
}
