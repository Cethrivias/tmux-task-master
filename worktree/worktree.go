package worktree

import (
	"fmt"
	"os/exec"
	"ttm/config"
)

type Worktree struct {
	TaskName     string
	WorktreeName string
	Fullpath     string
}

func New(taskName, worktreeName string) *Worktree {
	return &Worktree{
		TaskName:     taskName,
		WorktreeName: worktreeName,
		Fullpath:     config.Config.TasksPath + "/" + taskName + "/" + worktreeName,
	}
}

func (w *Worktree) Branch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	cmd.Dir = w.Fullpath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(output[:len(output)-1]), nil
}

func (w *Worktree) Create() error {
	cmd := exec.Command("git", "worktree", "add", "-B", w.TaskName, w.Fullpath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		return err
	}

	return nil
}
