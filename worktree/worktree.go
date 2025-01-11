package worktree

import (
	"errors"
	"os/exec"
	"ttm/config"
)

type Worktree struct {
	TaskName     string
	WorktreeName string
	Fullpath     string
}

var (
	GetGitBranchError   = errors.New("Could not get Git branch")
	CreateWorktreeError = errors.New("Could not create a worktree")
	DeleteWorktreeError = errors.New("Could not delete a worktree")
)

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
		return "", errors.Join(GetGitBranchError, errors.New(string(output)), err)
	}

	return string(output[:len(output)-1]), nil
}

func (w *Worktree) Create() error {
	cmd := exec.Command("git", "worktree", "add", "-B", w.TaskName, w.Fullpath)

	if output, err := cmd.CombinedOutput(); err != nil {
		return errors.Join(CreateWorktreeError, errors.New(string(output)), err)
	}

	return nil
}

func (w *Worktree) Delete() error {
	cmd := exec.Command("git", "worktree", "remove", w.Fullpath)
	cmd.Dir = w.Fullpath
	if output, err := cmd.CombinedOutput(); err != nil {
		return errors.Join(DeleteWorktreeError, errors.New(string(output)), err)
	}

	return nil
}
