package report

import (
	"fmt"
	"log"
)

// ImporterReport Stores statistics of an import process.
type ImporterReport struct {
	invalidEntryCountByFile map[string]int
}

// NewImporterReport Instantiates a new ImporterReport object.
func NewImporterReport() *ImporterReport {

	var invalidEntryCountByFile = make(map[string]int)

	return &ImporterReport{invalidEntryCountByFile}
}

// GetInvalidEntryCount Gets the number of invalid entries in the given file.
func (ir *ImporterReport) GetInvalidEntryCount(filename string) int {

	return ir.invalidEntryCountByFile[filename]
}

// IncreaseInvalidEntryCount Increments the invalid entry count for the given file by one.
func (ir *ImporterReport) IncreaseInvalidEntryCount(filename string) {

	ir.invalidEntryCountByFile[filename]++
}

// LogSummaryForFile Prints a summary report for the given file to the log.
func (ir *ImporterReport) LogSummaryForFile(filename string) {

	numberOfInvalidLines := ir.invalidEntryCountByFile[filename]
	if numberOfInvalidLines != 0 {
		message := fmt.Sprintf("There is/are %d invalid line(s) in %s.", numberOfInvalidLines, filename)
		log.Println(message)
	}
}
