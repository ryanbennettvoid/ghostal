package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetURLScheme_Postgres(t *testing.T) {
	input := "postgresql://admin:pw@localhost/main?sslmode=disable"
	output, err := GetURLScheme(input)
	assert.NoError(t, err)
	assert.Equal(t, "postgresql", output)
}

func TestGetURLScheme_Mongo(t *testing.T) {
	input := "mongodb://admin:admin@localhost:27017/main?tls=false"
	output, err := GetURLScheme(input)
	assert.NoError(t, err)
	assert.Equal(t, "mongodb", output)
}

func TestGetURLScheme_Other(t *testing.T) {
	input := "amazing://localhost:1337"
	output, err := GetURLScheme(input)
	assert.NoError(t, err)
	assert.Equal(t, "amazing", output)
}
