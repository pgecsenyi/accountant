package bll

import (
	"bll/common"
	"bll/report"
	"dal"
	"path"
	"util"
)

// Verifier Stores settings related to verification.
type Verifier struct {
	Db             dal.Database
	InputChecksums string
	BasePath       string
	Report         *report.VerifierReport
}

// NewVerifier Instantiates a new Verifier object.
func NewVerifier(db dal.Database, inputChecksums string, basePath string) Verifier {

	basePath = util.NormalizePath(basePath)
	report := report.NewVerifierReport()

	return Verifier{db, inputChecksums, basePath, report}
}

// Verify Verifies checksums in the given file.
func (verifier *Verifier) Verify(verifyNamesOnly bool) {

	verifier.Db.LoadFingerprints()
	verifier.Report.CountAll = verifier.Db.GetFingerprints().Len()
	verifier.verifyEntries(verifyNamesOnly)
	verifier.Report.LogSummary(verifyNamesOnly)
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
		verifier.Report.AddMissingFile(fingerprint.Filename)
	} else if !verifyNameOnly {
		verifier.verifyChecksum(fingerprint, fullPath)
	}
}

func (verifier *Verifier) verifyChecksum(fingerprint *dal.Fingerprint, fullPath string) {

	hasher := common.NewHasher(fingerprint.Algorithm)
	checksum := hasher.CalculateChecksum(fullPath)

	if !util.CompareByteSlices(checksum, fingerprint.Checksum) {
		verifier.Report.AddCorruptFile(fingerprint.Filename)
	}
}
