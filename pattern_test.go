package trie

import (
	"fmt"
	"testing"
)

func HandleError(t *testing.T, err error) {
	t.Error(err)
}

func TestNewPattern(t *testing.T) {
	tests := []struct {
		PatternStr string
		RegexpStr  string
	}{
		{"a.b.c", "a.b.c"},
		{".b.c.<name>", `.b.c.(\w+)`},
		{"a.b.c.<num:int>", `a.b.c.(\d+)`},
		{"a.b.<name:*>", `a.b.(.+)`},
	}

	for _, test := range tests {
		pattern := NewPattern(test.PatternStr)
		if pattern.regexpStr != test.RegexpStr {
			t.Errorf("NewPattern(%s) Parse Error! got %s", test.PatternStr, pattern.regexpStr)
		}
	}
}

func TestRegisterPattern(t *testing.T) {
	r := `^((\\+86)|(86))?(1(([35][0-9])|[8][0-9]|[7][06789]|[4][579]))\\d{8}$`
	err := RegisterPattern("mobile", r)
	HandleError(t, err)

	s := "a.b.<num:mobile>"
	pattern := NewPattern(s)
	if fmt.Sprintf(`a.b.(%s)`, r) != pattern.regexpStr {
		t.Errorf("NewPattern(%s) Parse Error! got %s", s, pattern.regexpStr)
	}

}
