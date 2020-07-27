package Gee

import (
	"net/http"
)

type HandlerFunc func(*Context)

type Engine struct {
	route *route
}

func New() *Engine {
	return &Engine{route: newRoute()}
}

func (e *Engine) addRoute (method string, pattern string, handle HandlerFunc)  {
	e.route.addRoute(method, pattern, handle)
}

func (e *Engine) GET (pattern string, handle HandlerFunc)  {
	e.addRoute("GET", pattern, handle)
}

func (e *Engine) POST (pattern string, handle HandlerFunc)  {
	e.addRoute("POST", pattern, handle)
}

/**
 实现了此方法，即可接受消息
 */
func (e *Engine) ServeHTTP (w http.ResponseWriter, req *http.Request)  {
	c := newContext(w, req)
	e.route.handler(c)
}

func (e *Engine) Run (addr string) error {
	return http.ListenAndServe(addr, e)
}