package main

import (
	"context"
	"encoding/base64"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"io"
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

	blocks := make(map[string][]byte)

	for {
		block := make([]byte, 4096)
		n, err := f.Read(block)
		if err == io.EOF {
			hash := base64.StdEncoding.EncodeToString(block[:n])
			blocks[hash] = block
			break
		} else if err != nil {
			logrus.Fatalf("unable to read file, %v", err)
		}

		hash := base64.StdEncoding.EncodeToString(block[:n])
		blocks[hash] = block
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