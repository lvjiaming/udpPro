package Gee

import (
	"fmt"
	"net/http"
)

type Context struct {
	Write http.ResponseWriter
	Req *http.Request

	Method string
	Path string

	Params map[string]string
	StateCode int
}


func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Write:     w,
		Req:       req,
		Method:    req.Method,
		Path:      req.URL.Path,
	}
}

func (c *Context) Param (key string) string {
	res, _ := c.Params[key]
	return res
}

func (c *Context) Status (code int)  {
	c.StateCode = code
	c.Write.WriteHeader(code)
}

func (c *Context) SetHeader (key, val string)  {
	c.Write.Header().Set(key, val)
}

func (c *Context) String (code int, format string, values ...interface{})  {
	c.Status(code)
	c.SetHeader("Content-Type", "text/plain")
	c.Write.Write([]byte(fmt.Sprintf(format, values...)))
}