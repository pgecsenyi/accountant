package bll

import (
	"checksum"
	"util"
)

// Calculator Stores settings related to checksum calculation.
type Calculator struct {
	InputDirectory  string
	OutputChecksums string
	BasePath        string
}

// RecordChecksumsForDirectory Calculates and stores checksums for the files in the given directory.
func (calculator *Calculator) RecordChecksumsForDirectory(hasher *checksum.FileHasher, algorithm string) {

	files := util.ListDirectoryRecursively(calculator.InputDirectory)
	fingerprints := checksum.CalculateChecksumsForFiles(calculator.InputDirectory, files, calculator.BasePath, algorithm)
	hasher.Fingerprints = fingerprints
	hasher.SaveCsv(calculator.OutputChecksums)
}
