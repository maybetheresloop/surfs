package block

import (
	"crypto/sha256"
	"encoding/base64"
	"io"
)

// A pair of a block's contents and its Base64-encoded SHA256 hash.
type Block struct {
	Block []byte
	Hash  string
}

var DefaultBlockSize uint64 = 4096

// Calculates the Base64-encoded SHA256 hash of the specified block.
func blockHash(block []byte) string {
	sha := sha256.Sum256(block)
	return base64.StdEncoding.EncodeToString(sha[:])
}

// Divides the contents of the specified reader into blocks of the specified size. A slice of block-hash pairs
// are returned.
func blocksWithSize(r io.Reader, size uint64) ([]Block, error) {
	blocks := make([]Block, 0, 1)
	for {
		block := make([]byte, size)
		n, err := io.ReadFull(r, block)
		if err != nil {
			if err == io.EOF {
				return blocks, nil
			} else if err == io.ErrUnexpectedEOF {
				blocks = append(blocks, Block{
					Block: block[:n],
					Hash:  blockHash(block[:n]),
				})
				return blocks, nil
			} else {
				return nil, err
			}
		}

		blocks = append(blocks, Block{
			Block: block,
			Hash:  blockHash(block),
		})
	}
}

// Divides the contents of the specified reader into blocks of 4KB. A slice of block-hash pairs
// are returned.
func Blocks(r io.Reader) ([]Block, error) {
	return blocksWithSize(r, DefaultBlockSize)
}
