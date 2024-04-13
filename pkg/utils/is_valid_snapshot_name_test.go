package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnit_IsValidSnapshotName_ValidAlphanumeric(t *testing.T) {
	input := "ttt123"
	output := IsValidSnapshotName(input)
	assert.True(t, output)
}

func TestUnit_IsValidSnapshotName_ValidUnderscore(t *testing.T) {
	input := "t_tt"
	output := IsValidSnapshotName(input)
	assert.True(t, output)
}

func TestUnit_IsValidSnapshotName_Invalid(t *testing.T) {
	input := "t-tt"
	output := IsValidSnapshotName(input)
	assert.False(t, output)
}
