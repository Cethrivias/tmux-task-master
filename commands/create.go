package commands

import (
	"errors"
	"fmt"
	"ttm/config"
	"ttm/task"
	"ttm/tmux"
)

func Create(args []string) error {
	if len(args) == 1 {
		return missingTaskName
	}

	name := args[1]
	if err := task.New(name).Create(); err != nil {
		return errors.Join(CouldNotCreateTask, err)
	}

	if err := tmux.NewSession(name); err != nil {
		switch {
		case errors.Is(err, tmux.DuplicateSessionError):
			fmt.Printf("Tmux session '%s' already exists\n", name)
		default:
			return errors.Join(CouldNotCreateTask, err)
		}
	}

	if err := tmux.NewWindow(name, "cmd"); err != nil {
		return errors.Join(CouldNotCreateTask, err)
	}

	err := tmux.SetWindowDir(name, "cmd", config.Config.TasksPath+"/"+name)
	if err != nil {
		return errors.Join(CouldNotCreateTask, err)
	}

	fmt.Printf("Created task '%s'\n", name)

	return nil
}
