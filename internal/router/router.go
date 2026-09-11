package router

import (
	"net/http"
	"strings"
)

// Middleware transforms a handler before it reaches route logic.
type Middleware func(http.Handler) http.Handler

const routeKeySep = "\x00"

// node is one segment in the dynamic-route trie.
// Static children are O(1); a param child captures one variable segment.
type node struct {
	static    map[string]*node
	param     *node
	paramName string
	handler   http.Handler
}

// Router resolves requests via O(1) static map lookup, then trie fallback
// for patterns like /users/:id.
type Router struct {
	static     map[string]http.Handler
	dynamic    map[string]*node
	middleware []Middleware
	notFound   http.Handler
}

// New creates an empty router with a default 404 handler.
func New() *Router {
	return &Router{
		static:  make(map[string]http.Handler),
		dynamic: make(map[string]*node),
		notFound: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "route not found", http.StatusNotFound)
		}),
	}
}

// Use appends global middleware executed before every matched route.
func (r *Router) Use(mw ...Middleware) {
	r.middleware = append(r.middleware, mw...)
}

// Handle registers a handler for a method + path pattern.
// Static paths (no ":") are stored in a map for O(1) lookup.
// Dynamic paths (e.g. /users/:id) are stored in a per-method trie.
func (r *Router) Handle(method, pattern string, handler http.HandlerFunc) {
	method = strings.ToUpper(method)

	if strings.Contains(pattern, ":") {
		r.addDynamic(method, pattern, handler)
		return
	}

	r.static[staticKey(method, pattern)] = handler
}

// GET is a shortcut for Handle("GET", ...).
func (r *Router) GET(pattern string, handler http.HandlerFunc) {
	r.Handle(http.MethodGet, pattern, handler)
}

// POST is a shortcut for Handle("POST", ...).
func (r *Router) POST(pattern string, handler http.HandlerFunc) {
	r.Handle(http.MethodPost, pattern, handler)
}

// PUT is a shortcut for Handle("PUT", ...).
func (r *Router) PUT(pattern string, handler http.HandlerFunc) {
	r.Handle(http.MethodPut, pattern, handler)
}

// DELETE is a shortcut for Handle("DELETE", ...).
func (r *Router) DELETE(pattern string, handler http.HandlerFunc) {
	r.Handle(http.MethodDelete, pattern, handler)
}

// Param reads a captured path parameter from the request.
// Example: pattern /users/:id with path /users/42 -> Param(r, "id") == "42"
func Param(r *http.Request, name string) string {
	return r.PathValue(name)
}

// ServeHTTP implements http.Handler.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	method := req.Method
	path := req.URL.Path

	// O(1) static route lookup.
	if handler, ok := r.static[staticKey(method, path)]; ok {
		r.serve(w, req, handler, nil)
		return
	}

	// Trie lookup for dynamic routes such as /users/:id.
	if handler, params := r.matchDynamic(method, path); handler != nil {
		r.serve(w, req, handler, params)
		return
	}

	r.notFound.ServeHTTP(w, req)
}

func (r *Router) addDynamic(method, pattern string, handler http.Handler) {
	root, ok := r.dynamic[method]
	if !ok {
		root = &node{static: make(map[string]*node)}
		r.dynamic[method] = root
	}

	curr := root
	for _, segment := range splitPath(pattern) {
		if strings.HasPrefix(segment, ":") {
			name := segment[1:]
			if curr.param == nil {
				curr.param = &node{
					static:    make(map[string]*node),
					paramName: name,
				}
			}
			curr = curr.param
			continue
		}

		child, ok := curr.static[segment]
		if !ok {
			child = &node{static: make(map[string]*node)}
			curr.static[segment] = child
		}
		curr = child
	}

	curr.handler = handler
}

func (r *Router) matchDynamic(method, path string) (http.Handler, map[string]string) {
	root, ok := r.dynamic[method]
	if !ok {
		return nil, nil
	}

	params := make(map[string]string)
	handler := walkTrie(root, splitPath(path), 0, params)
	if handler == nil {
		return nil, nil
	}

	return handler, params
}

func walkTrie(n *node, segments []string, index int, params map[string]string) http.Handler {
	if index == len(segments) {
		return n.handler
	}

	segment := segments[index]

	// Prefer static segments over param capture for specificity.
	if child, ok := n.static[segment]; ok {
		if handler := walkTrie(child, segments, index+1, params); handler != nil {
			return handler
		}
	}

	if n.param != nil {
		params[n.param.paramName] = segment
		if handler := walkTrie(n.param, segments, index+1, params); handler != nil {
			return handler
		}
		delete(params, n.param.paramName)
	}

	return nil
}

func (r *Router) serve(w http.ResponseWriter, req *http.Request, handler http.Handler, params map[string]string) {
	for key, value := range params {
		req.SetPathValue(key, value)
	}

	for i := len(r.middleware) - 1; i >= 0; i-- {
		handler = r.middleware[i](handler)
	}

	handler.ServeHTTP(w, req)
}

func staticKey(method, path string) string {
	return method + routeKeySep + path
}

func splitPath(path string) []string {
	path = strings.Trim(path, "/")
	if path == "" {
		return nil
	}
	return strings.Split(path, "/")
}
