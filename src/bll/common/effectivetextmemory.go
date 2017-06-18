package common

import (
	"encoding/hex"
	"hash"
	"hash/fnv"
)

// EffectiveTextMemory Stores hashes of strings and helps deciding whether it has already seen a string or not.
type EffectiveTextMemory struct {
	CountAll        int
	CountCollisions int
	UseCache        bool
	hashFunc        hash.Hash
	storage         map[string]bool
	textCache       map[string]bool
}

// NewEffectiveTextMemory Instantiates a new EffectiveTextMemory object.
func NewEffectiveTextMemory() *EffectiveTextMemory {

	storage := make(map[string]bool)
	textCache := make(map[string]bool)

	return &EffectiveTextMemory{0, 0, true, fnv.New32a(), storage, textCache}
}

// ClearCache Clears the cache.
func (etm *EffectiveTextMemory) ClearCache() {

	etm.textCache = make(map[string]bool)
}

// ContainsText Checks if the given string is memorized.
func (etm *EffectiveTextMemory) ContainsText(text string) bool {

	textHash := etm.calculateHash(text)

	return etm.storage[textHash]
}

// Write Memorizes the given string.
func (etm *EffectiveTextMemory) Write(text string) {

	if etm.UseCache {
		if etm.textCache[text] {
			return
		}
		etm.textCache[text] = true
	}

	textHash := etm.calculateHash(text)
	etm.storeHash(textHash)
}

func (etm *EffectiveTextMemory) calculateHash(text string) string {

	textBytes := []byte(text)
	textHashBytes := etm.hashFunc.Sum(textBytes)
	textHash := hex.EncodeToString(textHashBytes)

	return textHash
}

func (etm *EffectiveTextMemory) storeHash(textHash string) {

	etm.CountAll++

	if etm.storage[textHash] {
		etm.CountCollisions++
	} else {
		etm.storage[textHash] = true
	}
}
