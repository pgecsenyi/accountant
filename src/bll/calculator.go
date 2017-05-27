package bll

import (
	"container/list"
	"dal"
	"hash/fnv"
	"path"
	"util"
)

// Calculator Stores settings related to checksum calculation.
type Calculator struct {
	InputDirectory  string
	OutputChecksums string
	BasePath        string
	InputChecksums  string
}

// RecordChecksumsForDirectory Calculates and stores checksums for the files in the given directory.
func (calculator *Calculator) RecordChecksumsForDirectory(db *dal.Db, algorithm string) {

	files := util.ListDirectoryRecursively(calculator.InputDirectory)
	hasher := NewHasher(algorithm)
	effectiveBasePath := calculator.getEffectiveBasePath()
	var fingerprints *list.List
	if calculator.InputChecksums == "" {
		fingerprints = hasher.CalculateChecksumsForFiles(calculator.InputDirectory, effectiveBasePath, files)
	} else {
		fingerprints = calculator.calculateChecksumsForMissingFiles(db, hasher, effectiveBasePath, files)
	}
	db.Fingerprints = fingerprints
	db.SaveCsv(calculator.OutputChecksums)
}

func (calculator *Calculator) getEffectiveBasePath() string {

	return util.TrimPath(calculator.InputDirectory, calculator.BasePath)
}

func (calculator *Calculator) calculateChecksumsForMissingFiles(
	db *dal.Db, hasher Hasher, effectiveBasePath string, files []string) *list.List {

	fingerprints := list.New()
	nhs := newNameHashStorage(fnv.New32a())
	db.LoadNamesFromCsv(calculator.InputChecksums, nhs)
	nhs.ClearCache()
	for _, file := range files {
		fullPath := path.Join(effectiveBasePath, file)
		if !nhs.ContainsName(fullPath) {
			hasher.recordChecksumForFile(calculator.InputDirectory, effectiveBasePath, file, fingerprints)
		}
	}
	nhs.ClearCache()

	return fingerprints
}
