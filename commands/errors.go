package commands

import "errors"

var (
	missingTaskName = errors.New("You need to specify a task name")
)