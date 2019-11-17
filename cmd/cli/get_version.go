package main

import (
	"context"
	"errors"
	"fmt"
	"surfs/internal/meta"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

func GetVersion(c *cli.Context) error {

	fp := c.Args().First()
	if fp == "" {
		log.Errorf("Filepath must be specified.")
		return errors.New("must specify a file")
	}

	conf, err := getConfig(c)
	if err != nil {
		return err
	}

	addr := fmt.Sprintf("%s:%d", conf.MetadataConf.Host, conf.MetadataConf.Port)
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}

	client := meta.NewMetadataStoreClient(conn)

	req := &meta.GetVersionRequest{Filename: fp}
	res, err := client.GetVersion(context.Background(), req)
	if err != nil {
		return err
	}

	fmt.Println(res.Version)
	return nil
}
