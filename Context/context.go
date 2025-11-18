package Context

import (
	"net/http"
)

// type H map[string]interface{}
type HandlerFunc func(*Context)
type Context struct {
	// 基础请求信息
	Writer  http.ResponseWriter
	Request *http.Request

	// 请求相关数据
	Path   string
	Method string
	Params map[string]string // 动态路由参数

	// 响应状态
	StatusCode int
	// 中间件相关
	Handlers []HandlerFunc          // 存储中间件和主业务函数
	data     map[string]interface{} // 存储中间件及主业务函数之间传递的数据
	Index    int                    // 当前执行到的中间件索引
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Path:    r.URL.Path,
		Method:  r.Method,
		Writer:  w,
		Request: r,
		Index:   -1,
		data:    make(map[string]interface{}),
		Params:  make(map[string]string),
	}
}

// Next 继续执行下一个中间件或主业务函数
func (c *Context) Next() {
	c.Index++
	for c.Index < len(c.Handlers) {
		c.Handlers[c.Index](c)
		c.Index++
	}
}

// Abort 中断执行
func (c *Context) Abort() {
	c.Index = len(c.Handlers)
}

// Status 设置响应状态码
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// SetHeader 设置响应头
func (c *Context) SetHeader(key, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) Set(key string, value interface{}) {
	c.data[key] = value
}

func (c *Context) Get(key string) (value interface{}, ok bool) {
	value, ok = c.data[key]
	return
}

func (c *Context) SetHandlers(hs []HandlerFunc) {
	c.Handlers = hs
}
