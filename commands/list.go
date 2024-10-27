package commands

import (
	"fmt"
	"os"
	"ttm/config"
	"ttm/worktree"
)

func List(args []string) error {
	if len(args) == 1 {
		return listTasks()
	}
	taskName := args[1]
	return listWorktrees(taskName)
}

func listTasks() error {
	dirs, err := os.ReadDir(config.Config.TasksPath)
	if err != nil {
		return err
	}

	fmt.Println("Tasks:")
	for _, dir := range dirs {
		fmt.Println(" - " + dir.Name())
	}

	return nil
}

func listWorktrees(taskName string) error {
	dirs, err := os.ReadDir(config.Config.TasksPath + "/" + taskName)
	if err != nil {
		return err
	}

	fmt.Printf("Task '%s' projects:\n", taskName)
	for _, dir := range dirs {
		wt := worktree.New(taskName, dir.Name())
		branch, err := wt.Branch()
		if err != nil {
			return err
		}
		fmt.Printf(" - %s (%s)\n", dir.Name(), branch)
	}

	return nil
}
