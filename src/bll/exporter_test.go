package bll

import (
	"container/list"
	"fmr/bll/common"
	"fmr/bll/testutil"
	"fmr/dal"
	"testing"
)

var exportTestFingerprint1 = testutil.CreateSparseFingerprint(
	"textfile.txt", "15dfaa952a85ad9a458013fa2fc3bdc807d34e7f", "sha1")
var exportTestFingerprint2 = testutil.CreateSparseFingerprint(
	"compressed.tar.gz", "845178f3c9e7ec71f23e01e2187a1867", "md5")
var exportTestFingerprint3 = testutil.CreateSparseFingerprint(
	"animage.jpg", "afb25773", "crc32")
var exportTestFingerprint4 = testutil.CreateSparseFingerprint(
	"source.c", "357ad3058f7b5b71e0488df08ed1f6dfcdde722f298bdd9a903b1c8121d9db50", "sha256")
var exportTestFingerprint5 = testutil.CreateSparseFingerprint(
	"source2.c",
	"312c3581a742881b03a7b8f4311a67744e36152a6494806046154e005cd4230a9c7c439e273c4ab811e897f97bf92fa4136bab895b101c8792a7f0e05ecf5d41",
	"sha512")

func TestExporter(t *testing.T) {

	setupExporterTests()

	t.Run("Convert", testExporterConvert)
	t.Run("Convert_EmptyFilter", testExporterConvertEmptyFilter)
	t.Run("Convert_NameFilter", testExporterConvertFilterName)
	t.Run("Convert_NameAlgFilter", testExporterConvertFilterNameAlg)

	tearDownExporterTests()
}

func setupExporterTests() {

	testHelper.CreateTestRootDirectory()

	testHelper.CreateTestDirectory("tmp")
}

func testExporterConvert(t *testing.T) {

	fpFilter := common.NewFingerprintFilter("")
	expectedFingerprints := getFingerprintsToExport()
	testExporterWithFilter(t, fpFilter, expectedFingerprints)
}

func testExporterConvertEmptyFilter(t *testing.T) {

	fpFilter := common.NewFingerprintFilter(":")
	expectedFingerprints := getFingerprintsToExport()
	testExporterWithFilter(t, fpFilter, expectedFingerprints)
}

func testExporterConvertFilterName(t *testing.T) {

	fpFilter := common.NewFingerprintFilter("source")
	expectedFingerprints := getExpectedFingerprintsForExportWithNameFilter()
	testExporterWithFilter(t, fpFilter, expectedFingerprints)
}

func testExporterConvertFilterNameAlg(t *testing.T) {

	fpFilter := common.NewFingerprintFilter("source:sha256")
	expectedFingerprints := getExpectedFingerprintsForExportWithNameAlgFilter()
	testExporterWithFilter(t, fpFilter, expectedFingerprints)
}

func tearDownExporterTests() {

	testHelper.CleanUp()
}

func testExporterWithFilter(t *testing.T, fpFilter common.FingerprintFilter, expectedFingerprints *list.List) {

	// Arrange.
	memoryDatabase1 := dal.NewMemoryDatabase()
	memoryDatabase2 := dal.NewMemoryDatabase()
	memoryDatabase1.AddFingerprints(getFingerprintsToExport())
	fieldsToCheck := testutil.NewFingerprintFieldsToCheck(false, false, false)
	testPath := testHelper.GetTestPath("tmp")
	outputChecksums := testHelper.GetTestPath("out.csv")
	exporter := NewExporter(memoryDatabase1, testPath, "")
	importer := NewImporter(memoryDatabase2, testPath, outputChecksums)

	// Act.
	exporter.Convert(fpFilter)
	importer.Convert()
	testHelper.RemoveTestDirectory("tmp")
	testHelper.CreateTestDirectory("tmp")

	// Assert.
	if memoryDatabase2.GetFingerprints().Len() != expectedFingerprints.Len() {
		t.Errorf(
			"Wrong number of database entries: %d (expected: %d).",
			memoryDatabase2.GetFingerprints().Len(),
			expectedFingerprints.Len())
	}
	assertInvalidEntryCountInExportedFile(t, importer, "Checksum.crc32")
	assertInvalidEntryCountInExportedFile(t, importer, "Checksum.md5")
	assertInvalidEntryCountInExportedFile(t, importer, "Checksum.sha1")
	assertInvalidEntryCountInExportedFile(t, importer, "Checksum.sha256")
	assertInvalidEntryCountInExportedFile(t, importer, "Checksum.sha512")
	testutil.AssertContainsFingerprints(t, memoryDatabase2.GetFingerprints(), expectedFingerprints, fieldsToCheck)
}

func getFingerprintsToExport() *list.List {

	expectedFingerprints := testutil.CreateList(
		exportTestFingerprint1,
		exportTestFingerprint2,
		exportTestFingerprint3,
		exportTestFingerprint4,
		exportTestFingerprint5)

	return expectedFingerprints
}

func getExpectedFingerprintsForExportWithNameFilter() *list.List {

	expectedFingerprints := testutil.CreateList(exportTestFingerprint4, exportTestFingerprint5)

	return expectedFingerprints
}

func getExpectedFingerprintsForExportWithNameAlgFilter() *list.List {

	expectedFingerprints := testutil.CreateList(exportTestFingerprint4)

	return expectedFingerprints
}

func assertInvalidEntryCountInExportedFile(t *testing.T, importer Importer, filename string) {

	if importer.Report.GetInvalidEntryCount(filename) != 0 {
		t.Errorf("There should not be any invalid entries in the exported file (\"%s\").", filename)
	}
}
