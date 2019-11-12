package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Author = "maybetheresloop"
	app.Email = "maybetheresloop@gmail.com"
	app.Name = "surfs-cli"
	app.Usage = "The Surfs command-line interface."
	app.Version = "0.1.0"

	app.Commands = []cli.Command{
		{
			Name:   "create",
			Usage:  "Upload a file to the store",
			Action: Create,
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatalf("error, %v", err)
	}

}
