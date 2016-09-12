HTTP Router
===========

## Usage

```go
package main

import (
	"fmt"
	"net/http"
	"os"

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
