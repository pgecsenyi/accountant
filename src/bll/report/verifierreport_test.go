package report

import (
	"container/list"
	"testing"
	"util"
)

var verifierReportTestHelper = util.NewTestHelper()

func TestVerifierReport(t *testing.T) {

	t.Run("AddCorruptFile", testVrAddCorruptFile)
	t.Run("AddMissingFile", testVrAddMissingFile)
	t.Run("AddValidFile", testVrAddValidFile)
}

func testVrAddCorruptFile(t *testing.T) {

	vr := NewVerifierReport()
	testItem1 := "somedirectory/sumesubdirectory/somefile.txt"
	testItem2 := "someotherfile.txt"

	vr.AddCorruptFile(testItem1)
	vr.AddCorruptFile(testItem2)

	assertAllCount(t, vr, 2)
	assertListLength(t, vr.MissingFiles, "missing", 0)
	if !verifierReportTestHelper.HasStringItems(vr.CorruptFiles, testItem1, testItem2) {
		t.Error("The list of corrupt files is incomplete.")
	}
}

func testVrAddMissingFile(t *testing.T) {

	vr := NewVerifierReport()
	testItem1 := "somedirectory/sumesubdirectory/somefile.txt"
	testItem2 := "someotherfile.txt"

	vr.AddMissingFile(testItem1)
	vr.AddMissingFile(testItem2)

	assertAllCount(t, vr, 2)
	assertListLength(t, vr.CorruptFiles, "corrupt", 0)
	if !verifierReportTestHelper.HasStringItems(vr.MissingFiles, testItem1, testItem2) {
		t.Error("The list of missing files is incomplete.")
	}
}

func testVrAddValidFile(t *testing.T) {

	vr := NewVerifierReport()
	testItem1 := "somedirectory/sumesubdirectory/somefile.txt"
	testItem2 := "someotherfile.txt"

	vr.AddValidFile(testItem1)
	vr.AddValidFile(testItem2)

	assertAllCount(t, vr, 2)
	assertListLength(t, vr.CorruptFiles, "corrupt", 0)
	assertListLength(t, vr.MissingFiles, "missing", 0)
}

func assertAllCount(t *testing.T, vr *VerifierReport, expectedCount int) {

	if vr.CountAll != expectedCount {
		t.Errorf("Wrong number of all files: %d (expected: %d).", vr.CountAll, expectedCount)
	}
}

func assertListLength(t *testing.T, l *list.List, name string, expectedLength int) {

	if l.Len() != 0 {
		t.Errorf("Wrong number of %s files: %d (expected: %d).", name, l.Len(), expectedLength)
	}
}
