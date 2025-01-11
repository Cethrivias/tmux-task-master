package utils

import (
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"ttm/config"
)

func GenerateTestTaskName() string {
	return "ttm-test-task-" + strconv.Itoa(rand.Int())
}

func GetWorktreeName() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	worktreeName := filepath.Base(cwd)

	return worktreeName, err
}

func GetTaskPath(taskName string) string {
	return config.Config.TasksPath + "/" + taskName
}

func GetWorktreePath(taskName, worktreeName string) string {
	return GetTaskPath(taskName) + "/" + worktreeName
}

func SetupTask(taskName string) error {
	return os.MkdirAll(GetTaskPath(taskName), os.ModePerm)
}

func TeardownTask(taskName string) ([]byte, error) {
	taskPath := GetTaskPath(taskName)
	worktreeName, err := GetWorktreeName()
	if err != nil {
		return []byte{}, err
	}
	worktreePath := GetWorktreePath(taskName, worktreeName)

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
