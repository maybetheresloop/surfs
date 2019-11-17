package block

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func tempBlock(dir string, pattern string, content []byte) (string, error) {
	f, err := ioutil.TempFile(dir, pattern)
	if err != nil {
		return "", nil
	}

	defer f.Close()

	_, err = f.Write(content)
	if err != nil {
		return "", nil
	}

	return f.Name(), nil
}

func TestStore_HasBlock(t *testing.T) {
	fp, err := tempBlock(os.TempDir(), "surfs", []byte("block1"))
	assert.Nil(t, err)

	defer os.Remove(fp)

	s := &Store{
		blocks: map[string]datafile{
			"hash1": {
				path: fp,
			},
		},
		dataDir: "",
		counter: 0,
	}

	req := &HasBlockRequest{
		Hash: "hash1",
	}

	res, err := s.HasBlock(context.Background(), req)
	assert.Nil(t, err)
	assert.True(t, res.Success)

	req = &HasBlockRequest{
		Hash: "hash2",
	}

	res, err = s.HasBlock(context.Background(), req)
	assert.Nil(t, err)
	assert.False(t, res.Success)
}

func TestStore_GetBlock(t *testing.T) {
	fp, err := tempBlock(os.TempDir(), "surfs", []byte("block1"))
	assert.Nil(t, err)

	defer os.Remove(fp)

	s := &Store{
		blocks: map[string]datafile{
			"hash1": {
				path: fp,
			},
		},
		dataDir: "",
		counter: 0,
	}

	req := &GetBlockRequest{
		Hash: "hash1",
	}

	res, err := s.GetBlock(context.Background(), req)
	assert.Nil(t, err)
	assert.True(t, res.Success)
	assert.Equal(t, []byte("block1"), res.Block)

	req = &GetBlockRequest{
		Hash: "hash2",
	}

	res, err = s.GetBlock(context.Background(), req)
	assert.Nil(t, err)
	assert.False(t, res.Success)
	assert.Equal(t, []byte(nil), res.Block)
}
