package dal

import (
	"bytes"
	"container/list"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"util"
)

// CsvDatabase Logic for calculating checksums.
type CsvDatabase struct {
	inputPath    string
	outputPath   string
	fingerprints *list.List
}

// NewCsvDatabase Instantiates a new CsvDatabase object.
func NewCsvDatabase(inputPath string, outputPath string) *CsvDatabase {
	return &CsvDatabase{inputPath, outputPath, list.New()}
}

// AddFingerprint Adds a fingerprint to the database.
func (db *CsvDatabase) AddFingerprint(fingerprint *Fingerprint) {

	if fingerprint != nil {
		db.fingerprints.PushFront(fingerprint)
	}
}

// GetFingerprints Returns stored fingerprints.
func (db *CsvDatabase) GetFingerprints() *list.List {

	return db.fingerprints
}

// Load Loads fingerprints from the given CSV file.
func (db *CsvDatabase) Load() {

	content := readFileContent(db.inputPath)
	reader := csv.NewReader(bytes.NewReader(content))

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		util.CheckErrDontPanic(err, "")
		db.addFingerprint(record)
	}
}

// LoadNames Loads the filenames and forwards it to the given StringWriter.
func (db *CsvDatabase) LoadNames(writer util.StringWriter) {

	content := readFileContent(db.inputPath)
	reader := csv.NewReader(bytes.NewReader(content))

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		util.CheckErrDontPanic(err, "")
		writer.Write(record[0])
	}
}

// Save Saves fingerprints to the output CSV file.
func (db *CsvDatabase) Save() {

	records := db.createCsvRecords()

	file, err := os.Create(db.outputPath)
	util.CheckErrDontPanic(err, fmt.Sprintf("Cannot write file %s.", db.outputPath))
	defer file.Close()

	writeCsv(records, file)
}

// SetFingerprints Sets stored fingerprints.
func (db *CsvDatabase) SetFingerprints(fingerprints *list.List) {

	db.fingerprints = fingerprints
}

func (db *CsvDatabase) addFingerprint(record []string) {

	checksumBytes, err := hex.DecodeString(record[1])
	util.CheckErrDontPanic(err, "")

	fingerprint := createFingerprint(record, checksumBytes)
	db.fingerprints.PushFront(fingerprint)
}

func (db *CsvDatabase) createCsvRecords() [][]string {

	records := make([][]string, db.fingerprints.Len())
	index := 0

	for element := db.fingerprints.Front(); element != nil; element = element.Next() {
		fingerprint := element.Value.(*Fingerprint)
		records[index] = createCsvRecord(fingerprint)
		index++
	}

	return records
}

func readFileContent(filename string) []byte {

	content, err := ioutil.ReadFile(filename)
	util.CheckErr(err, fmt.Sprintf("Cannot read file %s.", filename))

	return content
}

func createFingerprint(record []string, checksumBytes []byte) *Fingerprint {

	fingerprint := new(Fingerprint)
	fingerprint.Filename = record[0]
	fingerprint.Checksum = checksumBytes
	fingerprint.Algorithm = record[2]
	fingerprint.CreatedAt = record[3]
	fingerprint.Creator = record[4]
	fingerprint.Note = record[5]

	return fingerprint
}

func createCsvRecord(fingerprint *Fingerprint) []string {

	return []string{
		fingerprint.Filename, hex.EncodeToString(fingerprint.Checksum), fingerprint.Algorithm,
		fingerprint.CreatedAt, fingerprint.Creator, fingerprint.Note}
}

func writeCsv(records [][]string, destination io.Writer) {

	writer := csv.NewWriter(destination)

	// Calls Flush internally.
	writer.WriteAll(records)

	util.CheckErrDontPanic(writer.Error(), "Error writing CSV.")
}
