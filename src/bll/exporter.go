package bll

import (
	"bll/common"
	"dal"
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"util"
)

// Exporter Exports checksums from CSV.
type Exporter struct {
	Db              dal.Database
	OutputDirectory string
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
func NewExporter(db dal.Database, outputDirectory string, basePath string) Exporter {

	basePath = util.NormalizePath(basePath)
	fileHandlers := fileHandlers{nil, nil, nil, nil, nil}

	return Exporter{db, outputDirectory, basePath, fileHandlers}
}

// Convert Converts checksum data to formats that third party utilities understand.
func (exporter *Exporter) Convert(fpFilter common.FingerprintFilter) {

	exporter.Db.LoadFingerprints()
	defer exporter.closeFiles()
	exporter.exportChecksums(fpFilter)
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

func (exporter *Exporter) exportChecksums(fpFilter common.FingerprintFilter) {

	fingerprints := exporter.Db.GetFingerprints()

	for element := fingerprints.Front(); element != nil; element = element.Next() {
		fingerprint := element.Value.(*dal.Fingerprint)
		if fpFilter.FilterFingerprint(fingerprint) {
			checksum := hex.EncodeToString(fingerprint.Checksum)
			fullPath := path.Join(exporter.BasePath, fingerprint.Filename)
			exporter.exportChecksum(fullPath, checksum, fingerprint.Algorithm)
		}
	}
}

func (exporter *Exporter) exportChecksum(filename string, hash string, algorithm string) {

	if algorithm == dal.CRC32 {
		entry := fmt.Sprintf("%s %s\n", filename, hash)
		exporter.saveEntry(&exporter.fileWriters.fCrc32, dal.CRC32EXT, entry)
	} else {
		entry := fmt.Sprintf("%s *%s\n", hash, filename)
		if algorithm == dal.MD5 {
			exporter.saveEntry(&exporter.fileWriters.fMd5, dal.MD5EXT, entry)
		} else if algorithm == dal.SHA1 {
			exporter.saveEntry(&exporter.fileWriters.fSha1, dal.SHA1EXT, entry)
		} else if algorithm == dal.SHA256 {
			exporter.saveEntry(&exporter.fileWriters.fSha256, dal.SHA256EXT, entry)
		} else if algorithm == dal.SHA512 {
			exporter.saveEntry(&exporter.fileWriters.fSha512, dal.SHA512EXT, entry)
		}
	}
}

func (exporter *Exporter) saveEntry(writer **os.File, extension string, entry string) {

	exporter.openOutputFile(writer, extension)
	(*writer).WriteString(entry)
}

func (exporter *Exporter) openOutputFile(writer **os.File, extension string) {

	if *writer == nil {
		fullPath := path.Join(exporter.OutputDirectory, "Checksum"+extension)
		newWriter, err := os.Create(fullPath)
		util.CheckErr(err, "Failed to open output file "+fullPath)

		*writer = newWriter
	}
}
