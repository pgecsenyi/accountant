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
	"time"
	"util"
)

// Hasher Logic for calculating checksums.
type Hasher struct {
	algorithm string
	hashFunc  hash.Hash
}

// NewHasher Instantiates a new Hasher object.
func NewHasher(algorithm string) Hasher {

	hashFunc := createHashFunc(algorithm)

	return Hasher{algorithm, hashFunc}
}

// CalculateChecksum Calculates the checksum of the given file.
func (hasher *Hasher) CalculateChecksum(filename string) []byte {

	return hasher.calculateChecksum(filename)
}

// CalculateFingerprint Calculates fingerprint for the given file.
func (hasher *Hasher) CalculateFingerprint(basePath string, effectiveBasePath string, file string) *dal.Fingerprint {

	currentTime := time.Now().Format(time.RFC3339)
	fingerprint := hasher.calculateFingerprint(basePath, effectiveBasePath, file, currentTime)

	return fingerprint
}

// CalculateFingerprints Calculates fingerprint for each file in the given list.
func (hasher *Hasher) CalculateFingerprints(basePath string, effectiveBasePath string, files []string) *list.List {

	currentTime := time.Now().Format(time.RFC3339)
	fingerprints := list.New()

	for _, file := range files {
		fingerprint := hasher.calculateFingerprint(basePath, effectiveBasePath, file, currentTime)
		fingerprints.PushFront(fingerprint)
	}

	return fingerprints
}

func (hasher *Hasher) calculateChecksum(filename string) []byte {

	file, err := os.Open(filename)
	util.CheckErr(err, "Cannot read file "+filename+".")
	defer file.Close()

	io.Copy(hasher.hashFunc, file)
	checksum := hasher.hashFunc.Sum(nil)[:]
	hasher.hashFunc.Reset()

	return checksum
}

func (hasher *Hasher) calculateFingerprint(
	basePath string, effectiveBasePath string, file string, currentTime string) *dal.Fingerprint {

	fullPath := path.Join(basePath, file)
	checksum := hasher.calculateChecksum(fullPath)
	effectivePath := util.NormalizePath(path.Join(effectiveBasePath, file))
	fingerprint := hasher.createFingerprint(effectivePath, checksum, currentTime)

	return fingerprint
}

func (hasher *Hasher) createFingerprint(file string, checksum []byte, currentTime string) *dal.Fingerprint {

	fp := new(dal.Fingerprint)
	fp.Filename = file
	fp.Checksum = checksum
	fp.Algorithm = hasher.algorithm
	fp.CreatedAt = currentTime
	fp.Creator = util.RuntimeVersion
	fp.Note = ""

	return fp
}

func createHashFunc(algorithm string) hash.Hash {

	if algorithm == dal.CRC32 {
		return crc32.NewIEEE()
	} else if algorithm == dal.MD5 {
		return md5.New()
	} else if algorithm == dal.SHA256 {
		return sha256.New()
	} else if algorithm == dal.SHA512 {
		return sha512.New()
	}

	return sha1.New()
}
