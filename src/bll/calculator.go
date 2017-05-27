package bll

import (
	"container/list"
	"dal"
	"path"
	"util"
)

// Calculator Stores settings related to checksum calculation.
type Calculator struct {
	InputDirectory    string
	OutputChecksums   string
	BasePath          string
	InputChecksums    string
	hasher            Hasher
	effectiveBasePath string
}

// NewCalculator Instantiates a new Calculator object.
func NewCalculator(
	inputDirectory string, algorithm string, outputChecksums string, basePath string, inputChecksums string) Calculator {

	hasher := NewHasher(algorithm)
	effectiveBasePath := util.TrimPath(inputDirectory, basePath)

	return Calculator{inputDirectory, outputChecksums, basePath, inputChecksums, hasher, effectiveBasePath}
}

// Calculate Calculates and stores checksums for the files in the given directory.
func (calculator *Calculator) Calculate(db *dal.Db) {

	files := util.ListDirectoryRecursively(calculator.InputDirectory)
	fingerprints := calculator.calculateFingerprints(db, files)

	db.Fingerprints = fingerprints
	db.SaveCsv(calculator.OutputChecksums)
}

func (calculator *Calculator) calculateFingerprints(db *dal.Db, files []string) *list.List {

	if calculator.InputChecksums == "" {
		return calculator.hasher.CalculateFingerprints(calculator.InputDirectory, calculator.effectiveBasePath, files)
	}

	return calculator.calculateFingerprintsForMissingFiles(db, files)
}

func (calculator *Calculator) calculateFingerprintsForMissingFiles(db *dal.Db, files []string) *list.List {

	fingerprints := list.New()
	etm := calculator.loadMissingNames(db)

	for _, file := range files {
		calculator.addFingerprintIfFileIsMissing(file, etm, fingerprints)
	}

	return fingerprints
}

func (calculator *Calculator) loadMissingNames(db *dal.Db) *effectiveTextMemory {

	etm := newEffectiveTextMemory()
	db.LoadNamesFromCsv(calculator.InputChecksums, etm)
	etm.ClearCache()

	return etm
}

func (calculator *Calculator) addFingerprintIfFileIsMissing(
	file string, etm *effectiveTextMemory, fingerprints *list.List) {

	fullPath := path.Join(calculator.effectiveBasePath, file)
	if !etm.ContainsText(fullPath) {
		fp := calculator.hasher.CalculateFingerprint(calculator.InputDirectory, calculator.effectiveBasePath, file)
		fingerprints.PushFront(fp)
	}
}
