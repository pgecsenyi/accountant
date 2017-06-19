package report

import (
	"testing"
	"util"
)

var comparerReportTestHelper = util.NewTestHelper()

func TestComparerReport(t *testing.T) {

	t.Run("AddMissingFile", testCrAddMissingFile)
	t.Run("AddNewFile", testCrAddNewFile)
}

func testCrAddMissingFile(t *testing.T) {

	cr := NewComparerReport()
	testItem := "somedirectory/sumesubdirectory/somefile.txt"

	cr.AddMissingFile(testItem)

	if !comparerReportTestHelper.HasStringItems(cr.MissingFiles, testItem) {
		t.Errorf("%s should be marked as missing.", testItem)
	}
}

func testCrAddNewFile(t *testing.T) {

	cr := NewComparerReport()
	testItem := "somedirectory/sumesubdirectory/somefile.txt"

	cr.AddNewFile(testItem)

	if !comparerReportTestHelper.HasStringItems(cr.NewFiles, testItem) {
		t.Errorf("%s should be marked as missing.", testItem)
	}
}
