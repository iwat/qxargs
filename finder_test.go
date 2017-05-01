package main

import (
	"strings"
	"testing"
)

func TestFind(t *testing.T) {
	finder := newFinder(true)
	matches, err := finder.Find("main", "go")
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 1 {
		t.Fatal("expected 1 match, got", len(matches))
	}
	if !strings.HasSuffix(matches[0], "main.go") {
		t.Fatal("expected main.go, got", matches[0])
	}
}
