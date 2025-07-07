package cli

import (
	"errors"
	"fmt"
	"os"
)

var commands = make([]Command, 0)
var description string

var (
	commandNotFoundError = errors.New("You need to specify a command")
	unknownCommandError  = errors.New("Unknown command")
)

func AddCommand(cmd Command) {
	commands = append(commands, cmd)
}

func Help() {
	help := fmt.Sprintf("%s\nAvailable commands:", description)

	for _, cmd := range commands {
		help = fmt.Sprintf("%s\n  %-15s - %s", help, cmd.Example, cmd.Description)
	}

	help = fmt.Sprintf("%s\n  %-15s - %s", help, "--help", "Prints this")

	fmt.Println(help)
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

	if os.Args[1] == "--help" {
		Help()
		return nil
	}

	return unknownCommandError
}

func SetDescription(desc string) { description = desc }

type Command struct {
	Name        string
	Action      func(args []string) error
	Example     string
	Description string
}
