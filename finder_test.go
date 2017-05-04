package main

import (
	"strings"
	"testing"
)

func TestFind(t *testing.T) {
	finder := NewFinder("main", "go")

	match := <-finder.Channel()
	if !strings.HasSuffix(match, "main.go") {
		t.Fatal("expected main.go, got", match)
	}

	defer func() {
		if r := recover(); r != nil {
			t.Log("recovered", r)
		}
	}()
	finder.Reset()
}
