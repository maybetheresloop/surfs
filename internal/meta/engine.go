package meta

import (
	"github.com/maybetheresloop/keychain"
)

// Engine is an interface to represent key-value storage for the
// filename-hash list mapping.
type engine interface {
	// Sets the metadata for the specified file.
	SetFileMetadata(filename string, stat Stat) error

	// Gets the metadata associated with the specified file.
	GetFileMetadata(filename string) (Stat, bool, error)
}

// Implementation of the Engine interface backed by a Keychain key-value
// store.
type keychainEngine keychain.Keychain

func (k *keychainEngine) Set(string) {

}

// Implementation of the Engine interface backed by a regular Go map.
type mapEngine map[string]Stat

func NewMapEngine() mapEngine {
	return make(map[string]Stat)
}

func (m mapEngine) SetFileMetadata(filename string, stat Stat) error {
	m[filename] = stat
	return nil
}

func (m mapEngine) GetFileMetadata(filename string) (Stat, bool, error) {
	stat, ok := m[filename]
	return stat, ok, nil
}
