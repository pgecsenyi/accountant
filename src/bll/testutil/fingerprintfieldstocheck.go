package testutil

// FingerprintFieldsToCheck Indicates which fields of a fingerprint have to be checked. Used for testing purposes.
type FingerprintFieldsToCheck struct {
	Filename  bool
	Checksum  bool
	Algorithm bool
	CreatedAt bool
	Creator   bool
	Note      bool
}

// NewFingerprintFieldsToCheck Instantiates a new FingerprintFieldsToCheck object.
func NewFingerprintFieldsToCheck(createdAt bool, creator bool, note bool) FingerprintFieldsToCheck {

	return FingerprintFieldsToCheck{true, true, true, createdAt, creator, note}
}
