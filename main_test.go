package main_test

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"ttm/config"
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

func TestCreate(t *testing.T) {
	// Arrange
	taskName := generateTaskName()

	// Act
	cmd := exec.Command("ttm", "create", taskName)
	output, err := cmd.Output()

	// Assert
	if err != nil {
		t.Errorf("Output:\n%s\nError:\n%s\n", string(output), err)
	}
	info, err := os.Stat(getTaskPath(taskName))
	if err != nil {
		t.Error(err)
	}
	if !info.IsDir() {
		t.Error("TaskPath is not a directory")
	}

	output, err = teardownTask(taskName)
	if err != nil {
		t.Errorf("Output:\n%s\nError:\n%s\n", string(output), err)
	}
}

func TestListTasks(t *testing.T) {
	// Arrange
	taskName := generateTaskName()
	err := setupTask(taskName)
	if err != nil {
		t.Error(err)
	}

	// Act
	cmd := exec.Command("ttm", "list")
	output, err := cmd.Output()

	// Assert
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(string(output), taskName) {
		t.Errorf("Output:\n%s\nDoes not contain task: %s\n", string(output), taskName)
	}

	output, err = teardownTask(taskName)
	if err != nil {
		t.Errorf("Output:\n%s\nError:\n%s\n", string(output), err)
	}
}

func TestListWorktrees(t *testing.T) {
	// Arrange
	taskName := generateTaskName()
	err := setupTask(taskName)
	if err != nil {
		t.Error(err)
	}
	worktreeName, err := getWorktreeName()
	if err != nil {
		t.Error(err)
	}
	worktreePath := getWorktreePath(taskName, worktreeName)
	cmd := exec.Command("git", "worktree", "add", "-B", taskName, worktreePath)
	output, err := cmd.Output()
	if err != nil {
		t.Errorf("Output:\n%s\nError:\n%s\n", string(output), err)
	}

	// Act
	output, err = exec.Command("ttm", "list", taskName).Output()

	// Assert
	if err != nil {
		t.Errorf("Output:\n%s\nError:\n%s\n", string(output), err)
	}
    expectedWorktree := fmt.Sprintf("%s (%s)", worktreeName, taskName)
	if !strings.Contains(string(output), expectedWorktree) {
		t.Errorf("Output:\n%s\nDoes not contain worktree: %s\n", string(output), expectedWorktree)
	}

	output, err = teardownTask(taskName)
	if err != nil {
		t.Errorf("Output:\n%s\nError:\n%s\n", string(output), err)
	}
}

func TestAdd(t *testing.T) {
	// arrange
	taskName := generateTaskName()

	// act
	cmd := exec.Command("ttm", "add", taskName)
	output, err := cmd.Output()

	// assert
	if err != nil {
		t.Errorf("Output:\n%s\nError:\n%s\n", string(output), err)
	}
	worktreeName, err := getWorktreeName()
	if err != nil {
		t.Error(err)
	}
	info, err := os.Stat(getWorktreePath(taskName, worktreeName))
	if err != nil {
		t.Error(err)
	}
	if !info.IsDir() {
		t.Error("WorktreePath is not a directory")
	}

	output, err = teardownTask(taskName)
	if err != nil {
		t.Errorf("Output:\n%s\nError:\n%s\n", string(output), err)
	}
}

func TestDeleteWorktree(t *testing.T) {
	// arrange
	taskName := generateTaskName()
	worktreeName, err := getWorktreeName()
	if err != nil {
		t.Error(err)
	}
	worktreePath := getWorktreePath(taskName, worktreeName)
	cmd := exec.Command("git", "worktree", "add", "-B", taskName, worktreePath)
	output, err := cmd.Output()
	if err != nil {
		t.Errorf("Output:\n%s\nError:\n%s\n", string(output), err)
	}

	// act
	cmd = exec.Command("ttm", "delete", taskName, worktreeName)
	output, err = cmd.Output()

	// assert
	if err != nil {
		t.Errorf("Output:\n%s\nError:\n%s\n", string(output), err)
	}
	_, err = os.Stat(worktreePath)
	if err != nil && !os.IsNotExist(err) {
		t.Error(err)
	}

	output, err = teardownTask(taskName)
	if err != nil {
		t.Errorf("Output:\n%s\nError:\n%s\n", string(output), err)
	}
}

func TestDeleteTask(t *testing.T) {
	// arrange
	taskName := generateTaskName()
	worktreeName, err := getWorktreeName()
	if err != nil {
		t.Error(err)
	}
	worktreePath := getWorktreePath(taskName, worktreeName)
	cmd := exec.Command("git", "worktree", "add", "-B", taskName, worktreePath)
	output, err := cmd.Output()
	if err != nil {
		t.Errorf("Output:\n%s\nError:\n%s\n", string(output), err)
	}

	// act
	cmd = exec.Command("ttm", "delete", taskName)
	pp, err := cmd.StdinPipe()
	if err != nil {
		t.Error(err)
	}
	pp.Write([]byte("y\n"))

	output, err = cmd.Output()

	// assert
	if err != nil {
		t.Errorf("Output:\n%s\nError:\n%s\n", string(output), err)
	}
	_, err = os.Stat(worktreePath)
	if err != nil && !os.IsNotExist(err) {
		t.Error(err)
	}
	_, err = os.Stat(getTaskPath(taskName))
	if err != nil && !os.IsNotExist(err) {
		t.Error(err)
	}

	output, err = teardownTask(taskName)
	if err != nil {
		t.Errorf("Output:\n%s\nError:\n%s\n", string(output), err)
	}
}

func generateTaskName() string {
	return "ttm-test-task-" + strconv.Itoa(rand.Int())
}

func getWorktreeName() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	worktreeName := filepath.Base(cwd)

	return worktreeName, err
}

func getTaskPath(taskName string) string {
	return config.Config.TasksPath + "/" + taskName
}

func getWorktreePath(taskName, worktreeName string) string {
	return getTaskPath(taskName) + "/" + worktreeName
}

func setupTask(taskName string) error {
	return os.MkdirAll(getTaskPath(taskName), os.ModePerm)
}

func teardownTask(taskName string) ([]byte, error) {
	taskPath := getTaskPath(taskName)
	worktreeName, err := getWorktreeName()
	if err != nil {
		return []byte{}, err
	}
	worktreePath := getWorktreePath(taskName, worktreeName)

	_, err = os.Stat(worktreePath)
	if err != nil && !os.IsNotExist(err) {
		// some random error
		return []byte{}, err
	}

	if !os.IsNotExist(err) {
		// dir exists. need to cleanup
		output, err := exec.Command("git", "worktree", "remove", worktreePath).Output()
		if err != nil {
			return output, err
		}
	}
	output, err := exec.Command("git", "branch", "-D", taskName).Output()
	if err != nil && (len(output) != 0 || err.Error() != "exit status 1") {
		// empty output + status 1 => branch does not exist
		return output, err
	}

	return []byte{}, os.RemoveAll(taskPath)
}
