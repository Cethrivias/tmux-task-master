package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"ttm/config"
	"ttm/worktree"
)

func main() {
	err := config.Init()
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
		if err := listWorktrees(taskName); err != nil {
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
		if len(os.Args) < 3 {
			log.Fatalln("You need to specify a task")
		}
		taskName := os.Args[2]

		if len(os.Args) == 3 {
			if err := deleteTask(taskName); err != nil {
				log.Fatalln(err)
			}
			return
		}

		worktreeName := os.Args[3]

		if err := deleteWorktree(taskName, worktreeName); err != nil {
			log.Fatal(err)
		}

		return
	}

	log.Fatalf("Unknown command '%s'\n", cmd)
}

func createTask(name string) error {
	taskPath := config.Config.TasksPath + "/" + name
	err := os.MkdirAll(taskPath, os.ModePerm)
	if err != nil {
		return err
	}

	fmt.Printf("Created task '%s'\n", name)

	return nil
}

func listTasks() error {
	dirs, err := os.ReadDir(config.Config.TasksPath)
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
	taskPath := config.Config.TasksPath + "/" + name
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	repoName := filepath.Base(cwd)
	worktreePath := taskPath + "/" + repoName
	cmd := exec.Command("git", "worktree", "add", "-B", name, worktreePath)

	output, err := cmd.Output()
	if err != nil {
		fmt.Println(string(output))
		return err
	}

	fmt.Printf("Added '%s' to task '%s'\n", repoName, name)

	return nil
}

func listWorktrees(taskName string) error {
	dirs, err := os.ReadDir(config.Config.TasksPath + "/" + taskName)
	if err != nil {
		return err
	}

	fmt.Printf("Task '%s' projects:\n", taskName)
	for _, dir := range dirs {
        wt := worktree.New(config.Config.TasksPath + "/" + taskName + "/" + dir.Name())
        branch, err := wt.Branch()
        if err != nil {
            return err
        }
		fmt.Printf(" - %s (%s)\n", dir.Name(), branch)
	}

	return nil
}

func deleteWorktree(taskName, worktreeName string) error {
	cmd := exec.Command("git", "worktree", "remove", config.Config.TasksPath+"/"+taskName+"/"+worktreeName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		return err
	}

	return err
}

func deleteTask(taskName string) error {
	worktreeDirs, err := os.ReadDir(config.Config.TasksPath + "/" + taskName)
	if err != nil {
		return err
	}
	if len(worktreeDirs) > 0 {
		input := ""
		fmt.Printf("This task contains %d projects. Do you want to delete it? (y/n)\n", len(worktreeDirs))
		_, err = fmt.Scan(&input)
		if err != nil {
			return err
		}
		if input != "y" {
			fmt.Println("Aborting")
			return nil
		}

		fmt.Printf("Deleting task '%s' projects:\n", taskName)
		for _, dir := range worktreeDirs {
			fmt.Printf(" - %s\n", dir.Name())
			cmd := exec.Command("git", "worktree", "remove", config.Config.TasksPath+"/"+taskName+"/"+dir.Name())
			output, err := cmd.Output()
			if err != nil {
				fmt.Println(string(output))
				return err
			}
		}
	}

	return os.Remove(config.Config.TasksPath + "/" + taskName)
}
