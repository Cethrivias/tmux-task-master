package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var home = os.Getenv("HOME")
var configPath = home + "/.config/ttm"
var config = TtmConfig{
	TasksPath: home + "/ttm",
}

func main() {
	err := readConfig()
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) < 2 {
		log.Fatalln("You need to specify a command")
	}

	cmd := os.Args[1]

	if cmd == "create" {
		if len(os.Args) < 3 {
			log.Fatalln("You need to specify task name")
		}
		name := os.Args[2]

		if err := createTask(name); err != nil {
			log.Fatal(err)
		}

		return
	}

	if cmd == "list" {
		if len(os.Args) == 2 {
			if err := listTasks(); err != nil {
				log.Fatal(err)
			}

			return
		}
		taskName := os.Args[2]
		if err := listProjects(taskName); err != nil {
			log.Fatal(err)
		}

		return
	}

	if cmd == "add" {
		if len(os.Args) < 3 {
			log.Fatalln("You need to specify task name")
		}
		name := os.Args[2]

		if err := addToTask(name); err != nil {
			log.Fatal(err)
		}

		return
	}

    if cmd == "delete" {
        if len(os.Args) < 4 {
            log.Fatalln("You need to specify task and project name")
        }
        taskName := os.Args[2]
        projectName := os.Args[3]

		if err := deleteProjectWorktree(taskName, projectName); err != nil {
			log.Fatal(err)
		}

		return
    }

	log.Fatalf("Unknown command '%s'\n", cmd)
}

func createTask(name string) error {
	taskPath := config.TasksPath + "/" + name
	err := os.MkdirAll(taskPath, os.ModePerm)
	if err != nil {
		return err
	}

	fmt.Printf("Created task '%s'\n", name)

	return nil
}

func listTasks() error {
	dirs, err := os.ReadDir(config.TasksPath)
	if err != nil {
		return err
	}

	fmt.Println("Tasks:")
	for _, dir := range dirs {
		fmt.Println(" - " + dir.Name())
	}

	return nil
}

func addToTask(name string) error {
	taskPath := config.TasksPath + "/" + name
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	repoName := filepath.Base(cwd)
	worktreePath := taskPath + "/" + repoName
	cmd := exec.Command("git", "worktree", "add", "-B", name, worktreePath)

	if err = cmd.Run(); err != nil {
		return err
	}

	fmt.Printf("Added '%s' to task '%s'\n", repoName, name)

	return nil
}

func listProjects(taskName string) error {
	dirs, err := os.ReadDir(config.TasksPath + "/" + taskName)
	if err != nil {
		return err
	}

	fmt.Printf("Task '%s' projects:\n", taskName)
	for _, dir := range dirs {
		fmt.Println(" - " + dir.Name())
	}

	return nil
}

func deleteProjectWorktree(taskName, projectName string) error {
    return nil
}

func readConfig() error {
	file, err := os.Open(configPath + "/ttm.json")
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(configPath, os.ModePerm)
			if err != nil {
				return err
			}
			file, err = os.Create(configPath + "/ttm.json")
			if err != nil {
				return err
			}
			encoder := json.NewEncoder(file)
			encoder.SetIndent("", "    ")
			err = encoder.Encode(&config)
			if err != nil {
				return err
			}
		}
	}

	decoder := json.NewDecoder(file)
	decoder.Decode(&config)

	return nil
}

type TtmConfig struct {
	TasksPath string `json:"tasksPath"`
}
