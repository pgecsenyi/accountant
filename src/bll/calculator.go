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
	fingerprints := hasher.CalculateChecksumsForFiles(calculator.InputDirectory, files, calculator.BasePath)
	db.Fingerprints = fingerprints
	db.SaveCsv(calculator.OutputChecksums)
}
