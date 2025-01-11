package task

import (
	"errors"
	"os"
	"ttm/config"
)

type Task struct {
	Name     string
	Fullpath string
}

var (
	CreateDirError = errors.New("Could not create task dir")
	DeleteDirError = errors.New("Could not delete task dir")
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
		return errors.Join(CreateDirError, err)
	}

	return nil
}

func (task *Task) Delete() error {
	if err := os.Remove(config.Config.TasksPath + "/" + task.Name); err != nil {
		return errors.Join(DeleteDirError, err)
	}

	return nil
}
