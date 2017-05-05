package internal

import (
	"strings"
)

type Engine struct {
	grepper *_Grepper
}

func NewEngine() *Engine {
	return &Engine{newGrepper()}
}

func (e *Engine) Query(args ...string) []string {
	var findArgs []string
	var grepArgs []string

	for _, arg := range args {
		if strings.HasPrefix(arg, "?") {
			grepArgs = append(grepArgs, arg[1:])
		} else {
			findArgs = append(findArgs, arg)
		}
	}

	results := []string(nil)
	finder := newFinder(findArgs...)
	for result := range finder.channel() {
		matched, err := e.grepper.grep(result, grepArgs...)
		if err != nil {
			continue
		}
		if matched {
			results = append(results, result)
		}
		if len(results) >= 10 {
			finder.reset()
			break
		}
	}

	return results
}
