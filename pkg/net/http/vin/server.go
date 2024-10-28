// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package vin

import (
	"context"
	"flag"
	"html/template"
	"net"
	"net/http"
	"os"
	"path"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"

	"ascale/pkg/conf/dsn"
	"ascale/pkg/conf/env"
	"ascale/pkg/log"
	"ascale/pkg/net/http/vin/bytesconv"
	"ascale/pkg/net/http/vin/render"
	"ascale/pkg/net/ip"
	"ascale/pkg/net/metadata"
	"ascale/pkg/stat"
	"ascale/pkg/xtime"
)

func init() {
	addFlag(flag.CommandLine)
}

func addFlag(fs *flag.FlagSet) {
	// tcp://0.0.0.0:8000/?timeout=1s
	v := os.Getenv("HTTP")
	fs.StringVar(&_httpDSN, "http", v, "listen http dsn, or use HTTP env variable.")
}

var (
	_httpDSN string

	stats = stat.HTTPServer
)

func parseDSN(rawdsn string) *ServerConfig {
	if rawdsn == "" {
		return nil
	}
	d, err := dsn.Parse(rawdsn)
	if err != nil {
		panic(errors.Wrapf(err, "vin: invalid dsn: %s", rawdsn))
	}
	conf := new(ServerConfig)
	if _, err = d.Bind(conf); err != nil {
		panic(errors.Wrapf(err, "vin: invalid dsn: %s", rawdsn))
	}
	return conf
}

const defaultMultipartMemory = 32 << 20 // 32 MB

var (
	default404Body   = []byte("404 page not found")
	default405Body   = []byte("405 method not allowed")
	defaultAppEngine bool
)

// HandlerFunc defines the handler used by vin middleware as return value.
type HandlerFunc func(*Context)

// HandlersChain defines a HandlerFunc array.
type HandlersChain []HandlerFunc

// Last returns the last handler in the chain. ie. the last handler is the main one.
func (c HandlersChain) Last() HandlerFunc {
	if length := len(c); length > 0 {
		return c[length-1]
	}
	return nil
}

// RouteInfo represents a request route's specification which contains method and path and its handler.
type RouteInfo struct {
	Method      string
	Path        string
	Handler     string
	HandlerFunc HandlerFunc
}

// RoutesInfo defines a RouteInfo array.
type RoutesInfo []RouteInfo

type ServerConfig struct {
	Network      string         `dsn:"network"`
	Address      string         `dsn:"address"`
	Timeout      xtime.Duration `dsn:"query.timeout"`
	ReadTimeout  xtime.Duration `dsn:"query.readTimeout"`
	WriteTimeout xtime.Duration `dsn:"query.writeTimeout"`
}

// Engine is the framework's instance, it contains the muxer, middleware and configuration settings.
// Create an instance of Engine, by using New() or Default()
type Engine struct {
	RouterGroup

	// Enables automatic redirection if the current route can't be matched but a
	// handler for the path with (without) the trailing slash exists.
	// For example if /foo/ is requested but a route only exists for /foo, the
	// client is redirected to /foo with http status code 301 for GET requests
	// and 307 for all other request methods.
	RedirectTrailingSlash bool

	// If enabled, the router tries to fix the current request path, if no
	// handle is registered for it.
	// First superfluous path elements like ../ or // are removed.
	// Afterwards the router does a case-insensitive lookup of the cleaned path.
	// If a handle can be found for this route, the router makes a redirection
	// to the corrected path with status code 301 for GET requests and 307 for
	// all other request methods.
	// For example /FOO and /..//Foo could be redirected to /foo.
	// RedirectTrailingSlash is independent of this option.
	RedirectFixedPath bool

	// If enabled, the router checks if another method is allowed for the
	// current route, if the current request can not be routed.
	// If this is the case, the request is answered with 'Method Not Allowed'
	// and HTTP status code 405.
	// If no other Method is allowed, the request is delegated to the NotFound
	// handler.
	HandleMethodNotAllowed bool
	ForwardedByClientIP    bool

	// If enabled, the url.RawPath will be used to find parameters.
	UseRawPath bool

	// If true, the path value will be unescaped.
	// If UseRawPath is false (by default), the UnescapePathValues effectively is true,
	// as url.Path gonna be used, which is already unescaped.
	UnescapePathValues bool

	// Value of 'maxMemory' param that is given to http.Request's ParseMultipartForm
	// method call.
	MaxMultipartMemory int64

	// RemoveExtraSlash a parameter can be parsed from the URL even with extra slashes.
	// See the PR #1817 and issue #1644
	RemoveExtraSlash bool

	HTMLRender  render.HTMLRender
	FuncMap     template.FuncMap
	allNoRoute  HandlersChain
	allNoMethod HandlersChain
	noRoute     HandlersChain
	noMethod    HandlersChain
	pool        sync.Pool
	trees       methodTrees

	address string

	metastore map[string]map[string]interface{} // metastore is the path as key and the metadata of this path as value, it export via /metadata
	server    atomic.Value                      // store *http.Server

	lock sync.RWMutex
	conf *ServerConfig
}

