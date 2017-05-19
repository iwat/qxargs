package internal

import (
	"strings"
	"testing"
)

func TestFind(t *testing.T) {
	finder := newFinder("find", "go")

	match := <-finder.channel()
	if !strings.HasSuffix(match, "finder.go") {
		t.Fatal("expected finder.go, got", match)
	}

	match = <-finder.channel()
	if !strings.HasSuffix(match, "finder_test.go") {
		t.Fatal("expected finder_test.go, got", match)
	}

	finder.reset()
}
