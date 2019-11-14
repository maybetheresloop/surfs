package main

import (
	"context"
	"errors"
	"os"
	"surfs/internal/block"
	"surfs/internal/meta"

	"google.golang.org/grpc"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
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
		log.Fatalf("unable to open %s, %v", src, err)
	}

	defer f.Close()

	// Open the file and split it into blocks. This currently reads the whole file into memory.
	_, hashes, err := block.MakeBlocks(f)
	if err != nil {
		return nil
	}

	// Create a client to interact with the block store service.
	conn, err := grpc.Dial("localhost:5678", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}

	defer conn.Close()

	client := meta.NewMetadataStoreClient(conn)

	log.Debug("calling ReadFile() RPC...")

	readReq := &meta.ReadFileRequest{Filename: dest}
	readRes, err := client.ReadFile(context.Background(), readReq)
	if err != nil {
		return err
	}

	modReq := &meta.ModifyFileRequest{
		Filename: dest,
		Version:  readRes.Version + 1,
		HashList: hashes,
	}

	modRes, err := client.ModifyFile(context.Background(), modReq)
	if err != nil {
		return err
	}

	//log.Info("establishing connection with server...")
	//conn, err := grpc.Dial("localhost:5678", grpc.WithInsecure(), grpc.WithBlock())
	//if err != nil {
	//	return err
	//}
	//defer conn.Close()
	//client := block.NewStoreClient(conn)
	//
	//// Send requests to store all the blocks of the file.
	//log.WithFields(log.Fields{
	//	"src":  src,
	//	"dest": dest,
	//}).Info("creating file")
	//
	//for _, blk := range blks {
	//	if _, err := client.StoreBlock(context.Background(), &block.StoreBlockRequest{
	//		Block: blk.Block,
	//		Hash:  blk.Hash,
	//	}); err != nil {
	//		return err
	//	}
	//}
	//
	//log.Info("successfully stored block")

	return nil
}
