package report

import (
	"container/list"
	"fmt"
	"log"
)

// VerificationReport Stores statistics of a verification process.
type VerificationReport struct {
	CountAll     int
	CorruptFiles *list.List
	MissingFiles *list.List
}

// NewVerificationReport Instantiates a new VerificationReport object.
func NewVerificationReport() *VerificationReport {

	return &VerificationReport{0, list.New(), list.New()}
}

// AddCorruptFile Adds the given file to the list of corrupt files.
func (vr *VerificationReport) AddCorruptFile(filename string) {

	vr.CorruptFiles.PushFront(filename)
	vr.CountAll++
	log.Println(fmt.Sprintf("Corrupt: %s", filename))
}

// AddMissingFile Adds the given file to the list of missing files.
func (vr *VerificationReport) AddMissingFile(filename string) {

	vr.MissingFiles.PushFront(filename)
	vr.CountAll++
	log.Println(fmt.Sprintf("Missing: %s", filename))
}

// AddValidFile Logs that the given file is valid.
func (vr *VerificationReport) AddValidFile(filename string) {

	vr.CountAll++
}

// LogSummary Prints a summary report to the log.
func (vr *VerificationReport) LogSummary(displayCorruptCount bool) {

	countCorrupt := vr.CorruptFiles.Len()
	countMissing := vr.MissingFiles.Len()
	countValid := vr.CountAll - countCorrupt - countMissing

	if displayCorruptCount {
		log.Println(fmt.Sprintf(
			"Summary: %d/%d valid, %d missing, %d corrupt.",
			countValid, vr.CountAll, countMissing, countCorrupt))
	} else {
		log.Println(fmt.Sprintf(
			"Summary: %d/%d exist(s), %d missing.",
			countValid, vr.CountAll, countMissing))
	}
}
