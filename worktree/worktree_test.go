package worktree_test

import (
	"log"
	"os"
	"testing"
	"ttm/config"
	"ttm/utils"
	"ttm/worktree"
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

func TestWorktreeDelete(t *testing.T) {
	taskName := utils.GenerateTestTaskName()
	worktreeName, err := utils.GetWorktreeName()
	if err != nil {
		t.Error(err)
		return
	}
	wt := worktree.New(taskName, worktreeName)

	t.Run("Can delete a worktree from any cwd", func(t *testing.T) {
		if err = wt.Create(); err != nil {
			t.Error(err)
			return
		}
		prevCwd, err := os.Getwd()
		if err != nil {
			t.Error("Could not get cwd", err)
			return
		}
		home, err := os.UserHomeDir()
		if err != nil {
			t.Error("Could not get user HOME", err)
			return
		}
		if err = os.Chdir(home); err != nil {
			t.Error("Could not chdir to user's HOME", err)
			return
		}

		out, err := wt.Delete()
		if err != nil {
			t.Error(out)
			t.Error(err)
		}
		if err = os.Chdir(prevCwd); err != nil {
			t.Error("Could not restore previous working directory", err)
			return
		}
	})

	out, err := utils.TeardownTask(taskName)
	if err != nil {
		t.Error(out)
		t.Error(err)
	}
}
