package internal

import (
	"testing"
)

func TestGrep(t *testing.T) {
	grepper := newGrepper()
	match, err := grepper.grep("grepper_test.go", "package", "func", "NewGrepper")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(match)
	if !match {
		t.Fatal("expected main.go to match")
	}
}
