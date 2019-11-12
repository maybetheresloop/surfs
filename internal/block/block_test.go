package block

import (
	"bytes"
	"testing"
)

func TestBlocksWithSize(t *testing.T) {
	r := bytes.NewReader([]byte("abcdefghijklmn"))

	blks, err := blocksWithSize(r, 4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(blks) != 4 {
		t.Fatalf("incorrect # blocks: expected =%d, got = %d", 4, len(blks))
	}
}
