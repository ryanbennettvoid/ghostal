package app

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApp_ParseProgramArgs(t *testing.T) {

	app := NewApp(GetGlobalLogger())

	input := []string{"someCommand", "option1", "option2"}
	output, err := app.parseProgramArgs(input)
	assert.NoError(t, err)
	assert.Equal(t, "someCommand", string(output.Command))
	{
		o, err := output.Options.Get(0, "")
		assert.NoError(t, err)
		assert.Equal(t, "option1", o)
	}
	{
		o, err := output.Options.Get(1, "")
		assert.NoError(t, err)
		assert.Equal(t, "option2", o)
	}
	{
		_, err := output.Options.Get(2, "xyz")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "xyz")
	}

}
