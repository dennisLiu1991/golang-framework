package igin

type RouterGroup struct {
	// 所有group都持有全局同一个engine的实例
	e      *Engine
	prefix string
	parent *RouterGroup

	middlewares []httpHandler
}

func NewGroup(e *Engine) *RouterGroup {
	return &RouterGroup{
		e: e,
	}
}

func (g *RouterGroup) Use(m ...httpHandler) {
	g.middlewares = append(g.middlewares, m...)
}

func (g *RouterGroup) Group(path string) *RouterGroup {
	childG := &RouterGroup{
		e:      g.e,
		parent: g,
		prefix: g.prefix + path,
	}
	g.e.groups = append(g.e.groups, childG)
	return childG
}

func (g *RouterGroup) GET(relativePath string, h httpHandler) {
	g.addHandler("GET", relativePath, h)
}

func (g *RouterGroup) POST(relativePath string, h httpHandler) {
	g.addHandler("POST", relativePath, h)
}

func (g *RouterGroup) addHandler(method, path string, h httpHandler) {
	path = g.prefix + path
	g.e.router.addHandler(method, path, h)
}
