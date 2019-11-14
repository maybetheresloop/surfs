package block

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func createTempFileWithContent(b []byte) (string, error) {
	f, err := ioutil.TempFile("/tmp", "datafile-test")
	defer f.Close()
	if err != nil {
		return "", err
	}

	_, err = f.Write(b)
	if err != nil {
		return "", err
	}

	return f.Name(), nil
}

func TestReadAll(t *testing.T) {
	name, err := createTempFileWithContent([]byte("asdf"))
	if err != nil {
		t.Fatalf("unable to create temp file, %v", err)
	}

	defer os.Remove(name)

	d := datafile{path: name}
	got, err := d.readAll()
	if err != nil {
		t.Fatalf("failed to read from temp file")
	}

	if bytes.Compare([]byte("asdf"), got) != 0 {
		t.Fatalf("datafile content incorrect: expected =%s, got =%s", []byte("asdf"), got)
	}

}
