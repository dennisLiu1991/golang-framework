package igin

import (
	"fmt"
	"net/http"
	"strings"
)

type httpHandler func(c *Context)

// Engine 是web服务的核心结构，负责启动服务、注册路由、监听端口、接受http请求等
type Engine struct {
	*RouterGroup

	groups []*RouterGroup
	router *router
}

func NewEngine() *Engine {
	e := &Engine{
		router: newRouter(),
	}
	e.RouterGroup = NewGroup(e)
	// 把自己也作为group放进去，后面可以统一处理
	e.groups = append(e.groups, e.RouterGroup)
	return e
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	fmt.Printf("receive request,[%s]path = %s\n", req.Method, path)

	key := genHandlerKey(req.Method, path)

	h, ok := e.router.handlerMap[key]
	if !ok {
		handler404(path, w)
		return
	}

	c := NewContext(req, w)
	c.handlers = e.findMWByPath(path)
	c.handlers = append(c.handlers, h)
	c.Next()
}

func (e *Engine) findMWByPath(path string) []httpHandler {
	var mds []httpHandler

	for _, g := range e.groups {
		if strings.HasPrefix(path, g.prefix) {
			mds = append(mds, g.middlewares...)
		}
	}
	// 需要带上engine里的middleware
	// mds = append(mds, e.middlewares...)
	return mds
}

func handler404(path string, w http.ResponseWriter) {
	fmt.Printf("can't handle this path : %s,return 404\n", path)
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 not found\n"))
}

// Run 启动服务
func (e *Engine) Run(addr string) {
	fmt.Printf("igin run on port %s ...\n", addr)
	e.printDebugInfo()
	http.ListenAndServe(addr, e)
}

func (e *Engine) printDebugInfo() {
	fmt.Println("============================debug info============================")
	handlerMap := e.router.handlerMap
	fmt.Println("register handler: ")
	for k := range handlerMap {
		arr := strings.Split(k, "-")
		fmt.Printf("[%s]%s\n", arr[0], arr[1])
	}

	fmt.Println("============================debug end============================")
}
