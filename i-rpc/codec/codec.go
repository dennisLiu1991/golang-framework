package codec

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
)

type CodeC interface {
	Write(v interface{}) error
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

func (d *defaultCodeC) Write(v interface{}) (err error) {
	defer func() {
		// flush之后才会将buf里的数据写入conn
		_ = d.buf.Flush()
		if err != nil {
			_ = d.Close()
		}
	}()

	return d.ec.Encode(v)
}

func (d *defaultCodeC) Read(v interface{}) error {
	return d.dc.Decode(v)
}

func (d *defaultCodeC) Close() error {
	return d.conn.Close()
}
