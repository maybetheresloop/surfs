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

const MaxAttempts = 10

var SrcRequired = errors.New("must specify a source file")
var DestRequired = errors.New("must specify a destination path")
var ExceededMaxRetries = errors.New("exceeded max retries")

// Create creates a file in the Surfs
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

	// Clean the destination filepath.

	f, err := os.Open(src)
	if err != nil {
		return err
	}

	defer f.Close()

	// Open the file and split it into blocks. This currently reads the whole file into memory.
	blockMap, err := block.MakeBlocks(f)
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

	// Retrieve the current version of the file to be created/updated from the metadata store.
	// For us to be able to update the file, we must send a request specifying the version number
	// to be exactly one more than the version number retrieved from the metadata store.
	log.Debug("Checking current file version...")

	readReq := &meta.ReadFileRequest{Filename: dest}
	readRes, err := metaClient.ReadFile(context.Background(), readReq)
	if err != nil {
		return err
	}

	modReq := &meta.ModifyFileRequest{
		Filename: dest,
		Version:  readRes.Version + 1,
		HashList: blockMap.Hashes,
	}

	attempts := 0

	for {
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
			}).Debug("Successfully created file.")

			return nil
		}

		attempts++

		if attempts >= MaxAttempts {
			log.Errorf("Failed to upload file after %d attempts, aborting.", MaxAttempts)
			return ExceededMaxRetries
		}

		// Otherwise, the ModifyFile() RPC returns a list of hashes whose corresponding
		// blocks are missing from the block store. If the list is nil, then we have the wrong
		// number and leave it to the user to try again.
		if modRes.MissingHashList == nil {
			log.Errorf("Version conflict, please try again.")
			return VersionConflict
		}

		// If the list is not empty, we upload the required blocks to the block store and call
		// ModifyFile again.
		log.Debugf("Block store is missing %d blocks, uploading them...", len(modRes.MissingHashList))

		for _, hash := range modRes.MissingHashList {
			req := &block.StoreBlockRequest{
				Block: blockMap.Blocks[hash],
				Hash:  hash,
			}

			_, err := blockClient.StoreBlock(context.Background(), req)
			if err != nil {
				return err
			}
		}

	}

}
