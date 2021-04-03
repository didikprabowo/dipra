package dipra

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
)

var (
	version = `v1.0.6`
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

	defaultHandler = func(c *Context) error {
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
	}
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
		Route map[string][]Route

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
		Route:   map[string][]Route{},
		Node:    Node{},
		IsDebug: true,
	}
	e.Pool.New = func() interface{} {
		return e.InitialContext(nil, nil)
	}

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

	ok := e.findRouter(method, path, handler)
	if !ok {
		e.Route[method] = append(e.Route[method], Route{
			Path:    e.Prefix + path,
			Method:  method,
			Handler: handler,
		})
	} else {
		log.Printf(fmt.Sprintf("path %s %s already exist, please use another path", path, method))
		os.Exit(1)
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

// Any is used request with all method
func (e *Engine) Any(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	for _, m := range allowMethod {
		e.AddRoute(path, m, handler, middleware...)
	}
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

	err = c.JSON(r["code"].(int), M{
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
	werrx = &WrapError{}
	for _, v := range e.Route[c.Method] {
		e.Node.SetNode(c, v)
		url, err := e.Node.ReserverURI()
		if err != nil {
			switch err.Error() {
			case http.StatusText(http.StatusNotFound):
				werrx.Code = http.StatusNotFound
			default:
				werrx.Code = http.StatusInternalServerError
			}
			werrx.Message = err.Error()

			return v, werrx
		}
		uriCtx := c.URL.EscapedPath()[1:len(c.URL.EscapedPath())]
		if uriCtx == url {
			if c.Method != v.Method {
				return r, Err405
			}
			return v, nil
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
