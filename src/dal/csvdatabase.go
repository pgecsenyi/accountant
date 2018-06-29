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
	fpInputPath        string
	fpOutputPath       string
	namePairOutputPath string
	fingerprints       *list.List
	namePairs          *list.List
}

// NewCsvDatabase Instantiates a new CsvDatabase object.
func NewCsvDatabase(fpInputPath string, fpOutputPath string, namePairOutputPath string) *CsvDatabase {

	return &CsvDatabase{fpInputPath, fpOutputPath, namePairOutputPath, list.New(), list.New()}
}

// AddFingerprint Adds a fingerprint to the database.
func (db *CsvDatabase) AddFingerprint(fingerprint *Fingerprint) {

	if fingerprint != nil {
		db.fingerprints.PushFront(fingerprint)
	}
}

// AddFingerprints Adds a list of fingerprints to the database.
func (db *CsvDatabase) AddFingerprints(fingerprints *list.List) {

	for element := fingerprints.Front(); element != nil; element = element.Next() {
		fingerprint := element.Value.(*Fingerprint)
		db.AddFingerprint(fingerprint)
	}
}

// AddNamePair Adds a name pair to the database.
func (db *CsvDatabase) AddNamePair(namePair *NamePair) {

	if namePair != nil {
		db.namePairs.PushFront(namePair)
	}
}

// Clear Removes all entries from the database.
func (db *CsvDatabase) Clear() {

	db.fingerprints.Init()
	db.namePairs.Init()
}

// GetFingerprints Returns stored fingerprints.
func (db *CsvDatabase) GetFingerprints() *list.List {

	return db.fingerprints
}

// LoadFingerprints Loads fingerprints from the given CSV file.
func (db *CsvDatabase) LoadFingerprints() {

	content := readFileContent(db.fpInputPath)
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

// LoadNamesFromFingeprints Loads the filenames and forwards it to the given StringWriter.
func (db *CsvDatabase) LoadNamesFromFingeprints(writer util.StringWriter) {

	content := readFileContent(db.fpInputPath)
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

// SaveFingerprints Saves fingerprints to the output CSV file.
func (db *CsvDatabase) SaveFingerprints() {

	records := db.createCsvRecords()

	file, err := os.Create(db.fpOutputPath)
	util.CheckErrDontPanic(err, fmt.Sprintf("Cannot write file %s.", db.fpOutputPath))
	defer file.Close()

	writeCsv(records, file)
}

// SaveNamePairs Saves name pairs to a text file.
func (db *CsvDatabase) SaveNamePairs() {

	outputFile, err := os.OpenFile(db.namePairOutputPath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0660)
	util.CheckErr(err, fmt.Sprintf("Cannot write name pairs to %s.", db.namePairOutputPath))
	defer outputFile.Close()

	for element := db.namePairs.Front(); element != nil; element = element.Next() {
		namePair := element.Value.(*NamePair)
		writeNamePair(namePair, outputFile)
	}
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

func writeNamePair(namePair *NamePair, outputFile *os.File) {

	outputFile.WriteString(namePair.NewName + "\r\n")
	outputFile.WriteString("    " + namePair.OldName + "\r\n")
	outputFile.WriteString("    \r\n")
	outputFile.WriteString("    \r\n")
}
