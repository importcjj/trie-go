Trie.go
=======
[![GoDoc](https://godoc.org/github.com/importcjj/trie-go?status.svg)](https://godoc.org/github.com/importcjj/trie-go)
[![Build Status](https://travis-ci.org/importcjj/trie-go.svg?branch=master)](https://travis-ci.org/importcjj/trie-go)

## Usage

```go
tree := trie.New()
// Put(pattern string, value interface())
tree.Put("/", "root")
tree.Put("/<id:int>", "name pattern")

// Has(pattern string) bool
duplicated := tree.Has("/")

// Match(key string) bool, result
ok, result := tree.Match("/")
// ok is true
// result.Params is nil
// result.Value.(string) is "root"

ok, result = tree.Match("/123")
// ok is true
// result.Params is {"id": 123}
// result.Value.(string) is "name pattern"

ok, result = tree.Match("/hi")
// ok is false

ok, node := tree.GetNode("/<id:int>")
if ok {
	node.Data["put_data"] = "hello"
}
```

## Examples

#### A HTTP router base on it.

```go
package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/importcjj/trie-go/router"
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

// BasicAuth is a Midwares
func BasicAuth(ctx *router.Context) {
	fmt.Fprintln(os.Stderr, ctx.Request.URL, "Call Basic Auth.")
}

// BeforeMetric mark a time point when the request start.
func BeforeMetric(ctx *router.Context) {
	// just a example, so use the params map to
	// record the time.
	ctx.Params["time"] = time.Now().Format("Mon Jan 2 15:04:05 -0700 MST 2006")
}

// AfterMetric log the time spent to handle the requeset.
func AfterMetric(ctx *router.Context) {
	start, _ := time.Parse("Mon Jan 2 15:04:05 -0700 MST 2006", ctx.Params["time"])
	dur := time.Since(start)
	fmt.Fprintf(os.Stderr, "%s spent", dur.String())
}

var r = router.New()

func init() {
	r.Get("/hello/world", Helloworld)
	r.Get("/hi/<username:str>", ParamHandler)
	// restful api style, this pattern can match such as
	// "/page/hi.html" "/page/static/inde.html" eta.
	r.Router("/page/<filepath:*>", PageResource)

	r.Before("/", BasicAuth, BeforeMetric)
	r.After("/", AfterMetric)
}

func main() {
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	server.ListenAndServe()
}

```
