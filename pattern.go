package trie

import (
	"fmt"
	"regexp"
	"strings"
)

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

type PatternStore struct {
	Patterns       map[string]string
	DefaultPattern func() string
}

func NewPatternStore() *PatternStore {
	return &PatternStore{
		Patterns: make(map[string]string),
	}
}

func (store *PatternStore) Register(name string, pattern string) error {
	if _, ok := store.Patterns[name]; ok {
		return ErrDuplicatedPatternName
	}
	store.Patterns[name] = fmt.Sprintf(`(%s)`, pattern)
	return nil
}

func (store *PatternStore) GetPattern(name string) string {
	if pattern, ok := store.Patterns[name]; ok {
		return pattern
	}
	return store.DefaultPattern()
}

var defaultPatternStore = NewPatternStore()

type Pattern struct {
	pattern         *regexp.Regexp
	params          []string
	patternName     string
	patternStr      string
	regexpStr       string
	IsRegexpPattern bool
}

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

func (pattern *Pattern) EqualStr(str string) bool {
	return str == pattern.patternStr
}

func (pattern *Pattern) MatchEverything() bool {
	return pattern.patternName == "*"
}
