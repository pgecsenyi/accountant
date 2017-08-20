package testutil

import (
	"container/list"
	"dal"
	"encoding/hex"
	"testing"
	"util"
)

// CreateSparseFingerprint Creates a fingerprint with the given name, checksum and algorithm.
func CreateSparseFingerprint(filename string, checksum string, algorithm string) *dal.Fingerprint {

	return CreateFingerprint(filename, checksum, algorithm, "", "", "")
}

// CreateFingerprint Creates a fingerprint with the given fields.
func CreateFingerprint(
	filename string, checksum string,
	algorithm string, createdAt string,
	creator string, note string) *dal.Fingerprint {

	checksumBytes, err := hex.DecodeString(checksum)
	util.CheckErr(err, "Unable to convert checksum from string.")

	return &dal.Fingerprint{filename, checksumBytes, algorithm, createdAt, creator, note}
}

// CreateList Creates a list containing the given items.
func CreateList(items ...interface{}) *list.List {

	result := list.New()
	for _, item := range items {
		result.PushFront(item)
	}

	return result
}

// GetExpectedFingerprintsForBasicCalculation Creates fingerprints for testing calculation logic.
func GetExpectedFingerprintsForBasicCalculation() *list.List {

	fp1 := CreateSparseFingerprint("test.txt", "1c291ca3", "crc32")
	fp2 := CreateSparseFingerprint("dir1/test.txt", "6b24cc6a", "crc32")
	expectedFingerprints := CreateList(fp1, fp2)

	return expectedFingerprints
}

// AssertContainsFingerprints Checks whether the given fingerprints are in the given list. Filtering can be customized
// using the "fieldsToCheck" parameter.
func AssertContainsFingerprints(
	t *testing.T, fingerprints *list.List,
	expectedFingerprints *list.List, fieldsToCheck FingerprintFieldsToCheck) {

	for element := fingerprints.Front(); element != nil; element = element.Next() {
		fingerprint := element.Value.(*dal.Fingerprint)
		if !isFingerprintInList(fingerprint, expectedFingerprints, fieldsToCheck) {
			t.Errorf("Fingerprint for file \"%s\" is unexpected.", fingerprint.Filename)
		}
	}
}

func isFingerprintInList(
	fingerprint *dal.Fingerprint, fingerprints *list.List,
	fieldsToCheck FingerprintFieldsToCheck) bool {

	for element := fingerprints.Front(); element != nil; element = element.Next() {
		fp := element.Value.(*dal.Fingerprint)

		if (!fieldsToCheck.Filename || fingerprint.Filename == fp.Filename) &&
			(!fieldsToCheck.Algorithm || fingerprint.Algorithm == fp.Algorithm) &&
			(!fieldsToCheck.CreatedAt || fingerprint.CreatedAt == fp.CreatedAt) &&
			(!fieldsToCheck.Creator || fingerprint.Creator == fp.Creator) &&
			(!fieldsToCheck.Note || fingerprint.Note == fp.Note) &&
			(!fieldsToCheck.Checksum || util.CompareByteSlices(fingerprint.Checksum, fp.Checksum)) {
			return true
		}
	}

	return false
}
