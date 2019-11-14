package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"surfs/internal/meta"

	"google.golang.org/grpc"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func run(c *cli.Context) error {
	port := c.Uint64("port")
	if port == 0 {
		return errors.New("must specify a valid port")
	}

	blockStoreHost := c.String("block-store-hostname")
	if blockStoreHost == "" {
		return errors.New("must specify a hostname for the block store")
	}

	blockStorePort := c.Uint64("block-store-port")
	if blockStorePort == 0 {
		return errors.New("must specify a port for the block store")
	}

	addr := fmt.Sprintf(":%d", port)
	blockStoreAddr := fmt.Sprintf("%s:%d", blockStoreHost, blockStorePort)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	store, err := meta.NewStore(blockStoreAddr)
	if err != nil {
		return err
	}
	defer store.Close()

	s := grpc.NewServer()
	meta.RegisterMetadataStoreServer(s, store)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return err
	}

	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "surfs-meta"
	app.Author = "maybetheresloop"
	app.Email = "maybetheresloop@gmail.com"
	app.Usage = "Start the Surfs metadata service."
	app.Version = "0.1.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "datadir, D",
			Usage: "Specifies the `DIR` where the metadata store files are located (default: ./meta).",
			Value: "./meta",
		},
		cli.UintFlag{
			Name:  "port, p",
			Usage: "Specifies the `PORT` the to listen on (default: 5679)",
			Value: 5679,
		},
		cli.StringFlag{
			Name:  "block-store-hostname, H",
			Usage: "Specifies the `HOSTNAME` of the Surfs block store service (default: localhost).",
			Value: "localhost",
		},
		cli.UintFlag{
			Name:  "block-store-port, P",
			Usage: "Specifies the `PORT` of the Surfs block store service (default: 5678).",
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
		log.Fatalf("error running metadata service, %v", err)
	}

}
