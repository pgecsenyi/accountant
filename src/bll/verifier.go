package bll

import (
	"dal"
	"fmt"
	"log"
	"path"
	"util"
)

// Verifier Stores settings related to verification.
type Verifier struct {
	Db             dal.Database
	InputChecksums string
	BasePath       string
	countInvalid   int
	countMissing   int
}

// NewVerifier Instantiates a new Verifier object.
func NewVerifier(db dal.Database, inputChecksums string, basePath string) Verifier {

	basePath = util.NormalizePath(basePath)

	return Verifier{db, inputChecksums, basePath, 0, 0}
}

// Verify Verifies checksums in the given file.
func (verifier *Verifier) Verify(verifyNamesOnly bool) {

	verifier.Db.LoadFingerprints()
	verifier.verifyEntries(verifyNamesOnly)
	verifier.printSummary(verifyNamesOnly)
}

func (verifier *Verifier) verifyEntries(verifyNamesOnly bool) {

	fingerprints := verifier.Db.GetFingerprints()

	for element := fingerprints.Front(); element != nil; element = element.Next() {
		fingerprint := element.Value.(*dal.Fingerprint)
		verifier.verifyEntry(fingerprint, verifyNamesOnly)
	}
}

func (verifier *Verifier) verifyEntry(fingerprint *dal.Fingerprint, verifyNameOnly bool) {

	fullPath := path.Join(verifier.BasePath, fingerprint.Filename)

	if !util.CheckIfFileExists(fullPath) {
		log.Println(fmt.Sprintf("Missing: %s", fingerprint.Filename))
		verifier.countMissing++
	} else if !verifyNameOnly {
		verifier.verifyChecksum(fingerprint, fullPath)
	}
}

func (verifier *Verifier) verifyChecksum(fingerprint *dal.Fingerprint, fullPath string) {

	hasher := NewHasher(fingerprint.Algorithm)
	checksum := hasher.CalculateChecksum(fullPath)

	if !compareByteSlices(checksum, fingerprint.Checksum) {
		log.Println(fmt.Sprintf("Invalid: %s", fingerprint.Filename))
		verifier.countInvalid++
	}
}

func (verifier *Verifier) printSummary(verifyNamesOnly bool) {

	countAll := verifier.Db.GetFingerprints().Len()
	countValid := countAll - verifier.countMissing

	if verifyNamesOnly {
		log.Println(fmt.Sprintf(
			"Summary: %d/%d exist(s), %d missing.",
			countValid, countAll, verifier.countMissing))
	} else {
		log.Println(fmt.Sprintf(
			"Summary: %d/%d valid, %d missing, %d invalid.",
			countValid, countAll, verifier.countMissing, verifier.countInvalid))
	}
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
