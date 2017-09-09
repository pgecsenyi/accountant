package dal

import (
	"container/list"
	"testing"
	"util"
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

func TestAddFingerprint(t *testing.T) {

	checksum := []byte{12, 23, 34, 45}
	fingerprint := &Fingerprint{"simple.txt", checksum, "sha1", "", "", ""}
	memoryDatabase := NewMemoryDatabase()

	memoryDatabase.AddFingerprint(fingerprint)
	actualFingerprints := memoryDatabase.GetFingerprints()

	assertStoredFingerprintsAreValid(t, actualFingerprints)
}

func TestAddNamePair(t *testing.T) {

	namePair := &NamePair{"new", "old"}
	memoryDatabase := NewMemoryDatabase()

	memoryDatabase.AddNamePair(namePair)
	actualNamePairs := memoryDatabase.GetNamePairs()

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

func TestLoadNamesFromFingeprints(t *testing.T) {

	fingerprint := &Fingerprint{"simple.txt", nil, "", "", "", ""}
	otfReadTester := newOnTheFlyFingerprintReadTester(t)
	memoryDatabase := NewMemoryDatabase()
	memoryDatabase.AddFingerprint(fingerprint)

	memoryDatabase.LoadNamesFromFingeprints(otfReadTester)
}

func TestSetFingerprint(t *testing.T) {

	checksum := []byte{12, 23, 34, 45}
	fingerprint := &Fingerprint{"simple.txt", checksum, "sha1", "", "", ""}
	fingerprints := list.New()
	fingerprints.PushFront(fingerprint)
	memoryDatabase := NewMemoryDatabase()

	memoryDatabase.SetFingerprints(fingerprints)
	actualFingerprints := memoryDatabase.GetFingerprints()

	assertStoredFingerprintsAreValid(t, actualFingerprints)
}

func assertStoredFingerprintsAreValid(t *testing.T, actualFingerprints *list.List) {

	if actualFingerprints.Len() != 1 {
		t.Error("There is a wrong number of Fingerprints in the list.")
	}
	actualFingerprint := actualFingerprints.Front().Value.(*Fingerprint)
	assertStoredFingerprintIsValid(t, actualFingerprint)
}

func assertStoredFingerprintIsValid(t *testing.T, actualFingerprint *Fingerprint) {

	expectedChecksum := []byte{12, 23, 34, 45}

	if actualFingerprint.Filename != "simple.txt" ||
		actualFingerprint.Algorithm != "sha1" ||
		!util.CompareByteSlices(expectedChecksum, actualFingerprint.Checksum) {
		t.Error("Wrong Fingerprint is in the database.")
	}
}
