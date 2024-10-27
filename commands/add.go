package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"ttm/worktree"
)

func Add(args []string) error {
	if len(args) < 2 {
		return missingTaskName
	}

	taskName := args[1]
	return addToTask(taskName)
}

func addToTask(taskName string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	repoName := filepath.Base(cwd)
	wt := worktree.New(taskName, repoName)
	err = wt.Create()
	if err != nil {
		return err
	}

	fmt.Printf("Added '%s' to task '%s'\n", repoName, taskName)

	return nil
}
