package bll

import (
	"encoding/hex"
	"hash"
)

type nameHashStorage struct {
	CountAll        int
	CountCollisions int
	hashFunc        hash.Hash
	storage         map[string]bool
	nameCache       map[string]bool
}

// NewNameHashStorage Instantiates a new NameHashStorage object.
func newNameHashStorage(hashFunc hash.Hash) *nameHashStorage {

	storage := make(map[string]bool)
	cache := make(map[string]bool)
	return &nameHashStorage{0, 0, hashFunc, storage, cache}
}

func (nhs *nameHashStorage) ClearCache() {

	nhs.nameCache = make(map[string]bool)
}

func (nhs *nameHashStorage) ContainsName(name string) bool {

	nameHash := nhs.calculateHash(name)

	return nhs.storage[nameHash]
}

func (nhs *nameHashStorage) Write(name string) {

	if nhs.nameCache[name] {
		return
	}

	nhs.nameCache[name] = true
	nameHash := nhs.calculateHash(name)
	nhs.storeHash(nameHash)
}

func (nhs *nameHashStorage) calculateHash(name string) string {

	nameBytes := []byte(name)
	nameHashBytes := nhs.hashFunc.Sum(nameBytes)
	nameHash := hex.EncodeToString(nameHashBytes)

	return nameHash
}

func (nhs *nameHashStorage) storeHash(hash string) {

	nhs.CountAll++

	if nhs.storage[hash] {
		nhs.CountCollisions++
	} else {
		nhs.storage[hash] = true
	}
}
