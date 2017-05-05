package internal

import (
	"testing"
)

func TestMatcher(t *testing.T) {
	input := "/go/src/github.com/iwat/qxargs/main.go"
	m1, err := newMatcher("main")
	if err != nil {
		t.Fatal(err)
	}

	if !m1.Matches(input) {
		t.Fatal("input should be matched by", m1)
	}
}
