package Middleware

import (
	"MyGin/Context"
	"fmt"
	"time"
)

const (
	green   = "\033[97;42m"
	white   = "\033[90;47m"
	yellow  = "\033[97;43m"
	red     = "\033[97;41m"
	blue    = "\033[97;44m"
	magenta = "\033[97;45m"
	cyan    = "\033[97;46m"
	reset   = "\033[0m"
)

// 日志中间件（记录每个HTTP请求的详细信息）
func Logger() Context.HandlerFunc {
	return func(c *Context.Context) {
		// 1. 记录请求开始时间
		start := time.Now()
		path := c.Path
		method := c.Method

		// 2.继续执行业务
		c.Next()

		// 3.记录结束时间和状态
		end := time.Since(start)
		status := c.StatusCode
		if status == 0 {
			status = 200
		}

		// 4.打印日志
		codeColor := statusColor(status)
		methodColor := methodColor(method)

		fmt.Printf("%s %d %s | %v | %s %s %s\n",
			codeColor, status, reset,
			end,
			methodColor, method, reset,
			path,
		)
	}
}

// 状态码颜色选择
func statusColor(code int) string {
	switch {
	case code >= 200 && code < 300:
		return green
	case code >= 300 && code < 400:
		return white
	case code >= 400 && code < 500:
		return yellow
	default:
		return red
	}
}

// 请求方法颜色选择
func methodColor(m string) string {
	switch m {
	case "GET":
		return cyan
	case "POST":
		return green
	case "PUT":
		return yellow
	case "DELETE":
		return red
	default:
		return blue
	}
}
