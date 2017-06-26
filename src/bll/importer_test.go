package bll

import (
	"bll/testutil"
	"container/list"
	"dal"
	"testing"
)

func TestImporter(t *testing.T) {

	setupImporterTests()

	t.Run("Convert", testImporterConvert)

	tearDownImporterTests()
}

func setupImporterTests() {

	testHelper.CreateTestRootDirectory()

	testHelper.CreateTestDirectory("subdir")
	testHelper.CreateTestFileWithContent(
		"test.sha",
		"15dfaa952a85ad9a458013fa2fc3bdc807d34e7f *textfile.txt"+
			"\r\n1a0041decc7147a86a01652e92a9027775d472c4 *presentation.odp")
	testHelper.CreateTestFileWithContent("import/almost-empty.md5", " ")
	testHelper.CreateTestFileWithContent(
		"subdir/md5.md5",
		"845178f3c9e7ec71f23e01e2187a1867  compressed.tar.gz\n8d5f2e17f783cc066de6e02adc74566e  executable")
	testHelper.CreateTestFileWithContent(
		"subdir/something.sfv",
		"; some header\r\nanimage.jpg AFB25773\r\nanotherimage.png D7B3144F")
	testHelper.CreateTestFileWithContent(
		"CHK.sha256",
		"357ad3058f7b5b71e0488df08ed1f6dfcdde722f298bdd9a903b1c8121d9db50 *source.c")
	testHelper.CreateTestFileWithContent(
		"subdir/sh.sha512",
		"312c3581a742881b03a7b8f4311a67744e36152a6494806046154e005cd4230a9c7c439e273c4ab811e897f97bf92fa4136bab895b101c8792a7f0e05ecf5d41 *important.odt"+
			"\r\n3c3581a742881b03a7b8f4311a67744e36152a6494806046154e005cd4230a9c7c439e273c4ab811e897f97bf92fa4136bab895b101c8792a7f0e05ecf5 *wrongentry.7z"+
			"\r\nsomething"+
			"\r\n312s3581a742881b03a7b8f4311a67744e36152a6494806046154e005cd4230a9c7c439e273c4ab811e897f97bf92fa4136bab895b101c8792a7f0e05ecf5d41 *another_wrong_entry.bad")
}

func testImporterConvert(t *testing.T) {

	// Arrange.
	expectedFingerprints := getExpectedFingerprintsForImport()
	memoryDatabase := dal.NewMemoryDatabase()
	testPath := testHelper.GetTestRootDirectory()
	outputChecksums := testHelper.GetTestPath("out.csv")
	importer := NewImporter(memoryDatabase, testPath, outputChecksums)

	// Act.
	importer.Convert()

	// Assert.
	if memoryDatabase.GetFingerprints().Len() != 8 {
		t.Errorf("Wrong number of database entries: %d (expected: %d).", memoryDatabase.GetFingerprints().Len(), 8)
	}
	testFile := testHelper.GetTestPath("subdir/sh.sha512")
	if importer.Report.GetInvalidEntryCount(testFile) != 3 {
		t.Errorf("Wrong number of invalid entries for file \"%s\".", testFile)
	}
	testutil.AssertContainsFingerprints(t, memoryDatabase.GetFingerprints(), expectedFingerprints)
}

func tearDownImporterTests() {

	testHelper.CleanUp()
}

func getExpectedFingerprintsForImport() *list.List {

	fp1 := testutil.CreateFingerprint("textfile.txt", "15dfaa952a85ad9a458013fa2fc3bdc807d34e7f", "sha1")
	fp2 := testutil.CreateFingerprint("presentation.odp", "1a0041decc7147a86a01652e92a9027775d472c4", "sha1")
	fp3 := testutil.CreateFingerprint("compressed.tar.gz", "845178f3c9e7ec71f23e01e2187a1867", "md5")
	fp4 := testutil.CreateFingerprint("executable", "8d5f2e17f783cc066de6e02adc74566e", "md5")
	fp5 := testutil.CreateFingerprint("animage.jpg", "afb25773", "crc32")
	fp6 := testutil.CreateFingerprint("anotherimage.png", "d7b3144f", "crc32")
	fp7 := testutil.CreateFingerprint("source.c", "357ad3058f7b5b71e0488df08ed1f6dfcdde722f298bdd9a903b1c8121d9db50", "sha256")
	fp8 := testutil.CreateFingerprint(
		"important.odt",
		"312c3581a742881b03a7b8f4311a67744e36152a6494806046154e005cd4230a9c7c439e273c4ab811e897f97bf92fa4136bab895b101c8792a7f0e05ecf5d41",
		"sha512")
	expectedFingerprints := testutil.CreateList(fp1, fp2, fp3, fp4, fp5, fp6, fp7, fp8)

	return expectedFingerprints
}
