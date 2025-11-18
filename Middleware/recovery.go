package Middleware

import (
	"MyGin/Context"
	"fmt"
	"net/http"
	"runtime/debug"
)

func Recovery() Context.HandlerFunc {
	return func(c *Context.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 打印错误与堆栈
				fmt.Printf("\n[Recovery] panic: %v\n%s\n", err, debug.Stack())

				// 返回状态码 500
				c.Status(http.StatusInternalServerError)
				c.SetHeader("Content-Type", "text/plain; charset=utf-8")
				_, _ = c.Writer.Write([]byte("Internal Server Error"))

				// 中断剩余中间件，不继续执行
				c.Abort()
			}
		}()

		c.Next()
	}
}
