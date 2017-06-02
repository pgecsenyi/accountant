package bll

import (
	"container/list"
	"fmt"
	"log"
)

// ComparerReport Stores the otherwise not persisted result of comparison.
type ComparerReport struct {
	MissingFiles *list.List
	NewFiles     *list.List
}

// NewComparerReport Instantiates a new ComparerReport object.
func NewComparerReport() *ComparerReport {

	return &ComparerReport{list.New(), list.New()}
}

// AddMissingFile Adds the given file to the list of missing files.
func (cr *ComparerReport) AddMissingFile(filename string) {

	cr.MissingFiles.PushFront(filename)
	log.Println(fmt.Sprintf("Missing: %s", filename))
}

// AddNewFile Adds the given file to the list of new files.
func (cr *ComparerReport) AddNewFile(filename string) {

	cr.NewFiles.PushFront(filename)
	log.Println(fmt.Sprintf("New: %s", filename))
}
