package commands

import (
	"errors"
	"fmt"
	"ttm/task"
	"ttm/tmux"
)

func Create(args []string) error {
	if len(args) == 1 {
		return missingTaskName
	}

	name := args[1]
	if err := task.New(name).Create(); err != nil {
		return errors.Join(couldNotCreateTask, err)
	}

	if err := tmux.NewSession(name); err != nil {
		switch {
		case errors.Is(err, tmux.DuplicateSessionError):
			fmt.Printf("Tmux session '%s' already exists\n", name)
		default:
			return errors.Join(couldNotCreateTask, err)
		}
	}

	fmt.Printf("Created task '%s'\n", name)

	return nil
}
