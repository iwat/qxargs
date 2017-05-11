package internal

import (
	"errors"
	"regexp"
	"strings"
)

type _Matcher interface {
	Matches(input string) bool
}

type _StringMatcher struct {
	pattern string
}

func newMatcher(pattern string) (matcher _Matcher, err error) {
	negate := false

	if strings.HasPrefix(pattern, "-") {
		if len(pattern) == 1 {
			return nil, errors.New("incomplete negate matcher")
		}
		negate = true
		pattern = pattern[1:]
	}

	if len(pattern) > 2 && strings.HasPrefix(pattern, "/") && strings.HasSuffix(pattern, "/") {
		matcher, err = newRegexpMatcher(pattern)
	} else {
		matcher = newStringMatcher(pattern)
	}

	if negate {
		matcher = newNegateMatcher(matcher)
	}

	return matcher, err
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
	pattern = pattern[1 : len(pattern)-1]
	p, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return &_RegexpMatcher{p}, nil
}

func (m _RegexpMatcher) Matches(input string) bool {
	return m.pattern.FindStringSubmatch(input) != nil
}

type _NegateMatcher struct {
	original _Matcher
}

func newNegateMatcher(matcher _Matcher) *_NegateMatcher {
	return &_NegateMatcher{matcher}
}

func (m _NegateMatcher) Matches(input string) bool {
	return !m.original.Matches(input)
}
