package app

import "errors"

type ProgramArgs struct {
	Command Command
	Options string
}

func ParseProgramArgs(args []string) (ProgramArgs, error) {
	if len(args) == 0 || len(args[0]) == 0 {
		return ProgramArgs{}, errors.New("invalid args")
	}
	command := Command(args[0])
	options := ""
	if len(args) > 1 {
		options = args[1]
	}
	return ProgramArgs{
		Command: command,
		Options: options,
	}, nil
}
