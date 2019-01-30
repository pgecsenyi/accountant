package dal

import (
	"container/list"
	"fmr/util"
)

// Database Interface for database implementations.
type Database interface {
	AddFingerprint(fingerprint *Fingerprint)
	AddFingerprints(fingerprints *list.List)
	AddNamePair(namePair *NamePair)
	Clear()
	GetFingerprints() *list.List
	GetNamePairs() *list.List
	LoadFingerprints()
	LoadNamesFromFingeprints(writer util.StringWriter)
	SaveFingerprints()
	SaveNamePairs()
}
