package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGrep(t *testing.T) {
	grepper := newGrepper()

	match, err := grepper.grep("grepper_test.go", "package", "func", "NewGrepper")
	assert.Nil(t, err)
	assert.True(t, match)

	match, err = grepper.grep("grepper_test.go", "package", "func", "-NewGrepper")
	assert.Nil(t, err)
	assert.False(t, match)
}
