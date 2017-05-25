package bll

import (
	"checksum"
	"container/list"
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"strings"
	"util"
)

// Exporter Exports checksums from CSV.
type Exporter struct {
	InputChecksums  string
	OutputDirectory string
	Filter          string
	BasePath        string
	fileWriters     fileHandlers
}

type fileHandlers struct {
	fCrc32  *os.File
	fMd5    *os.File
	fSha1   *os.File
	fSha256 *os.File
	fSha512 *os.File
}

// NewExporter Instantiates a new Exporter object.
func NewExporter(inputChecksums string, outputDirectory string, filter string, basePath string) Exporter {

	fileHandlers := fileHandlers{nil, nil, nil, nil, nil}
	return Exporter{inputChecksums, outputDirectory, filter, basePath, fileHandlers}
}

// Convert Converts checksum data to formats that third party utilities understand.
func (exporter *Exporter) Convert(hasher *checksum.FileHasher) {

	hasher.LoadCsv(exporter.InputChecksums)
	defer exporter.closeFiles()
	exporter.exportChecksums(hasher.Fingerprints)
}

func (exporter *Exporter) closeFiles() {

	fw := exporter.fileWriters
	if fw.fCrc32 != nil {
		fw.fCrc32.Close()
	}
	if fw.fMd5 != nil {
		fw.fMd5.Close()
	}
	if fw.fSha1 != nil {
		fw.fSha1.Close()
	}
	if fw.fSha256 != nil {
		fw.fSha256.Close()
	}
	if fw.fSha512 != nil {
		fw.fSha512.Close()
	}
}

func (exporter *Exporter) exportChecksums(fingerprints *list.List) {

	for element := fingerprints.Front(); element != nil; element = element.Next() {
		meta := element.Value.(*checksum.Fingerprint)
		if strings.Contains(meta.Filename, exporter.Filter) {
			checksum := hex.EncodeToString(meta.Checksum)
			fullPath := path.Join(exporter.BasePath, meta.Filename)
			exporter.exportChecksum(fullPath, checksum, meta.Algorithm)
		}
	}
}

func (exporter *Exporter) exportChecksum(filename string, hash string, algorithm string) {

	if algorithm == checksum.CRC32 {
		exporter.openOutputFile(&exporter.fileWriters.fCrc32, checksum.CRC32EXT)
		exporter.fileWriters.fCrc32.WriteString(fmt.Sprintf("%s %s\n", filename, hash))
	} else if algorithm == checksum.MD5 {
		exporter.openOutputFile(&exporter.fileWriters.fMd5, checksum.MD5EXT)
		exporter.fileWriters.fMd5.WriteString(fmt.Sprintf("%s *%s\n", hash, filename))
	} else if algorithm == checksum.SHA1 {
		exporter.openOutputFile(&exporter.fileWriters.fSha1, checksum.SHA1EXT)
		exporter.fileWriters.fSha1.WriteString(fmt.Sprintf("%s *%s\n", hash, filename))
	} else if algorithm == checksum.SHA256 {
		exporter.openOutputFile(&exporter.fileWriters.fSha256, checksum.SHA256EXT)
		exporter.fileWriters.fSha256.WriteString(fmt.Sprintf("%s *%s\n", hash, filename))
	} else if algorithm == checksum.SHA512 {
		exporter.openOutputFile(&exporter.fileWriters.fSha512, checksum.SHA512EXT)
		exporter.fileWriters.fSha512.WriteString(fmt.Sprintf("%s *%s\n", hash, filename))
	}
}

func (exporter *Exporter) openOutputFile(writer **os.File, extension string) {

	if *writer == nil {
		fullPath := path.Join(exporter.OutputDirectory, "Checksum"+extension)
		newWriter, err := os.Create(fullPath)
		*writer = newWriter
		util.CheckErr(err, "Failed to open output file "+fullPath)
	}
}
