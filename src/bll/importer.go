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
	InputDirectory   string
	OutputChecksums  string
	patterns         importEntryPatterns
	fingerprintProto *dal.Fingerprint
}

type importEntryPatterns struct {
	patternCrc32  *regexp.Regexp
	patternMd5    *regexp.Regexp
	patternSha1   *regexp.Regexp
	patternSha256 *regexp.Regexp
	patternSha512 *regexp.Regexp
}

// NewImporter Instantiates a new Importer object.
func NewImporter(inputDirectory string, outputChecksums string) Importer {

	patterns := importEntryPatterns{nil, nil, nil, nil, nil}
	fingerprintProto := new(dal.Fingerprint)

	return Importer{inputDirectory, outputChecksums, patterns, fingerprintProto}
}

// Convert Converts checksum data produced by third party utilities to CSV.
func (importer *Importer) Convert(db *dal.Db) {

	files := util.ListDirectoryRecursively(importer.InputDirectory)

	for _, file := range files {
		fullPath := path.Join(importer.InputDirectory, file)
		importer.updateProtoTime(fullPath)
		importer.loadDataFromFile(db, fullPath)
	}

	db.SaveCsv(importer.OutputChecksums)
}

func (importer *Importer) updateProtoTime(filePath string) {

	fileInfo, err := os.Stat(filePath)
	util.CheckErrDontPanic(err, "Unable to get the file modification time for "+filePath+".")
	importer.fingerprintProto.CreatedAt = fileInfo.ModTime().UTC().Format(time.RFC3339)
}

func (importer *Importer) loadDataFromFile(db *dal.Db, filePath string) {

	extension := path.Ext(filePath)

	if extension == dal.CRC32EXT {
		importer.fingerprintProto.Algorithm = dal.CRC32
		compilePattern(&importer.patterns.patternCrc32, dal.PATTERNCRC32, dal.CRC32LEN)
		importer.parseFile(db, filePath, importer.patterns.patternCrc32, ';')
	} else if extension == dal.MD5EXT {
		importer.fingerprintProto.Algorithm = dal.MD5
		compilePattern(&importer.patterns.patternMd5, dal.PATTERNCOMMON, dal.MD5LEN)
		importer.parseFile(db, filePath, importer.patterns.patternMd5, '*')
	} else if extension == dal.SHA1EXT {
		importer.fingerprintProto.Algorithm = dal.SHA1
		compilePattern(&importer.patterns.patternSha1, dal.PATTERNCOMMON, dal.SHA1LEN)
		importer.parseFile(db, filePath, importer.patterns.patternSha1, '*')
	} else if extension == dal.SHA256EXT {
		importer.fingerprintProto.Algorithm = dal.SHA256
		compilePattern(&importer.patterns.patternSha256, dal.PATTERNCOMMON, dal.SHA256LEN)
		importer.parseFile(db, filePath, importer.patterns.patternSha256, '*')
	} else if extension == dal.SHA512EXT {
		importer.fingerprintProto.Algorithm = dal.SHA512
		compilePattern(&importer.patterns.patternSha512, dal.PATTERNCOMMON, dal.SHA512LEN)
		importer.parseFile(db, filePath, importer.patterns.patternSha512, '*')
	}
}

func (importer *Importer) parseFile(db *dal.Db, filePath string, pattern *regexp.Regexp, commentChar byte) {

	file, err := os.Open(filePath)
	util.CheckErr(err, "Cannot open file "+filePath+".")
	defer file.Close()

	idxFilename, idxChecksum := getFilenameChecksumIndices(pattern.SubexpNames())

	numberOfInvalidLines := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if !importer.parseLine(db, scanner.Text(), pattern, commentChar, idxFilename, idxChecksum) {
			numberOfInvalidLines++
		}
	}

	util.CheckErr(scanner.Err(), "Error reading file "+filePath+".")

	if numberOfInvalidLines != 0 {
		message := fmt.Sprintf("There is/are %d invalid line(s) in %s.", numberOfInvalidLines, filePath)
		log.Println(message)
	}
}

func (importer *Importer) parseLine(
	db *dal.Db, line string,
	pattern *regexp.Regexp, commentChar byte,
	idxFilename int, idxChecksum int) bool {

	// This line is a comment or whitespace.
	if len(strings.TrimSpace(line)) <= 0 || line[0] == commentChar {
		return true
	}

	// Try to filename and checksum.
	matches := pattern.FindStringSubmatch(line)
	if matches == nil || len(matches) < 3 {
		return false
	}

	// Add fingerprint to the database.
	if !importer.addEntry(db, matches[idxFilename], matches[idxChecksum]) {
		return false
	}

	return true
}

func (importer *Importer) addEntry(db *dal.Db, file string, checksum string) bool {

	checksumBytes, err := hex.DecodeString(checksum)
	if err != nil {
		return false
	}

	fingerprint := importer.cloneFingerprintProto(file, checksumBytes)
	db.Fingerprints.PushFront(fingerprint)

	return true
}

func (importer *Importer) cloneFingerprintProto(file string, checksum []byte) *dal.Fingerprint {

	clone := new(dal.Fingerprint)
	clone.Filename = util.NormalizePath(file)
	clone.Checksum = checksum
	clone.Algorithm = importer.fingerprintProto.Algorithm
	clone.CreatedAt = importer.fingerprintProto.CreatedAt

	return clone
}

func compilePattern(re **regexp.Regexp, pattern string, length int) {

	if *re == nil {
		completePattern := fmt.Sprintf(pattern, length)
		*re = regexp.MustCompile(completePattern)
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
