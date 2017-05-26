package bll

import (
	"container/list"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"dal"
	"hash"
	"hash/crc32"
	"io"
	"os"
	"path"
	"runtime"
	"time"
	"util"
)

var currentTime = time.Now().Format(time.RFC3339)
var runTimeVersion = runtime.Version()

// Hasher Logic for calculating checksums.
type Hasher struct {
	algorithm string
	hashFunc  hash.Hash
}

// NewHasher Instantiates a new Hasher object.
func NewHasher(algorithm string) Hasher {

	hashFunc := getHashFunc(algorithm)

	return Hasher{algorithm, hashFunc}
}

// CalculateChecksumForFile Calculates checksum for the given file using the given algorithm.
func (hasher *Hasher) CalculateChecksumForFile(filename string) []byte {

	file, err := os.Open(filename)
	util.CheckErr(err, "Cannot read file "+filename+".")
	defer file.Close()

	io.Copy(hasher.hashFunc, file)
	checksum := hasher.hashFunc.Sum(nil)[:]
	hasher.hashFunc.Reset()

	return checksum
}

// CalculateChecksumsForFiles Calculates checksum for each file in the given list.
func (hasher *Hasher) CalculateChecksumsForFiles(basePath string, effectiveBasePath string, files []string) *list.List {

	fingerprints := list.New()
	for _, file := range files {
		hasher.recordChecksumForFile(basePath, effectiveBasePath, file, fingerprints)
	}

	return fingerprints
}

func (hasher *Hasher) recordChecksumForFile(
	basePath string, effectiveBasePath string, filePath string, fingerprints *list.List) {

	fullPath := path.Join(basePath, filePath)
	checksum := hasher.CalculateChecksumForFile(fullPath)

	fp := new(dal.Fingerprint)
	fp.Filename = util.NormalizePath(path.Join(effectiveBasePath, filePath))
	fp.Checksum = checksum
	fp.Algorithm = hasher.algorithm
	fp.CreatedAt = currentTime
	fp.Creator = runTimeVersion
	fp.Note = ""

	fingerprints.PushFront(fp)
}

func getHashFunc(algorithm string) hash.Hash {

	if algorithm == dal.CRC32 {
		return crc32.NewIEEE()
	} else if algorithm == dal.MD5 {
		return md5.New()
	} else if algorithm == dal.SHA1 {
		return sha1.New()
	} else if algorithm == dal.SHA256 {
		return sha256.New()
	} else if algorithm == dal.SHA512 {
		return sha512.New()
	}

	return nil
}
