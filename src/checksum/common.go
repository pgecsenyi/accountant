package checksum

import (
	"container/list"
	"io"
	"os"
	"path"
	"runtime"
	"time"
	"util"
)

var currentTime = time.Now().Format(time.RFC3339)
var runTimeVersion = runtime.Version()

// CalculateChecksumForFile Calculates checksum for the given file using the given algorithm.
func CalculateChecksumForFile(filename string, algorithm string) []byte {

	file, err := os.Open(filename)
	util.CheckErr(err, "Cannot read file "+filename+".")
	defer file.Close()

	calculator := CreateCalculator(algorithm)
	io.Copy(calculator, file)
	checksum := calculator.Sum(nil)[:]

	return checksum
}

// CalculateChecksumsForFiles Calculates checksum for each file in the given list.
func CalculateChecksumsForFiles(basePath string, files []string, prefixToRemove string, algorithm string) *list.List {

	fingerprints := list.New()
	for _, file := range files {
		recordChecksumForFile(basePath, file, algorithm, prefixToRemove, fingerprints)
	}

	return fingerprints
}

func recordChecksumForFile(
	basePath string, filePath string,
	algorithm string, prefixToRemove string,
	fingerprints *list.List) {

	fullPath := path.Join(basePath, filePath)
	checksum := CalculateChecksumForFile(fullPath, algorithm)
	normalizedPath := util.NormalizePath(fullPath)[len(prefixToRemove):]

	fp := new(Fingerprint)
	fp.Filename = normalizedPath
	fp.Checksum = checksum
	fp.Algorithm = algorithm
	fp.CreatedAt = currentTime
	fp.Creator = runTimeVersion
	fp.Note = ""

	fingerprints.PushFront(fp)
}
