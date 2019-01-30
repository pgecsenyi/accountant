package dal

import (
	"container/list"
	"fmr/util"
	"testing"
)

type onTheFlyFingerprintReadTester struct {
	t       *testing.T
	counter int
}

func newOnTheFlyFingerprintReadTester(t *testing.T) *onTheFlyFingerprintReadTester {
	return &onTheFlyFingerprintReadTester{t, 0}
}

// Write Memorizes the given string.
func (otfReadTester *onTheFlyFingerprintReadTester) Write(text string) {

	if otfReadTester.counter == 1 {
		otfReadTester.t.Error("More than one Fingerprint is in the database.")
	}
	if text != "simple.txt" {
		otfReadTester.t.Error("Wrong Fingerprint name is in the database.")
	}
	otfReadTester.counter++
}

var testHelper = util.NewTestHelper()

func assertStoredFingerprintIsValid(t *testing.T, actualFingerprints *list.List) {

	if actualFingerprints.Len() != 1 {
		t.Error("There is a wrong number of Fingerprints in the list.")
	}

	actualFingerprint := actualFingerprints.Front().Value.(*Fingerprint)
	expectedChecksum := []byte{12, 23, 34, 45}

	if actualFingerprint.Filename != "simple.txt" ||
		actualFingerprint.Algorithm != "sha1" ||
		!util.CompareByteSlices(expectedChecksum, actualFingerprint.Checksum) {
		t.Error("Wrong Fingerprint is in the database.")
	}
}

func testDatabaseAddFingerprint(t *testing.T, database Database) {

	checksum := []byte{12, 23, 34, 45}
	fingerprint := &Fingerprint{"simple.txt", checksum, "sha1", "", "", ""}

	database.AddFingerprint(fingerprint)
	actualFingerprints := database.GetFingerprints()

	assertStoredFingerprintIsValid(t, actualFingerprints)
}

func testDatabaseAddFingerprints(t *testing.T, database Database) {

	checksum := []byte{12, 23, 34, 45}
	fingerprint := &Fingerprint{"simple.txt", checksum, "sha1", "", "", ""}
	fingerprints := list.New()
	fingerprints.PushFront(fingerprint)

	database.AddFingerprints(fingerprints)
	actualFingerprints := database.GetFingerprints()

	assertStoredFingerprintIsValid(t, actualFingerprints)
}

func testDatabaseAddNamePair(t *testing.T, database Database) {

	namePair := &NamePair{"new", "old"}

	database.AddNamePair(namePair)
	actualNamePairs := database.GetNamePairs()

	if actualNamePairs.Len() != 1 {
		t.Errorf("The is a wrong number of NamePairs in the list.")
	}
	actualNamePair := actualNamePairs.Front().Value.(*NamePair)
	if actualNamePair.NewName != "new" {
		t.Errorf("Wrong new name: %s.", actualNamePair.NewName)
	}
	if actualNamePair.OldName != "old" {
		t.Errorf("Wrong old name: %s.", actualNamePair.OldName)
	}
}

func testDatabaseClear(t *testing.T, database Database) {

	fingerprint := &Fingerprint{"simple.txt", nil, "", "", "", ""}
	namePair := &NamePair{"apple", "orange"}

	database.AddFingerprint(fingerprint)
	database.AddNamePair(namePair)
	database.Clear()
	fingerprintCount := database.GetFingerprints().Len()
	namePairCount := database.GetNamePairs().Len()

	if fingerprintCount != 0 {
		t.Errorf("Wrong number of fingerprints in the database: %d.", fingerprintCount)
	}
	if namePairCount != 0 {
		t.Errorf("Wrong number of name pairs in the database: %d.", namePairCount)
	}
}

func testDatabaseLoadNamesFromFingerprints(t *testing.T, database Database) {

	fingerprint := &Fingerprint{"simple.txt", nil, "", "", "", ""}
	otfReadTester := newOnTheFlyFingerprintReadTester(t)

	database.AddFingerprint(fingerprint)
	database.SaveFingerprints()

	database.LoadNamesFromFingeprints(otfReadTester)
}
