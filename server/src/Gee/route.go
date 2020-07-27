package Gee

import (
	"fmt"
	"net/http"
	"strings"
)

type route struct {
	handlers map[string]HandlerFunc
	roots map[string]*node
}

func newRoute() *route {
	return &route{
		handlers: make(map[string]HandlerFunc),
		roots: make(map[string]*node),
	}
}

func (r *route) parsePattern (pattern string) []string {
	strList := strings.Split(pattern, "/")
	list := make([]string, 0)
	for _, item := range strList{
		if item != "" {
			list = append(list, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return list
}

func (r *route) addRoute (method string, pattern string, handler HandlerFunc)  {
	parts := r.parsePattern(pattern)

	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	key := method + "-" + pattern
	r.handlers[key] = handler
	//fmt.Println("测试树：", r.roots)

	//var printChild func(*node)
	//printChild = func(n *node) {
	//	fmt.Println(fmt.Sprintf("pattern: %s part: %s ", n.pattern, n.part))
	//	if len(n.children) > 0 {
	//		for _, child := range n.children{
	//			printChild(child)
	//		}
	//	}
	//}
	//
	//for key, val := range r.roots{
	//	fmt.Println(key)
	//	//fmt.Println(*val)
	//	printChild(val)
	//}
}

func (r *route) getRoute (method, path string) (*node, map[string]string) {
	strList := r.parsePattern(path)
	fmt.Println("path: ", path)
	params := make(map[string]string)
	res, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	searchRes := res.search(strList, 0)
	if searchRes != nil {
		fmt.Println("结果：", searchRes.pattern)
		parts := r.parsePattern(searchRes.pattern)
		for key, part := range parts{
			if part[0] == ':' {
				params[part[1:]] = strList[key]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(strList[key:], "/")
				break
			}
		}
		return searchRes, params
	}
	return nil, nil
}

func (r *route) handler (c *Context)  {
	n, params := r.getRoute(c.Method, c.Path)
	fmt.Println("params: ", params)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}

