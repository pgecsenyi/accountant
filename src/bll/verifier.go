package bll

import (
	"checksum"
	"container/list"
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

// VerifyRecords Verifies checksums in the given file.
func (verifier *Verifier) Verify(hasher *checksum.FileHasher) {

	hasher.LoadCsv(verifier.InputChecksums)
	verifier.verifyEntries(hasher.Fingerprints)
	verifier.printSummary(hasher.Fingerprints)
}

func (verifier *Verifier) verifyEntries(fingerprints *list.List) {

	for element := fingerprints.Front(); element != nil; element = element.Next() {
		fingerprint := element.Value.(*checksum.Fingerprint)
		verifier.verifyEntry(fingerprint)
	}
}

func (verifier *Verifier) verifyEntry(fingerprint *checksum.Fingerprint) {

	fullPath := path.Join(verifier.BasePath, fingerprint.Filename)

	if !util.CheckIfFileExists(fullPath) {
		fmt.Println(fmt.Sprintf("%s does not exist", fingerprint.Filename))
		verifier.countMissing++
	} else {
		checksum := checksum.CalculateChecksumForFile(fullPath, fingerprint.Algorithm)
		if !util.Compare(checksum, fingerprint.Checksum) {
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
