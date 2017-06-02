package bll

import (
	"container/list"
	"dal"
	"encoding/hex"
	"util"
)

// Comparer Stores settings related to comparison.
type Comparer struct {
	Db             dal.Database
	InputDirectory string
	BasePath       string
	Report         *ComparerReport
}

// NewComparer Instantiates a new Comparer object.
func NewComparer(db dal.Database, inputDirectory string, basePath string) Comparer {

	comparerReport := NewComparerReport()

	return Comparer{db, inputDirectory, basePath, comparerReport}
}

// Compare Verifies and stores changes in the given directory based on the checksums calculated earlier.
func (comparer *Comparer) Compare(algorithm string) {

	oldFingerprints := comparer.loadOldFingerprints()
	newFingerprints := comparer.calculateNewFingerprints(algorithm)

	comparer.Db.SetFingerprints(newFingerprints)
	comparer.Db.SaveFingerprints()

	comparer.compareAndSaveResults(oldFingerprints, newFingerprints)
	comparer.Db.SaveNamePairs()
}

func (comparer *Comparer) loadOldFingerprints() *list.List {

	comparer.Db.LoadFingerprints()
	oldFingerprints := comparer.Db.GetFingerprints()

	return oldFingerprints
}

func (comparer *Comparer) calculateNewFingerprints(algorithm string) *list.List {

	hasher := NewHasher(algorithm)
	effectiveBasePath := comparer.getEffectiveBasePath()
	files := util.ListFilesRecursively(comparer.InputDirectory)
	newFingerprints := hasher.CalculateFingerprints(comparer.InputDirectory, effectiveBasePath, files)

	return newFingerprints
}

func (comparer *Comparer) getEffectiveBasePath() string {

	return util.TrimPath(comparer.InputDirectory, comparer.BasePath)
}

func (comparer *Comparer) compareAndSaveResults(oldFingerprints *list.List, newFingerprints *list.List) {

	cache := buildFingerprintCache(oldFingerprints)
	foundFingerprints := make(map[string]bool)

	for element := newFingerprints.Front(); element != nil; element = element.Next() {
		fingerprint := element.Value.(*dal.Fingerprint)
		checksum := hex.EncodeToString(fingerprint.Checksum)
		matchingFingerprint := cache[checksum]
		comparer.processMatch(fingerprint, checksum, matchingFingerprint, foundFingerprints)
	}

	comparer.collectMissingFiles(cache, foundFingerprints)
}

func (comparer *Comparer) processMatch(
	fingerprint *dal.Fingerprint, checksum string,
	matchingFingerprint *dal.Fingerprint, foundFingerprints map[string]bool) {

	if matchingFingerprint == nil {
		comparer.Report.AddNewFile(fingerprint.Filename)
	} else {
		comparer.Db.AddNamePair(&dal.NamePair{fingerprint.Filename, matchingFingerprint.Filename})
		fingerprint.CreatedAt = matchingFingerprint.CreatedAt
		foundFingerprints[checksum] = true
	}
}

func (comparer *Comparer) collectMissingFiles(oldFingerprints map[string]*dal.Fingerprint, foundFingerprints map[string]bool) {

	for hash, fingerprint := range oldFingerprints {
		hasFound := foundFingerprints[hash]
		if !hasFound {
			comparer.Report.AddMissingFile(fingerprint.Filename)
		}
	}
}

func buildFingerprintCache(fingerprints *list.List) map[string]*dal.Fingerprint {

	var cache = make(map[string]*dal.Fingerprint)

	for element := fingerprints.Front(); element != nil; element = element.Next() {
		fingerprint := element.Value.(*dal.Fingerprint)
		checksum := hex.EncodeToString(fingerprint.Checksum)
		cache[checksum] = fingerprint
	}

	return cache
}
