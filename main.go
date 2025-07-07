package main

import (
	"fmt"
	"ttm/cli"
	"ttm/commands"
	"ttm/config"
)

var (
	version string
	build   string
)

func main() {
	cli.SetDescription(fmt.Sprintf("Tmux Task Master\nVersion: %s (Build: %s)", version, build))
	cli.AddCommand(cli.Command{
		Name:        "create",
		Action:      commands.Create,
		Example:     "create <task>",
		Description: "Creates a dir and a tmux session for the <task>",
	})
	cli.AddCommand(cli.Command{
		Name:        "add",
		Action:      commands.Add,
		Example:     "add <task>",
		Description: "Creates a worktree for the current repo in the <task> dir",
	})
	cli.AddCommand(cli.Command{
		Name:        "list",
		Action:      commands.List,
		Example:     "list [task]",
		Description: "Without argument - lists all tasks. With argument - lists worktrees of in the <task> dir",
	})
	cli.AddCommand(cli.Command{
		Name:        "delete",
		Action:      commands.Delete,
		Example:     "delete <task>",
		Description: "Deletes the <task> dir with all worktrees and kills the tmux session",
	})

	err := config.Init()
	if err != nil {
		fmt.Println(err)
		cli.Help()
		return
	}

	err = cli.Run()
	if err != nil {
		fmt.Println(err)
		cli.Help()
	}
}
