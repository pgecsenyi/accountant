package dal

import (
	"container/list"
	"util"
)

// MemoryDatabase Logic for calculating checksums.
type MemoryDatabase struct {
	Fingerprints *list.List
	NamePairs    *list.List
}

// NewMemoryDatabase Instantiates a new MemoryDatabase object.
func NewMemoryDatabase() *MemoryDatabase {

	return &MemoryDatabase{list.New(), list.New()}
}

// AddFingerprint Adds a fingerprint to the database.
func (db *MemoryDatabase) AddFingerprint(fingerprint *Fingerprint) {

	if fingerprint != nil {
		db.Fingerprints.PushFront(fingerprint)
	}
}

// AddNamePair Adds a name pair to the database.
func (db *MemoryDatabase) AddNamePair(namePair *NamePair) {

	if namePair != nil {
		db.NamePairs.PushFront(namePair)
	}
}

// GetFingerprints Returns stored fingerprints.
func (db *MemoryDatabase) GetFingerprints() *list.List {

	return db.Fingerprints
}

// LoadFingerprints Does nothing, there's nothing to load.
func (db *MemoryDatabase) LoadFingerprints() {
}

// LoadNamesFromFingeprints Passes filenames to the given StringWriter.
func (db *MemoryDatabase) LoadNamesFromFingeprints(writer util.StringWriter) {

	for element := db.Fingerprints.Front(); element != nil; element = element.Next() {
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

// SetFingerprints Sets stored fingerprints.
func (db *MemoryDatabase) SetFingerprints(fingerprints *list.List) {

	db.Fingerprints = fingerprints
}
