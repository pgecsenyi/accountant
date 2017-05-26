package bll

import (
	"bufio"
	"dal"
	"encoding/hex"
	"fmt"
	"log"
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
func (importer *Importer) Convert(db *dal.Db) {

	fpPrototype := new(dal.Fingerprint)
	fpPrototype.Note = ""

	files := util.ListDirectoryRecursively(importer.InputDirectory)
	for _, file := range files {
		fullPath := path.Join(importer.InputDirectory, file)
		setCreatedAtTime(fpPrototype, fullPath)
		importer.loadDataFromFile(db, fullPath, fpPrototype)
	}

	db.SaveCsv(importer.OutputChecksums)
}

func setCreatedAtTime(fpPrototype *dal.Fingerprint, filePath string) {

	fileInfo, err := os.Stat(filePath)
	util.CheckErrDontPanic(err, "Unable to get the file modification time for "+filePath+".")
	fpPrototype.CreatedAt = fileInfo.ModTime().UTC().Format(time.RFC3339)
}

func (importer *Importer) loadDataFromFile(db *dal.Db, filePath string, fpPrototype *dal.Fingerprint) {

	extension := path.Ext(filePath)

	if extension == dal.CRC32EXT {
		fpPrototype.Algorithm = dal.CRC32
		compilePattern(&importer.patternCrc32, dal.PATTERNCRC32, dal.CRC32LEN)
		parseFile(db, fpPrototype, filePath, importer.patternCrc32, ';')
	} else if extension == dal.MD5EXT {
		fpPrototype.Algorithm = dal.MD5
		compilePattern(&importer.patternMd5, dal.PATTERNCOMMON, dal.MD5LEN)
		parseFile(db, fpPrototype, filePath, importer.patternMd5, '*')
	} else if extension == dal.SHA1EXT {
		fpPrototype.Algorithm = dal.SHA1
		compilePattern(&importer.patternSha1, dal.PATTERNCOMMON, dal.SHA1LEN)
		parseFile(db, fpPrototype, filePath, importer.patternSha1, '*')
	} else if extension == dal.SHA256EXT {
		fpPrototype.Algorithm = dal.SHA256
		compilePattern(&importer.patternSha256, dal.PATTERNCOMMON, dal.SHA256LEN)
		parseFile(db, fpPrototype, filePath, importer.patternSha256, '*')
	} else if extension == dal.SHA512EXT {
		fpPrototype.Algorithm = dal.SHA512
		compilePattern(&importer.patternSha512, dal.PATTERNCOMMON, dal.SHA512LEN)
		parseFile(db, fpPrototype, filePath, importer.patternSha512, '*')
	}
}

func compilePattern(re **regexp.Regexp, pattern string, length int) {

	if *re == nil {
		completePattern := fmt.Sprintf(pattern, length)
		*re = regexp.MustCompile(completePattern)
	}
}

func parseFile(
	db *dal.Db, fpPrototype *dal.Fingerprint,
	filePath string, pattern *regexp.Regexp, commentChar byte) {

	file, err := os.Open(filePath)
	util.CheckErr(err, "Cannot open file "+filePath+".")
	defer file.Close()

	idxFilename, idxChecksum := getFilenameChecksumIndices(pattern.SubexpNames())

	numberOfInvalidLines := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if !parseLine(db, fpPrototype, scanner.Text(), pattern, commentChar, idxFilename, idxChecksum) {
			numberOfInvalidLines++
		}
	}

	util.CheckErr(scanner.Err(), "Error reading file "+filePath+".")

	if numberOfInvalidLines != 0 {
		message := fmt.Sprintf("There is/are %d invalid line(s) in %s.", numberOfInvalidLines, filePath)
		log.Println(message)
	}
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
	db *dal.Db, fpPrototype *dal.Fingerprint,
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
	newPrototype.Filename = util.NormalizePath(matches[idxFilename])
	newPrototype.Checksum = checksumBytes
	db.Fingerprints.PushFront(newPrototype)

	return true
}

func cloneFingerprintPrototype(fpPrototype *dal.Fingerprint) *dal.Fingerprint {

	clone := new(dal.Fingerprint)
	clone.Algorithm = fpPrototype.Algorithm
	clone.CreatedAt = fpPrototype.CreatedAt

	return clone
}
