package report

import "testing"

func TestImporterReport(t *testing.T) {

	t.Run("IncreaseInvalidEntryCount", testIrIncreaseInvalidEntryCount)
}

func testIrIncreaseInvalidEntryCount(t *testing.T) {

	ir := NewImporterReport()
	testItem1 := "somedirectory/sumesubdirectory/sometext.txt"
	testItem2 := "somedirectory/someimage.png"
	testItem3 := "somedoc.odt"

	ir.IncreaseInvalidEntryCount(testItem1)
	ir.IncreaseInvalidEntryCount(testItem1)
	ir.IncreaseInvalidEntryCount(testItem2)

	assertInvalidEntryCount(t, ir, testItem1, 2)
	assertInvalidEntryCount(t, ir, testItem2, 1)
	assertInvalidEntryCount(t, ir, testItem3, 0)
}

func assertInvalidEntryCount(t *testing.T, ir *ImporterReport, filename string, expectedCount int) {

	actualCount := ir.GetInvalidEntryCount(filename)
	if actualCount != expectedCount {
		t.Errorf("Wrong invalid entry count for file %s: %d. Expected: %d.", filename, actualCount, expectedCount)
	}
}
