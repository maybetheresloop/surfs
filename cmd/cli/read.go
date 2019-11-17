package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"surfs/internal/block"
	"surfs/internal/meta"

	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

func read(c *cli.Context) error {

	conf, err := getConfig(c)
	if err != nil {
		return err
	}

	src := c.Args().First()
	if src == "" {
		return errors.New("must specify a file to copy")
	}

	dest := c.Args().Get(1)
	if dest == "" {
		return errors.New("must specify a destination")
	}

	// Validate destination path. If the destination path is a directory, the downloaded
	// file will be placed in that directory. Otherwise, the destination path must not already
	// be taken by another file.
	stat, err := os.Stat(dest)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else if stat.IsDir() {
		dest = path.Join(dest, filepath.Base(src))
	}

	// Set up metadata store client.
	addr := fmt.Sprintf("%s:%d", conf.MetadataConf.Host, conf.MetadataConf.Port)
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}

	defer conn.Close()

	client := meta.NewMetadataStoreClient(conn)

	// Set up block store client.
	blockAddr := fmt.Sprintf("%s:%d", conf.BlockConf.Host, conf.BlockConf.Port)
	blockConn, err := grpc.Dial(blockAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}

	defer blockConn.Close()

	blockClient := block.NewStoreClient(blockConn)

	readReq := &meta.ReadFileRequest{
		Filename: src,
	}

	readRes, err := client.ReadFile(context.Background(), readReq)
	if err != nil {
		return err
	}

	if readRes.HashList == nil {
		fmt.Println(NotFound)
		return NotFound
	}

	f, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer f.Close()

	wr := bufio.NewWriter(f)

	// Download all the blocks corresponding to the file and write them to the
	// destination file.
	for _, hash := range readRes.HashList {
		getReq := &block.GetBlockRequest{Hash: hash}
		getRes, err := blockClient.GetBlock(context.Background(), getReq)
		if err != nil {
			return err
		}

		if !getRes.Success {
			return errors.New("block missing from block store")
		}

		_, err = wr.Write(getRes.Block)
		if err != nil {
			return err
		}
	}

	return wr.Flush()
}
