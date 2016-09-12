package trie

import (
	"fmt"
	"testing"
)

func TestNewPattern(t *testing.T) {
	p := NewPattern("<name:str>is<age:int>")
	ok, params := p.Match("jiajuis100")
	if !ok {
		t.Error("They should be matched.")
	} else {
		fmt.Printf("%#v\n", params)
	}
}
