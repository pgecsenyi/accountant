package dal

import (
	"container/list"
	"util"
)

// Database Interface for database implementations.
type Database interface {
	AddFingerprint(fingerprint *Fingerprint)
	GetFingerprints() *list.List
	Load()
	LoadNames(writer util.StringWriter)
	Save()
	SetFingerprints(*list.List)
}
