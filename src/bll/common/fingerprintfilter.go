package common

import (
	"dal"
	"strings"
)

// FingerprintFilter Stores Fingerprint filter criteria.
type FingerprintFilter struct {
	filenameFilter  string
	algorithmFilter string
}

// NewFingerprintFilter Instantiates a new FingerprintFilter object.
func NewFingerprintFilter(filter string) FingerprintFilter {

	filenameFilter, algorithmFilter := getFilterParts(filter)

	return FingerprintFilter{filenameFilter, algorithmFilter}
}

// FilterFingerprint Checks whether the given Fingerprint object matches the saved filters.
func (ff *FingerprintFilter) FilterFingerprint(fingerprint *dal.Fingerprint) bool {

	return (ff.filenameFilter == "" || strings.Contains(fingerprint.Filename, ff.filenameFilter)) &&
		(ff.algorithmFilter == "" || fingerprint.Algorithm == ff.algorithmFilter)
}

func getFilterParts(filter string) (string, string) {

	if filter == "" {
		return "", ""
	}

	separatorIndex := strings.Index(filter, ":")
	if separatorIndex == -1 {
		return filter, ""
	}

	return filter[:separatorIndex], filter[separatorIndex+1:]
}
