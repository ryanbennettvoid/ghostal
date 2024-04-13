package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnit_ToPointer(t *testing.T) {
	n := ToPointer(3)
	assert.IsType(t, new(int), n)
}
