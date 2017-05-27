package bll

import (
	"container/list"
	"dal"
	"encoding/hex"
	"log"
	"os"
	"util"
)

// Comparer Stores settings related to comparison.
type Comparer struct {
	Db             dal.Database
	InputDirectory string
	OutputNames    string
	BasePath       string
}

// Compare Verifies and stores changes in the given directory based on the checksums calculated earlier.
func (comparer *Comparer) Compare(algorithm string) {

	oldFingerprints := comparer.loadOldFingerprints()
	newFingerprints := comparer.calculateNewFingerprints(algorithm)

	comparer.Db.SetFingerprints(newFingerprints)
	comparer.Db.Save()

	compareAndSaveResults(oldFingerprints, newFingerprints, comparer.OutputNames)
}

func (comparer *Comparer) loadOldFingerprints() *list.List {

	comparer.Db.Load()
	oldFingerprints := comparer.Db.GetFingerprints()

	return oldFingerprints
}

func (comparer *Comparer) calculateNewFingerprints(algorithm string) *list.List {

	hasher := NewHasher(algorithm)
	effectiveBasePath := comparer.getEffectiveBasePath()
	files := util.ListDirectoryRecursively(comparer.InputDirectory)
	newFingerprints := hasher.CalculateFingerprints(comparer.InputDirectory, effectiveBasePath, files)

	return newFingerprints
}

func (comparer *Comparer) getEffectiveBasePath() string {

	return util.TrimPath(comparer.InputDirectory, comparer.BasePath)
}

func compareAndSaveResults(oldFingerprints *list.List, newFingerprints *list.List, outputPath string) {

	cache := buildFingerprintCache(oldFingerprints)

	outputFile, err := os.OpenFile(outputPath, os.O_CREATE|os.O_TRUNC|os.O_RDWR|os.O_APPEND, 0660)
	util.CheckErr(err, "Cannot open output for name pairs: "+outputPath)
	defer outputFile.Close()

	compareFingerprints(cache, newFingerprints, outputFile)
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

func compareFingerprints(oldFingerprints map[string]*dal.Fingerprint, newFingerprints *list.List, output *os.File) {

	foundFingerprints := make(map[string]bool)

	for element := newFingerprints.Front(); element != nil; element = element.Next() {
		fingerprint := element.Value.(*dal.Fingerprint)
		checksum := hex.EncodeToString(fingerprint.Checksum)
		matchingFingerprint := oldFingerprints[checksum]
		processMatch(fingerprint, checksum, matchingFingerprint, output, foundFingerprints)
	}

	printRemovedFiles(oldFingerprints, foundFingerprints)
}

func processMatch(
	fingerprint *dal.Fingerprint, checksum string, matchingFingerprint *dal.Fingerprint,
	output *os.File, foundFingerprints map[string]bool) {

	if matchingFingerprint == nil {
		log.Println("New: " + fingerprint.Filename + ".")
	} else {
		saveNamePair(matchingFingerprint.Filename, fingerprint.Filename, output)
		fingerprint.CreatedAt = matchingFingerprint.CreatedAt
		foundFingerprints[checksum] = true
	}
}

func saveNamePair(oldName string, newName string, output *os.File) {

	output.WriteString(newName + "\r\n")
	output.WriteString("    " + oldName + "\r\n")
	output.WriteString("    \r\n")
	output.WriteString("    \r\n")
}

func printRemovedFiles(oldFingerprints map[string]*dal.Fingerprint, foundFingerprints map[string]bool) {

	for hash, fingerprint := range oldFingerprints {
		hasFound := foundFingerprints[hash]
		if !hasFound {
			log.Println("Missing: " + fingerprint.Filename)
		}
	}
}
