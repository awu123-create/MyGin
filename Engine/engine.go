package Engine

import (
	"MyGin/Context"
	"MyGin/Middleware"
	"MyGin/Router"
	"net/http"
	"strings"
	"sync"
)

// 实现分组控制
type RouterGroup struct {
	prefix      string
	middlewares []Context.HandlerFunc
	engine      *Engine // 保存当前路由组属于哪个 Engine 实例
}

// Engine 入口
type Engine struct {
	*RouterGroup
	router *Router.Router
	//middlewares []Context.HandlerFunc
	groups      []*RouterGroup
	contextPool sync.Pool
}

func NewEngine() *Engine {
	engine := &Engine{
		router: Router.NewRouter(),
	}
	engine.RouterGroup = &RouterGroup{
		prefix: "",
		engine: engine,
	}

	// engine 有自己的全局中间件和全局前缀
	engine.groups = []*RouterGroup{engine.RouterGroup}
	engine.contextPool.New = func() any {
		return &Context.Context{} // 定义当池子里没东西时，怎么新建一个对象
	}
	return engine
}

func Default() *Engine {
	engine := &Engine{
		router: Router.NewRouter(),
	}
	engine.RouterGroup = &RouterGroup{
		prefix: "",
		engine: engine,
	}

	engine.groups = []*RouterGroup{engine.RouterGroup}
	engine.groups[0].middlewares = []Context.HandlerFunc{
		Middleware.Logger(),
		Middleware.Recovery(),
	}
	engine.contextPool.New = func() any {
		return &Context.Context{} // 定义当池子里没东西时，怎么新建一个对象
	}
	return engine
}

func (r *RouterGroup) Group(prefix string) *RouterGroup {
	engine := r.engine
	newGroup := &RouterGroup{
		prefix: r.prefix + prefix,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// Use 注册中间件
func (r *RouterGroup) Use(middlewares ...Context.HandlerFunc) {
	r.middlewares = append(r.middlewares, middlewares...)
}

func (r *RouterGroup) GET(path string, handler Context.HandlerFunc) {
	fullpath := r.prefix + path
	r.engine.router.AddRoute("GET", fullpath, handler)
}

func (r *RouterGroup) POST(path string, handler Context.HandlerFunc) {
	r.engine.router.AddRoute("POST", path, handler)
}

// Run 启动服务
func (e *Engine) Run(address string) error {
	return http.ListenAndServe(address, e) // address 是监听端口
}

/*
1. 从 Pool 获取 Context
2. 初始化 Context
3. 从 Router 查找路由 handler
4. 合并中间件 + handler
5. 执行 c.Next()
6. 回收 context
*/
func (e *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	c := e.contextPool.Get().(*Context.Context)
	c.Writer = writer
	c.Request = request
	c.Method = request.Method
	c.Path = request.URL.Path
	c.Index = -1
	c.Handlers = nil

	// 收集匹配的 Group 中间件
	for _, group := range e.groups {
		if strings.HasPrefix(c.Path, group.prefix) {
			c.Handlers = append(c.Handlers, group.middlewares...)
		}
	}
	handler, params := e.router.GetRoute(request.Method, request.URL.Path)

	if handler != nil {
		c.Params = params

		c.Handlers = append(c.Handlers, handler)

		// 启动执行链
		c.Next()
	} else {
		c.Handlers = []Context.HandlerFunc{
			func(c *Context.Context) {
				c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
			},
		}
		c.Next()
	}

	// 回收
	c.SetHandlers(nil)
	c.Params = nil
	e.contextPool.Put(c)
}
