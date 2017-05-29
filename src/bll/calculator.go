package bll

import (
	"container/list"
	"dal"
	"path"
	"util"
)

// Calculator Stores settings related to checksum calculation.
type Calculator struct {
	Db                dal.Database
	InputDirectory    string
	BasePath          string
	hasher            Hasher
	effectiveBasePath string
}

// NewCalculator Instantiates a new Calculator object.
func NewCalculator(db dal.Database, inputDirectory string, algorithm string, basePath string) Calculator {

	hasher := NewHasher(algorithm)
	effectiveBasePath := util.TrimPath(inputDirectory, basePath)

	return Calculator{db, inputDirectory, basePath, hasher, effectiveBasePath}
}

// Calculate Calculates and stores checksums for the files in the given directory.
func (calculator *Calculator) Calculate(missingOnly bool) {

	files := util.ListFilesRecursively(calculator.InputDirectory)
	fingerprints := calculator.calculateFingerprints(files, missingOnly)

	calculator.Db.SetFingerprints(fingerprints)
	calculator.Db.Save()
}

func (calculator *Calculator) calculateFingerprints(files []string, missingOnly bool) *list.List {

	if missingOnly {
		return calculator.calculateFingerprintsForMissingFiles(files)
	}

	return calculator.hasher.CalculateFingerprints(calculator.InputDirectory, calculator.effectiveBasePath, files)
}

func (calculator *Calculator) calculateFingerprintsForMissingFiles(files []string) *list.List {

	fingerprints := list.New()
	etm := calculator.loadMissingNames()

	for _, file := range files {
		calculator.addFingerprintIfFileIsMissing(file, etm, fingerprints)
	}

	return fingerprints
}

func (calculator *Calculator) loadMissingNames() *effectiveTextMemory {

	etm := newEffectiveTextMemory()
	calculator.Db.LoadNames(etm)
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
