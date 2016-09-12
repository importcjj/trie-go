package main

import (
	"fmt"
	"net/http"

	"github.com/importcjj/trie.go/router"
)

var r = router.New()

var HelloWorld = &router.Handler{
	OnGet: func(ctx *router.Context) {
		ctx.WriteString("Hello, world!")
	},
}

func init() {
	r.Router("/helloworld", HelloWorld)
	r.Get("/hi/<username>", func(ctx *router.Context) {
		username := ctx.ParamString("username")
		text := fmt.Sprintf("Hi, %s", username)
		ctx.WriteString(text)
	})

	r.Get("/file", func(ctx *router.Context) {
		ctx.WriteString("this is file")
	})

	r.Get("/file/<filename:*>", func(ctx *router.Context) {
		filename := ctx.ParamString("filename")
		ctx.WriteString(filename)
	})
}

func main() {
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	server.ListenAndServe()
}
