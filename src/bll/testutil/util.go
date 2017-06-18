package testutil

import (
	"container/list"
	"dal"
	"encoding/hex"
	"testing"
	"util"
)

func CreateFingerprint(filename string, checksum string, algorithm string) *dal.Fingerprint {

	checksumBytes, err := hex.DecodeString(checksum)
	util.CheckErr(err, "Unable to convert checksum from string.")

	return &dal.Fingerprint{filename, checksumBytes, algorithm, "", "", ""}
}

func CreateList(items ...interface{}) *list.List {

	result := list.New()
	for _, item := range items {
		result.PushFront(item)
	}

	return result
}

func GetExpectedFingerprintsForBasicCalculation() *list.List {

	fp1 := CreateFingerprint("test.txt", "1c291ca3", "crc32")
	fp2 := CreateFingerprint("dir1/test.txt", "6b24cc6a", "crc32")
	expectedFingerprints := CreateList(fp1, fp2)

	return expectedFingerprints
}

func AssertContainsFingerprints(t *testing.T, fingerprints *list.List, expectedFingerprints *list.List) {

	for element := fingerprints.Front(); element != nil; element = element.Next() {
		fingerprint := element.Value.(*dal.Fingerprint)
		if !isFingerprintInList(fingerprint, expectedFingerprints) {
			t.Errorf("Fingerprint for file \"%s\" (filename or checksum) is unexpected.", fingerprint.Filename)
		}
	}
}

func isFingerprintInList(fingerprint *dal.Fingerprint, fingerprints *list.List) bool {

	for element := fingerprints.Front(); element != nil; element = element.Next() {
		fp := element.Value.(*dal.Fingerprint)
		if fingerprint.Filename == fp.Filename &&
			fingerprint.Algorithm == fp.Algorithm &&
			util.CompareByteSlices(fingerprint.Checksum, fp.Checksum) {
			return true
		}
	}

	return false
}
