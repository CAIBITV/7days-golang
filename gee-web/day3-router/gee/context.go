package gee

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Path   string
	Method string
	Params map[string]string
	// response info
	StatusCode int
}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value

}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}

// PostForm 获取 Post 请求的表单某个键值对的值
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// Query 获取 Get 请求的参数的值
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) String(code int, format string, values ...interface{}) {
	// 注意顺序 需要先 Header().Set() 再 WriteHeader() 再 Write()
	// 在 WriteHeader() 后调用 Header().Set 是不会生效的
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	// 创建缓冲区
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	// encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		// 这里存在一个问题 因为前面已经 Header().Set() 和 WriteHeader() 了
		// encoder.Encode(obj) 相当于 Write() 了
		// 下面的 http.Error() 也相当于 是不会生效的
		// 我们通过缓冲区 buf 来解决这个问题
		//http.Error(c.Writer, err.Error(), 500)
		http.Error(c.Writer, err.Error(), 500)
	}
	c.Writer.Write(buf.Bytes())
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}
