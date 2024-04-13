package app

import (
	"fmt"
)

type Options []string

func (o Options) Get(idx int, description string) (string, error) {
	if idx > len(o)-1 {
		return "", fmt.Errorf("option at position %d (%s) not found", idx, description)
	}
	return o[idx], nil
}

type ProgramArgs struct {
	Command Command
	Options Options
}
