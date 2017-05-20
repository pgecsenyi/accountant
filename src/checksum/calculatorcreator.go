package checksum

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"hash/crc32"
)

// CRC32 Identifies the CRC32 algorithm.
const CRC32 string = "crc32"

// MD5 Identifies the MD5 algorithm.
const MD5 string = "md5"

// SHA1 Identifies the SHA-1 algorithm.
const SHA1 string = "sha1"

// SHA256 Identifies the SHA-256 algorithm.
const SHA256 string = "sha256"

// SHA512 Identifies the SHA-512 algorithm.
const SHA512 string = "sha512"

// CreateCalculator Creates the appropriate ChecksumCalculator for the given algorithm.
func CreateCalculator(algorithm string) hash.Hash {

	if algorithm == CRC32 {
		return crc32.NewIEEE()
	} else if algorithm == MD5 {
		return md5.New()
	} else if algorithm == SHA1 {
		return sha1.New()
	} else if algorithm == SHA256 {
		return sha256.New()
	} else if algorithm == SHA512 {
		return sha512.New()
	}

	return nil
}
