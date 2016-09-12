package main

import (
	"net/http"

	"github.com/importcjj/trie.go/router"
	"github.com/importcjj/trie.go/router/examples/handlers"
)

func main() {
	r := router.New()
	r.Router("/helloworld", handlers.HelloWorld)

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	server.ListenAndServe()
}
