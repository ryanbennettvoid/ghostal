package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestToRelativeTime(t *testing.T) {

	now := time.Now()

	{
		pastTime := now.Add(-time.Second * 5)
		output := ToRelativeTime(pastTime, now)
		assert.Equal(t, "5s", output)
	}
	{
		pastTime := now.Add(-time.Hour * 6)
		output := ToRelativeTime(pastTime, now)
		assert.Equal(t, "6h", output)
	}
	{
		pastTime := now.Add(-time.Hour * 24 * 7)
		output := ToRelativeTime(pastTime, now)
		assert.Equal(t, "7d", output)
	}
	{
		pastTime := now.Add(-time.Hour * 24 * 800)
		output := ToRelativeTime(pastTime, now)
		assert.Equal(t, "800d", output)
	}

}
