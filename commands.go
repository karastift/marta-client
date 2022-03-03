package main

import (
	"errors"
	"strings"
)

// use a map so we can check for existing commands by accessing a key
// store number of needed arguments as value
var commands map[string]int = map[string]int{
	"ping": 0,
	"info": 0,
}

const cmdPrefix string = "!"

type Command struct {
	CmdType string
	Args    []string
}

func NewCommand(raw string) (*Command, error) {

	cmd := Command{}

	raw = strings.TrimSpace(raw)

	if !strings.HasPrefix(raw, cmdPrefix) {
		return nil, errors.New("failed to parse command: command prefix is missing")
	}
	raw = strings.TrimPrefix(raw, cmdPrefix)

	splitted := strings.Split(raw, " ")

	cmd.CmdType = splitted[0]
	cmd.Args = splitted[1:]

	neededArguments, cmdExists := commands[cmd.CmdType]

	if !cmdExists {
		return nil, errors.New("failed to parse command: command does not exist")
	}

	if len(cmd.Args) < neededArguments {
		return nil, errors.New("failed to parse command: not enough arguments")
	}

	return &cmd, nil
}
