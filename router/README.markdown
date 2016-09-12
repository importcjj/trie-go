HTTP Router
===========

## Usage

```go
package main

import (
    "fmt"
    "github.com/importcjj/trie.go/router"
    "net/http"
)

var r = router.New()

var HelloWorld = &router.NewHanlder{
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
}

func main() {
    r := router.New()
    server := &http.Server{
        Addr: ":8080",
        Handler: r,
    }
    server.ListenAndServe()
}

```
