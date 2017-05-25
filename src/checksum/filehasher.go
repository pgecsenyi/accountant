package checksum

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

// FileHasher Logic for calculating checksums.
type FileHasher struct {
	Fingerprints *list.List
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

// NewFileHasher Instantiates a new FileHasher object.
func NewFileHasher() FileHasher {
	return FileHasher{list.New()}
}

// LoadCsv Loads fingerprints from the given CSV file.
func (fh *FileHasher) LoadCsv(filename string) {

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

// SaveCsv Exports fingerprints to the given CSV file.
func (fh *FileHasher) SaveCsv(filename string) {

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
