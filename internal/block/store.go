package block

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
)

const FilePrefix = "blk_"

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Store struct {
	// Mapping of filename to file blocks.
	blocks map[string]datafile

	// Key-value engine for storing file metadata.
	engine engine

	// The data directory for this instance of the metadata store.
	dataDir string

	// Assigns sequential IDs to blocks in the block store.
	counter uint64
}

func NewStore(dataDir string) (*Store, error) {

	engine, err := openKeychainEngine("block.keychain")
	if err != nil {
		return nil, err
	}

	return &Store{
		blocks:  make(map[string]datafile),
		engine:  engine,
		dataDir: dataDir,
		counter: 0,
	}, nil
}

func (s *Store) StoreBlock(ctx context.Context, req *StoreBlockRequest) (*StoreBlockResponse, error) {

	// This could totally cause collisions, need to fix this later.
	base := fmt.Sprintf("%s%d", FilePrefix, rand.Uint64())
	path := filepath.Join(s.dataDir, base)

	log.WithFields(log.Fields{
		"hash": req.Hash,
		"path": path,
		"size": len(req.Block),
	}).Debug("Storing block...")

	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	if _, err := f.Write(req.Block); err != nil {
		return nil, err
	}

	if err := f.Sync(); err != nil {
		return nil, err
	}

	s.blocks[req.Hash] = datafile{
		path: path,
	}

	return &StoreBlockResponse{
		Success: true,
	}, nil
}

func (s *Store) HasBlock(ctx context.Context, req *HasBlockRequest) (*HasBlockResponse, error) {
	_, ok := s.blocks[req.Hash]
	return &HasBlockResponse{
		Success: ok,
	}, nil
}

func (s *Store) GetBlock(ctx context.Context, req *GetBlockRequest) (*GetBlockResponse, error) {
	df, ok := s.blocks[req.Hash]

	if !ok {
		return &GetBlockResponse{
			Success: ok,
			Block:   nil,
		}, nil
	}

	b, err := df.readAll()
	if err != nil {
		return nil, err
	}

	return &GetBlockResponse{
		Success: true,
		Block:   b,
	}, nil

}
