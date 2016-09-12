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
