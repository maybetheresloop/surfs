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

	app.Commands = []cli.Command{
		{
			Name:   "create",
			Flags:  flags,
			Usage:  "Upload a file to the store",
			Action: Create,
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatalf("error, %v", err)
	}
}
