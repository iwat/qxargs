package internal

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type Finder struct {
	matchers []_Matcher
	channel  chan string
	reset    chan bool
}

func NewFinder(queries ...string) *Finder {
	matchers := make([]_Matcher, 0, len(queries))
	for _, query := range queries {
		matcher, err := newMatcher(query)
		if err != nil {
			println(err)
			continue
		}
		matchers = append(matchers, matcher)
	}

	finder := &Finder{
		matchers: matchers,
		channel:  make(chan string),
		reset:    make(chan bool),
	}

	go func() {
		_ = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			if finder.shouldSkip(info.Name()) {
				if info.IsDir() {
					return filepath.SkipDir
				} else {
					return nil
				}
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
			case <-finder.reset:
				return errors.New("reset")
			case finder.channel <- rel:
			}

			return nil
		})
		close(finder.channel)
		close(finder.reset)
	}()

	return finder
}

func (f *Finder) Channel() <-chan string {
	return f.channel
}

func (f *Finder) Reset() {
	f.reset <- true
}

func (f *Finder) shouldSkip(basematcher string) bool {
	if basematcher == "." || basematcher == ".." {
		return false
	}

	if strings.HasPrefix(basematcher, ".") {
		return true
	}

	return false
}

func (f *Finder) matches(rel string) bool {
	for _, matcher := range f.matchers {
		if !matcher.Matches(rel) {
			return false
		}
	}

	return true
}
