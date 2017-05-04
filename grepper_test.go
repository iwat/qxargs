package main

import (
	"testing"
)

func TestGrep(t *testing.T) {
	grepper := NewGrepper()
	match, err := grepper.Grep("main.go", "package", "func", "newGrepper")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(match)
	if !match {
		t.Fatal("expected main.go to match")
	}
}
