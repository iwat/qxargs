package internal

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type _Finder struct {
	matchers []_Matcher
	walker   chan string
	done     chan bool
}

func newFinder(queries ...string) *_Finder {
	matchers := make([]_Matcher, 0, len(queries))
	for _, query := range queries {
		matcher, err := newMatcher(query)
		if err != nil {
			continue
		}
		matchers = append(matchers, matcher)
	}

	finder := &_Finder{
		matchers: matchers,
		walker:   make(chan string),
		done:     make(chan bool, 1),
	}

	go func() {
		defer close(finder.walker)
		_ = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			if finder.shouldSkip(info.Name()) {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}

			if info.IsDir() {
				return nil
			}

			rel, err := filepath.Rel(".", path)
			if err != nil {
				panic(err)
			}

			if !finder.matches(rel) {
				return nil
			}

			select {
			case <-finder.done:
				return errors.New("done")
			case finder.walker <- rel:
			}

			return nil
		})
	}()

	return finder
}

func (f *_Finder) channel() <-chan string {
	return f.walker
}

func (f *_Finder) reset() {
	f.done <- true
	close(f.done)
}

func (f *_Finder) shouldSkip(basematcher string) bool {
	if basematcher == "." || basematcher == ".." {
		return false
	}

	if strings.HasPrefix(basematcher, ".") {
		return true
	}

	return false
}

func (f *_Finder) matches(rel string) bool {
	for _, matcher := range f.matchers {
		if !matcher.Matches(rel) {
			return false
		}
	}

	return true
}
