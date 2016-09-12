package trie

import (
	"fmt"
	"testing"
)

func TestSplitPattern(t *testing.T) {
	var (
		p1 = "/"
		p2 = "a/b"
		p3 = "/a/b/c/"
		p4 = "a/b/"
	)

	trie := New()
	fmt.Println(trie.splitPattern(p1))
	fmt.Println(trie.splitPattern(p2))
	fmt.Println(trie.splitPattern(p3))
	fmt.Println(trie.splitPattern(p4))
}

func TestAddPattern(t *testing.T) {
	trie := New()
	trie.Put("/", "root")
	trie.Put("/hi", "hi")
	trie.Put("a/b/c", "a-b-c")
}

func TestHasPattern(t *testing.T) {
	trie := New()
	trie.Put("/", "root")
	if has := trie.Has("/"); !has {
		t.Errorf("The trie already has pattern %s", "/")
	}
}

func TestRouter(t *testing.T) {
	trie := New()
	trie.Put("/", "root")
	if ok, value, params := trie.Match("/"); !ok {
		t.Error("they should matched")
	} else {
		fmt.Println(value)
		fmt.Printf("%#v\n", params)
	}

	trie.Put("/hello/world", "hellworld")
	if ok, value, params := trie.Match("/hello/world"); !ok {
		t.Error("they should matched")
	} else {
		fmt.Println(value)
		fmt.Printf("%#v\n", params)
	}

	trie.Put("/hello/<name:str>", "hellworld")
	if ok, value, params := trie.Match("/hello/jiaju"); !ok {
		t.Error("they should matched")
	} else {
		fmt.Println(value)
		fmt.Printf("%#v\n", params)
	}
}
