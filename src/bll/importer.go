package bll

import (
	"bll/report"
	"bufio"
	"dal"
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
	Db               dal.Database
	InputDirectory   string
	OutputChecksums  string
	Report           *report.ImportReport
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
func NewImporter(db dal.Database, inputDirectory string, outputChecksums string) Importer {

	patterns := importEntryPatterns{nil, nil, nil, nil, nil}
	fingerprintProto := new(dal.Fingerprint)
	report := report.NewImportReport()

	return Importer{db, inputDirectory, outputChecksums, report, patterns, fingerprintProto}
}

// Convert Converts checksum data produced by third party utilities to CSV.
func (importer *Importer) Convert() {

	files := util.ListFilesRecursively(importer.InputDirectory)

	for _, file := range files {
		fullPath := path.Join(importer.InputDirectory, file)
		importer.updateProtoTime(fullPath)
		importer.loadDataFromFile(fullPath)
	}

	importer.Db.SaveFingerprints()
}

func (importer *Importer) updateProtoTime(filePath string) {

	fileInfo, err := os.Stat(filePath)
	util.CheckErrDontPanic(err, "Unable to get the file modification time for "+filePath+".")
	importer.fingerprintProto.CreatedAt = fileInfo.ModTime().UTC().Format(time.RFC3339)
}

func (importer *Importer) loadDataFromFile(filePath string) {

	extension := path.Ext(filePath)

	if extension == dal.CRC32EXT {
		importer.fingerprintProto.Algorithm = dal.CRC32
		compilePattern(&importer.patterns.patternCrc32, dal.PATTERNCRC32, dal.CRC32LEN)
		importer.parseFile(filePath, importer.patterns.patternCrc32, ';')
	} else if extension == dal.MD5EXT {
		importer.fingerprintProto.Algorithm = dal.MD5
		compilePattern(&importer.patterns.patternMd5, dal.PATTERNCOMMON, dal.MD5LEN)
		importer.parseFile(filePath, importer.patterns.patternMd5, '*')
	} else if extension == dal.SHA1EXT {
		importer.fingerprintProto.Algorithm = dal.SHA1
		compilePattern(&importer.patterns.patternSha1, dal.PATTERNCOMMON, dal.SHA1LEN)
		importer.parseFile(filePath, importer.patterns.patternSha1, '*')
	} else if extension == dal.SHA256EXT {
		importer.fingerprintProto.Algorithm = dal.SHA256
		compilePattern(&importer.patterns.patternSha256, dal.PATTERNCOMMON, dal.SHA256LEN)
		importer.parseFile(filePath, importer.patterns.patternSha256, '*')
	} else if extension == dal.SHA512EXT {
		importer.fingerprintProto.Algorithm = dal.SHA512
		compilePattern(&importer.patterns.patternSha512, dal.PATTERNCOMMON, dal.SHA512LEN)
		importer.parseFile(filePath, importer.patterns.patternSha512, '*')
	}
}

func (importer *Importer) parseFile(filePath string, pattern *regexp.Regexp, commentChar byte) {

	file, err := os.Open(filePath)
	util.CheckErr(err, "Cannot open file "+filePath+".")
	defer file.Close()

	idxFilename, idxChecksum := getFilenameChecksumIndices(pattern.SubexpNames())

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if !importer.parseLine(scanner.Text(), pattern, commentChar, idxFilename, idxChecksum) {
			importer.Report.IncreaseInvalidEntryCount(filePath)
		}
	}

	util.CheckErr(scanner.Err(), "Error reading file "+filePath+".")

	importer.Report.LogSummaryForFile(filePath)
}

func (importer *Importer) parseLine(
	line string,
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
	if !importer.addEntry(matches[idxFilename], matches[idxChecksum]) {
		return false
	}

	return true
}

func (importer *Importer) addEntry(file string, checksum string) bool {

	checksumBytes, err := hex.DecodeString(checksum)
	if err != nil {
		return false
	}

	fingerprint := importer.cloneFingerprintProto(file, checksumBytes)
	importer.Db.AddFingerprint(fingerprint)

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
