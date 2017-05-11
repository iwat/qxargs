package internal

import (
	"bufio"
	"os"
)

type _Grepper struct {
}

func newGrepper() *_Grepper {
	return &_Grepper{}
}

func (g *_Grepper) grep(file string, keywords ...string) (bool, error) {
	matchers := make(map[string]_Matcher)
	negateMatchers := make(map[string]_Matcher)
	for _, keyword := range keywords {
		matcher, err := newMatcher(keyword)
		if err != nil {
			continue
		}
		if _, ok := matcher.(*_NegateMatcher); ok {
			negateMatchers[keyword] = matcher
		} else {
			matchers[keyword] = matcher
		}
	}

	f, err := os.Open(file)
	if err != nil {
		return false, err
	}

	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		for keyword, matcher := range matchers {
			if matcher.Matches(scanner.Text()) {
				delete(matchers, keyword)
			}
		}
		for _, negateMatcher := range negateMatchers {
			if !negateMatcher.Matches(scanner.Text()) {
				return false, nil
			}
		}
		if len(matchers) == 0 && len(negateMatchers) == 0 {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return len(matchers) == 0, nil
}
