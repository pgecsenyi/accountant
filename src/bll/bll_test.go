package bll

import (
	"dal"
	"encoding/hex"
	"testing"
	"time"
	"util"
)

var testHelper = util.NewTestHelper()

func Test_Hasher(t *testing.T) {

	setupHasherTests()
	t.Run("CalculateChecksum_Crc32", test_CalculateChecksum_Crc32)
	t.Run("CalculateChecksum_Md5", test_CalculateChecksum_Md5)
	t.Run("CalculateChecksum_Sha1", test_CalculateChecksum_Sha1)
	t.Run("CalculateChecksum_Sha256", test_CalculateChecksum_Sha256)
	t.Run("CalculateChecksum_Sha512", test_CalculateChecksum_Sha512)
	t.Run("CalculateFingerprint", test_CalculateFingerprint)
	t.Run("CalculateFingerprints", test_CalculateFingerprints)
	tearDownHasherTests()
}

func setupHasherTests() {

	testHelper.CreateTestRootDirectory()
	testHelper.CreateTestDirectory("dir1")
	testHelper.CreateTestFileWithContent("test.txt", "Hello World!")
	testHelper.CreateTestFileWithContent("dir1/test.txt", "Lorem ipsum, dolor sit amet.")
}

func test_CalculateChecksum_Crc32(t *testing.T) {

	testChecksumCalculation(t, "crc32", "1c291ca3")
}

func test_CalculateChecksum_Md5(t *testing.T) {

	testChecksumCalculation(t, "md5", "ed076287532e86365e841e92bfc50d8c")
}

func test_CalculateChecksum_Sha1(t *testing.T) {

	testChecksumCalculation(t, "sha1", "2ef7bde608ce5404e97d5f042f95f89f1c232871")
}

func test_CalculateChecksum_Sha256(t *testing.T) {

	testChecksumCalculation(t, "sha256", "7f83b1657ff1fc53b92dc18148a1d65dfc2d4b1fa3d677284addd200126d9069")
}

func test_CalculateChecksum_Sha512(t *testing.T) {

	testChecksumCalculation(
		t,
		"sha512",
		"861844d6704e8573fec34d967e20bcfef3d424cf48be04e6dc08f2bd58c729743371015ead891cc3cf1c9d34b49264b510751b1ff9e537937bc46b5d6ff4ecc8")
}

func test_CalculateFingerprint(t *testing.T) {

	// Arrange
	hasher := NewHasher("crc32")

	// Act.
	startTime := time.Now()
	fingerprint := hasher.CalculateFingerprint(testHelper.GetTestRootDirectory(), "", "test.txt")
	endTime := time.Now()

	// Assert.
	checksum := hex.EncodeToString(fingerprint.Checksum)
	createdAt, _ := time.Parse(time.RFC3339, fingerprint.CreatedAt)

	if fingerprint.Filename != "test.txt" {
		t.Errorf("Wrong filename: %s.", fingerprint.Filename)
	}
	if checksum != "1c291ca3" {
		t.Errorf("Wrong checksum: %s.", checksum)
	}
	if fingerprint.Algorithm != "crc32" {
		t.Errorf("Wrong algorithm: %s.", fingerprint.Algorithm)
	}
	if endTime.Unix() < createdAt.Unix() || createdAt.Unix() < startTime.Unix() {
		t.Errorf("Wrong timestamp: %s.", fingerprint.CreatedAt)
	}
	if fingerprint.Creator != util.RuntimeVersion {
		t.Errorf("Wrong creator: %s.", fingerprint.Creator)
	}
	if fingerprint.Note != "" {
		t.Errorf("Wrong note: %s.", fingerprint.Note)
	}
}

func test_CalculateFingerprints(t *testing.T) {

	// Arrange.
	hasher := NewHasher("crc32")

	// Act.
	fingerprints := hasher.CalculateFingerprints(
		testHelper.GetTestRootDirectory(),
		"",
		[]string{"test.txt", "dir1/test.txt"})

	// Assert.
	if fingerprints.Len() != 2 {
		t.Errorf("Wrong number of items in result set: %d.", fingerprints.Len())
	}

	for element := fingerprints.Front(); element != nil; element = element.Next() {
		fingerprint := element.Value.(*dal.Fingerprint)
		checksum := hex.EncodeToString(fingerprint.Checksum)
		expectedChecksum := ""

		if fingerprint.Filename == "test.txt" {
			expectedChecksum = "1c291ca3"
		} else if fingerprint.Filename == "dir1/test.txt" {
			expectedChecksum = "6b24cc6a"
		} else {
			t.Errorf("Wrong filename: %s.", fingerprint.Filename)
		}

		if checksum != expectedChecksum {
			t.Errorf("Wrong %s checksum for file %s.", checksum, fingerprint.Filename)
		}
	}
}

func tearDownHasherTests() {

	testHelper.CleanUp()
}

func testChecksumCalculation(t *testing.T, algorithm string, expectedChecksum string) {

	hasher := NewHasher(algorithm)
	checksumBytes := hasher.CalculateChecksum(testHelper.GetTestPath("test.txt"))
	checksum := hex.EncodeToString(checksumBytes)

	if checksum != expectedChecksum {
		t.Errorf("Wrong %s checksum: %s.", algorithm, checksum)
	}
}
