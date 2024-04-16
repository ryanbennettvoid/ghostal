package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnit_StringAsBool_True(t *testing.T) {
	value, err := StringAsBool("true")
	assert.NoError(t, err)
	assert.True(t, value)
}

func TestUnit_StringAsBool_False(t *testing.T) {
	value, err := StringAsBool("false")
	assert.NoError(t, err)
	assert.False(t, value)
}

func TestUnit_StringAsBool_Error(t *testing.T) {
	_, err := StringAsBool("blah")
	assert.Error(t, err)
}
