package block

import (
	"strings"

	"github.com/maybetheresloop/keychain"
)

type engine interface {
	SetHashList(filename string, hashList []string) error
	GetHashList(filename string) ([]string, error)
}

type keychainEngine struct {
	keys *keychain.Keychain
}

func (k *keychainEngine) SetHashList(filename string, hashList []string) error {
	var value []byte
	if hashList != nil {
		value = []byte(strings.Join(hashList, ","))
	}

	return k.keys.Set([]byte(filename), value)
}

func (k *keychainEngine) GetHashList(filename string) ([]string, error) {
	hashes, err := k.keys.Get([]byte(filename))
	if err != nil {
		return nil, err
	}

	if hashes == nil {
		return nil, nil
	}

	return strings.Split(string(hashes), ","), nil
}

func openKeychainEngine(filename string) (*keychainEngine, error) {
	keys, err := keychain.OpenConf(filename, &keychain.Conf{Sync: true})
	if err != nil {
		return nil, err
	}

	return &keychainEngine{keys: keys}, nil
}
