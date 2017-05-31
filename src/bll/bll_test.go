package bll

import (
	"container/list"
	"dal"
	"encoding/hex"
	"testing"
	"time"
	"util"
)

var testHelper = util.NewTestHelper()

func Test_EffectiveTextMemory(t *testing.T) {

	t.Run("Cache", test_Etm_Default)
	t.Run("ClearCache", test_Etm_ClearCache)
	t.Run("NoCache", test_Etm_NoCache)
}

func test_Etm_Default(t *testing.T) {

	etm := newEffectiveTextMemory()

	etm.Write("test/a.txt")
	etm.Write("test/a.txt")
	etm.Write("test/b.txt")

	if etm.CountAll != 2 {
		t.Errorf("Wrong number of all entries: %d.", etm.CountAll)
	}
	if etm.CountCollisions != 0 {
		t.Errorf("Wrong number of collisions: %d.", etm.CountCollisions)
	}
	testIfEtmContainsText(t, etm, "test/a.txt", true)
	testIfEtmContainsText(t, etm, "test/b.txt", true)
	testIfEtmContainsText(t, etm, "test/something/c.txt", false)
}

func test_Etm_ClearCache(t *testing.T) {

	etm := newEffectiveTextMemory()

	etm.Write("test/a.txt")
	etm.ClearCache()
	etm.Write("test/a.txt")
	etm.Write("test/b.txt")

	if etm.CountAll != 3 {
		t.Errorf("Wrong number of all entries: %d.", etm.CountAll)
	}
	if etm.CountCollisions != 1 {
		t.Errorf("Wrong number of collisions: %d.", etm.CountCollisions)
	}
}

func test_Etm_NoCache(t *testing.T) {

	etm := newEffectiveTextMemory()
	etm.UseCache = false

	etm.Write("test/a.txt")
	etm.Write("test/a.txt")
	etm.Write("test/b.txt")

	if etm.CountAll != 3 {
		t.Errorf("Wrong number of all entries: %d.", etm.CountAll)
	}
	if etm.CountCollisions != 1 {
		t.Errorf("Wrong number of collisions: %d.", etm.CountCollisions)
	}
}

func Test_Other(t *testing.T) {

	setupOtherTests()
	t.Run("Hasher_CalculateChecksum_Crc32", test_Hasher_CalculateChecksum_Crc32)
	t.Run("Hasher_CalculateChecksum_Md5", test_Hasher_CalculateChecksum_Md5)
	t.Run("Hasher_CalculateChecksum_Sha1", test_Hasher_CalculateChecksum_Sha1)
	t.Run("Hasher_CalculateChecksum_Sha256", test_Hasher_CalculateChecksum_Sha256)
	t.Run("Hasher_CalculateChecksum_Sha512", test_Hasher_CalculateChecksum_Sha512)
	t.Run("Hasher_CalculateFingerprint", test_Hasher_CalculateFingerprint)
	t.Run("Hasher_CalculateFingerprints", test_Hasher_CalculateFingerprints)
	t.Run("Calculator_Calculate_All", test_Calculator_All)
	t.Run("Calculator_Calculate_MissingOnly", test_Calculator_MissingOnly)
	t.Run("Comparer_Compare", test_Comparer_Compare)
	tearDownOtherTests()
}

func setupOtherTests() {

	testHelper.CreateTestRootDirectory()

	testHelper.CreateTestDirectory("calculation")
	testHelper.CreateTestDirectory("calculation/dir1")
	testHelper.CreateTestFileWithContent("calculation/test.txt", "Hello World!")
	testHelper.CreateTestFileWithContent("calculation/dir1/test.txt", "Lorem ipsum, dolor sit amet.")

	testHelper.CreateTestDirectory("comparison")
	testHelper.CreateTestDirectory("comparison/dir1")
	testHelper.CreateTestFileWithContent("comparison/test2.txt", "Hello World!")
	testHelper.CreateTestFileWithContent("comparison/dir1/test.txt", "Lorem ipsum, dolor sit amet.")
}

func test_Hasher_CalculateChecksum_Crc32(t *testing.T) {

	testChecksumCalculation(t, "crc32", "1c291ca3")
}

func test_Hasher_CalculateChecksum_Md5(t *testing.T) {

	testChecksumCalculation(t, "md5", "ed076287532e86365e841e92bfc50d8c")
}

func test_Hasher_CalculateChecksum_Sha1(t *testing.T) {

	testChecksumCalculation(t, "sha1", "2ef7bde608ce5404e97d5f042f95f89f1c232871")
}

func test_Hasher_CalculateChecksum_Sha256(t *testing.T) {

	testChecksumCalculation(
		t,
		"sha256",
		"7f83b1657ff1fc53b92dc18148a1d65dfc2d4b1fa3d677284addd200126d9069")
}

func test_Hasher_CalculateChecksum_Sha512(t *testing.T) {

	testChecksumCalculation(
		t,
		"sha512",
		"861844d6704e8573fec34d967e20bcfef3d424cf48be04e6dc08f2bd58c729743371015ead891cc3cf1c9d34b49264b510751b1ff9e537937bc46b5d6ff4ecc8")
}

