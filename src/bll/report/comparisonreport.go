package report

import (
	"container/list"
	"fmt"
	"log"
)

// ComparisonReport Stores statistics of a comparison process.
type ComparisonReport struct {
	MissingFiles *list.List
	NewFiles     *list.List
}

// NewComparisonReport Instantiates a new ComparisonReport object.
func NewComparisonReport() *ComparisonReport {

	return &ComparisonReport{list.New(), list.New()}
}

// AddMissingFile Adds the given file to the list of missing files.
func (cr *ComparisonReport) AddMissingFile(filename string) {

	cr.MissingFiles.PushFront(filename)
	log.Println(fmt.Sprintf("Missing: %s", filename))
}

// AddNewFile Adds the given file to the list of new files.
func (cr *ComparisonReport) AddNewFile(filename string) {

	cr.NewFiles.PushFront(filename)
	log.Println(fmt.Sprintf("New: %s", filename))
}
