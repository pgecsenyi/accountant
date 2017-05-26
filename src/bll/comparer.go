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
	InputDirectory  string
	InputChecksums  string
	OutputNames     string
	OutputChecksums string
	BasePath        string
}

// RecordNameChangesForDirectory Verifies and stores changes in the given directory based on the checksums calculated earlier.
func (comparer *Comparer) RecordNameChangesForDirectory(db *dal.Db, algorithm string) {

	db.LoadCsv(comparer.InputChecksums)
	oldFingerprints := db.Fingerprints

	hasher := NewHasher(algorithm)
	files := util.ListDirectoryRecursively(comparer.InputDirectory)
	effectiveBasePath := comparer.getEffectiveBasePath()
	newFingerprints := hasher.CalculateChecksumsForFiles(comparer.InputDirectory, effectiveBasePath, files)
	db.Fingerprints = newFingerprints

	recordDifferences(oldFingerprints, newFingerprints, comparer.OutputNames)
	db.SaveCsv(comparer.OutputChecksums)
}

func (comparer *Comparer) getEffectiveBasePath() string {

	return util.TrimPath(comparer.InputDirectory, comparer.BasePath)
}

func recordDifferences(oldFingerprints *list.List, newFingerprints *list.List, outputPath string) {

	cache := buildFingerprintCache(oldFingerprints)

	f, err := os.OpenFile(outputPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
	util.CheckErr(err, "Cannot open output for name pairs: "+outputPath)
	defer f.Close()

	compareFingerprints(cache, newFingerprints, f)
}

func buildFingerprintCache(fingerprints *list.List) map[string]*dal.Fingerprint {

	var cache = make(map[string]*dal.Fingerprint)
	for element := fingerprints.Front(); element != nil; element = element.Next() {
		fp := element.Value.(*dal.Fingerprint)
		fpString := hex.EncodeToString(fp.Checksum)
		cache[fpString] = fp
	}

	return cache
}

func compareFingerprints(oldFingerprints map[string]*dal.Fingerprint, newFingerprints *list.List, output *os.File) {

	foundFingerprints := make(map[string]bool)

	for element := newFingerprints.Front(); element != nil; element = element.Next() {
		newFp := element.Value.(*dal.Fingerprint)
		newFpString := hex.EncodeToString(newFp.Checksum)
		oldFp := oldFingerprints[newFpString]
		if oldFp == nil {
			log.Println("New file: " + newFp.Filename + ".")
		} else {
			saveNameChange(oldFp.Filename, newFp.Filename, output)
			newFp.CreatedAt = oldFp.CreatedAt
			oldFpString := hex.EncodeToString(oldFp.Checksum)
			foundFingerprints[oldFpString] = true
		}
	}

	printRemovedFiles(oldFingerprints, foundFingerprints)
}

func saveNameChange(oldName string, newName string, output *os.File) {

	output.WriteString(newName + "\r\n")
	output.WriteString("    " + oldName + "\r\n")
	output.WriteString("    \r\n")
	output.WriteString("    \r\n")
}

func printRemovedFiles(oldFingerprints map[string]*dal.Fingerprint, foundFingerprints map[string]bool) {

	for hash, fingerprint := range oldFingerprints {
		hasFound := foundFingerprints[hash]
		if !hasFound {
			log.Println("Missing file: " + fingerprint.Filename)
		}
	}
}
