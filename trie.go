package trie

import (
	"errors"
	"strings"
)

const (
	defalutDelimeter = "/"
)

// Errors
var (
	ErrDuplicatedPatternName = errors.New("trie.go: the pattern name already exists")
)

// Trie is a tree
type Trie struct {
	Delimeter string
	Root      *Node
}

// New returns a new Trie object.
func New() *Trie {
	return &Trie{
		Delimeter: defalutDelimeter,
		Root:      NewNode(),
	}
}

// SetDelimeter sets the delimeter of the trie object.
func (trie *Trie) SetDelimeter(delimeter string) {
	trie.Delimeter = delimeter
}

func (trie *Trie) splitPattern(pattern string) []string {
	p := strings.TrimRight(pattern, trie.Delimeter)
	parts := strings.Split(p, trie.Delimeter)
	if parts[0] == "" {
		parts[0] = string(pattern[0])
	}
	return parts
}

// Put add a new pattern to the trie tree.
func (trie *Trie) Put(pattern string, value interface{}) error {
	// duplicated pattern.
	if trie.Has(pattern) {
		return ErrDuplicatedPatternName
	}
	parts := trie.splitPattern(pattern)
	t := trie.Root
	for i, part := range parts {
		if node, ok := t.Children[part]; ok {
			t = node
			continue
		}
		// make a new node and put it to t's children nodes.
		node := NewNode(part)
		node.Value = value
		t.Children[part] = node
		if i == len(parts)-1 {
			node.patternEnding = true
		}
		t = node
	}

	return nil
}

// Has returns true if the pattern is duplicated
// in the trie tree. Otherwise, returns false.
func (trie *Trie) Has(pattern string) bool {
	parts := trie.splitPattern(pattern)
	t := trie.Root
	for _, part := range parts {
		if _, ok := t.Children[part]; !ok {
			return false
		}
		t = t.Children[part]
	}
	return true
}

// Match try to match the key with a pattern, if
// successful, it will return true and the value
// which was put with the matched pattern. If the
// pattern is't a certain string, this function
// will also return the params matched by this pattern.
func (trie *Trie) Match(v string) (bool, *Result) {
	var result = NewResult()
	parts := trie.splitPattern(v)
	length := len(parts)
	var paramMaps []map[string]string

	t := trie.Root
	for i, part := range parts {
		isMatched := false
		for _, node := range t.Children {
			if ok, params := node.Pattern.Match(part); ok {
				if i == length-1 && !node.isPatternEnding() {
					continue
				}

				t = node
				isMatched = true

				result.ChainData = append(result.ChainData, node.Data)

				if node.Pattern.MatchEverything() {
					for k, v := range params {
						seg := []string{v}
						params[k] = strings.Join(append(seg, parts[i+1:]...), defalutDelimeter)
					}
					paramMaps = append(paramMaps, params)
					goto finish
				}
				paramMaps = append(paramMaps, params)
				break
			}
		}
		if !isMatched {
			return false, nil
		}
	}
finish:
	result.Value = t.Value
	var m = make(map[string]string)
	for _, params := range paramMaps {
		for k, v := range params {
			m[k] = v
		}
	}
	result.Params = m

	return true, result
}

// GetNode allows you use the origin pattern string which was used in Put
// to get the node which points it.
func (trie *Trie) GetNode(v string) (ok bool, node *Node) {
	parts := trie.splitPattern(v)
	t := trie.Root
	for _, part := range parts {
		isMatched := false
		for _, node := range t.Children {
			if ok := node.Pattern.EqualToStr(part); ok {
				t = node
				isMatched = true
				break
			}
		}
		if !isMatched {
			return false, nil
		}
	}

	return true, t
}

// Node is the tree node of the Trie.
type Node struct {
	Pattern       *Pattern
	Value         interface{}
	Data          map[string]interface{}
	Children      map[string]*Node
	patternEnding bool
}

// NewNode returns a new Node object.
func NewNode(str ...string) *Node {
	node := &Node{
		Children: make(map[string]*Node),
		Data:     make(map[string]interface{}),
	}
	if len(str) > 0 {
		node.Pattern = NewPattern(str[0])
		// add regexpPattern
	}
	return node
}

func (node *Node) isPatternEnding() bool {
	return node.patternEnding
}

// Result is return by Match and GetNode.
type Result struct {
	Params    map[string]string
	Value     interface{}
	ChainData []map[string]interface{}
}

// NewResult returns a new Result object.
func NewResult() *Result {
	return &Result{}
}
