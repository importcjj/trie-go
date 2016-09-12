package router

import (
	"net/http"
	"strconv"

	"github.com/importcjj/trie.go"
)

// HTTP methods
const (
	MethodGet    = "GET"
	MethodPost   = "POST"
	MethodPut    = "PUT"
	MethodPatch  = "PATCH"
	MethodDelete = "DELETE"
)

type Router struct {
	trie     *trie.Trie
	handlers map[string]*Handler
}

func New() *Router {
	return &Router{
		trie: trie.New(),
	}
}

func (router *Router) Router(pattern string, handler HandlerInterface) {
	router.trie.Put(pattern, handler)
}

func (router *Router) Get(pattern string, handlefunc func(*Context)) {
	exist, h := router.trie.Get(pattern)
	var handler HandlerInterface
	if exist {
		handler = h.(HandlerInterface)
	} else {
		handler = NewHanlder()
	}
	handler.Get(handlefunc)
	router.trie.Put(pattern, handler)
}

func (router *Router) Post(pattern string, handlefunc func(*Context)) {
	exist, h := router.trie.Get(pattern)
	var handler HandlerInterface
	if !exist {
		handler = h.(HandlerInterface)
	} else {
		handler = &Handler{}
	}
	handler.Post(handlefunc)
	router.trie.Put(pattern, handler)
}

func (router *Router) Put(pattern string, handlefunc func(*Context)) {
	exist, h := router.trie.Get(pattern)
	var handler HandlerInterface
	if !exist {
		handler = h.(HandlerInterface)
	} else {
		handler = &Handler{}
	}
	handler.Put(handlefunc)
	router.trie.Put(pattern, handler)
}

func (router *Router) Patch(pattern string, handlefunc func(*Context)) {
	exist, h := router.trie.Get(pattern)
	var handler HandlerInterface
	if !exist {
		handler = h.(HandlerInterface)
	} else {
		handler = &Handler{}
	}
	handler.Patch(handlefunc)
	router.trie.Put(pattern, handler)
}

func (router *Router) Delete(pattern string, handlefunc func(*Context)) {
	exist, h := router.trie.Get(pattern)
	var handler HandlerInterface
	if !exist {
		handler = h.(HandlerInterface)
	} else {
		handler = &Handler{}
	}
	handler.Delete(handlefunc)
	router.trie.Put(pattern, handler)
}

type Handler struct {
	funcs    map[string]func(*Context)
	OnGet    func(context *Context)
	OnPost   func(context *Context)
	OnPut    func(context *Context)
	OnPatch  func(context *Context)
	OnDelete func(context *Context)
}

func NewHanlder() *Handler {
	return &Handler{
		funcs: make(map[string]func(*Context)),
	}
}

type HandlerInterface interface {
	Get(func(*Context))
	Post(func(*Context))
	Put(func(*Context))
	Patch(func(*Context))
	Delete(func(*Context))

	DoGet(*Context)
	DoPost(*Context)
	DoPut(*Context)
	DoPatch(*Context)
	DoDelete(*Context)
}

func (handler *Handler) Get(handleFunc func(*Context)) {
	handler.OnGet = handleFunc
}

func (handler *Handler) Post(handleFunc func(*Context)) {
	handler.OnPost = handleFunc
}

func (handler *Handler) Put(handleFunc func(*Context)) {
	handler.OnPut = handleFunc
}

func (handler *Handler) Patch(handleFunc func(*Context)) {
	handler.OnPatch = handleFunc
}

func (handler *Handler) Delete(handleFunc func(*Context)) {
	handler.OnDelete = handleFunc
}

func (handler *Handler) DoGet(context *Context) {
	// handlerFunc, ok := handler.funcs[MethodGet]
	if handler.OnGet == nil {
		status := http.StatusMethodNotAllowed
		text := http.StatusText(status)
		context.ResponseWriter.Write([]byte(text))
		context.ResponseWriter.WriteHeader(status)
		return
	}
	handler.OnGet(context)
}

func (handler *Handler) DoPost(context *Context) {
	if handler.OnPost == nil {
		status := http.StatusMethodNotAllowed
		text := http.StatusText(status)
		context.ResponseWriter.Write([]byte(text))
		context.ResponseWriter.WriteHeader(status)
		return
	}
	handler.OnPost(context)
}

func (handler *Handler) DoPut(context *Context) {
	if handler.OnPut == nil {
		status := http.StatusMethodNotAllowed
		text := http.StatusText(status)
		context.ResponseWriter.Write([]byte(text))
		context.ResponseWriter.WriteHeader(status)
		return
	}
	handler.OnPut(context)
}

func (handler *Handler) DoPatch(context *Context) {
	if handler.OnPatch == nil {
		status := http.StatusMethodNotAllowed
		text := http.StatusText(status)
		context.ResponseWriter.Write([]byte(text))
		context.ResponseWriter.WriteHeader(status)
		return
	}
	handler.OnPatch(context)
}

func (handler *Handler) DoDelete(context *Context) {
	if handler.OnDelete == nil {
		status := http.StatusMethodNotAllowed
		text := http.StatusText(status)
		context.ResponseWriter.Write([]byte(text))
		context.ResponseWriter.WriteHeader(status)
		return
	}
	handler.OnDelete(context)
}

func (router *Router) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	ok, i, params := router.trie.Match(r.RequestURI)
	if !ok {
		rw.Write([]byte(http.StatusText(http.StatusNotFound)))
		// rw.WriteHeader(404)
		return
	}
	ctx := NewContent(rw, r)
	ctx.Params = params
	handler := i.(HandlerInterface)
	switch r.Method {
	case MethodGet:
		handler.DoGet(ctx)
	case MethodPost:
		handler.DoPost(ctx)
	case MethodPatch:
		handler.DoPatch(ctx)
	case MethodPut:
		handler.DoPut(ctx)
	case MethodDelete:
		handler.DoDelete(ctx)
	default:
		status := http.StatusMethodNotAllowed
		text := http.StatusText(status)
		rw.Write([]byte(text))
		rw.WriteHeader(status)
		return
	}

}

type Context struct {
	Params         map[string]string
	ResponseWriter http.ResponseWriter
	Request        *http.Request
}

func NewContent(rw http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		ResponseWriter: rw,
		Request:        r,
	}
}

func (context *Context) WriteString(v string) {
	context.ResponseWriter.Write([]byte(v))
}

func (context *Context) ParamString(key string, d ...string) string {
	if param, ok := context.Params[key]; ok {
		return param
	}
	if len(d) > 0 {
		return d[0]
	}

	return ""
}

func (context *Context) ParamInt(key string, d ...int) (int, error) {
	if param, ok := context.Params[key]; ok {
		v, err := strconv.Atoi(param)
		if err != nil {
			return 0, err
		}
		return v, nil
	}
	if len(d) > 0 {
		return d[0], nil
	}

	return 0, nil
}
