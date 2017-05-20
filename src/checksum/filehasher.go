package checksum

import (
	"bytes"
	"container/list"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"time"
	"util"
)

// FileHasher Logic for calculating checksums.
type FileHasher struct {
	Fingerprints *list.List
	algorithm    string
}

// Fingerprint Stores the necessary data to identify a file and a bit more.
type Fingerprint struct {
	Filename  string
	Checksum  []byte
	Algorithm string
	CreatedAt string
	Creator   string
	Note      string
}

var currentTime = time.Now().Format(time.RFC3339)
var runTimeVersion = runtime.Version()

// NewFileHasher Instantiates a new FileHasher object.
func NewFileHasher(algorithm string) FileHasher {
	return FileHasher{list.New(), algorithm}
}

// CalculateChecksumsForFiles Calculates checksum for each file in the given list.
func (fh *FileHasher) CalculateChecksumsForFiles(basePath string, files []string, prefixToRemove string) {

	for _, file := range files {
		fh.recordChecksumForFile(basePath, file, fh.algorithm, prefixToRemove)
	}
}

// ExportToCsv Exports fingerprints to the given CSV file.
func (fh *FileHasher) ExportToCsv(filename string) {

	records := createStringArrayFromFingerprints(fh.Fingerprints)
	writeChecksumsToCsvFile(records, filename)
}

// ImportFromCsv Imports fingerprints from the given CSV file.
func (fh *FileHasher) ImportFromCsv(filename string) {

	content := readFileContent(filename)
	reader := csv.NewReader(bytes.NewReader(content))

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		util.CheckErrDontPanic(err, "")

		checksumBytes, err := hex.DecodeString(record[1])
		util.CheckErrDontPanic(err, "")
		fh.addFingerprint(record, checksumBytes)
	}
}

// Reset Clears fingerprints.
func (fh *FileHasher) Reset() {

	(*fh).Fingerprints = list.New()
}

// VerifyFiles Verifies checksums in the stored list applying the given base path.
func (fh *FileHasher) VerifyFiles(basePath string) {

	numberOfInvalid := 0
	numberOfNonExisting := 0

	for element := fh.Fingerprints.Front(); element != nil; element = element.Next() {
		meta := element.Value.(*Fingerprint)
		fullPath := path.Join(basePath, meta.Filename)
		if !util.CheckIfFileExists(fullPath) {
			fmt.Println(fmt.Sprintf("%s does not exist", meta.Filename))
			numberOfNonExisting++
		} else {
			checksum := calculateChecksumForFile(fullPath, meta.Algorithm)
			if !util.Compare(checksum, meta.Checksum) {
				fmt.Println(fmt.Sprintf("%s has an invalid checksum", meta.Filename))
				numberOfInvalid++
			}
		}
	}

	numberOfAll := fh.Fingerprints.Len()
	fmt.Println()
	fmt.Printf(
		"Valid: %d/%d, missing: %d, invalid: %d.",
		numberOfAll-numberOfNonExisting-numberOfInvalid,
		numberOfAll, numberOfNonExisting, numberOfInvalid)
}

func (fh *FileHasher) recordChecksumForFile(basePath string, filePath string, algorithm string, prefixToRemove string) {

	fullPath := path.Join(basePath, filePath)
	checksum := calculateChecksumForFile(fullPath, algorithm)
	normalizedPath := util.NormalizePath(fullPath)[len(prefixToRemove):]

	fp := new(Fingerprint)
	fp.Filename = normalizedPath
	fp.Checksum = checksum
	fp.Algorithm = algorithm
	fp.CreatedAt = currentTime
	fp.Creator = runTimeVersion
	fp.Note = ""
	fh.Fingerprints.PushFront(fp)
}

func readFileContent(filename string) []byte {

	content, err := ioutil.ReadFile(filename)
	util.CheckErr(err, "Cannot read file "+filename+".")

	return content
}

func (fh *FileHasher) addFingerprint(record []string, checksumBytes []byte) {

	fp := new(Fingerprint)
	fp.Filename = record[0]
	fp.Checksum = checksumBytes
	fp.Algorithm = record[2]
	fp.CreatedAt = record[3]
	fp.Creator = record[4]
	fp.Note = record[5]
	fh.Fingerprints.PushFront(fp)
}

func calculateChecksumForFile(filename string, algorithm string) []byte {

	file, err := os.Open(filename)
	util.CheckErr(err, "Cannot read file "+filename+".")
	defer file.Close()

	calculator := CreateCalculator(algorithm)
	io.Copy(calculator, file)
	checksum := calculator.Sum(nil)[:]

	return checksum
}

func createStringArrayFromFingerprints(fingerprints *list.List) [][]string {

	records := make([][]string, fingerprints.Len())
	index := 0
	for element := fingerprints.Front(); element != nil; element = element.Next() {
		meta := element.Value.(*Fingerprint)
		records[index] = []string{
			meta.Filename, hex.EncodeToString(meta.Checksum), meta.Algorithm,
			meta.CreatedAt, meta.Creator, meta.Note}
		index++
	}

	return records
}

func writeChecksumsToCsvFile(records [][]string, filename string) {

	file, err := os.Create(filename)
	util.CheckErrDontPanic(err, "Cannot read file "+filename)

	writeCsv(records, file)

	defer file.Close()
}

func writeCsv(records [][]string, destination io.Writer) {

	writer := csv.NewWriter(destination)

	// Calls Flush internally.
	writer.WriteAll(records)

	util.CheckErrDontPanic(writer.Error(), "Error writing CSV.")
}
