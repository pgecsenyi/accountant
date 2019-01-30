package report

import (
	"fmr/util"
	"testing"
)

var comparisonReportTestHelper = util.NewTestHelper()

func TestComparisonReport(t *testing.T) {

	t.Run("AddMissingFile", testCrAddMissingFile)
	t.Run("AddNewFile", testCrAddNewFile)
}

func testCrAddMissingFile(t *testing.T) {

	cr := NewComparisonReport()
	testItem := "somedirectory/sumesubdirectory/somefile.txt"

	cr.AddMissingFile(testItem)

	if !comparisonReportTestHelper.HasStringItems(cr.MissingFiles, testItem) {
		t.Errorf("%s should be marked as missing.", testItem)
	}
}

func testCrAddNewFile(t *testing.T) {

	cr := NewComparisonReport()
	testItem := "somedirectory/sumesubdirectory/somefile.txt"

	cr.AddNewFile(testItem)

	if !comparisonReportTestHelper.HasStringItems(cr.NewFiles, testItem) {
		t.Errorf("%s should be marked as missing.", testItem)
	}
}
