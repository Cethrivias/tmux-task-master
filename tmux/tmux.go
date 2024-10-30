package tmux

import (
	"errors"
	"os/exec"
	"strings"
	"ttm/config"
)

var (
	DuplicateSessionError = errors.New("Duplicated tmux sessions")
	UnexpectedError       = errors.New("Unexpected tmux error")
	WindowSetupError      = errors.New("Window setup error")
)

func NewSession(name string) error {
	out, err := exec.Command("tmux", "new-session", "-d", "-s", name, "-n", "vim").CombinedOutput()
	if err != nil {
		return handleError(errors.Join(errors.New(string(out)), err))
	}

	err = sendKeys(name, "vim", "cd "+config.Config.TasksPath+"/"+name)
	if err != nil {
		return errors.Join(WindowSetupError, err)
	}

	return nil
}

func sendKeys(session, window, command string) error {
	out, err := exec.Command("tmux", "send-keys", "-t", session+":"+window, command, "C-m").CombinedOutput()
	if err != nil {
		return errors.Join(errors.New(string(out)), err)
	}

	out, err = exec.Command("tmux", "send-keys", "-t", session+":"+window, "clear", "C-m").CombinedOutput()
	if err != nil {
		return errors.Join(errors.New(string(out)), err)
	}

	return nil
}

func handleError(err error) error {
	if strings.HasPrefix(err.Error(), "duplicate session") {
		return errors.Join(DuplicateSessionError, err)
	}
	return errors.Join(UnexpectedError, err)
}
