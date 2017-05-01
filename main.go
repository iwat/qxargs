package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	findArgs := make([]string, 0)
	grepArgs := make([]string, 0)

	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "?") {
			grepArgs = append(grepArgs, arg[1:])
		} else {
			findArgs = append(findArgs, arg)
		}
	}

	finder := newFinder(true)
	results, err := finder.Find(findArgs...)
	if err != nil {
		panic(err)
	}

	grepper := newGrepper()
	results, err = grepper.Grep(results, grepArgs...)
	if err != nil {
		panic(err)
	}

	fmt.Println(results)
}
