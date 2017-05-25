package dal

// CRC32 Identifies the CRC32 algorithm.
const CRC32 string = "crc32"

// CRC32EXT Stores the extension of an external file containing CRC32 checksums.
const CRC32EXT string = ".sfv"

// CRC32LEN Stores the length of a CRC32 value.
const CRC32LEN int = 8

// MD5 Identifies the MD5 algorithm.
const MD5 string = "md5"

// MD5EXT Stores the extension of an external file containing MD5 checksums.
const MD5EXT string = ".md5"

// MD5LEN Stores the length of an MD5 hash.
const MD5LEN int = 32

// SHA1 Identifies the SHA-1 algorithm.
const SHA1 string = "sha1"

// SHA1EXT Stores the extension of an external file containing SHA-1 checksums.
const SHA1EXT string = ".sha"

// SHA1LEN Stores the length of an SHA-1 hash.
const SHA1LEN int = 40

// SHA256 Identifies the SHA-256 algorithm.
const SHA256 string = "sha256"

// SHA256EXT Stores the extension of an external file containing SHA-256 checksums.
const SHA256EXT string = ".sha256"

// SHA256LEN Stores the length of an SHA-256 hash.
const SHA256LEN int = 64

// SHA512 Identifies the SHA-512 algorithm.
const SHA512 string = "sha512"

// SHA512EXT Stores the extension of an external file containing SHA-512 checksums.
const SHA512EXT string = ".sha512"

// SHA512LEN Stores the length of an SHA-512 hash.
const SHA512LEN int = 128

// PATTERNCOMMON The Regular Expression for the common file types.
const PATTERNCOMMON string = "^(?P<hash>[a-fA-F0-9]{%d}) ( |\\*)(?P<file>.+)$"

// PATTERNCRC32 The Regular Expression for the CRC32 file types.
const PATTERNCRC32 string = "^(?P<file>.+) (?P<hash>[a-fA-F0-9]{%d})$"
