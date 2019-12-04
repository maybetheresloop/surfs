package meta

import (
	"context"
	"surfs/internal/block"

	"google.golang.org/grpc"

	log "github.com/sirupsen/logrus"
)

type MetadataStore struct {
	// Key-value storage engine for file metadata.
	engine engine

	// gRPC client to the block store.
	client block.StoreClient

	// Underlying gRPC connection for block store client.
	conn *grpc.ClientConn
}

// Creates a new Metadata store service.
func NewStore(blockStoreAddr string) (*MetadataStore, error) {

	log.Debugf("Connecting to block store at %s...", blockStoreAddr)

	conn, err := grpc.Dial(blockStoreAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	client := block.NewStoreClient(conn)

	log.Debug("Connected to block store.")

	return &MetadataStore{
		//files:  make(map[string]Stat),
		conn:   conn,
		client: client,
		engine: newMapEngine(),
	}, nil
}

// Closes the store. This simply closes the underlying gRPC connection.
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
	}).Debug("Reading file...")

	// Even if the file metadata is not found, returning the zero value still works.
	//st, _ := s.files[req.Filename]
	st, _, err := s.engine.getFileMetadata(req.Filename)
	if err != nil {
		return nil, err
	}

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
	}).Debug("Modifying file...")

	// The new version number must be exactly one more than the current version number. If it is not,
	// then we reject the modification.
	//oldVersion := s.files[req.Filename].version
	st, _, err := s.engine.getFileMetadata(req.Filename)
	if err != nil {
		return nil, err
	}

	oldVersion := st.version

	if req.Version != oldVersion+1 {

		log.WithFields(log.Fields{
			"filename":   req.Filename,
			"newVersion": req.Version,
			"oldVersion": oldVersion,
		}).Debug("Did not modify file; invalid new file version.")

		return &ModifyFileResponse{Success: false}, nil
	}

	// Check for missing blocks. If there are any blocks missing in the block store, return a list of
	// those missing blocks to the client. Otherwise, we have all the required blocks, and it is safe
	// for us to modify the file metadata to point to the new list of blocks.
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
		}).Debug("Modified file successfully.")

		//s.files[req.Filename] = Stat{
		//	hashList: req.HashList,
		//	version:  req.Version,
		//}
		if err := s.engine.setFileMetadata(req.Filename, Stat{
			hashList: req.HashList,
			version:  req.Version,
		}); err != nil {
			return nil, err
		}

		return &ModifyFileResponse{Success: true}, nil
	}

	log.WithFields(log.Fields{
		"filename": req.Filename,
		"version":  req.Version,
	}).Debugf("Did not modify file successfully; missing %d block(s).", len(missing))

	return &ModifyFileResponse{Success: false, MissingHashList: missing}, nil
}

// Deletes the specified file.
func (s *MetadataStore) DeleteFile(ctx context.Context, req *DeleteFileRequest) (*DeleteFileResponse, error) {
	log.WithFields(log.Fields{
		"filename": req.Filename,
		"version":  req.Version,
	}).Debug("Deleting file...")

	// The new version number must be exactly one more than the current version number. If it is not,
	// then we reject the deletion.
	//oldVersion := s.files[req.Filename].version
	st, _, err := s.engine.getFileMetadata(req.Filename)
	if err != nil {
		return nil, err
	}

	oldVersion := st.version

	if req.Version != oldVersion+1 {

		log.WithFields(log.Fields{
			"filename":   req.Filename,
			"newVersion": req.Version,
			"oldVersion": oldVersion,
		}).Debug("Did not delete file: invalid new file version.")

		return &DeleteFileResponse{Success: false}, nil
	}

	// Deleting the file simply consists of setting its hash list to a nil slice, which is automatically set by
	// the zero value of stat.
	//s.files[req.Filename] = Stat{
	//	version: req.Version,
	//}
	if err := s.engine.setFileMetadata(req.Filename, Stat{
		version: req.Version,
	}); err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"filename": req.Filename,
		"version":  req.Version,
	}).Debug("Deleted file successfully.")

	return &DeleteFileResponse{Success: true}, nil
}

func (s *MetadataStore) GetVersion(ctx context.Context, req *GetVersionRequest) (*GetVersionResponse, error) {

	//version := s.files[req.Filename].version
	st, _, err := s.engine.getFileMetadata(req.Filename)
	if err != nil {
		return nil, err
	}

	version := st.version

	log.WithFields(log.Fields{
		"version": version,
	}).Debug("Found file version.")

	// If the filename doesn't exist, we are fine because the zero value of stat{} will have version = 0.
	return &GetVersionResponse{
		Version: version,
	}, nil
}
