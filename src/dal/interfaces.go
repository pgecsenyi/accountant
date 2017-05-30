package dal

import (
	"container/list"
	"util"
)

// Database Interface for database implementations.
type Database interface {
	AddFingerprint(fingerprint *Fingerprint)
	AddNamePair(namePair *NamePair)
	GetFingerprints() *list.List
	LoadFingerprints()
	LoadNamesFromFingeprints(writer util.StringWriter)
	SaveFingerprints()
	SaveNamePairs()
	SetFingerprints(*list.List)
}
