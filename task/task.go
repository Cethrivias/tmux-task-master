package task

import (
	"errors"
	"fmt"
	"os"
	"ttm/config"
)

type Task struct {
	Name     string
	Fullpath string
}

var (
	couldNotCreateTaskDir = errors.New("Could not create task dir")
)

func New(name string) *Task {
	return &Task{
		Name:     name,
		Fullpath: config.Config.TasksPath + "/" + name,
	}
}

func (task *Task) Create() error {
	err := os.MkdirAll(task.Fullpath, os.ModePerm)
	if err != nil {
		return errors.Join(couldNotCreateTaskDir, err)
	}

	fmt.Printf("Created task '%s'\n", task.Name)

	return nil
}

func (task *Task) Delete() error {
	return os.Remove(config.Config.TasksPath + "/" + task.Name)
}
