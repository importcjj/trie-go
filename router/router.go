package router

import (
	"net/http"
	"strconv"

	"github.com/importcjj/trie-go"
)

// HTTP methods
const (
	MethodGet    = "GET"
	MethodPost   = "POST"
	MethodPut    = "PUT"
	MethodPatch  = "PATCH"
	MethodDelete = "DELETE"
)

// Router is a ServeMux.
type Router struct {
	trie     *trie.Trie
	handlers map[string]*Handler
}

// New returns a new Router object.
func New() *Router {
	return &Router{
		trie: trie.New(),
	}
}

// Router binds a handler to the pattern. When a requet come in, try to
// call the matched handle func to handle the request.
func (router *Router) Router(pattern string, handler HandlerInterface) {
	router.trie.Put(pattern, handler)
}

// Before binds midwares to a specific pattern. when a HTTP who's URL is matched with
// this pattern, call these midwares at first.
func (router *Router) Before(pattern string, midwares ...func(*Context)) {
	exist, node := router.trie.GetNode(pattern)
	if exist {
		if beforeHooks, ok := node.Data["_before_hooks"]; ok {
			hooks := beforeHooks.([]func(*Context))
			node.Data["_before_hooks"] = append(hooks, midwares...)
			return
		}
		node.Data["_before_hooks"] = midwares
		return
	}
	// warn log
}

// After binds midwares to a specific pattern. when a HTTP who's URL is matched with
// this pattern, call these midwares at last.
func (router *Router) After(pattern string, midwares ...func(context *Context)) {
	exist, node := router.trie.GetNode(pattern)
	if exist {
		if afterHooks, ok := node.Data["_after_hooks"]; ok {
			hooks := afterHooks.([]func(*Context))
			node.Data["_after_hooks"] = append(hooks, midwares...)
			return
		}
		node.Data["_after_hooks"] = midwares
		return
	}
	// warn log
}

// Get Binds the handlefunc which just handle the GET request to the pattern.
func (router *Router) Get(pattern string, handlefunc func(*Context)) {
	exist, node := router.trie.GetNode(pattern)
	var handler HandlerInterface
	if exist {
		h := node.Value
		handler = h.(HandlerInterface)
	} else {
		handler = NewHanlder()
	}
	handler.Get(handlefunc)
	router.trie.Put(pattern, handler)
}

// Post Binds the handlefunc which just handle the POST request to the pattern.
func (router *Router) Post(pattern string, handlefunc func(*Context)) {
	exist, node := router.trie.GetNode(pattern)
	var handler HandlerInterface
	if exist {
		h := node.Value
		handler = h.(HandlerInterface)
	} else {
		handler = &Handler{}
	}
	handler.Post(handlefunc)
	router.trie.Put(pattern, handler)
}

// Put Binds the handlefunc which just handle the PUT request to the pattern.
func (router *Router) Put(pattern string, handlefunc func(*Context)) {
	exist, node := router.trie.GetNode(pattern)
	var handler HandlerInterface
	if exist {
		h := node.Value
		handler = h.(HandlerInterface)
	} else {
		handler = &Handler{}
	}
	handler.Put(handlefunc)
	router.trie.Put(pattern, handler)
}

// Patch Binds the handlefunc which just handle the PATCH request to the pattern.
func (router *Router) Patch(pattern string, handlefunc func(*Context)) {
	exist, node := router.trie.GetNode(pattern)
	var handler HandlerInterface
	if exist {
		h := node.Value
		handler = h.(HandlerInterface)
	} else {
		handler = &Handler{}
	}
	handler.Patch(handlefunc)
	router.trie.Put(pattern, handler)
}

// Delete Binds the handlefunc which just handle the DELETE request to the pattern.
func (router *Router) Delete(pattern string, handlefunc func(*Context)) {
	exist, node := router.trie.GetNode(pattern)
	var handler HandlerInterface
	if exist {
		h := node.Value
		handler = h.(HandlerInterface)
	} else {
		handler = &Handler{}
	}
	handler.Delete(handlefunc)
	router.trie.Put(pattern, handler)
}

// Handler is a HTTP handler.
type Handler struct {
	funcs    map[string]func(*Context)
	OnGet    func(context *Context)
	OnPost   func(context *Context)
	OnPut    func(context *Context)
	OnPatch  func(context *Context)
	OnDelete func(context *Context)
}

// NewHanlder returns a new Handler object.
func NewHanlder() *Handler {
	return &Handler{
		funcs: make(map[string]func(*Context)),
	}
}

// HandlerInterface is a Interface.
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

// Get updates OnGet method.
func (handler *Handler) Get(handleFunc func(*Context)) {
	handler.OnGet = handleFunc
}

// Post updates OnPost method.
func (handler *Handler) Post(handleFunc func(*Context)) {
	handler.OnPost = handleFunc
}

// Put updates OnPut method.
func (handler *Handler) Put(handleFunc func(*Context)) {
	handler.OnPut = handleFunc
}

// Patch updates OnPatch method.
func (handler *Handler) Patch(handleFunc func(*Context)) {
	handler.OnPatch = handleFunc
}

// Delete updates OnDelete method.
func (handler *Handler) Delete(handleFunc func(*Context)) {
	handler.OnDelete = handleFunc
}

// DoGet handle the GET request.
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

// DoPost handle the POST request.
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

// DoPut handle the PUT request.
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

// DoPatch handle the PATCH request.
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

// DoDelete handle the DELETE request.
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

func (router *Router) mapMidwares(context *Context, midwares ...func(*Context)) {
	for _, midware := range midwares {
		midware(context)
	}
}

func (router *Router) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	ok, result := router.trie.Match(r.RequestURI)
	if !ok {
		rw.Write([]byte(http.StatusText(http.StatusNotFound)))
		// rw.WriteHeader(404)
		return
	}
	ctx := NewContent(rw, r)
	ctx.Params = result.Params
	h := result.Value
	handler := h.(HandlerInterface)

	// before hooks
	for _, data := range result.ChainData {
		hooks := data["_before_hooks"]
		if hooks != nil {
			for _, hook := range hooks.([]func(*Context)) {
				hook(ctx)
			}
		}
	}

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
	}
	// after hooks
	for _, data := range result.ChainData {
		hooks := data["_after_hooks"]
		if hooks != nil {
			for _, hook := range hooks.([]func(*Context)) {
				hook(ctx)
			}
		}
	}
}

// Context store the context of a request.
type Context struct {
	Params         map[string]string
	Data           map[string]interface{}
	ResponseWriter http.ResponseWriter
	Request        *http.Request
}

// NewContent returns a new Context object.
func NewContent(rw http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		ResponseWriter: rw,
		Data:           make(map[string]interface{}),
		Request:        r,
	}
}

// WriteString write string to ResponseWriter.
func (context *Context) WriteString(v string) {
	context.ResponseWriter.Write([]byte(v))
}

// ParamString Get param string field with the specific key.
func (context *Context) ParamString(key string, d ...string) string {
	if param, ok := context.Params[key]; ok {
		return param
	}
	if len(d) > 0 {
		return d[0]
	}

	return ""
}

// ParamInt Get param int field with the specific key.
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
