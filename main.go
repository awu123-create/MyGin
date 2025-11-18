package main

import (
	"MyGin/Context"
	"MyGin/Engine"
)

func test() Context.HandlerFunc {
	return func(c *Context.Context) {
		c.String(200, "this is a test middleware\n")
		c.Next()
	}
}

func main() {
	e := Engine.Default()

	// 使用路由组
	v1 := e.Group("/hello")
	v1.Use(test())
	{
		v1.GET("/world", func(c *Context.Context) {
			c.JSON(200, "name:black myth WuKong")
		})
	}

	e.Run(":9090")
}
