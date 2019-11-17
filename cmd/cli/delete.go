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

func Delete(c *cli.Context) error {

	fp := c.Args().First()

	if fp == "" {
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

	defer conn.Close()

	client := meta.NewMetadataStoreClient(conn)

	readReq := &meta.ReadFileRequest{
		Filename: fp,
	}

	readRes, err := client.ReadFile(context.Background(), readReq)
	if err != nil {
		return err
	}

	if readRes.HashList == nil {
		log.Error("File not found.")
		fmt.Println("Not found")
		return NotFound
	}

	delReq := &meta.DeleteFileRequest{
		Filename: fp,
		Version:  readRes.Version + 1,
	}

	delRes, err := client.DeleteFile(context.Background(), delReq)
	if err != nil {
		return err
	}

	if !delRes.Success {
		log.Error("Version conflict; please try again.")
		return VersionConflict
	}

	log.Debug("Deleted file successfully.")
	fmt.Println("OK")

	return nil
}
