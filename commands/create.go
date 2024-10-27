package commands

import (
	"errors"
	"ttm/task"
)

func Create(args []string) error {
	if len(args) == 1 {
		return errors.New("You need to specify task name")
	}

	name := args[1]
    return task.New(name).Create()
}

