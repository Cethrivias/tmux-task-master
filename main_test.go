package main_test

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"slices"
	"strings"
	"testing"
	"ttm/config"
	"ttm/tmux"
	"ttm/utils"
)

func TestMain(m *testing.M) {
	// Global setup
	err := config.Init()
	if err != nil {
		log.Fatalln(err)
	}

	// Run all tests
	exitCode := m.Run()

	// Global teardown
	// teardown()

	// Exit with the test exit code
	os.Exit(exitCode)
}

func TestHappyPath(t *testing.T) {
	taskName := utils.GenerateTestTaskName()
	worktreeName, err := utils.GetWorktreeName()
	if err != nil {
		t.Error(err)
	}
	worktreePath := utils.GetWorktreePath(taskName, worktreeName)

	t.Run("Creates task", func(t *testing.T) {
		// act
		out, err := utils.TtmCreate(taskName)

		// assert
		if err != nil {
			t.Errorf("Output:\n%s\nError:\n%s\n", out, err)
		}
		info, err := os.Stat(utils.GetTaskPath(taskName))
		if err != nil {
			t.Error(err)
		}
		if !info.IsDir() {
			t.Error("TaskPath is not a directory")
		}
	})

	t.Run("Lists tasks", func(t *testing.T) {
		// act
		cmd := exec.Command("ttm", "list")
		output, err := cmd.Output()

		// assert
		if err != nil {
			t.Error(err)
		}
		if !strings.Contains(string(output), taskName) {
			t.Errorf("Output:\n%s\nDoes not contain task: %s\n", string(output), taskName)
		}
	})

	t.Run("Adds worktree", func(t *testing.T) {
		// act
		cmd := exec.Command("ttm", "add", taskName)
		output, err := cmd.Output()

		// assert
		if err != nil {
			t.Errorf("Output:\n%s\nError:\n%s\n", string(output), err)
		}
		info, err := os.Stat(utils.GetWorktreePath(taskName, worktreeName))
		if err != nil {
			t.Error(err)
		}
		if !info.IsDir() {
			t.Error("WorktreePath is not a directory")
		}
	})

	t.Run("Lists worktrees", func(t *testing.T) {
		// act
		output, err := exec.Command("ttm", "list", taskName).Output()

		// Assert
		if err != nil {
			t.Errorf("Output:\n%s\nError:\n%s\n", string(output), err)
		}
		expectedWorktree := fmt.Sprintf("%s (%s)", worktreeName, taskName)
		if !strings.Contains(string(output), expectedWorktree) {
			t.Errorf(
				"Output:\n%s\nDoes not contain worktree: %s\n",
				string(output),
				expectedWorktree,
			)
		}
	})

	t.Run("Deletes worktree", func(t *testing.T) {
		// act
		cmd := exec.Command("ttm", "delete", taskName, worktreeName)
		output, err := cmd.Output()

		// assert
		if err != nil {
			t.Errorf("Output:\n%s\nError:\n%s\n", string(output), err)
		}
		_, err = os.Stat(worktreePath)
		if err != nil && !os.IsNotExist(err) {
			t.Error(err)
		}
	})

	t.Run("Deletes task", func(t *testing.T) {
		// arrange
		out, err := utils.TtmAdd(taskName)
		if err != nil {
			t.Errorf("Output:\n%s\nError:\n%s\n", out, err)
		}

		// act
		cmd := exec.Command("ttm", "delete", taskName)
		pp, err := cmd.StdinPipe()
		if err != nil {
			t.Error(err)
		}
		pp.Write([]byte("y\n"))
		out, err = utils.TtmDelete(taskName)

		// assert
		if err != nil {
			t.Errorf("Output:\n%s\nError:\n%s\n", out, err)
		}
		_, err = os.Stat(worktreePath)
		if err != nil && !os.IsNotExist(err) {
			t.Error(err)
		}
		_, err = os.Stat(utils.GetTaskPath(taskName))
		if err != nil && !os.IsNotExist(err) {
			t.Error(err)
		}
	})

	out, err := utils.TeardownTask(taskName)
	if err != nil {
		t.Errorf("Output:\n%s\nError:\n%s\n", string(out), err)
	}
}

func TestTmux(t *testing.T) {
	taskName := utils.GenerateTestTaskName()

	t.Run("Creates session", func(t *testing.T) {
		// act
		out, err := utils.TtmCreate(taskName)
		if err != nil {
			t.Errorf("Output:\n%s\nError:\n%s\n", out, err)
		}

		// assert
		sessions, err := tmux.ListSessions()
		if err != nil {
			t.Errorf("Could not list tmux sessions. Error:\n%s\n", err)
		}

		hasSession := slices.ContainsFunc(
			sessions,
			func(session string) bool { return strings.HasPrefix(session, taskName) },
		)

		if !hasSession {
			t.Errorf("Tmux session '%s' was not created\n", taskName)
		}
	})

	t.Run("Creates windows", func(t *testing.T) {
		wins, err := tmux.ListWindows(taskName)
		if err != nil {
			t.Error(err)
		}
		if !slices.ContainsFunc(wins, func(win string) bool { return strings.Contains(win, "vim") }) {
			t.Error("Window 'vim' was not created")
		}
		if !slices.ContainsFunc(wins, func(win string) bool { return strings.Contains(win, "cmd") }) {
			t.Error("Window 'cmd' was not created")
		}
	})

	t.Run("Kills session", func(t *testing.T) {
		out, err := utils.TtmDelete(taskName)
		if err != nil {
			t.Errorf("Output:\n%s\nError:\n%s\n", out, err)
		}

		sessions, err := tmux.ListSessions()

		hasSession := slices.ContainsFunc(
			sessions,
			func(session string) bool { return strings.HasPrefix(session, taskName) },
		)

		if hasSession {
			t.Errorf("Tmux session '%s' was not killed\n", taskName)
		}
	})

	out, err := utils.TeardownTask(taskName)
	if err != nil {
		t.Errorf("Output:\n%s\nError:\n%s\n", string(out), err)
	}
}
