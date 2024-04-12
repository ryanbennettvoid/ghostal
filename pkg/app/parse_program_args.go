package app

import (
	"errors"
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

func ParseProgramArgs(args []string) (ProgramArgs, error) {
	if len(args) == 0 || len(args[0]) == 0 {
		return ProgramArgs{}, errors.New("invalid args")
	}
	command := Command(args[0])
	options := make([]string, 0)
	if len(args) > 1 {
		options = args[1:]
	}
	return ProgramArgs{
		Command: command,
		Options: options,
	}, nil
}
