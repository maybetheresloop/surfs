package main

import (
	"context"
	"errors"
	"fmt"
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

	log.SetLevel(log.DebugLevel)

	conf, err := getConfig(c)
	if err != nil {
		return err
	}

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
		return err
	}

	defer f.Close()

	// Open the file and split it into blocks. This currently reads the whole file into memory.
	blocks, err := block.MakeBlocks(f)
	if err != nil {
		return err
	}

	// Create a client to interact with the metadata store.
	metaAddr := fmt.Sprintf("%s:%d", conf.MetadataConf.Host, conf.MetadataConf.Port)
	metaConn, err := grpc.Dial(metaAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}

	defer metaConn.Close()

	metaClient := meta.NewMetadataStoreClient(metaConn)

	// Create a client to interact with the block store.
	blockAddr := fmt.Sprintf("%s:%d", conf.BlockConf.Host, conf.BlockConf.Port)
	blockConn, err := grpc.Dial(blockAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}

	defer blockConn.Close()

	blockClient := block.NewStoreClient(blockConn)

	log.Debug("calling ReadFile() RPC...")

	readReq := &meta.ReadFileRequest{Filename: dest}
	readRes, err := metaClient.ReadFile(context.Background(), readReq)
	if err != nil {
		return err
	}

	hashes := make([]string, 0, len(blocks))
	for k, _ := range blocks {
		hashes = append(hashes, k)
	}

	modReq := &meta.ModifyFileRequest{
		Filename: dest,
		Version:  readRes.Version + 1,
		HashList: hashes,
	}

	// Call the metadata store's ModifyFile() RPC to request an update of the hash list.
	// If the ModifyFile RPC() returns with success, then we are done.
	modRes, err := metaClient.ModifyFile(context.Background(), modReq)
	if err != nil {
		return err
	}

	if modRes.Success {
		log.WithFields(log.Fields{
			"src":  src,
			"dest": dest,
		}).Debug("successfully created file")

		return nil
	}

	// Otherwise, the ModifyFile() RPC returns a list of hashes whose corresponding
	// blocks are missing from the block store. We upload those blocks to the block
	// store and call ModifyFile again.
	log.Debugf("block store is missing %d blocks, uploading them", len(modRes.MissingHashList))

	for _, hash := range modRes.MissingHashList {
		req := &block.StoreBlockRequest{
			Block: blocks[hash],
			Hash:  hash,
		}

		_, err := blockClient.StoreBlock(context.Background(), req)
		if err != nil {
			return err
		}
	}

	// All blocks have been stored, try again.
	//modRes, err = metaClient.ModifyFile(context.Background(), modReq)
	//if err != nil {
	//	return err
	//}
	//
	//if modRes.Success {
	//	log.WithFields(log.Fields{
	//		"src":  src,
	//		"dest": dest,
	//	}).Debug("successfully created file")
	//
	//	return nil
	//}

	return nil
}
