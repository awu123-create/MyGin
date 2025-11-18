package Context

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ---------输出响应---------
// String 以String形式输出响应
func (c *Context) String(code int, format string, values ...interface{}) {
	c.Status(code)
	c.SetHeader("Content-Type", "text/plain; charset=utf-8")
	_, _ = fmt.Fprintf(c.Writer, format, values...)
}

// JSON 以JSON形式输出响应
func (c *Context) JSON(code int, obj interface{}) {
	c.Status(code)
	c.SetHeader("Content-Type", "application/json; charset=utf-8")
	encoder := json.NewEncoder(c.Writer)
	encode := encoder.Encode(obj)
	if encode != nil {
		http.Error(c.Writer, "Invalid", http.StatusInternalServerError)
	}
}

// HTML 以HTML形式输出响应
func (c *Context) HTML(code int, html string) {
	c.Status(code)
	c.SetHeader("Content-Type", "text/html; charset=utf-8")
	_, _ = c.Writer.Write([]byte(html))
}
