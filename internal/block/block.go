package block

import (
	"crypto/sha256"
	"encoding/base64"
	"io"
)

type Block struct {
	Block []byte
	Hash string
}

var DefaultBlockSize = 4096

func Blocks(r io.Reader) ([]Block, error) {

	blocks := make([]Block, 1)

	for {
		block := make([]byte, DefaultBlockSize)

		n, err := r.Read(block)
		if err == io.EOF {
			sha := sha256.Sum256(block[:n])
			hash := base64.StdEncoding.EncodeToString(sha[:])
			blocks = append(blocks, Block{
				Block: block,
				Hash: hash,
			})
			break
		} else if err != nil {
			return nil, err
		}

		sha := sha256.Sum256(block[:n])
		hash := base64.StdEncoding.EncodeToString(sha[:])
		blocks = append(blocks, Block{
			Block: block,
			Hash: hash,
		})
	}

	return blocks, nil
}