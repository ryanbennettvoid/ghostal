package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSanitizeDBURL(t *testing.T) {

	input := "postgresql://admin:pw@localhost/main?sslmode=disable"
	output, err := SanitizeDBURL(input)
	assert.NoError(t, err)
	assert.Equal(t, `postgresql://admin:xxxxx@localhost/main?sslmode=disable`, output)

}
