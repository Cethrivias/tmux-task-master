package main_test

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
	"ttm/config"
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
		cmd := exec.Command("ttm", "create", taskName)
		output, err := cmd.Output()

		// assert
		if err != nil {
			t.Errorf("Output:\n%s\nError:\n%s\n", string(output), err)
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
		out, err := exec.Command("ttm", "add", taskName).CombinedOutput()
		if err != nil {
			t.Errorf("Output:\n%s\nError:\n%s\n", string(out), err)
		}

		// act
		cmd := exec.Command("ttm", "delete", taskName)
		pp, err := cmd.StdinPipe()
		if err != nil {
			t.Error(err)
		}
		pp.Write([]byte("y\n"))
		out, err = cmd.CombinedOutput()

		// assert
		if err != nil {
			t.Errorf("Output:\n%s\nError:\n%s\n", string(out), err)
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
