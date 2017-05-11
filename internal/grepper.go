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
	matchers := make(map[string]_Matcher, len(keywords))
	for _, keyword := range keywords {
		matcher, err := newMatcher(keyword)
		if err != nil {
			continue
		}
		matchers[keyword] = matcher
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
		if len(matchers) == 0 {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return len(matchers) == 0, nil
}
