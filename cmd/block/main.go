package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"surfs/internal/block"

	"google.golang.org/grpc"

	"github.com/urfave/cli"

	log "github.com/sirupsen/logrus"
)

// Checks that a given path is a directory.
func validateDataDir(dataDir string) error {
	stat, err := os.Stat(dataDir)
	if err != nil {
		return err
	} else if !stat.IsDir() {
		return errors.New("must specify a valid directory")
	}

	return nil
}

func run(c *cli.Context) error {

	port := c.Uint("port")
	if port == 0 {
		return errors.New("must specify a valid port")
	}

	if c.Bool("V") {
		log.SetLevel(log.DebugLevel)
	}

	if c.Bool("VV") {
		log.SetLevel(log.TraceLevel)
	}

	addr := fmt.Sprintf(":%d", port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	log.Debugf("starting block service on %s", addr)

	dataDir := c.String("datadir")
	if err := validateDataDir(dataDir); err != nil {
		return err
	}

	log.Debugf("using data directory: %s", dataDir)

	store := block.NewStore(dataDir)

	s := grpc.NewServer()
	block.RegisterStoreServer(s, store)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	return nil
}

func main() {

	app := cli.NewApp()

	app.Name = "surfs-block"
	app.Version = "0.1.0"
	app.Usage = "Start the Surfs block service."

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "datadir, D",
			Usage: "Specifies the `DIR` where the block store files are located (default: ./data).",
			Value: "./data",
		},
		cli.UintFlag{
			Name:  "port, p",
			Usage: "Specifies the `PORT` the block service is to listen on (default: 5678)",
			Value: 5678,
		},
		cli.BoolFlag{
			Name:  "V",
			Usage: "Enables verbose output",
		},
		cli.BoolFlag{
			Name:  "VV",
			Usage: "Enables even more verbose output",
		},
	}

	app.Action = run

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("error running service, %v", err)
	}
}
