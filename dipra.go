package dipra

import (
	"fmt"
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
)

type (

	// HandlerFunc ro running HandlerFunc context
	HandlerFunc func(*Context) error

	// MiddlewareFunc to handle middleware
	MiddlewareFunc func(HandlerFunc) HandlerFunc

	// Engine core of dipra
	Engine struct {
		route      route
		middleware []MiddlewareFunc
		pool       sync.Pool
		IsDebug    bool
	}

	// Route for handler routing
	route struct {
		trees map[string]*node
		pool  sync.Pool
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
		route: route{
			trees: map[string]*node{},
		},
	}
	e.pool.New = func() interface{} {
		return e.InitialContext(nil, nil)
	}
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
		params:  &params{},
		Binding: Binding{Request: r},
	}

	return c
}

// AddRoute is used for set routing(path,mehtod, handle),
// By besides be able to set middleware
func (e *Engine) AddRoute(path, method string, h HandlerFunc, mh ...MiddlewareFunc) {

	isPathValid(path)

	nm := e.route.trees[method]
	if nm == nil {
		nm = &node{}
		e.route.trees[method] = nm
	}
	if !e.route.allowMethods(method) {
		panic("dipra : method not allowed")
	}

	if len(mh) > 0 {
		h = applyMiddleware(h, mh...)
	}

	nm.insert(method, path, h)
}

// Use is used for add handlefuncs
func (e *Engine) Use(middleware ...MiddlewareFunc) {
	e.middleware = append(e.middleware, middleware...)
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

// ANY is used HTTP Request with ALL METHOD
func (e *Engine) ANY(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	for _, v := range allowMethod {
		e.AddRoute(path, v, handler, middleware...)
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

	ctx := e.pool.Get().(*Context)
	ctx.Reset(w, r)
	e.HandlerRoute(ctx)
	e.pool.Put(ctx)
}

// HandlerRoute is used running context http
func (e *Engine) HandlerRoute(c *Context) {

	reqMethod := c.Method
	reqURL := c.URL.Path
	if reqURL[len(reqURL)-1] == '/' && len(reqURL) > 1 {
		reqURL = reqURL[:len(reqURL)-2]
	}

	isPathValid(reqURL)

	h := HandlerFunc(func(c *Context) error {
		return Err404
	})

	root := e.route.trees[reqMethod]
	if root != nil {
		_, params, hr := root.find(reqURL, c.params)
		if hr == nil {
			h = defaultErrorHandler(hr, Err404)
		} else {
			h = hr
			if params != nil {
				c.params.putParams(params)
			}
		}
	}

	if e.middleware != nil {
		h = applyMiddleware(h, e.middleware...)
	}

	if err := h(c); err != nil {
		e.HandlerError(err, c)
	}

	return

}

func applyMiddleware(h HandlerFunc, middleware ...MiddlewareFunc) HandlerFunc {
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h)
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
