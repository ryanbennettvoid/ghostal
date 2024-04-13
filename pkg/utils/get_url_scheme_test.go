package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnit_GetURLScheme_Postgres(t *testing.T) {
	input := "postgresql://admin:pw@localhost/main?sslmode=disable"
	output, err := GetURLScheme(input)
	assert.NoError(t, err)
	assert.Equal(t, "postgresql", output)
}

func TestUnit_GetURLScheme_Mongo(t *testing.T) {
	input := "mongodb://admin:admin@localhost:27017/main?tls=false"
	output, err := GetURLScheme(input)
	assert.NoError(t, err)
	assert.Equal(t, "mongodb", output)
}

func TestUnit_GetURLScheme_Other(t *testing.T) {
	input := "amazing://localhost:1337"
	output, err := GetURLScheme(input)
	assert.NoError(t, err)
	assert.Equal(t, "amazing", output)
}
