package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"os"
	"surfs/internal/block"
)

var SrcRequired = errors.New("must specify a source file")
var DestRequired = errors.New("must specify a destination path")

func Create(c *cli.Context) error {
	args := c.Args()

	src := args.Get(0)
	if src == "" {
		return SrcRequired
	}

	dest := args.Get(1)
	if dest == "" {
		return DestRequired
	}

	f, err := os.Open(src)
	if err != nil {
		logrus.Fatalf("unable to open %s, %v", src, err)
	}

	defer f.Close()

	blocks, err := block.Blocks(f)
	for _, blk := range blocks {
		fmt.Printf("block: len =%d, hash =%s", len(blk.Block), blk.Hash)
	}

	logrus.Info("establishing connection with server...")
	conn, err := grpc.Dial("localhost:7878", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logrus.Fatalf("unable to connect, %v", err)
	}
	defer conn.Close()

	logrus.WithFields(logrus.Fields{
		"src": src,
		"dest": dest,
	}).Info("creating file")

	client := block.NewStoreClient(conn)

	_, err = client.StoreBlock(context.Background(), &block.StoreBlockRequest{
		Block:                []byte("hello"),
		Hash:                 "asdf",
	})
	if err != nil {
		logrus.Fatalf("unable to store block, %v", err)
	}

	return nil
}