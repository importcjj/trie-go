package handlers

import "github.com/importcjj/trie.go/router"

var HelloWorld = &router.Handler{
	OnGet: func(ctx *router.Context) {
		ctx.WriteString("Get /helloworld")
	},

	OnPost: func(ctx *router.Context) {
		ctx.WriteString("Post /helloworld")
	},
}
