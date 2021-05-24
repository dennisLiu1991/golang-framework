package igin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Context struct {
	req *http.Request
	rw  http.ResponseWriter

	URL *url.URL

	hIndex   int
	handlers []httpHandler
}

func NewContext(r *http.Request, rw http.ResponseWriter) *Context {
	return &Context{
		req: r,
		rw:  rw,
		URL: r.URL,
		// 设置为-1 为了能够从c.Next()开始
		hIndex: -1,
	}
}

func (c *Context) Next() {
	c.hIndex++
	if c.hIndex > len(c.handlers) {
		return
	}
	c.handlers[c.hIndex](c)
}

func (c *Context) GetRequest() *http.Request {
	return c.req
}

func (c *Context) GetResponseWriter() http.ResponseWriter {
	return c.rw
}

func (c *Context) String(content string) {
	c.rw.Write([]byte(content))
}

func (c *Context) JSON(jsonMap map[string]interface{}) {
	d, err := json.Marshal(jsonMap)
	if err != nil {
		fmt.Printf("marshal json : %v", err)
		return
	}

	c.rw.Write(d)
}
