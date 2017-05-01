package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]
	finder := newFinder(true)
	results, err := finder.Find(args...)
	if err != nil {
		panic(err)
	}

	fmt.Println(results)
}
