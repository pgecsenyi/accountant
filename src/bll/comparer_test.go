package bll

import (
	"bll/testutil"
	"container/list"
	"dal"
	"testing"
)

func TestComparer(t *testing.T) {

	setupComparerTests()

	t.Run("Compare", testComparerCompare)

	tearDownComparerTests()
}

func setupComparerTests() {

	testHelper.CreateTestRootDirectory()

	testHelper.CreateTestDirectory("")
	testHelper.CreateTestDirectory("dir1")
	testHelper.CreateTestFileWithContent("test2.txt", "Hello World!")
	testHelper.CreateTestFileWithContent("test3.txt", "Go is an open source programming language")
	testHelper.CreateTestFileWithContent("dir1/test.txt", "Lorem ipsum, dolor sit amet.")
}

func testComparerCompare(t *testing.T) {

	// Arrange.
	expectedFingerprints := getExpectedFingerprintsForComparison()
	fp1 := testutil.CreateFingerprint("test.txt", "1c291ca3", "crc32")
	fp2 := testutil.CreateFingerprint("dir1/test.txt", "6b24cc6a", "crc32")
	fp3 := testutil.CreateFingerprint("some-deleted-file", "f32ab44c", "crc32")
	memoryDatabase := dal.NewMemoryDatabase()
	memoryDatabase.AddFingerprint(fp1)
	memoryDatabase.AddFingerprint(fp2)
	memoryDatabase.AddFingerprint(fp3)
	testPath := testHelper.GetTestRootDirectory()
	comparer := NewComparer(memoryDatabase, testPath, testPath)

	// Act.
	comparer.Compare("crc32")

	// Assert.
	testutil.AssertContainsFingerprints(t, memoryDatabase.Fingerprints, expectedFingerprints)
	assertComparerMissingFiles(t, &comparer)
	assertComparerNewFiles(t, &comparer)
}

func tearDownComparerTests() {

	testHelper.CleanUp()
}

func getExpectedFingerprintsForComparison() *list.List {

	fp1 := testutil.CreateFingerprint("test2.txt", "1c291ca3", "crc32")
	fp2 := testutil.CreateFingerprint("test3.txt", "1881d07b", "crc32")
	fp3 := testutil.CreateFingerprint("dir1/test.txt", "6b24cc6a", "crc32")
	expectedFingerprints := testutil.CreateList(fp1, fp2, fp3)

	return expectedFingerprints
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
