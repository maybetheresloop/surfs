package meta

import (
	"context"
	"fmt"
	"os"
	"surfs/internal/block"

	"google.golang.org/grpc"

	log "github.com/sirupsen/logrus"
)

type MetadataStore struct {
	// Map from filename to hash list.
	files map[string]stat

	conn *grpc.ClientConn

	// gRPC client to the block store.
	client block.StoreClient
}

// Creates a new Metadata store service.
func NewStore(blockStoreAddr string) (*MetadataStore, error) {
	fmt.Fprintf(os.Stderr, "connecting to block store")

	conn, err := grpc.Dial(blockStoreAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	client := block.NewStoreClient(conn)

	fmt.Fprintf(os.Stderr, "connected to block store")

	return &MetadataStore{
		files:  make(map[string]stat),
		conn:   conn,
		client: client,
	}, nil
}

func (s *MetadataStore) Close() error {
	if s.conn == nil {
		return nil
	}

	return s.conn.Close()
}

// Reads a file from the metadata store. In reality, this RPC returns the hashes of the blocks corresponding to the
// desired file. It is the responsibility of the calling client to then contact the block store service and retrieve
// from it the blocks corresponding to the hashes.
func (s *MetadataStore) ReadFile(ctx context.Context, req *ReadFileRequest) (*ReadFileResponse, error) {
	log.WithFields(log.Fields{
		"filename": req.Filename,
	}).Debug("reading file")

	// Even if the file metadata is not found, returning the zero value still works.
	st, _ := s.files[req.Filename]

	return &ReadFileResponse{
		Version:  st.version,
		HashList: st.hashList,
	}, nil
}

// Modifies the specified file.
func (s *MetadataStore) ModifyFile(ctx context.Context, req *ModifyFileRequest) (*ModifyFileResponse, error) {
	log.WithFields(log.Fields{
		"filename": req.Filename,
		"version":  req.Version,
	}).Debug("modifying file")

	missing := make([]string, 0, 16)
	for _, hash := range req.HashList {
		res, err := s.client.HasBlock(context.Background(), &block.HasBlockRequest{
			Hash: hash,
		})
		if err != nil {
			return nil, err
		}

		if !res.Success {
			missing = append(missing, hash)
		}
	}

	if len(missing) == 0 {
		log.WithFields(log.Fields{
			"filename": req.Filename,
			"version":  req.Version,
		}).Debug("modified file successfully")

		s.files[req.Filename] = stat{
			hashList: req.HashList,
		}

		return &ModifyFileResponse{Success: true}, nil
	}

	log.WithFields(log.Fields{
		"filename": req.Filename,
		"version":  req.Version,
	}).Debugf("did not modify file successfully; missing %d blocks", len(missing))

	return &ModifyFileResponse{Success: false, MissingHashList: missing}, nil
}

// Deletes the specified file.
func (s *MetadataStore) DeleteFile(ctx context.Context, req *DeleteFileRequest) (*DeleteFileResponse, error) {
	log.WithFields(log.Fields{
		"filename": req.Filename,
		"version":  req.Version,
	}).Debug("deleting file")

	return &DeleteFileResponse{Success: true}, nil
}
