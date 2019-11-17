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

var DefaultBlockSize uint64 = 64

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

// Divides the contents of the specified reader into blocks of the specified size. A map of hashes to blocks is
// returned.
func makeBlocksWithSize(r io.Reader, size uint64) (map[string][]byte, []string, error) {
	blocks := make(map[string][]byte)
	hashes := make([]string, 0, 64)
	for {
		block := make([]byte, size)
		n, err := io.ReadFull(r, block)
		if err != nil {
			if err == io.EOF {
				return blocks, hashes, nil
			} else if err == io.ErrUnexpectedEOF {
				hash := blockHash(block[:n])
				blocks[hash] = block[:n]
				hashes = append(hashes, hash)
				return blocks, hashes, nil
			} else {
				return nil, nil, err
			}
		}

		hash := blockHash(block)
		blocks[hash] = block
		hashes = append(hashes, hash)
	}
}

// Divides the contents of the specified reader into blocks of 4KB. A map of hashes to block data
// and a slice of the hashes in order are returned.
func MakeBlocks(r io.Reader) (map[string][]byte, []string, error) {
	return makeBlocksWithSize(r, DefaultBlockSize)
}

func BlocksWithSizeHint(r io.Reader) ([]Block, error) {
	return nil, nil
}
