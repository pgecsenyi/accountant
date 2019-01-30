package bll

import (
	"fmr/bll/common"
	"fmr/bll/report"
	"fmr/dal"
	"fmr/util"
	"path"
)

// Verifier Stores settings related to verification.
type Verifier struct {
	Db       dal.Database
	BasePath string
	Report   *report.VerificationReport
}

// NewVerifier Instantiates a new Verifier object.
func NewVerifier(db dal.Database, basePath string) Verifier {

	basePath = util.NormalizePath(basePath)
	report := report.NewVerificationReport()

	return Verifier{db, basePath, report}
}

// Verify Verifies checksums in the given file.
func (verifier *Verifier) Verify(verifyNamesOnly bool, fpFilter common.FingerprintFilter) {

	verifier.Db.LoadFingerprints()
	verifier.verifyEntries(verifyNamesOnly, fpFilter)
	verifier.Report.LogSummary(!verifyNamesOnly)
}

func (verifier *Verifier) verifyEntries(verifyNamesOnly bool, fpFilter common.FingerprintFilter) {

	fingerprints := verifier.Db.GetFingerprints()

	for element := fingerprints.Front(); element != nil; element = element.Next() {
		fingerprint := element.Value.(*dal.Fingerprint)
		if fpFilter.FilterFingerprint(fingerprint) {
			verifier.verifyEntry(fingerprint, verifyNamesOnly)
		}
	}
}

func (verifier *Verifier) verifyEntry(fingerprint *dal.Fingerprint, verifyNameOnly bool) {

	fullPath := path.Join(verifier.BasePath, fingerprint.Filename)

	if !util.CheckIfFileExists(fullPath) {
		verifier.Report.AddMissingFile(fingerprint.Filename)
	} else if !verifyNameOnly {
		verifier.verifyChecksum(fingerprint, fullPath)
	} else {
		verifier.Report.AddValidFile(fingerprint.Filename)
	}
}

func (verifier *Verifier) verifyChecksum(fingerprint *dal.Fingerprint, fullPath string) {

	hasher := common.NewHasher(fingerprint.Algorithm)
	checksum := hasher.CalculateChecksum(fullPath)

	if util.CompareByteSlices(checksum, fingerprint.Checksum) {
		verifier.Report.AddValidFile(fingerprint.Filename)
	} else {
		verifier.Report.AddCorruptFile(fingerprint.Filename)
	}
}
