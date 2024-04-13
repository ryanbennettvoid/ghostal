package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsValidSnapshotName_ValidAlphanumeric(t *testing.T) {
	input := "ttt123"
	output := IsValidSnapshotName(input)
	assert.True(t, output)
}

func TestIsValidSnapshotName_ValidUnderscore(t *testing.T) {
	input := "t_tt"
	output := IsValidSnapshotName(input)
	assert.True(t, output)
}

func TestIsValidSnapshotName_Invalid(t *testing.T) {
	input := "t-tt"
	output := IsValidSnapshotName(input)
	assert.False(t, output)
}
