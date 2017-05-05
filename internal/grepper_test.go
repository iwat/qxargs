package internal

import (
	"testing"
)

func TestGrep(t *testing.T) {
	grepper := NewGrepper()
	match, err := grepper.Grep("grepper_test.go", "package", "func", "NewGrepper")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(match)
	if !match {
		t.Fatal("expected main.go to match")
	}
}