var _ IRouter = &Engine{}

// Start listen and serve bm engine by given DSN.
func (engine *Engine) Start() error {
	conf := engine.conf
	l, err := net.Listen(conf.Network, conf.Address)
	if err != nil {
		errors.Wrapf(err, "gin: listen tcp: %s", conf.Address)
		return err
	}

	log.Infof("gin: start http listen addr: %s", conf.Address)
	server := &http.Server{
		ReadTimeout:  time.Duration(conf.ReadTimeout),
		WriteTimeout: time.Duration(conf.WriteTimeout),
	}
	go func() {
		if err := engine.RunServer(server, l); err != nil {
			if errors.Cause(err) == http.ErrServerClosed {
				log.Info("gin: server closed")
				return
			}
			panic(errors.Wrapf(err, "vin: engine.ListenServer(%+v, %+v)", server, l))
		}
	}()

	return nil
}

// New returns a new blank Engine instance without any middleware attached.
// By default the configuration is:
// - RedirectTrailingSlash:  true
// - RedirectFixedPath:      false
// - HandleMethodNotAllowed: false
// - ForwardedByClientIP:    true
// - UseRawPath:             false
// - UnescapePathValues:     true
func New() *Engine {
	engine := &Engine{
		RouterGroup: RouterGroup{
			Handlers: nil,
			basePath: "/",
			root:     true,
		},
		metastore: make(map[string]map[string]interface{}),
		conf: &ServerConfig{
			Timeout: xtime.Duration(time.Second),
		},
		FuncMap:                template.FuncMap{},
		RedirectTrailingSlash:  true,
		RedirectFixedPath:      false,
		HandleMethodNotAllowed: false,
		ForwardedByClientIP:    true,
		UseRawPath:             false,
		RemoveExtraSlash:       false,
		UnescapePathValues:     true,
		address:                ip.InternalIP(),
		MaxMultipartMemory:     defaultMultipartMemory,
		trees:                  make(methodTrees, 0, 9),
	}
	engine.RouterGroup.engine = engine
	engine.pool.New = func() interface{} {
		return engine.allocateContext()
	}

	engine.addRoute("GET", "/metrics", HandlersChain{monitor()})
	engine.addRoute("GET", "/metadata", HandlersChain{engine.metadata()})

	startPerf()
	return engine
}

func NewServer(conf *ServerConfig) *Engine {
	if !flag.Parsed() {
		log.Warn("[vin] please call flag.Parse() before Init warden server, some configure may not effect.")
	}
	envConf := parseDSN(_httpDSN)
	if envConf != nil {
		conf = envConf
	}

	engine := &Engine{
		RouterGroup: RouterGroup{
			Handlers: nil,
			basePath: "/",
			root:     true,
		},
		metastore:              make(map[string]map[string]interface{}),
		conf:                   conf,
		FuncMap:                template.FuncMap{},
		RedirectTrailingSlash:  true,
		RedirectFixedPath:      false,
		HandleMethodNotAllowed: false,
		ForwardedByClientIP:    true,
		UseRawPath:             false,
		RemoveExtraSlash:       false,
		UnescapePathValues:     true,
		address:                ip.InternalIP(),
		MaxMultipartMemory:     defaultMultipartMemory,
		trees:                  make(methodTrees, 0, 9),
	}
	if err := engine.SetConfig(conf); err != nil {
		panic(err)
	}
	engine.RouterGroup.engine = engine
	engine.addRoute("GET", "/metrics", HandlersChain{monitor()})
	engine.addRoute("GET", "/metadata", HandlersChain{engine.metadata()})
	startPerf()
	engine.pool.New = func() interface{} {
		return engine.allocateContext()
	}
	return engine
}

