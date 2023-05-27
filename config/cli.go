package config

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func InitCommands(db Database) {
	InitialDB()
	cmdApp := cli.NewApp()
	cmdApp.Commands = []*cli.Command{
		{
			Name: "db:migrate",
			Action: func(c *cli.Context) error {
				Migration(db)
				return nil
			},
		},
	}

	err := cmdApp.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}
