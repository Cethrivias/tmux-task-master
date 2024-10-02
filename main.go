package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var home = os.Getenv("HOME")
var tasksPath = home + "/Projects/tasks"

func main() {
	fmt.Println("Arguments:", os.Args)

	cmd := os.Args[1]

	if cmd == "create" {
		if len(os.Args) < 3 {
			log.Fatalf("You need to specify task name")
		}
		name := os.Args[2]

		fmt.Printf("Creating task '%s'\n", name)

		if err := createTask(name); err != nil {
			log.Fatal(err)
		}

		return
	}

	if cmd == "add" {
		if len(os.Args) < 3 {
			log.Fatalf("You need to specify task name")
		}
		name := os.Args[2]

		if err := addToTask(name); err != nil {
			log.Fatal(err)
		}

		return
	}

	log.Fatalf("Unknown command '%s'\n", cmd)
}

func createTask(name string) error {
	taskPath := fmt.Sprintf("%s/%s", tasksPath, name)
	err := os.MkdirAll(taskPath, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func addToTask(name string) error {
	taskPath := fmt.Sprintf("%s/%s", tasksPath, name)
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
    repoName := filepath.Base(cwd)
	worktreePath := fmt.Sprintf("%s/%s", taskPath, repoName)
	cmd := exec.Command("git", "worktree", "add", "-B", name, worktreePath)

	if err = cmd.Run(); err != nil {
		return err
	}

	return nil
}
