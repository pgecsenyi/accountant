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
func (calculator *Calculator) RecordChecksumsForDirectory(hasher *checksum.FileHasher) {

	files := util.ListDirectoryRecursively(calculator.InputDirectory)
	hasher.CalculateChecksumsForFiles(calculator.InputDirectory, files, calculator.BasePath)
	hasher.ExportToCsv(calculator.OutputChecksums)
}
