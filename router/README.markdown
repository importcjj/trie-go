HTTP Router
===========

## Usage

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/importcjj/trie.go/router"
)

func Helloworld(ctx *router.Context) {
	ctx.WriteString("hello, world!")
}

func ParamHandler(ctx *router.Context) {
	username := ctx.ParamString("username")
	text := fmt.Sprintf("hi, %s", username)
	ctx.WriteString(text)
}

var PageResource = &router.Handler{
	OnGet: func(ctx *router.Context) {
		filepath := ctx.ParamString("filepath")
		text := fmt.Sprintf("Get page %s", filepath)
		ctx.WriteString(text)
	},

	OnPost: func(ctx *router.Context) {
		filepath := ctx.ParamString("filepath")
		text := fmt.Sprintf("Post page %s", filepath)
		ctx.WriteString(text)
	},
}

var r = router.New()

func init() {
	r.Get("/hello/world", Helloworld)
	r.Get("/hi/<username:str>", ParamHandler)
	// restful api style, this pattern can match such as
	// "/page/hi.html" "/page/static/inde.html" eta.
	r.Router("/page/<filepath:*>", PageResource)
}

func main() {
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	server.ListenAndServe()
}

```
