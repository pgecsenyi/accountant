package bll

import (
	"bufio"
	"checksum"
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
	"util"
)

// Importer Stores settings related to import.
type Importer struct {
	InputDirectory  string
	OutputChecksums string
	patternCrc32    *regexp.Regexp
	patternMd5      *regexp.Regexp
	patternSha1     *regexp.Regexp
	patternSha256   *regexp.Regexp
	patternSha512   *regexp.Regexp
}

// NewImporter Instantiates a new Importer object.
func NewImporter(inputDirectory string, outputChecksums string) Importer {

	return Importer{inputDirectory, outputChecksums, nil, nil, nil, nil, nil}
}

// Convert Converts checksum data produced by third party utilities to CSV.
func (importer *Importer) Convert(hasher *checksum.FileHasher) {

	fpPrototype := new(checksum.Fingerprint)
	fpPrototype.Note = ""

	files := util.ListDirectoryRecursively(importer.InputDirectory)
	for _, file := range files {
		fullPath := path.Join(importer.InputDirectory, file)
		setCreatedAtTime(fpPrototype, fullPath)
		importer.loadDataFromFile(hasher, fullPath, fpPrototype)
	}

	hasher.ExportToCsv(importer.OutputChecksums)
}

func setCreatedAtTime(fpPrototype *checksum.Fingerprint, filePath string) {

	fileInfo, err := os.Stat(filePath)
	util.CheckErrDontPanic(err, "Unable to get the file modification time for "+filePath+".")
	fpPrototype.CreatedAt = fileInfo.ModTime().UTC().Format(time.RFC3339)
}

func (importer *Importer) loadDataFromFile(
	hasher *checksum.FileHasher, filePath string, fpPrototype *checksum.Fingerprint) {

	extension := path.Ext(filePath)

	if extension == checksum.CRC32EXT {
		fpPrototype.Algorithm = checksum.CRC32
		compilePattern(&importer.patternCrc32, checksum.PATTERNCRC32, checksum.CRC32LEN)
		parseFile(hasher, fpPrototype, filePath, importer.patternCrc32, ';')
	} else if extension == checksum.MD5EXT {
		fpPrototype.Algorithm = checksum.MD5
		compilePattern(&importer.patternMd5, checksum.PATTERNCOMMON, checksum.MD5LEN)
		parseFile(hasher, fpPrototype, filePath, importer.patternMd5, '*')
	} else if extension == checksum.SHA1EXT {
		fpPrototype.Algorithm = checksum.SHA1
		compilePattern(&importer.patternSha1, checksum.PATTERNCOMMON, checksum.SHA1LEN)
		parseFile(hasher, fpPrototype, filePath, importer.patternSha1, '*')
	} else if extension == checksum.SHA256EXT {
		fpPrototype.Algorithm = checksum.SHA256
		compilePattern(&importer.patternSha256, checksum.PATTERNCOMMON, checksum.SHA256LEN)
		parseFile(hasher, fpPrototype, filePath, importer.patternSha256, '*')
	} else if extension == checksum.SHA512EXT {
		fpPrototype.Algorithm = checksum.SHA512
		compilePattern(&importer.patternSha512, checksum.PATTERNCOMMON, checksum.SHA512LEN)
		parseFile(hasher, fpPrototype, filePath, importer.patternSha512, '*')
	}
}

func compilePattern(re **regexp.Regexp, pattern string, length int) {

	if *re == nil {
		completePattern := fmt.Sprintf(pattern, length)
		*re = regexp.MustCompile(completePattern)
	}
}

func parseFile(
	hasher *checksum.FileHasher, fpPrototype *checksum.Fingerprint,
	filePath string, pattern *regexp.Regexp, commentChar byte) int {

	file, err := os.Open(filePath)
	util.CheckErr(err, "Cannot open file "+filePath+".")
	defer file.Close()

	idxFilename, idxChecksum := getFilenameChecksumIndices(pattern.SubexpNames())

	numberOfInvalidLines := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if !parseLine(hasher, fpPrototype, scanner.Text(), pattern, commentChar, idxFilename, idxChecksum) {
			numberOfInvalidLines++
		}
	}

	util.CheckErr(scanner.Err(), "Error reading file "+filePath+".")

	if numberOfInvalidLines != 0 {
		fmt.Printf("There is/are %d invalid line(s) in %s.", numberOfInvalidLines, filePath)
		fmt.Println()
	}

	return 0
}

func getFilenameChecksumIndices(indexNames []string) (int, int) {

	idxChecksum := 0
	idxFilename := 0

	for idx, item := range indexNames {
		if item == "file" {
			idxFilename = idx
		} else if item == "hash" {
			idxChecksum = idx
		}
	}

	return idxFilename, idxChecksum
}

func parseLine(
	hasher *checksum.FileHasher, fpPrototype *checksum.Fingerprint,
	line string, pattern *regexp.Regexp, commentChar byte, idxFilename int, idxChecksum int) bool {

	// This line is a comment or whitespace.
	if len(strings.TrimSpace(line)) <= 0 || line[0] == commentChar {
		return true
	}

	// Try to get the checksum.
	matches := pattern.FindStringSubmatch(line)
	if matches == nil || len(matches) < 3 {
		return false
	}

	// Add fingerprint to the database.
	checksumBytes, err := hex.DecodeString(matches[idxChecksum])
	if err != nil {
		return false
	}

	newPrototype := cloneFingerprintPrototype(fpPrototype)
	newPrototype.Filename = matches[idxFilename]
	newPrototype.Checksum = checksumBytes
	hasher.Fingerprints.PushFront(newPrototype)

	return true
}

func cloneFingerprintPrototype(fpPrototype *checksum.Fingerprint) *checksum.Fingerprint {

	clone := new(checksum.Fingerprint)
	clone.Algorithm = fpPrototype.Algorithm
	clone.CreatedAt = fpPrototype.CreatedAt

	return clone
}
