package bll

import (
	"checksum"
)

// Verifier Stores settings related to verification.
type Verifier struct {
	InputChecksums string
	BasePath       string
}

// VerifyRecords Verifies checksums in the given file.
func (verifier *Verifier) VerifyRecords(hasher *checksum.FileHasher) {

	hasher.LoadFromCsv(verifier.InputChecksums)
	hasher.VerifyFiles(verifier.BasePath)
}
