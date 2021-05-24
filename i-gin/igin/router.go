package igin

import "fmt"

type router struct {
	handlerMap map[string]httpHandler
}

func newRouter() *router {
	return &router{
		handlerMap: make(map[string]httpHandler),
	}
}

// genHandlerKey 生成key，格式例如 GET-/a/b/c ，依据该格式实现请求路由
func genHandlerKey(reqFunc, path string) string {
	return fmt.Sprintf("%s-%s", reqFunc, path)
}

func (r *router) addHandler(reqFunc, path string, h httpHandler) {
	key := genHandlerKey(reqFunc, path)
	r.handlerMap[key] = h
}
