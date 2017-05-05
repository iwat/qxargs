package internal

import (
	"strings"
	"testing"
)

func TestFind(t *testing.T) {
	finder := NewFinder("find", "go")

	match := <-finder.Channel()
	if !strings.HasSuffix(match, "finder.go") {
		t.Fatal("expected finder.go, got", match)
	}

	match = <-finder.Channel()
	if !strings.HasSuffix(match, "finder_test.go") {
		t.Fatal("expected finder_test.go, got", match)
	}

	defer func() {
		if r := recover(); r != nil {
			t.Log("recovered", r)
		}
	}()
	finder.Reset()
}