func test_Hasher_CalculateFingerprint(t *testing.T) {

	// Arrange
	hasher := NewHasher("crc32")

	// Act.
	startTime := time.Now()
	fingerprint := hasher.CalculateFingerprint(testHelper.GetTestPath("calculation"), "", "test.txt")
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

func test_Hasher_CalculateFingerprints(t *testing.T) {

	// Arrange.
	expectedFingerprints := getExpectedFingerprintsForBasicCalculation()
	hasher := NewHasher("crc32")

	// Act.
	fingerprints := hasher.CalculateFingerprints(
		testHelper.GetTestPath("calculation"),
		"",
		[]string{"test.txt", "dir1/test.txt"})

	// Assert.
	if fingerprints.Len() != 2 {
		t.Errorf("Wrong number of items in result set: %d.", fingerprints.Len())
	}
	testCalculatedFingerprints(t, fingerprints, expectedFingerprints)
}

func test_Calculator_All(t *testing.T) {

	// Arrange.
	expectedFingerprints := getExpectedFingerprintsForBasicCalculation()
	memoryDatabase := dal.NewMemoryDatabase()
	testPath := testHelper.GetTestPath("calculation")
	calculator := NewCalculator(memoryDatabase, testPath, "crc32", testPath)

	// Act.
	calculator.Calculate(false)

	// Assert.
	if memoryDatabase.Fingerprints.Len() != 2 {
		t.Errorf("Wrong number of items in result set: %d.", memoryDatabase.Fingerprints.Len())
	}
	testCalculatedFingerprints(t, memoryDatabase.Fingerprints, expectedFingerprints)
}

func test_Calculator_MissingOnly(t *testing.T) {

	// Arrange.
	expectedFingerprints := getExpectedFingerprintsForBasicCalculation()
	fp1 := createFingerprint("test.txt", "1c291ca3", "crc32")
	memoryDatabase := dal.NewMemoryDatabase()
	memoryDatabase.AddFingerprint(fp1)
	testPath := testHelper.GetTestPath("calculation")
	calculator := NewCalculator(memoryDatabase, testPath, "crc32", testPath)

	// Act.
	calculator.Calculate(true)

	// Assert.
	if memoryDatabase.Fingerprints.Len() != 1 {
		t.Errorf("Wrong number of items in result set: %d.", memoryDatabase.Fingerprints.Len())
	}
	testCalculatedFingerprints(t, memoryDatabase.Fingerprints, expectedFingerprints)
}

func test_Comparer_Compare(t *testing.T) {

	// Arrange.
	expectedFingerprints := getExpectedFingerprintsForComparison()
	fp1 := createFingerprint("test.txt", "1c291ca3", "crc32")
	fp2 := createFingerprint("dir1/test.txt", "6b24cc6a", "crc32")
	memoryDatabase := dal.NewMemoryDatabase()
	memoryDatabase.AddFingerprint(fp1)
	memoryDatabase.AddFingerprint(fp2)
	testPath := testHelper.GetTestPath("comparison")
	comparer := Comparer{memoryDatabase, testPath, testPath}

	// Act.
	comparer.Compare("crc32")

	// Assert.
	testCalculatedFingerprints(t, memoryDatabase.Fingerprints, expectedFingerprints)
}

func tearDownOtherTests() {

	testHelper.CleanUp()
}

func testIfEtmContainsText(t *testing.T, etm *effectiveTextMemory, text string, shouldContain bool) {

	if shouldContain && !etm.ContainsText(text) {
		t.Errorf("Should contain text: \"%s\".", text)
	} else if !shouldContain && etm.ContainsText(text) {
		t.Errorf("Should not contain text: \"%s\".", text)
	}
}

func testCalculatedFingerprints(t *testing.T, fingerprints *list.List, expectedFingerprints *list.List) {

	for element := fingerprints.Front(); element != nil; element = element.Next() {
		fingerprint := element.Value.(*dal.Fingerprint)
		if !isFingerprintInList(fingerprint, expectedFingerprints) {
			t.Errorf("Fingerprint for file \"%s\" (filename or checksum) is unexpected.", fingerprint.Filename)
		}
	}
}

func testChecksumCalculation(t *testing.T, algorithm string, expectedChecksum string) {

	hasher := NewHasher(algorithm)
	checksumBytes := hasher.CalculateChecksum(testHelper.GetTestPath("calculation/test.txt"))
	checksum := hex.EncodeToString(checksumBytes)

	if checksum != expectedChecksum {
		t.Errorf("Wrong %s checksum: %s.", algorithm, checksum)
	}
}

func createFingerprint(filename string, checksum string, algorithm string) *dal.Fingerprint {

	checksumBytes, err := hex.DecodeString(checksum)
	util.CheckErr(err, "Unable to convert checksum from string.")

	return &dal.Fingerprint{filename, checksumBytes, algorithm, "", "", ""}
}

func createList(items ...interface{}) *list.List {

	result := list.New()
	for _, item := range items {
		result.PushFront(item)
	}

	return result
}

func getExpectedFingerprintsForBasicCalculation() *list.List {

	fp1 := createFingerprint("test.txt", "1c291ca3", "crc32")
	fp2 := createFingerprint("dir1/test.txt", "6b24cc6a", "crc32")
	expectedFingerprints := createList(fp1, fp2)

	return expectedFingerprints
}

func getExpectedFingerprintsForComparison() *list.List {

	fp1 := createFingerprint("test2.txt", "1c291ca3", "crc32")
	fp2 := createFingerprint("dir1/test.txt", "6b24cc6a", "crc32")
	expectedFingerprints := createList(fp1, fp2)

	return expectedFingerprints
}

func isFingerprintInList(fingerprint *dal.Fingerprint, fingerprints *list.List) bool {

	for element := fingerprints.Front(); element != nil; element = element.Next() {
		fp := element.Value.(*dal.Fingerprint)
		if fingerprint.Filename == fp.Filename &&
			fingerprint.Algorithm == fp.Algorithm &&
			util.CompareByteSlices(fingerprint.Checksum, fp.Checksum) {
			return true
		}
	}

	return false
}
