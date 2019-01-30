package dal

import (
	"container/list"
	"fmr/util"
)

// MemoryDatabase Logic for calculating checksums.
type MemoryDatabase struct {
	fingerprints *list.List
	namePairs    *list.List
}

// NewMemoryDatabase Instantiates a new MemoryDatabase object.
func NewMemoryDatabase() *MemoryDatabase {

	return &MemoryDatabase{list.New(), list.New()}
}

// AddFingerprint Adds a fingerprint to the database.
func (db *MemoryDatabase) AddFingerprint(fingerprint *Fingerprint) {

	if fingerprint != nil {
		db.fingerprints.PushFront(fingerprint)
	}
}

// AddFingerprints Adds a list of fingerprints to the database.
func (db *MemoryDatabase) AddFingerprints(fingerprints *list.List) {

	for element := fingerprints.Front(); element != nil; element = element.Next() {
		fingerprint := element.Value.(*Fingerprint)
		db.AddFingerprint(fingerprint)
	}
}

// AddNamePair Adds a name pair to the database.
func (db *MemoryDatabase) AddNamePair(namePair *NamePair) {

	if namePair != nil {
		db.namePairs.PushFront(namePair)
	}
}

// Clear Removes all entries from the database.
func (db *MemoryDatabase) Clear() {

	db.fingerprints.Init()
	db.namePairs.Init()
}

// GetFingerprints Returns stored fingerprints.
func (db *MemoryDatabase) GetFingerprints() *list.List {

	return db.fingerprints
}

// GetNamePairs Returns stored name pairs.
func (db *MemoryDatabase) GetNamePairs() *list.List {

	return db.namePairs
}

// LoadFingerprints Does nothing, there's nothing to load.
func (db *MemoryDatabase) LoadFingerprints() {
}

// LoadNamesFromFingeprints Passes filenames to the given StringWriter.
func (db *MemoryDatabase) LoadNamesFromFingeprints(writer util.StringWriter) {

	for element := db.fingerprints.Front(); element != nil; element = element.Next() {
		fingerprint := element.Value.(*Fingerprint)
		writer.Write(fingerprint.Filename)
	}
}

// SaveFingerprints Does nothing, there is nowhere to save.
func (db *MemoryDatabase) SaveFingerprints() {
}

// SaveNamePairs Does nothing, there is nowhere to save.
func (db *MemoryDatabase) SaveNamePairs() {
}
