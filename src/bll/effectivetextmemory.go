package bll

import (
	"encoding/hex"
	"hash"
	"hash/fnv"
)

type effectiveTextMemory struct {
	CountAll        int
	CountCollisions int
	hashFunc        hash.Hash
	storage         map[string]bool
	textCache       map[string]bool
}

// NewNameHashStorage Instantiates a new NameHashStorage object.
func newEffectiveTextMemory() *effectiveTextMemory {

	storage := make(map[string]bool)
	textCache := make(map[string]bool)

	return &effectiveTextMemory{0, 0, fnv.New32a(), storage, textCache}
}

func (etm *effectiveTextMemory) ClearCache() {

	etm.textCache = make(map[string]bool)
}

func (etm *effectiveTextMemory) ContainsText(text string) bool {

	textHash := etm.calculateHash(text)

	return etm.storage[textHash]
}

func (etm *effectiveTextMemory) Write(text string) {

	if etm.textCache[text] {
		return
	}

	etm.textCache[text] = true
	textHash := etm.calculateHash(text)
	etm.storeHash(textHash)
}

func (etm *effectiveTextMemory) calculateHash(text string) string {

	textBytes := []byte(text)
	textHashBytes := etm.hashFunc.Sum(textBytes)
	textHash := hex.EncodeToString(textHashBytes)

	return textHash
}

func (etm *effectiveTextMemory) storeHash(textHash string) {

	etm.CountAll++

	if etm.storage[textHash] {
		etm.CountCollisions++
	} else {
		etm.storage[textHash] = true
	}
}
