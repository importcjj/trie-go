package trie

import (
	"fmt"
	"regexp"
	"strings"
)

// consts
const (
	DefaultPatternDelimeter = ":"
)

var triePattern = regexp.MustCompile(`<(?P<pattern>\w+?:*[\w|\*]+?)>`)

func init() {
	defaultPatternStore.Register("str", `\w+`)
	defaultPatternStore.Register("int", `\d+`)
	defaultPatternStore.Register("*", `.+`)
	defaultPatternStore.DefaultPattern = func() string {
		return `(\w+)`
	}
}

// PatternStore is a pattern register.
type PatternStore struct {
	Patterns       map[string]string
	DefaultPattern func() string
}

// NewPatternStore returns a new PatternStore object.
func NewPatternStore() *PatternStore {
	return &PatternStore{
		Patterns: make(map[string]string),
	}
}

// Register can adds a new customize pattern to the pattern store.
func (store *PatternStore) Register(name string, pattern string) error {
	if _, ok := store.Patterns[name]; ok {
		return ErrDuplicatedPatternName
	}
	store.Patterns[name] = fmt.Sprintf(`(%s)`, pattern)
	return nil
}

// GetPattern returns a pattern which is bound to name.
// If there is none, returns the default pattern. In the
// Default Pattern Store, it just returns then `str` pattern.
func (store *PatternStore) GetPattern(name string) string {
	if pattern, ok := store.Patterns[name]; ok {
		return pattern
	}
	return store.DefaultPattern()
}

var defaultPatternStore = NewPatternStore()

// RegisterPattern allow you add a customize pattern to the defaultPatternStore.
func RegisterPattern(name string, pattern string) error {
	return defaultPatternStore.Register(name, pattern)
}

// Pattern of trie node.
type Pattern struct {
	pattern         *regexp.Regexp
	params          []string
	patternName     string
	patternStr      string
	regexpStr       string
	IsRegexpPattern bool
}

// NewPattern returns a new Pattern object.
func NewPattern(str string) *Pattern {
	var params []string
	var subPatternCount int
	var subPatternName string
	regexpPatternStr := triePattern.ReplaceAllStringFunc(str, func(substr string) string {
		p := strings.Split(strings.Trim(substr, "<>"), DefaultPatternDelimeter)
		param := p[0]
		params = append(params, param)
		subPatternName = ""
		if len(p) > 1 {
			subPatternName = p[1]
		}
		subPatternCount++
		return defaultPatternStore.GetPattern(subPatternName)
	})

	var pattern = regexp.MustCompile(regexpPatternStr)
	var isRegexpPattern = (str != regexpPatternStr)
	p := &Pattern{
		pattern:         pattern,
		params:          params,
		patternStr:      str,
		regexpStr:       regexpPatternStr,
		IsRegexpPattern: isRegexpPattern,
	}
	if subPatternCount == 1 && subPatternName == "*" {
		p.patternName = subPatternName
	}

	return p
}

// Match matches the given string with self's regexpStr, if pattern matched, it
// will return true and the params it found. Otherwise, just returns
// false and nil.
func (pattern *Pattern) Match(str string) (bool, map[string]string) {
	if pattern.IsRegexpPattern {
		matches := pattern.pattern.FindAllStringSubmatch(str, -1)
		if len(matches) == 0 {
			return false, nil
		}
		var patternMap = make(map[string]string)
		for i, param := range pattern.params {
			patternMap[param] = matches[0][i+1]
		}
		return true, patternMap
	}
	return str == pattern.patternStr, nil
}

// EqualToStr returns true if the pattern's regexpStr just
// equal to self's patternStr
func (pattern *Pattern) EqualToStr(str string) bool {
	return str == pattern.patternStr
}

// MatchEverything returns true, if the pattern is the
// `*` pattern.
func (pattern *Pattern) MatchEverything() bool {
	return pattern.patternName == "*"
}
