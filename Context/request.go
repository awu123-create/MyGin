package Context

// 获取请求参数
// Query 获取 Query 参数
func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

// PostForm 获取表单参数
func (c *Context) PostForm(key string) string {
	return c.Request.FormValue(key)
}

// Param 获取动态路由参数
func (c *Context) Param(key string) string {
	return c.Params[key]
}
