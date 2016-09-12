Trie.go
=======

## Usage

```go
tree := trie.New()
# Put(pattern string, value interface())
tree.Put("/", "root")
tree.Put("/<id:int>", "name pattern")

# Has(pattern string) bool
duplicated := tree.Has("/")

# Match(key string) bool, map[string]interface{}, interface{}
ok, m, value := tree.Match("/")
# ok is true
# m is nil
# value.(string) is "root"

ok, m, value = tree.Match("/123")
# ok is true
# m is {"id": 123}
# value.(string) is "name pattern"

ok, m, value = tree.Match("/hi")
# ok is false
```

## Examples

#### A HTTP router base on it.

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
