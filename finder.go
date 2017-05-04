package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type _Finder struct {
	names   []string
	channel chan string
	reset   chan bool
}

func NewFinder(names ...string) *_Finder {
	lower_names := make([]string, 0, len(names))
	for _, name := range names {
		lower_names = append(lower_names, strings.ToLower(name))
	}

	finder := &_Finder{
		names:   lower_names,
		channel: make(chan string),
		reset:   make(chan bool),
	}

	go func() {
		filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
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

func (f *_Finder) Channel() <-chan string {
	return f.channel
}

func (f *_Finder) Reset() {
	f.reset <- true
}

func (f *_Finder) shouldSkip(basename string) bool {
	if basename == "." || basename == ".." {
		return false
	}

	if strings.HasPrefix(basename, ".") {
		return true
	}

	return false
}

func (f *_Finder) matches(rel string) bool {
	lower_rel := strings.ToLower(rel)
	for _, name := range f.names {
		if !strings.Contains(lower_rel, name) {
			return false
		}
	}

	return true
}
