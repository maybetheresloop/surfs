package block

import (
	"context"
	log "github.com/sirupsen/logrus"
	"os"
)

type Store struct {
	blocks map[string]*os.File
}

func NewStore() *Store {
	return &Store {
		blocks: make(map[string]*os.File),
	}
}

func (s *Store) StoreBlock(ctx context.Context, req *StoreBlockRequest) (*StoreBlockResponse, error) {
	log.WithFields(log.Fields{
		"hash": req.Hash,
		"size": len(req.Block),
	}).Debug("storing block")

	return &StoreBlockResponse{
		Success:              true,
	}, nil
}

func (s *Store) HasBlock(ctx context.Context, req *HasBlockRequest) (*HasBlockResponse, error) {
	return &HasBlockResponse{
		Success:              true,
	}, nil
}

func (s *Store) GetBlock(ctx context.Context, req *GetBlockRequest) (*GetBlockResponse, error) {
	return &GetBlockResponse{
		Success:              false,
		Block:                nil,
	}, nil
}