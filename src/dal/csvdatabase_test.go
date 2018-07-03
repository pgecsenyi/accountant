package dal

import (
	"testing"
)

func TestCsvDatabase(t *testing.T) {

	setupCsvDatabaseTests()

	t.Run("CsvDatabase_AddFingerprint", testCsvDatabaseAddFingerprint)
	t.Run("CsvDatabase_AddFingerprints", testCsvDatabaseAddFingerprints)
	t.Run("CsvDatabase_AddNamePair", testCsvDatabaseAddNamePair)
	t.Run("CsvDatabase_Clear", testCsvDatabaseClear)
	t.Run("CsvDatabase_LoadNamesFromFingerprints", testCsvDatabaseLoadNamesFromFingerprints)
	t.Run("CsvDatabase_SaveAndLoadFingerprints", testCsvDatabaseSaveAndLoadFingerprints)

	tearDownCsvDatabaseTests()
}

func setupCsvDatabaseTests() {

	testHelper.CreateTestRootDirectory()
}

func tearDownCsvDatabaseTests() {

	testHelper.CleanUp()
}

func testCsvDatabaseAddFingerprint(t *testing.T) {

	csvDatabase := NewCsvDatabase(
		testHelper.GetTestPath("fingerprints.csv"),
		testHelper.GetTestPath("fingerprints.csv"),
		testHelper.GetTestPath("namepairs.fm"))
	testDatabaseAddFingerprint(t, csvDatabase)
}

func testCsvDatabaseAddFingerprints(t *testing.T) {

	csvDatabase := NewCsvDatabase(
		testHelper.GetTestPath("fingerprints.csv"),
		testHelper.GetTestPath("fingerprints.csv"),
		testHelper.GetTestPath("namepairs.fm"))
	testDatabaseAddFingerprints(t, csvDatabase)
}

func testCsvDatabaseAddNamePair(t *testing.T) {

	csvDatabase := NewCsvDatabase(
		testHelper.GetTestPath("fingerprints.csv"),
		testHelper.GetTestPath("fingerprints.csv"),
		testHelper.GetTestPath("namepairs.fm"))
	testDatabaseAddNamePair(t, csvDatabase)
}

func testCsvDatabaseClear(t *testing.T) {

	csvDatabase := NewCsvDatabase(
		testHelper.GetTestPath("fingerprints.csv"),
		testHelper.GetTestPath("fingerprints.csv"),
		testHelper.GetTestPath("namepairs.fm"))
	testDatabaseClear(t, csvDatabase)
}

func testCsvDatabaseLoadNamesFromFingerprints(t *testing.T) {

	csvDatabase := NewCsvDatabase(
		testHelper.GetTestPath("fingerprints.csv"),
		testHelper.GetTestPath("fingerprints.csv"),
		testHelper.GetTestPath("namepairs.fm"))
	testDatabaseLoadNamesFromFingerprints(t, csvDatabase)
}

func testCsvDatabaseSaveAndLoadFingerprints(t *testing.T) {

	checksum := []byte{12, 23, 34, 45}
	fingerprint := &Fingerprint{"simple.txt", checksum, "sha1", "", "", ""}
	csvDatabase := NewCsvDatabase(
		testHelper.GetTestPath("fingerprints.csv"),
		testHelper.GetTestPath("fingerprints.csv"),
		testHelper.GetTestPath("namepairs.fm"))

	csvDatabase.AddFingerprint(fingerprint)
	csvDatabase.SaveFingerprints()
	csvDatabase.Clear()
	csvDatabase.LoadFingerprints()
	actualFingerprints := csvDatabase.GetFingerprints()

	assertStoredFingerprintIsValid(t, actualFingerprints)
}
