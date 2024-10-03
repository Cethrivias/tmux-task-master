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
    fmt.Println("Arguments:", os.Args)

	fmt.Println("Reading config")
	err := readConfig()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Config: %+v\n", config)

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
	taskPath := fmt.Sprintf("%s/%s", config.TasksPath, name)
	err := os.MkdirAll(taskPath, os.ModePerm)
	if err != nil {
		return err
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
