package dipra

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"
	"sync"
)

var (
	version = `v1.0.3`
	welcome = `
 Welcome to :
 ______  _                    
 |  _  \(_)                   
 | | | | _  _ __   _ __  __ _ 
 | | | || || '_ \ | '__|/ _  |
 | |/ / | || |_) || |  | (_| |
 |___/  |_|| .__/ |_|   \__,_|
 	| |                
 	|_|  %v             										  	
 Mini framework for build microservice, High speed and small size`
)

type (

	// HandlerFunc ro running HandlerFunc context
	HandlerFunc func(*Context) error

	// MiddlewareFunc to handle middleware
	MiddlewareFunc func(HandlerFunc) HandlerFunc

	// Engine core of dipra
	Engine struct {
		// Prefix
		Prefix string

		// Route
		Route []Route

		// MiddlewareFunc func
		HandleMiddleware []MiddlewareFunc

		// Sync.Pool
		Pool sync.Pool

		// Node
		Node Node

		// IsDebug...
		IsDebug bool
	}

	// Route for handler routing
	Route struct {
		Path    string
		Method  string
		Handler HandlerFunc
	}

	// M map[string]interface{}
	M map[string]interface{}
)

const (
	// DefaultPort is 8080
	DefaultPort string = ":8080"
)

// Default Engine
func Default() *Engine {
	e := &Engine{
		Route: []Route{},
		Node:  Node{},
	}
	e.Pool.New = func() interface{} {
		return e.InitialContext(nil, nil)
	}
	e.Route = append(e.Route, Route{
		Path:   "/",
		Method: http.MethodGet,
		Handler: func(c *Context) error {
			return c.JSON(http.StatusOK, M{
				"version":     version,
				"message":     "Welcome to dipra, have fun",
				"code":        http.StatusOK,
				"quick_start": "https://github.com/didikprabowo/dipra#Installation",
				"language":    "GO(Golang)",
				"author": M{
					"name":   "Didik Prabowo",
					"github": "https://github.com/didikprabowo",
				},
			})
		},
	})
	e.IsDebug = true
	return e
}

// Debug for used config debug
func (e *Engine) Debug(debug bool) {
	e.IsDebug = debug
}

// InitialContext to define
func (e *Engine) InitialContext(w http.ResponseWriter, r *http.Request) *Context {
	c := &Context{
		ResponseWriter: w,
		Request:        r,
		Writen: ResponseWriter{
			Response:   w,
			StatusCode: http.StatusOK,
		},
		Params:  Param{},
		Binding: Binding{Request: r},
	}

	return c
}

// AddRoute is used for set routing(path,mehtod, handle),
// By besides be able to set middleware
func (e *Engine) AddRoute(path, method string, handler HandlerFunc, middleware ...MiddlewareFunc) {

	if exist := e.findRouter(method, path, handler); !exist {
		e.Route = append(e.Route,
			Route{Path: e.Prefix + path,
				Method:  method,
				Handler: handler,
			})
	}

	e.HandleMiddleware = append(e.HandleMiddleware, middleware...)
}

// Use is used for add handlefuncs
func (e *Engine) Use(middleware ...MiddlewareFunc) {
	e.HandleMiddleware = append(e.HandleMiddleware, middleware...)
}

// Group is used for grouped route
func (e *Engine) Group(group string, m ...MiddlewareFunc) *Engine {
	e.Prefix = group
	e.HandleMiddleware = append(e.HandleMiddleware, m...)
	return e
}

// GET is used HTTP Request with GET METHOD
func (e *Engine) GET(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	e.AddRoute(path, http.MethodGet, handler, middleware...)
}

// POST is used HTTP Request with POST METHOD
func (e *Engine) POST(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	e.AddRoute(path, http.MethodPost, handler, middleware...)
}

// PUT is used HTTP Request with PUT METHOD
func (e *Engine) PUT(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	e.AddRoute(path, http.MethodPut, handler, middleware...)
}

// PATCH is used HTTP Request with PATCH METHOD
func (e *Engine) PATCH(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	e.AddRoute(path, http.MethodPatch, handler, middleware...)
}

// DELETE is used HTTP Request with DELETE METHOD
func (e *Engine) DELETE(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	e.AddRoute(path, http.MethodDelete, handler, middleware...)
}

// OPTION is used HTTP Request with OPTION METHOD
func (e *Engine) OPTION(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	e.AddRoute(path, http.MethodOptions, handler, middleware...)
}

// TRACE is used HTTP Request with TRACE METHOD
func (e *Engine) TRACE(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	e.AddRoute(path, http.MethodTrace, handler, middleware...)
}

