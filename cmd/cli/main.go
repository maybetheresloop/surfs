package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Author = "maybetheresloop"
	app.Email = "maybetheresloop@gmail.com"
	app.Name = "surfs-cli"
	app.Usage = "The Surfs command-line interface."
	app.Version = "0.1.0"

	flags := []cli.Flag{
		cli.StringFlag{
			Name:  "block-store-hostname",
			Usage: "Specifies the `HOSTNAME` of the Surfs block store service (default: localhost).",
		},
		cli.UintFlag{
			Name:  "block-store-port",
			Usage: "Specifies the `PORT` of the Surfs block store service (default: 5678).",
		},
		cli.StringFlag{
			Name:  "metadata-store-hostname",
			Usage: "Specifies the `HOSTNAME` of the Surfs metadata store service (default: localhost).",
		},
		cli.UintFlag{
			Name:  "metadata-store-port",
			Usage: "Specifies the `PORT` of the Surfs block store service (default: 5679).",
		},
		cli.StringFlag{
			Name:      "config, c",
			Usage:     "Specifies a configuration `FILE`",
			TakesFile: true,
			Value:     "conf/keychain.toml",
		},
	}

	app.Flags = flags

	app.Commands = []cli.Command{
		{
			Name:   "create",
			Usage:  "Upload a file to the store",
			Action: Create,
		},
		{
			Name:   "get-version",
			Usage:  "Get the current version of a file.",
			Action: GetVersion,
		},
		{
			Name:   "delete",
			Usage:  "Delete a file from the store.",
			Action: Delete,
		},
		{
			Name:   "read",
			Usage:  "Retrieve a file from Surfs.",
			Action: read,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("error, %v", err)
	}
}