func DefaultServer(conf *ServerConfig) *Engine {
	engine := NewServer(conf)
	engine.Use(Recovery(), Trace(env.AppID), Logger(), CORS(), UUID())
	return engine
}

// Default returns an Engine instance with the Logger and Recovery middleware already attached.
func Default() *Engine {
	engine := New()
	engine.Use(Recovery(), Trace(env.AppID), Logger(), CORS(), UUID())
	return engine
}

func (engine *Engine) SetConfig(conf *ServerConfig) (err error) {
	if conf.Timeout <= 0 {
		return errors.New("vin: config timeout must greater than 0")
	}
	if conf.Network == "" {
		conf.Network = "tcp"
	}
	engine.lock.Lock()
	engine.conf = conf
	engine.lock.Unlock()
	return
}

func (engine *Engine) allocateContext() *Context {
	return &Context{engine: engine, KeysMutex: &sync.RWMutex{}}
}

// SetFuncMap sets the FuncMap used for template.FuncMap.
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.FuncMap = funcMap
}

// NoRoute adds handlers for NoRoute. It return a 404 code by default.
func (engine *Engine) NoRoute(handlers ...HandlerFunc) {
	engine.noRoute = handlers
	engine.rebuild404Handlers()
}

// NoMethod sets the handlers called when... TODO.
func (engine *Engine) NoMethod(handlers ...HandlerFunc) {
	engine.noMethod = handlers
	engine.rebuild405Handlers()
}

// Use attaches a global middleware to the router. ie. the middleware attached though Use() will be
// included in the handlers chain for every single request. Even 404, 405, static files...
// For example, this is the right place for a logger or error management middleware.
func (engine *Engine) Use(middleware ...HandlerFunc) IRoutes {
	engine.RouterGroup.Use(middleware...)
	engine.rebuild404Handlers()
	engine.rebuild405Handlers()
	return engine
}

// Ping is used to set the general HTTP ping handler.
func (engine *Engine) Ping(handler HandlerFunc) {
	engine.GET("/monitor/ping", handler)
}

// Register is used to export metadata to discovery.
func (engine *Engine) Register(handler HandlerFunc) {
	engine.GET("/register", handler)
}
func (engine *Engine) rebuild404Handlers() {
	engine.allNoRoute = engine.combineHandlers(engine.noRoute)
}

func (engine *Engine) rebuild405Handlers() {
	engine.allNoMethod = engine.combineHandlers(engine.noMethod)
}

func (engine *Engine) addRoute(method, path string, handlers HandlersChain) {
	assert1(path[0] == '/', "path must begin with '/'")
	assert1(method != "", "HTTP method can not be empty")
	assert1(len(handlers) > 0, "there must be at least one handler")

	if _, ok := engine.metastore[path]; !ok {
		engine.metastore[path] = make(map[string]interface{})
	}
	engine.metastore[path]["method"] = method
	// debugPrintRoute(method, path, handlers)
	root := engine.trees.get(method)
	if root == nil {
		root = new(node)
		root.fullPath = "/"
		engine.trees = append(engine.trees, methodTree{method: method, root: root})
	}
	root.addRoute(path, handlers)
}

// Routes returns a slice of registered routes, including some useful information, such as:
// the http method, path and the handler name.
func (engine *Engine) Routes() (routes RoutesInfo) {
	for _, tree := range engine.trees {
		routes = iterate("", tree.method, routes, tree.root)
	}
	return routes
}

