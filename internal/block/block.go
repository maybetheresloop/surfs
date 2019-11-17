package block

import (
	"crypto/sha256"
	"encoding/base64"
	"io"
)

type Map struct {
	Hashes []string
	Blocks map[string][]byte
}

const DefaultBlockSize uint64 = 64

// Calculates the Base64-encoded SHA256 hash of the specified block.
func blockHash(block []byte) string {
	sha := sha256.Sum256(block)
	return base64.StdEncoding.EncodeToString(sha[:])
}

// Divides the contents of the specified reader into blocks of the specified size. A map of hashes to blocks is
// returned.
func makeBlocksWithSize(r io.Reader, size uint64) (*Map, error) {
	m := &Map{
		Blocks: make(map[string][]byte),
		Hashes: make([]string, 0, 64),
	}
	for {
		block := make([]byte, size)
		n, err := io.ReadFull(r, block)
		if err != nil {
			if err == io.EOF {
				return m, nil
			} else if err == io.ErrUnexpectedEOF {
				hash := blockHash(block[:n])
				m.Blocks[hash] = block[:n]
				m.Hashes = append(m.Hashes, hash)
				return m, nil
			} else {
				return nil, err
			}
		}

		hash := blockHash(block)
		m.Blocks[hash] = block
		m.Hashes = append(m.Hashes, hash)
	}
}

// Divides the contents of the specified reader into blocks of 4KB. A map of hashes to block data
// and a slice of the hashes in order are returned.
func MakeBlocks(r io.Reader) (*Map, error) {
	return makeBlocksWithSize(r, DefaultBlockSize)
}
