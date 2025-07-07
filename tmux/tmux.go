package tmux

import (
	"errors"
	"strings"
	"ttm/config"
	"ttm/utils"
)

// reference

// tmux new-window -c some_directory

//#!/bin/sh
//
// if [ "$#" -ne 1 ]; then
//    printf '%s\n' "exactly 1 argument needed" >&2
//    exit 1
// fi
//
// if [ "$CD_TO_PANE" ]; then
//    cd "$(tmux display -p '#{pane_current_path}')"
// fi
//
// tpth="$(realpath -e "$1//")" || exit 1
//
// tmux list-panes -s -F '#{window_id} #{pane_id} #{pane_current_path}' \
// | while IFS=' ' read -r wndw pne pth; do
//     [ "$pth" = "$tpth" ] && tmux select-pane -t "$pne" && tmux select-window -t "$wndw" && exit 0
// done || tmux new-window -c "$tpth"

var (
	DuplicateSessionError = errors.New("Duplicate session error")
	NewSessionError       = errors.New("Could not create tmux session")
	KillSessionError      = errors.New("Could not kill tmux session")
	ListSessionsError     = errors.New("Could not list tmux sessions")
	ListWindowsError      = errors.New("Could not list tmux windows")
	NewWindowError        = errors.New("Could not create tmux window")
	SetWindowDirError     = errors.New("Could not set tmux window dir")
)

func NewSession(session string) error {
	out, err := utils.Run("tmux", "new-session", "-d", "-s", session, "-n", "vim")
	if err != nil {
		if strings.HasPrefix(err.Error(), "duplicate session") {
			return errors.Join(DuplicateSessionError, errors.New(out), err)
		}
		return errors.Join(NewSessionError, errors.New(out), err)
	}

	err = SetWindowDir(session, "vim", config.Config.TasksPath+"/"+session)
	if err != nil {
		return errors.Join(NewSessionError, err)
	}

	return nil
}

func ListSessions() ([]string, error) {
	out, err := utils.Run("tmux", "list-sessions")
	if err != nil {
		return []string{}, errors.Join(ListSessionsError, err)
	}

	return strings.Split(out, "\n"), nil
}

func ListWindows(session string) ([]string, error) {
	out, err := utils.Run("tmux", "list-windows", "-t", session)

	if err != nil {
		return []string{}, errors.Join(ListWindowsError, err)
	}

	return strings.Split(out, "\n"), nil
}

func KillSession(session string) error {
	out, err := utils.Run("tmux", "kill-session", "-t", session)
	if err != nil {
		return errors.Join(KillSessionError, errors.New(out), err)
	}

	return nil
}

func NewWindow(session, window string) error {
	out, err := utils.Run("tmux", "new-window", "-t", session, "-n", window)
	if err != nil {
		return errors.Join(NewWindowError, errors.New(out), err)
	}

	return nil
}

func SetWindowDir(session, window, dir string) error {
	if err := sendKeys(session, window, "cd "+dir); err != nil {
		return errors.Join(SetWindowDirError, err)
	}

	return nil
}

func sendKeys(session, window, command string) error {
	out, err := utils.Run("tmux", "send-keys", "-t", session+":"+window, command, "C-m")
	if err != nil {
		return errors.Join(errors.New(out), err)
	}

	out, err = utils.Run("tmux", "send-keys", "-t", session+":"+window, "clear", "C-m")
	if err != nil {
		return errors.Join(errors.New(out), err)
	}

	return nil
}
