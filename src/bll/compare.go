package bll

import (
	"checksum"
	"container/list"
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
func (comparer *Comparer) RecordNameChangesForDirectory(hasher *checksum.FileHasher) {

	hasher.LoadFromCsv(comparer.InputChecksums)
	oldFingerprints := hasher.Fingerprints

	hasher.Reset()

	files := util.ListDirectoryRecursively(comparer.InputDirectory)
	hasher.CalculateChecksumsForFiles(comparer.InputDirectory, files, comparer.BasePath)
	newFingerprints := hasher.Fingerprints

	recordDifferences(oldFingerprints, newFingerprints, comparer.OutputNames)
	hasher.ExportToCsv(comparer.OutputChecksums)
}

func recordDifferences(oldFingerprints *list.List, newFingerprints *list.List, outputPath string) {

	cache := buildFingerprintCache(oldFingerprints)

	f, err := os.OpenFile(outputPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
	util.CheckErr(err, "Cannot open output for name pairs: "+outputPath)
	defer f.Close()

	compareFingerprints(cache, newFingerprints, f)
}

func buildFingerprintCache(fingerprints *list.List) map[string]*checksum.Fingerprint {

	var cache = make(map[string]*checksum.Fingerprint)
	for element := fingerprints.Front(); element != nil; element = element.Next() {
		fp := element.Value.(*checksum.Fingerprint)
		fpString := hex.EncodeToString(fp.Checksum)
		cache[fpString] = fp
	}

	return cache
}

func compareFingerprints(oldFingerprints map[string]*checksum.Fingerprint, newFingerprints *list.List, output *os.File) {

	foundFingerprints := make(map[string]bool)

	for element := newFingerprints.Front(); element != nil; element = element.Next() {
		newFp := element.Value.(*checksum.Fingerprint)
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

func printRemovedFiles(oldFingerprints map[string]*checksum.Fingerprint, foundFingerprints map[string]bool) {

	for hash, fingerprint := range oldFingerprints {
		hasFound := foundFingerprints[hash]
		if !hasFound {
			log.Println("Missing file: " + fingerprint.Filename)
		}
	}
}
