package block

import (
	"io/ioutil"
	"os"
)

// A wrapper around a file that holds the data of a corresponding block.
type datafile struct {
	path string
}

// Reads the entire file to a byte slice.
func (d *datafile) readAll() ([]byte, error) {
	f, err := os.Open(d.path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return b, nil
}