func iterate(path, method string, routes RoutesInfo, root *node) RoutesInfo {
	path += root.path
	if len(root.handlers) > 0 {
		handlerFunc := root.handlers.Last()
		routes = append(routes, RouteInfo{
			Method:      method,
			Path:        path,
			Handler:     nameOfFunction(handlerFunc),
			HandlerFunc: handlerFunc,
		})
	}
	for _, child := range root.children {
		routes = iterate(path, method, routes, child)
	}
	return routes
}

// Run attaches the router to a http.Server and starts listening and serving HTTP requests.
// It is a shortcut for http.ListenAndServe(addr, router)
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func (engine *Engine) Run(addr ...string) (err error) {
	// defer func() { debugPrintError(err) }()

	address := resolveAddress(addr)

	server := &http.Server{
		Addr:    address,
		Handler: engine,
	}
	engine.server.Store(server)
	if err = server.ListenAndServe(); err != nil {
		err = errors.Wrapf(err, "addrs: %v", addr)
	}
	return
}

// RunTLS attaches the router to a http.Server and starts listening and serving HTTPS (secure) requests.
// It is a shortcut for http.ListenAndServeTLS(addr, certFile, keyFile, router)
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func (engine *Engine) RunTLS(addr, certFile, keyFile string) (err error) {
	server := &http.Server{
		Addr:    addr,
		Handler: engine,
	}
	engine.server.Store(server)
	if err = server.ListenAndServeTLS(certFile, keyFile); err != nil {
		err = errors.Wrapf(err, "tls: %s/%s:%s", addr, certFile, keyFile)
	}
	return
}

// RunServer will serve and start listening HTTP requests by given server and listener.
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func (engine *Engine) RunServer(server *http.Server, l net.Listener) (err error) {
	server.Handler = engine
	engine.server.Store(server)
	if err = server.Serve(l); err != nil {
		err = errors.Wrapf(err, "listen server: %+v/%+v", server, l)
		return
	}
	return
}

// Server is used to load stored http server.
func (engine *Engine) Server() *http.Server {
	s, ok := engine.server.Load().(*http.Server)
	if !ok {
		return nil
	}
	return s
}

// Shutdown the http server without interrupting active connections.
func (engine *Engine) Shutdown(ctx context.Context) error {
	server := engine.Server()
	if server == nil {
		return errors.New("mars: no server")
	}
	return errors.WithStack(server.Shutdown(ctx))
}

func (engine *Engine) metadata() HandlerFunc {
	return func(c *Context) {
		c.JSON(engine.metastore, nil)
	}
}

// RunUnix attaches the router to a http.Server and starts listening and serving HTTP requests
// through the specified unix socket (ie. a file).
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func (engine *Engine) RunUnix(file string) (err error) {
	os.Remove(file)
	listener, err := net.Listen("unix", file)
	if err != nil {
		err = errors.Wrapf(err, "unix: %s", file)
		return
	}
	defer listener.Close()
	server := &http.Server{
		Handler: engine,
	}
	engine.server.Store(server)
	if err = server.Serve(listener); err != nil {
		err = errors.Wrapf(err, "unix: %s", file)
	}
	return
}

// RunListener attaches the router to a http.Server and starts listening and serving HTTP requests
// through the specified net.Listener
func (engine *Engine) RunListener(listener net.Listener) (err error) {
	// debugPrint("Listening and serving HTTP on listener what's bind with address@%s", listener.Addr())
	// defer func() { debugPrintError(err) }()
	err = http.Serve(listener, engine)
	return
}

