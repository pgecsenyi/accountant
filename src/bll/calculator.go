package bll

import (
	"dal"
	"util"
)

// Calculator Stores settings related to checksum calculation.
type Calculator struct {
	InputDirectory  string
	OutputChecksums string
	BasePath        string
}

// RecordChecksumsForDirectory Calculates and stores checksums for the files in the given directory.
func (calculator *Calculator) RecordChecksumsForDirectory(db *dal.Db, algorithm string) {

	files := util.ListDirectoryRecursively(calculator.InputDirectory)
	hasher := NewHasher(algorithm)
	effectiveBasePath := calculator.getEffectiveBasePath()
	fingerprints := hasher.CalculateChecksumsForFiles(calculator.InputDirectory, effectiveBasePath, files)
	db.Fingerprints = fingerprints
	db.SaveCsv(calculator.OutputChecksums)
}

func (calculator *Calculator) getEffectiveBasePath() string {

	return util.TrimPath(calculator.InputDirectory, calculator.BasePath)
}
