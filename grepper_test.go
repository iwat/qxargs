package main

import (
	"strings"
	"testing"
)

func TestGrep(t *testing.T) {
	grepper := newGrepper()
	results, err := grepper.Grep([]string{"main.go", "finder.go"}, "package", "func", "newGrepper")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(results)
	if len(results) != 1 {
		t.Fatal("expected 1 match, got", len(results))
	}
	if !strings.HasSuffix(results[0], "main.go") {
		t.Fatal("expected main.go, got", results[0])
	}
}
