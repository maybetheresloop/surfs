package block

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
)

type Block struct {
	Block []byte
	Hash string
}

var DefaultBlockSize uint64 = 4096

func blockHash(block []byte) string {
	sha := sha256.Sum256(block)
	return base64.StdEncoding.EncodeToString(sha[:])
}

func blocksWithSize(r io.Reader, size uint64) ([]Block, error) {
	blocks := make([]Block, 0, 1)
	for {
		block := make([]byte, size)
		n, err := io.ReadFull(r, block)
		fmt.Printf("block: %s\n", block)
		if err != nil {
			if err == io.EOF {
				return blocks, nil
			} else if err == io.ErrUnexpectedEOF {
				blocks = append(blocks, Block{
					Block: block[:n],
					Hash: blockHash(block[:n]),
				})
				return blocks, nil
			} else {
				return nil, err
			}
		}

		blocks = append(blocks, Block{
			Block: block,
			Hash: blockHash(block),
		})
	}
}

func Blocks(r io.Reader) ([]Block, error) {
	return blocksWithSize(r, DefaultBlockSize)
}