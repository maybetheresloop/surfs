package meta

import (
	"bytes"
	"encoding/binary"

	"github.com/maybetheresloop/keychain"
)

// Engine is an interface to represent key-value storage for the
// filename-hash list mapping.
type engine interface {
	// Sets the metadata for the specified file.
	setFileMetadata(filename string, stat Stat) error

	// Gets the metadata associated with the specified file.
	getFileMetadata(filename string) (Stat, bool, error)
}

// Implementation of the Engine interface backed by a Keychain key-value
// store.
type keychainEngine struct {
	inner *keychain.Keychain
}

// Creates a Keychain engine by opening the specified Keychain store file.
func openKeychainEngine(name string) (keychainEngine, error) {
	kc, err := keychain.Open(name)
	if err != nil {
		return keychainEngine{}, nil
	}

	return keychainEngine{inner: kc}, nil
}

// Sets the metadata for the specified file.
// TODO: research on how the hash list can be stored more compactly instead of wasting space
func (k keychainEngine) setFileMetadata(filename string, stat Stat) error {
	var buf bytes.Buffer

	var versionBytes [8]byte
	binary.BigEndian.PutUint64(versionBytes[:], stat.version)
	buf.Write(versionBytes[:])

	for _, hash := range stat.hashList {
		buf.WriteString(hash)
	}

	return k.inner.Set([]byte(filename), buf.Bytes())
}

func (k keychainEngine) getFileMetadata(filename string) (Stat, bool, error) {
	b, err := k.inner.Get([]byte(filename))
	if err != nil {
		return Stat{}, false, err
	}

	if b == nil {
		return Stat{}, false, nil
	}

	stat := Stat{}

	// Each hash is 32 bytes long, and the header is 8 bytes.
	numHashes := (len(b) - 8) / 32
	stat.version = binary.BigEndian.Uint64(b[:8])
	stat.hashList = make([]string, 0, numHashes)

	b = b[8:]
	for i := 0; i < numHashes; i++ {
		stat.hashList = append(stat.hashList, string(b[i*8:i*8+8]))
	}

	return Stat{}, true, nil
}

// Implementation of the Engine interface backed by a regular Go map.
type mapEngine map[string]Stat

func newMapEngine() mapEngine {
	return make(map[string]Stat)
}

func (m mapEngine) setFileMetadata(filename string, stat Stat) error {
	m[filename] = stat
	return nil
}

func (m mapEngine) getFileMetadata(filename string) (Stat, bool, error) {
	stat, ok := m[filename]
	return stat, ok, nil
}
