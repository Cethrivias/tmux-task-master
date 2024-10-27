package cli

import (
	"errors"
	"os"
)

var commands = make([]Command, 0)

var (
	commandNotFoundError = errors.New("You need to specify a command")
	unknownCommandError  = errors.New("Unknown command")
)

func AddCommand(cmd Command) {
	commands = append(commands, cmd)
}

func Run() error {
	if len(os.Args) < 2 {
		return commandNotFoundError
	}

	for _, cmd := range commands {
		if cmd.Name == os.Args[1] {
			return cmd.Action(os.Args[1:])
		}
	}

	return unknownCommandError
}

type Command struct {
	Name   string
	Action func(args []string) error
}
