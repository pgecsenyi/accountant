package report

import (
	"fmt"
	"log"
)

// ImportReport Stores statistics of an import process.
type ImportReport struct {
	invalidEntryCountByFile map[string]int
}

// NewImportReport Instantiates a new ImportReport object.
func NewImportReport() *ImportReport {

	var invalidEntryCountByFile = make(map[string]int)

	return &ImportReport{invalidEntryCountByFile}
}

// GetInvalidEntryCount Gets the number of invalid entries in the given file.
func (ir *ImportReport) GetInvalidEntryCount(filename string) int {

	return ir.invalidEntryCountByFile[filename]
}

// IncreaseInvalidEntryCount Increments the invalid entry count for the given file by one.
func (ir *ImportReport) IncreaseInvalidEntryCount(filename string) {

	ir.invalidEntryCountByFile[filename]++
}

// LogSummaryForFile Prints a summary report for the given file to the log.
func (ir *ImportReport) LogSummaryForFile(filename string) {

	numberOfInvalidLines := ir.invalidEntryCountByFile[filename]
	if numberOfInvalidLines != 0 {
		message := fmt.Sprintf("There is/are %d invalid line(s) in %s.", numberOfInvalidLines, filename)
		log.Println(message)
	}
}
