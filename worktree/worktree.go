package worktree

import (
	"os/exec"
)

type Worktree struct {
	Fullpath string
}

func New(fullpath string) *Worktree {
	return &Worktree{
		Fullpath: fullpath,
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
