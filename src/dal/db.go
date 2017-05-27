package dal

import (
	"bytes"
	"container/list"
	"encoding/csv"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
	"util"
)

// Db Logic for calculating checksums.
type Db struct {
	Fingerprints *list.List
}

// NewDb Instantiates a new Db object.
func NewDb() Db {
	return Db{list.New()}
}

// LoadCsv Loads fingerprints from the given CSV file.
func (fh *Db) LoadCsv(filename string) {

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

// LoadNamesFromCsv Loads the filenames from the given CSV and forwards it to the given StringWriter.
func (fh *Db) LoadNamesFromCsv(filename string, writer util.StringWriter) {

	content := readFileContent(filename)
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

// SaveCsv Exports fingerprints to the given CSV file.
func (fh *Db) SaveCsv(filename string) {

	records := createStringArrayFromFingerprints(fh.Fingerprints)

	file, err := os.Create(filename)
	util.CheckErrDontPanic(err, "Cannot read file "+filename)
	defer file.Close()

	writeCsv(records, file)
}

func readFileContent(filename string) []byte {

	content, err := ioutil.ReadFile(filename)
	util.CheckErr(err, "Cannot read file "+filename+".")

	return content
}

func (fh *Db) addFingerprint(record []string, checksumBytes []byte) {

	fp := new(Fingerprint)
	fp.Filename = record[0]
	fp.Checksum = checksumBytes
	fp.Algorithm = record[2]
	fp.CreatedAt = record[3]
	fp.Creator = record[4]
	fp.Note = record[5]

	fh.Fingerprints.PushFront(fp)
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

func writeCsv(records [][]string, destination io.Writer) {

	writer := csv.NewWriter(destination)

	// Calls Flush internally.
	writer.WriteAll(records)

	util.CheckErrDontPanic(writer.Error(), "Error writing CSV.")
}
