package bll

import (
	"container/list"
	"dal"
	"fmt"
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
		fmt.Println(fmt.Sprintf("%s does not exist", fingerprint.Filename))
		verifier.countMissing++
	} else {
		hasher := NewHasher(fingerprint.Algorithm)
		checksum := hasher.CalculateChecksumForFile(fullPath)
		if !compareByteSlices(checksum, fingerprint.Checksum) {
			fmt.Println(fmt.Sprintf("%s has an invalid checksum", fingerprint.Filename))
			verifier.countInvalid++
		}
	}
}

func (verifier *Verifier) printSummary(fingerprints *list.List) {

	countAll := fingerprints.Len()
	countValid := countAll - verifier.countMissing - verifier.countInvalid

	fmt.Println()
	fmt.Printf(
		"Valid: %d/%d, missing: %d, invalid: %d.",
		countValid, countAll, verifier.countMissing, verifier.countInvalid)
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
