package checksum

// CRC32LEN Stores the length of a CRC32 value.
const CRC32LEN int = 8

// MD5LEN Stores the length of an MD5 hash.
const MD5LEN int = 32

// SHA1LEN Stores the length of an SHA-1 hash.
const SHA1LEN int = 40

// SHA256LEN Stores the length of an SHA-256 hash.
const SHA256LEN int = 64

// SHA512LEN Stores the length of an SHA-512 hash.
const SHA512LEN int = 128

// PATTERNCOMMON The Regular Expression for the common file types.
const PATTERNCOMMON string = "^(?P<hash>[a-fA-F0-9]{%d}) \\*{0,1}(?P<file>.+)$"

// PATTERNCRC32 The Regular Expression for the CRC32 file types.
const PATTERNCRC32 string = "^(?P<file>.+) (?P<hash>[a-fA-F0-9]{%d})$"
