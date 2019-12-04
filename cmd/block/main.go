package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"surfs/internal/block"

	"github.com/BurntSushi/toml"

	"google.golang.org/grpc"

	"github.com/urfave/cli"

	log "github.com/sirupsen/logrus"
)

type blockStore struct {
	Host    string
	Port    uint
	DataDir string
}

type config struct {
	BlockStore blockStore `toml:"block-store"`
}

func defaultConf() config {
	return config{
		BlockStore: blockStore{
			Host:    "localhost",
			Port:    5678,
			DataDir: "data",
		},
	}
}

func parseConfigFromFile(name string, conf *config) error {
	_, err := toml.DecodeFile(name, &conf)
	return err
}

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

	var conf config
	confPath := c.Args().First()
	if confPath != "" {
		if err := parseConfigFromFile(confPath, &conf); err != nil {
			return err
		}
	}

	if port := c.Uint("port"); port != 0 {
		conf.BlockStore.Port = port
	}

	if dataDir := c.String("datadir"); dataDir != "" {
		if err := validateDataDir(dataDir); err != nil {
			return err
		}

		conf.BlockStore.DataDir = dataDir
	}

	if c.Bool("V") {
		log.SetLevel(log.DebugLevel)
	}

	if c.Bool("VV") {
		log.SetLevel(log.TraceLevel)
	}

	log.Debugf("using data directory: %s", conf.BlockStore.DataDir)

	addr := fmt.Sprintf(":%d", conf.BlockStore.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	log.Debugf("starting block store service on %s", addr)

	store, err := block.NewStore(conf.BlockStore.DataDir)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	block.RegisterStoreServer(s, store)

	if err := s.Serve(lis); err != nil {
		return err
	}

	return nil
}

func main() {

	app := cli.NewApp()

	app.Name = "surfs-block"
	app.Version = "0.1.0"
	app.Usage = "Start the Surfs block store service."

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "datadir, D",
			Usage: "Specifies the `DIR` where the block store files are located",
			Value: "./data",
		},
		cli.UintFlag{
			Name:  "port, p",
			Usage: "Specifies the `PORT` the block store service is to listen on",
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
