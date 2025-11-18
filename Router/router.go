package Router

import (
	"MyGin/Context"
	"strings"
)

// 实现基础路由注册和匹配功能
// 使用前缀树来支持动态路由
type node struct {
	//当前节点的路径部分
	part string

	//是否是通配符节点
	isParam    bool //:id
	isWildcard bool //*filepath

	//子节点
	children []*node

	//路由处理函数
	handler Context.HandlerFunc

	//路由中参数的名称
	paramName string
}

type Router struct {
	trees map[string]*node // method => root node
}

func NewRouter() *Router {
	return &Router{
		trees: make(map[string]*node),
	}
}

// 1.分割路径
func parsePath(path string) []string {
	parts := strings.Split(path, "/")
	res := make([]string, 0)
	for _, p := range parts {
		if p == "" {
			continue
		}
		res = append(res, p)
		if p[0] == '*' {
			break
		}
	}
	return res
}

/*
2.插入节点
(1)看看当前节点是否已有 part 子节点
(2)没有就创建一个
(3)跳到子节点，继续
*/
func (n *node) insert(parts []string, handler Context.HandlerFunc) {
	cur := n

	for _, part := range parts {
		var child *node
		for _, c := range cur.children {
			if c.part == part {
				child = c
				break
			}
		}

		if child == nil {
			child = &node{
				part:       part,
				paramName:  extractParamName(part),
				isParam:    part[0] == ':',
				isWildcard: part[0] == '*',
			}
			cur.children = append(cur.children, child)
		}

		cur = child
		if child.isWildcard {
			break
		}
	}
	cur.handler = handler
}

// 提取参数名称
func extractParamName(part string) string {
	if part[0] == ':' || part[0] == '*' {
		return part[1:]
	}
	return ""
}

// 3.匹配节点(优先级：静态路由 > 动态路由 > 通配符路由)
func (n *node) search(parts []string, params map[string]string) *node {
	// 递归终止条件
	if len(parts) == 0 {
		if n.handler != nil {
			return n
		}
		return nil
	}

	part := parts[0]
	next := parts[1:]

	// 首先匹配静态路由
	for _, child := range n.children {
		if child.part == part {
			res := child.search(next, params)
			if res != nil {
				return res
			}
		}
	}

	// 然后匹配动态路由
	for _, child := range n.children {
		if child.isParam {
			params[child.paramName] = part
			res := child.search(next, params)
			if res != nil {
				return res
			}
		}
	}

	// 最后匹配通配符路由
	for _, child := range n.children {
		if child.isWildcard {
			params[child.paramName] = strings.Join(parts, "/")
			return child
		}
	}

	return nil
}

// 4.实现注册路由部分
func (r *Router) AddRoute(method, path string, handler Context.HandlerFunc) {
	if _, ok := r.trees[method]; !ok {
		r.trees[method] = &node{}
	}

	parts := parsePath(path)
	r.trees[method].insert(parts, handler)
}

func (r *Router) GetRoute(method, path string) (Context.HandlerFunc, map[string]string) {
	root := r.trees[method]
	if root == nil {
		return nil, nil
	}

	parts := parsePath(path)
	param := make(map[string]string)
	node := root.search(parts, param) // param会被写入

	if node != nil {
		return node.handler, param
	}
	return nil, nil
}

// 实现分组控制
