package block

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

const FilePrefix = "blk_"

type Store struct {
	blocks  map[string]datafile
	dataDir string
	counter uint64
}

func NewStore(dataDir string) *Store {
	return &Store{
		blocks:  make(map[string]datafile),
		dataDir: dataDir,
		counter: 0,
	}
}

func (s *Store) StoreBlock(ctx context.Context, req *StoreBlockRequest) (*StoreBlockResponse, error) {
	log.WithFields(log.Fields{
		"hash": req.Hash,
		"size": len(req.Block),
	}).Debug("storing block")

	base := fmt.Sprintf("%s%d", FilePrefix, s.counter)
	path := filepath.Join(s.dataDir, base)

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
	return &GetBlockResponse{
		Success: false,
		Block:   nil,
	}, nil
}
