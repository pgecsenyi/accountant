package common

import (
	"dal"
	"testing"
)

func TestFingerprintFilter(t *testing.T) {

	t.Run("EmptyFilter", testEmptyFilter)
	t.Run("FilenameFilter_Match_Extension", testFilenameFilterMatchExtension)
	t.Run("FilenameFilter_Match_Name", testFilenameFilterMatchName)
	t.Run("FilenameFilter_Match_NameNoAlgorithm", testFilenameFilterMatchNameNoAlgorithm)
	t.Run("FilenameFilter_NoMatch", testFilenameFilterNoMatch)
	t.Run("AlgorithmFilter_Match", testAlgorithmFilterMatch)
	t.Run("AlgorithmFilter_NoMatch_DifferentAlgorithm", testAlgorithmFilterNoMatchDifferentAlgorithm)
	t.Run("AlgorithmFilter_NoMatch_NotWellDefinedAlgorithm", testAlgorithmFilterNoMatchNotWellDefinedAlgorithm)
	t.Run("AllFilters", testAllFilters)
}

func testEmptyFilter(t *testing.T) {

	fp := createFingerprintWithNameAndAlg("sample-file-with_aWeird-name.txt", "crc32")
	filter := ""
	fpFilter := NewFingerprintFilter(filter)

	result := fpFilter.FilterFingerprint(fp)

	assertMatch(t, fp.Filename, filter, true, result)
}

func testFilenameFilterMatchExtension(t *testing.T) {

	fp := createFingerprintWithNameAndAlg("sample-file-with_aWeird-name.txt", "crc32")
	filter := ".txt"
	fpFilter := NewFingerprintFilter(filter)

	result := fpFilter.FilterFingerprint(fp)

	assertMatch(t, fp.Filename, filter, true, result)
}

func testFilenameFilterMatchName(t *testing.T) {

	fp := createFingerprintWithNameAndAlg("sample-file-with_aWeird-name.txt", "crc32")
	filter := "file-with"
	fpFilter := NewFingerprintFilter(filter)

	result := fpFilter.FilterFingerprint(fp)

	assertMatch(t, fp.Filename, filter, true, result)
}

func testFilenameFilterMatchNameNoAlgorithm(t *testing.T) {

	fp := createFingerprintWithNameAndAlg("sample-file-with_aWeird-name.txt", "crc32")
	filter := "file-with:"
	fpFilter := NewFingerprintFilter(filter)

	result := fpFilter.FilterFingerprint(fp)

	assertMatch(t, fp.Filename, filter, true, result)
}

func testFilenameFilterNoMatch(t *testing.T) {

	fp := createFingerprintWithNameAndAlg("sample-file-with_aWeird-name.txt", "crc32")
	filter := "f1le-with"
	fpFilter := NewFingerprintFilter(filter)

	result := fpFilter.FilterFingerprint(fp)

	assertMatch(t, fp.Filename, filter, false, result)
}

func testAlgorithmFilterMatch(t *testing.T) {

	fp := createFingerprintWithNameAndAlg("sample-file-with_aWeird-name.txt", "sha512")
	filter := ":sha512"
	fpFilter := NewFingerprintFilter(filter)

	result := fpFilter.FilterFingerprint(fp)

	assertMatch(t, fp.Algorithm, filter, true, result)
}

func testAlgorithmFilterNoMatchDifferentAlgorithm(t *testing.T) {

	fp := createFingerprintWithNameAndAlg("sample-file-with_aWeird-name.txt", "sha512")
	filter := ":crc32"
	fpFilter := NewFingerprintFilter(filter)

	result := fpFilter.FilterFingerprint(fp)

	assertMatch(t, fp.Algorithm, filter, false, result)
}

func testAlgorithmFilterNoMatchNotWellDefinedAlgorithm(t *testing.T) {

	fp := createFingerprintWithNameAndAlg("sample-file-with_aWeird-name.txt", "sha512")
	filter := ":sha"
	fpFilter := NewFingerprintFilter(filter)

	result := fpFilter.FilterFingerprint(fp)

	assertMatch(t, fp.Algorithm, filter, false, result)
}

func testAllFilters(t *testing.T) {

	fp := createFingerprintWithNameAndAlg("sample-file-with_aWeird-name.txt", "sha512")
	filter := "sample:sha512"
	fpFilter := NewFingerprintFilter(filter)

	result := fpFilter.FilterFingerprint(fp)

	assertMatch(t, fp.Filename+" | "+fp.Algorithm, filter, true, result)
}

func createFingerprintWithNameAndAlg(filename string, algorithm string) *dal.Fingerprint {

	return &dal.Fingerprint{filename, nil, algorithm, "", "", ""}
}

func assertMatch(t *testing.T, filteredText string, filter string, shouldMatch bool, match bool) {

	if shouldMatch && !match {
		t.Errorf("\"%s\" should match filter \"%s\".", filteredText, filter)
	} else if !shouldMatch && match {
		t.Errorf("\"%s\" should not match filter \"%s\".", filteredText, filter)
	}
}
