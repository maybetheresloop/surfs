package meta

import (
	"context"
	"reflect"
	"surfs/internal/block"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

type mockClient struct {
	blocks map[string][]byte
}

func (m *mockClient) HasBlock(ctx context.Context, in *block.HasBlockRequest, opts ...grpc.CallOption) (*block.HasBlockResponse, error) {
	_, ok := m.blocks[in.Hash]
	return &block.HasBlockResponse{
		Success: ok,
	}, nil
}

func (m *mockClient) GetBlock(ctx context.Context, in *block.GetBlockRequest, opts ...grpc.CallOption) (*block.GetBlockResponse, error) {
	blk, ok := m.blocks[in.Hash]

	return &block.GetBlockResponse{
		Success: ok,
		Block:   blk,
	}, nil
}

func (m *mockClient) StoreBlock(ctx context.Context, in *block.StoreBlockRequest, opts ...grpc.CallOption) (*block.StoreBlockResponse, error) {
	m.blocks[in.Hash] = in.Block

	return &block.StoreBlockResponse{
		Success: true,
	}, nil
}

func expectReadFile(store *MetadataStore, filename string, expected *ReadFileResponse, t *testing.T) {
	req := &ReadFileRequest{
		Filename: filename,
	}

	res, err := store.ReadFile(context.Background(), req)
	assert.Nil(t, err)

	assert.Equal(t, res.Version, expected.Version)
	assert.ElementsMatch(t, expected.HashList, res.HashList)
}

func TestMetadataStore_ReadFile(t *testing.T) {
	mock := &mockClient{blocks: map[string][]byte{}}
	store := &MetadataStore{
		files: map[string]stat{
			"file1": {
				hashList: []string{"hash1", "hash2"},
				version:  1,
			},
		},
		conn:   nil,
		client: mock,
	}

	expectReadFile(store, "file1", &ReadFileResponse{HashList: []string{"hash1", "hash2"}, Version: 1}, t)
	expectReadFile(store, "file2", &ReadFileResponse{HashList: nil, Version: 0}, t)
}

func expectModifyFile(store *MetadataStore, req *ModifyFileRequest, expected *ModifyFileResponse, t *testing.T) {
	res, err := store.ModifyFile(context.Background(), req)
	if err != nil {
		t.Fatalf("failed to modify file, %v", err)
	}

	if expected.Success != res.Success {
		t.Fatalf("incorrect success status: expected =%t, got =%t", expected.Success, res.Success)
	}

	if !reflect.DeepEqual(expected.MissingHashList, res.MissingHashList) {
		t.Fatalf("incorrect missing hash list: expected =%v, got =%v", expected.MissingHashList, res.MissingHashList)
	}

	if expected.Success && !reflect.DeepEqual(store.files[req.Filename].hashList, req.HashList) {
		t.Fatalf("incorrect hash list after modify: expected =%v, got =%v", req.HashList, store.files[req.Filename].hashList)
	}
}

func TestMetadataStore_ModifyFile(t *testing.T) {

	mock := &mockClient{blocks: map[string][]byte{
		"hash1": []byte("block1"),
		"hash2": []byte("block2"),
	}}
	store := &MetadataStore{
		files: map[string]stat{
			"file1": {
				hashList: []string{"hash1"},
				version:  1,
			},
		},
		conn:   nil,
		client: mock,
	}

	// *********************************************
	// * TEST #1: Test modifying an existing file. *
	// *********************************************

	// Test all blocks found.
	expectModifyFile(store, &ModifyFileRequest{
		Filename: "file1",
		Version:  2,
		HashList: []string{"hash1", "hash2"},
	}, &ModifyFileResponse{Success: true, MissingHashList: nil}, t)

	// Expect correct modifications.
	stat := store.files["file1"]
	assert.Equal(t, stat.version, uint64(2))
	assert.Equal(t, stat.hashList, []string{"hash1", "hash2"})

	// ***************
	// * END TEST #1 *
	// ***************

	// **************************************
	// * TEST #2: Test creating a new file. *
	// **************************************

	// Test missing blocks in block store.
	expectModifyFile(store, &ModifyFileRequest{
		Filename: "file2",
		Version:  1,
		HashList: []string{"hash1", "hash2", "hash3", "hash4"},
	}, &ModifyFileResponse{Success: false, MissingHashList: []string{"hash3", "hash4"}}, t)

	// ***************
	// * END TEST #2 *
	// ***************
}

func expectDeleteFile(store *MetadataStore, req *DeleteFileRequest, expected *DeleteFileResponse, t *testing.T) {
	res, err := store.DeleteFile(context.Background(), req)

	assert.Nil(t, err)
	assert.Equal(t, expected.Success, res.Success)
}

func TestMetadataStore_DeleteFile(t *testing.T) {
	mock := &mockClient{blocks: map[string][]byte{}}
	store := &MetadataStore{
		files: map[string]stat{
			"file1": {
				version:  1,
				hashList: []string{"hash1"},
			},
			"file2": {
				version:  2,
				hashList: []string{"hash2"},
			},
		},
		conn:   nil,
		client: mock,
	}

	expectDeleteFile(store, &DeleteFileRequest{Filename: "file1", Version: 2}, &DeleteFileResponse{Success: true}, t)
	expectDeleteFile(store, &DeleteFileRequest{Filename: "file2", Version: 2}, &DeleteFileResponse{Success: false}, t)
}

func expectGetVersion(store *MetadataStore, req *GetVersionRequest, expected *GetVersionResponse, t *testing.T) {
	res, err := store.GetVersion(context.Background(), req)
	assert.Nil(t, err)

	assert.Equal(t, expected.Version, res.Version)
}

func TestMetadataStore_GetVersion(t *testing.T) {
	store := &MetadataStore{
		files: map[string]stat{
			"file1": {
				version:  1,
				hashList: nil,
			},
		},
		conn:   nil,
		client: nil,
	}

	expectGetVersion(store, &GetVersionRequest{Filename: "file1"}, &GetVersionResponse{Version: 1}, t)
	expectGetVersion(store, &GetVersionRequest{Filename: "file2"}, &GetVersionResponse{Version: 0}, t)
}
