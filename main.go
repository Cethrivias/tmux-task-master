package main

import (
	"log"
	"ttm/cli"
	"ttm/commands"
	"ttm/config"
)

func main() {
	err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	cli.AddCommand(cli.Command{
		Name:   "create",
		Action: commands.Create,
	})
	cli.AddCommand(cli.Command{
		Name:   "list",
		Action: commands.List,
	})
	cli.AddCommand(cli.Command{
		Name:   "add",
		Action: commands.Add,
	})
	cli.AddCommand(cli.Command{
		Name:   "delete",
		Action: commands.Delete,
	})

	err = cli.Run()
	if err != nil {
		log.Fatal(err)
	}
}
