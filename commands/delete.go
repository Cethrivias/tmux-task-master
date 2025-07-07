package commands

import (
	"fmt"
	"os"
	"ttm/task"
	"ttm/tmux"
	"ttm/worktree"
)

func Delete(args []string) error {
	if len(args) < 2 {
		return missingTaskName
	}

	taskName := args[1]
	if len(args) == 2 {
		return deleteTask(taskName)
	}

	worktreeName := args[2]
	return deleteWorktree(taskName, worktreeName)
}

func deleteWorktree(taskName, worktreeName string) error {
	if err := worktree.New(taskName, worktreeName).Delete(); err != nil {
		return err
	}

	return nil
}

func deleteTask(taskName string) error {
	t := task.New(taskName)
	worktreeDirs, err := os.ReadDir(t.Fullpath)
	if err != nil {
		return err
	}
	if len(worktreeDirs) > 0 {
		input := ""
		fmt.Printf("This task contains %d projects. Do you want to delete it? (y/n)\n", len(worktreeDirs))
		if _, err = fmt.Scan(&input); err != nil {
			return err
		}
		if input != "y" {
			fmt.Println("Aborting")
			return nil
		}

		fmt.Printf("Deleting task '%s' projects:\n", taskName)
		for _, dir := range worktreeDirs {
			fmt.Printf(" - %s\n", dir.Name())
			if err := worktree.New(taskName, dir.Name()).Delete(); err != nil {
				return err
			}
		}
	}

	if err := tmux.KillSession(taskName); err != nil {
		return err
	}

	return t.Delete()
}