// CONNECT is used HTTP Request with CONNECT METHOD
func (e *Engine) CONNECT(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	e.AddRoute(path, http.MethodConnect, handler, middleware...)
}

// Static is used define http to get file type
func (e *Engine) Static(prefix, root string) {
	cp := func(c *Context) error {
		cleanURL := path.Clean(c.Request.URL.String())
		name := path.Join(root, strings.ReplaceAll(cleanURL, prefix, ""))

		isExist, err := os.Stat(root)
		if err != nil || !isExist.IsDir() {
			return err
		}
		return c.File(name)
	}

	e.AddRoute(prefix+"/*", http.MethodGet, cp)
}

// defaulterrorHttp is used set error default
func defaulterrorHTTP(w http.ResponseWriter, code int, err error) MiddlewareFunc {
	return func(c HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			return c.JSON(code, M{
				"error": err,
			})
		}
	}
}

// defaultErrorHandler is used default handler
func defaultErrorHandler(c HandlerFunc, werrx *WrapError) HandlerFunc {
	return func(c *Context) error {
		return c.JSON(werrx.Code, M{
			"error": werrx,
		})
	}
}

// HandlerError for handler context error
func (e *Engine) HandlerError(err error, c *Context) {

	var (
		r          = M{}
		debug bool = true
	)

	if e != nil {
		debug = e.IsDebug
	}

	eStr, ok := err.(*WrapError)
	if !ok || !debug {
		r["code"] = http.StatusInternalServerError
		r["message"] = err.Error()
	} else {
		r["code"] = eStr.Code
		r["message"] = eStr.Message.(string)
	}

	err = c.JSON(200, M{
		"error": r,
	})
}

// ServeHTTP is used run http server
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	ctx := e.Pool.Get().(*Context)
	ctx.Reset(w, r)
	rt, err := e.HandlerRoute(ctx)
	if err != nil {
		rt.Handler = defaultErrorHandler(rt.Handler, err)
	}

	if e.HandleMiddleware != nil {
		rt.Handler = e.addMiddleware(rt.Handler)
	}

	if err := rt.Handler(ctx); err != nil {
		e.HandlerError(err, ctx)
	}

	e.Pool.Put(ctx)
}

// HandlerRoute is used running context http
func (e *Engine) HandlerRoute(c *Context) (r Route, werrx *WrapError) {

	sort.Slice(e.Route, func(i, j int) bool {
		return (e.Route[i].Path[1:len(e.Route[i].Path)] == c.URL.String()[1:len(c.URL.String())])
	})

	for _, rt := range e.Route {

		e.Node.SetNode(c, rt)
		url, err := e.Node.ReserverURI()
		if err != nil {
			werrx.Message = err.Error()
			werrx.Code = http.StatusInternalServerError
			return rt, werrx
		}
		uriCtx := c.URL.EscapedPath()[1:len(c.URL.EscapedPath())]
		if uriCtx == url {
			if c.Method != rt.Method {
				return rt, Err405
			}
			return rt, werrx
		}

	}

	return r, Err404
}

// WrapMiddleware is used wrapping with returns HandlerFunc
func (e *Engine) addMiddleware(h HandlerFunc) HandlerFunc {
	for i := len(e.HandleMiddleware) - 1; i >= 0; i-- {
		if e.HandleMiddleware[i](h) != nil {
			h = e.HandleMiddleware[i](h)
		}
	}
	return h
}

// Run Server with HTTP
func (e *Engine) Run(addr string) (err error) {
	if addr == "" {
		addr = DefaultPort
	}

	if e.IsDebug {
		fmt.Fprintf(os.Stdout, fmt.Sprintf(welcome, Blue+version+Reset))
		fmt.Fprintf(os.Stdout, fmt.Sprintf("\n Server started with http %v[::]:%v %v\n", Green, addr, Reset))
	}

	err = http.ListenAndServe(addr, e)
	if err != nil {
		return err
	}
	return err
}

// RunTLSaddr Server with HTTPS
func (e *Engine) RunTLSaddr(addr string, certFile, keyFile string) (err error) {
	if addr == "" {
		addr = DefaultPort
	}

	if e.IsDebug {
		fmt.Fprintf(os.Stdout, fmt.Sprintf(welcome, Blue+version+Reset))
		fmt.Fprintf(os.Stdout, fmt.Sprintf("\n Server started with http %v[::]:%v %v\n", Green, addr, Reset))
	}

	err = http.ListenAndServeTLS(addr, certFile, keyFile, e)
	if err != nil {
		return err
	}

	return err
}