// ServeHTTP conforms to the http.Handler interface.
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := engine.pool.Get().(*Context)
	c.writermem.reset(w)
	c.Request = req
	c.reset()

	engine.lock.RLock()
	tm := time.Duration(engine.conf.Timeout)
	engine.lock.RUnlock()

	if ctm := timeout(req); ctm > 0 && tm > ctm {
		tm = ctm
	}
	md := metadata.MD{
		metadata.Color:      color(req),
		metadata.RemoteIP:   remoteIP(req),
		metadata.RemotePort: remotePort(req),
		metadata.Caller:     caller(req),
		metadata.Mirror:     mirror(req),
	}
	ctx := metadata.NewContext(context.Background(), md)

	var cancel func()
	if tm > 0 {
		c.Context, cancel = context.WithTimeout(ctx, tm)
	} else {
		c.Context, cancel = context.WithCancel(ctx)
	}
	defer cancel()

	engine.handleHTTPRequest(c)

	engine.pool.Put(c)
}

func (engine *Engine) handleHTTPRequest(c *Context) {
	httpMethod := c.Request.Method
	rPath := c.Request.URL.Path
	unescape := false
	if engine.UseRawPath && len(c.Request.URL.RawPath) > 0 {
		rPath = c.Request.URL.RawPath
		unescape = engine.UnescapePathValues
	}

	if engine.RemoveExtraSlash {
		rPath = cleanPath(rPath)
	}

	// Find root of the tree for the given HTTP method
	t := engine.trees
	for i, tl := 0, len(t); i < tl; i++ {
		if t[i].method != httpMethod {
			continue
		}
		root := t[i].root
		// Find route in tree
		value := root.getValue(rPath, c.Params, unescape)
		if value.handlers != nil {
			c.handlers = value.handlers
			c.Params = value.params
			c.fullPath = value.fullPath
			c.Next()
			c.writermem.WriteHeaderNow()
			return
		}
		if httpMethod != "CONNECT" && rPath != "/" {
			if value.tsr && engine.RedirectTrailingSlash {
				redirectTrailingSlash(c)
				return
			}
			if engine.RedirectFixedPath && redirectFixedPath(c, root, engine.RedirectFixedPath) {
				return
			}
		}
		break
	}

	if engine.HandleMethodNotAllowed {
		for _, tree := range engine.trees {
			if tree.method == httpMethod {
				continue
			}
			if value := tree.root.getValue(rPath, nil, unescape); value.handlers != nil {
				c.handlers = engine.allNoMethod
				serveError(c, http.StatusMethodNotAllowed, default405Body)
				return
			}
		}
	}

	c.handlers = engine.allNoRoute
	serveError(c, http.StatusNotFound, default404Body)
}

var mimePlain = []string{MIMEPlain}

func serveError(c *Context, code int, defaultMessage []byte) {
	c.writermem.status = code
	c.Next()
	if c.writermem.Written() {
		return
	}
	if c.writermem.Status() == code {
		c.writermem.Header()["Content-Type"] = mimePlain
		_, err := c.Writer.Write(defaultMessage)
		if err != nil {
			// debugPrint("cannot write message to writer during serve error: %v", err)
		}
		return
	}
	c.writermem.WriteHeaderNow()
}

func redirectTrailingSlash(c *Context) {
	req := c.Request
	p := req.URL.Path
	if prefix := path.Clean(c.Request.Header.Get("X-Forwarded-Prefix")); prefix != "." {
		p = prefix + "/" + req.URL.Path
	}
	req.URL.Path = p + "/"
	if length := len(p); length > 1 && p[length-1] == '/' {
		req.URL.Path = p[:length-1]
	}
	redirectRequest(c)
}

func redirectFixedPath(c *Context, root *node, trailingSlash bool) bool {
	req := c.Request
	rPath := req.URL.Path

	if fixedPath, ok := root.findCaseInsensitivePath(cleanPath(rPath), trailingSlash); ok {
		req.URL.Path = bytesconv.BytesToString(fixedPath)
		redirectRequest(c)
		return true
	}
	return false
}

func redirectRequest(c *Context) {
	req := c.Request
	// rPath := req.URL.Path
	rURL := req.URL.String()

	code := http.StatusMovedPermanently // Permanent redirect, request with GET method
	if req.Method != http.MethodGet {
		code = http.StatusTemporaryRedirect
	}
	// debugPrint("redirecting request %d: %s --> %s", code, rPath, rURL)
	http.Redirect(c.Writer, req, rURL, code)
	c.writermem.WriteHeaderNow()
}
