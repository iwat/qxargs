package internal

import (
	"regexp"
	"strings"
)

type _Matcher interface {
	Matches(input string) bool
}

type _StringMatcher struct {
	pattern string
}

func newMatcher(pattern string) (_Matcher, error) {
	if len(pattern) > 2 && strings.HasPrefix(pattern, "/") && strings.HasSuffix(pattern, "/") {
		return newRegexpMatcher(pattern)
	}
	return newStringMatcher(pattern), nil
}

func newStringMatcher(pattern string) _StringMatcher {
	return _StringMatcher{strings.ToLower(pattern)}
}

func (m _StringMatcher) Matches(input string) bool {
	lower := strings.ToLower(input)

	return strings.Contains(lower, m.pattern)
}

type _RegexpMatcher struct {
	pattern *regexp.Regexp
}

func newRegexpMatcher(pattern string) (*_RegexpMatcher, error) {
	p, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return &_RegexpMatcher{p}, nil
}

func (m _RegexpMatcher) Matches(input string) bool {
	return m.pattern.FindStringSubmatch(input) != nil
}
