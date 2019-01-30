package bll

import (
	"container/list"
	"fmr/bll/testutil"
	"fmr/dal"
	"fmr/util"
	"testing"
)

func TestComparer(t *testing.T) {

	setupComparerTests()

	t.Run("Compare_AllFields", testComparerCompareAllFields)
	t.Run("Compare_NewAndMissingFiles", testComparerCompareNewAndMissingFiles)

	tearDownComparerTests()
}

func setupComparerTests() {

	testHelper.CreateTestRootDirectory()

	testHelper.CreateTestDirectory("alldata")
	testHelper.CreateTestDirectory("alldata/dir1")
	testHelper.CreateTestFileWithContent("alldata/test2.txt", "Hello World!")
	testHelper.CreateTestFileWithContent("alldata/orange.txt", "Go is an open source programming language")
	testHelper.CreateTestFileWithContent("alldata/dir1/test.txt", "Lorem ipsum, dolor sit amet.")

	testHelper.CreateTestDirectory("newandmissing")
	testHelper.CreateTestDirectory("newandmissing/dir1")
	testHelper.CreateTestFileWithContent("newandmissing/test2.txt", "Hello World!")
	testHelper.CreateTestFileWithContent("newandmissing/test3.txt", "Go is an open source programming language")
	testHelper.CreateTestFileWithContent("newandmissing/dir1/test.txt", "Lorem ipsum, dolor sit amet.")
}

func testComparerCompareAllFields(t *testing.T) {

	// Arrange.
	expectedFingerprints := getExpectedFingerprintsForAllFieldsComparison()
	fieldsToCheck := testutil.NewFingerprintFieldsToCheck(true, true, true)
	memoryDatabase := getBaselineDatabaseForAllFieldsComparison()
	testPath := testHelper.GetTestDirectory("alldata")
	comparer := NewComparer(memoryDatabase, testPath, testPath)

	// Act.
	comparer.Compare("crc32")

	// Assert.
	testutil.AssertContainsFingerprints(t, memoryDatabase.GetFingerprints(), expectedFingerprints, fieldsToCheck)
}

func testComparerCompareNewAndMissingFiles(t *testing.T) {

	// Arrange.
	expectedFingerprints := getExpectedFingerprintsForNewAndMissingComparison()
	fieldsToCheck := testutil.NewFingerprintFieldsToCheck(false, false, false)
	memoryDatabase := getBaselineDatabaseForNewAndMissingComparison()
	testPath := testHelper.GetTestDirectory("newandmissing")
	comparer := NewComparer(memoryDatabase, testPath, testPath)

	// Act.
	comparer.Compare("crc32")

	// Assert.
	testutil.AssertContainsFingerprints(t, memoryDatabase.GetFingerprints(), expectedFingerprints, fieldsToCheck)
	assertComparerMissingFiles(t, &comparer)
	assertComparerNewFiles(t, &comparer)
}

func tearDownComparerTests() {

	testHelper.CleanUp()
}

func getExpectedFingerprintsForAllFieldsComparison() *list.List {

	fp1 := testutil.CreateFingerprint(
		"test2.txt", "1c291ca3", "crc32", "2011-09-19T08:04:45Z", "Double Commander", "Lorem ipsum")
	fp2 := testutil.CreateFingerprint(
		"orange.txt", "1881d07b", "crc32", "2016-05-14T20:31:32+03:00", util.RuntimeVersion, "")
	fp3 := testutil.CreateFingerprint(
		"dir1/test.txt", "6b24cc6a", "crc32", "2017-04-28T09:57:26Z", "Total Commander", "")
	expectedFingerprints := testutil.CreateList(fp1, fp2, fp3)

	return expectedFingerprints
}

func getBaselineDatabaseForAllFieldsComparison() *dal.MemoryDatabase {

	fp1 := testutil.CreateFingerprint(
		"test.txt", "1c291ca3", "crc32", "2011-09-19T08:04:45Z", "Double Commander", "Lorem ipsum")
	fp2 := testutil.CreateFingerprint(
		"dir1/test.txt", "6b24cc6a", "crc32", "2017-04-28T09:57:26Z", "Total Commander", "")
	fp3 := testutil.CreateFingerprint(
		"apple.txt", "1881d07b", "crc32", "2016-05-14T20:31:32+03:00", util.RuntimeVersion, "")

	memoryDatabase := dal.NewMemoryDatabase()
	memoryDatabase.AddFingerprint(fp1)
	memoryDatabase.AddFingerprint(fp2)
	memoryDatabase.AddFingerprint(fp3)

	return memoryDatabase
}

func getExpectedFingerprintsForNewAndMissingComparison() *list.List {

	fp1 := testutil.CreateSparseFingerprint("test2.txt", "1c291ca3", "crc32")
	fp2 := testutil.CreateSparseFingerprint("test3.txt", "1881d07b", "crc32")
	fp3 := testutil.CreateSparseFingerprint("dir1/test.txt", "6b24cc6a", "crc32")
	expectedFingerprints := testutil.CreateList(fp1, fp2, fp3)

	return expectedFingerprints
}

func getBaselineDatabaseForNewAndMissingComparison() *dal.MemoryDatabase {

	fp1 := testutil.CreateSparseFingerprint("test.txt", "1c291ca3", "crc32")
	fp2 := testutil.CreateSparseFingerprint("dir1/test.txt", "6b24cc6a", "crc32")
	fp3 := testutil.CreateSparseFingerprint("some-deleted-file", "f32ab44c", "crc32")

	memoryDatabase := dal.NewMemoryDatabase()
	memoryDatabase.AddFingerprint(fp1)
	memoryDatabase.AddFingerprint(fp2)
	memoryDatabase.AddFingerprint(fp3)

	return memoryDatabase
}

func assertComparerMissingFiles(t *testing.T, comparer *Comparer) {

	if !testHelper.HasStringItems(comparer.Report.MissingFiles, "some-deleted-file") {
		t.Error("File should be marked as missing: \"some-deleted-file\".")
	}
}

func assertComparerNewFiles(t *testing.T, comparer *Comparer) {

	if !testHelper.HasStringItems(comparer.Report.NewFiles, "test3.txt") {
		t.Error("File should be marked as new: \"test3.txt\".")
	}
}
