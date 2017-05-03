package main

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type _Finder struct {
	debug bool
}

func newFinder(debug bool) *_Finder {
	return &_Finder{debug: debug}
}

func (f *_Finder) Find(names ...string) ([]string, error) {
	args := []string{".", "-type", "f", "-not", "-path", "*/\\.*"}
	for _, name := range names {
		args = append(args, "-ipath", "*"+name+"*")
	}

	if f.debug {
		fmt.Println(append(append([]string{}, "find"), args...))
	}
	cmd := exec.Command("find", args...)

	readCloser, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	defer func() { readCloser.Close() }()

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	matches := make([]string, 0)
	reader := bufio.NewReader(readCloser)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		matches = append(matches, strings.TrimSpace(line))
	}

	if err := cmd.Wait(); err != nil {
		return matches, err
	}

	return matches, nil
}
