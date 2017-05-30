package dal

// Fingerprint Stores the necessary data to identify a file and a bit more.
type Fingerprint struct {
	Filename  string
	Checksum  []byte
	Algorithm string
	CreatedAt string
	Creator   string
	Note      string
}

// NamePair Stores old name - new name pairs.
type NamePair struct {
	NewName string
	OldName string
}
