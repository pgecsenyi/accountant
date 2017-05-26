package bll

import (
	"container/list"
	"dal"
	"fmt"
	"log"
	"path"
	"util"
)

// Verifier Stores settings related to verification.
type Verifier struct {
	InputChecksums string
	BasePath       string
	countInvalid   int
	countMissing   int
}

// NewVerifier Instantiates a new Verifier object.
func NewVerifier(inputChecksums string, basePath string) Verifier {

	basePath = util.NormalizePath(basePath)

	return Verifier{inputChecksums, basePath, 0, 0}
}

// Verify Verifies checksums in the given file.
func (verifier *Verifier) Verify(db *dal.Db) {

	db.LoadCsv(verifier.InputChecksums)
	verifier.verifyEntries(db.Fingerprints)
	verifier.printSummary(db.Fingerprints)
}

func (verifier *Verifier) verifyEntries(fingerprints *list.List) {

	for element := fingerprints.Front(); element != nil; element = element.Next() {
		fingerprint := element.Value.(*dal.Fingerprint)
		verifier.verifyEntry(fingerprint)
	}
}

func (verifier *Verifier) verifyEntry(fingerprint *dal.Fingerprint) {

	fullPath := path.Join(verifier.BasePath, fingerprint.Filename)

	if !util.CheckIfFileExists(fullPath) {
		log.Println(fmt.Sprintf("Missing: %s", fingerprint.Filename))
		verifier.countMissing++
	} else {
		hasher := NewHasher(fingerprint.Algorithm)
		checksum := hasher.CalculateChecksumForFile(fullPath)
		if !compareByteSlices(checksum, fingerprint.Checksum) {
			log.Println(fmt.Sprintf("Corrupt: %s", fingerprint.Filename))
			verifier.countInvalid++
		}
	}
}

func (verifier *Verifier) printSummary(fingerprints *list.List) {

	countAll := fingerprints.Len()
	countValid := countAll - verifier.countMissing - verifier.countInvalid

	log.Println(fmt.Sprintf(
		"Summary: %d/%d valid, %d missing, %d invalid.",
		countValid, countAll, verifier.countMissing, verifier.countInvalid))
}

func compareByteSlices(slice1 []byte, slice2 []byte) bool {

	if (slice1 == nil && slice2 != nil) || (slice1 != nil && slice2 == nil) {
		return false
	}
	if len(slice1) != len(slice2) {
		return false
	}

	for i := 0; i < len(slice1); i++ {
		if slice1[i] != slice2[i] {
			return false
		}
	}

	return true
}
