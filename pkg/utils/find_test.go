package utils

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestUnit_FindString(t *testing.T) {
	result, err := Find([]string{"A", "B", "C"}, func(t string) bool {
		return strings.ToLower(t) == "a"
	})
	assert.NoError(t, err)
	assert.Equal(t, "A", result)
}

func TestUnit_FindStruct(t *testing.T) {
	type Item struct {
		ID   int
		Name string
	}
	items := []Item{
		{1, "XXX"},
		{2, "YYY"},
		{3, "ZZZ"},
	}
	result, err := Find(items, func(item Item) bool {
		return item.ID == 2
	})
	assert.NoError(t, err)
	assert.Equal(t, "YYY", result.Name)
}
