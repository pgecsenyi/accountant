package bll

import (
	"fmr/bll/common"
	"fmr/bll/testutil"
	"fmr/dal"
	"strings"
	"testing"
)

func TestVerifier(t *testing.T) {

	setupVerifierTests()

	t.Run("Verify", testVerifierVerify)
	t.Run("Verify_Filtered", testVerifierVerifyFiltered)
	t.Run("Verify_NamesOnly", testVerifierVerifyNamesOnly)

	tearDownVerifierTests()
}

func setupVerifierTests() {

	testHelper.CreateTestRootDirectory()

	testHelper.CreateTestDirectory("dir1")
	testHelper.CreateTestFileWithContent("test.txt", "Hello World!")
	testHelper.CreateTestFileWithContent("dir1/test.txt", "Lorem ipsum, dolor sit amet.")
}

func testVerifierVerify(t *testing.T) {

	expectedCorruptFiles := []string{"test.txt"}
	expectedMissingFiles := []string{"hello.world"}
	fpFilter := common.NewFingerprintFilter("")
	testVerifierWithFilter(t, fpFilter, expectedCorruptFiles, expectedMissingFiles)
}

func testVerifierVerifyFiltered(t *testing.T) {

	expectedCorruptFiles := []string{"test.txt"}
	expectedMissingFiles := []string{}
	fpFilter := common.NewFingerprintFilter(".txt")
	testVerifierWithFilter(t, fpFilter, expectedCorruptFiles, expectedMissingFiles)
}

func testVerifierVerifyNamesOnly(t *testing.T) {

	// Arrange.
	fp1 := testutil.CreateSparseFingerprint("test.txt", "1d291cf2", "crc32")
	fp2 := testutil.CreateSparseFingerprint("hello.world", "a1b2c3d4", "crc32")
	memoryDatabase := dal.NewMemoryDatabase()
	memoryDatabase.AddFingerprint(fp1)
	memoryDatabase.AddFingerprint(fp2)
	testPath := testHelper.GetTestRootDirectory()
	verifier := NewVerifier(memoryDatabase, testPath)
	fpFilter := common.NewFingerprintFilter("")

	// Act.
	verifier.Verify(true, fpFilter)

	// Assert.
	if verifier.Report.CorruptFiles.Len() != 0 {
		t.Error("Checksum verification is not supposed to be performed this time.")
	}
	if !testHelper.HasStringItems(verifier.Report.MissingFiles, "hello.world") {
		t.Error("File should be marked as missing: \"hello.world\".")
	}
}

func tearDownVerifierTests() {

	testHelper.CleanUp()
}

func testVerifierWithFilter(
	t *testing.T, fpFilter common.FingerprintFilter, expectedCorruptFiles []string, expectedMissingFiles []string) {

	// Arrange.
	fp1 := testutil.CreateSparseFingerprint("test.txt", "1d291cf2", "crc32")
	fp2 := testutil.CreateSparseFingerprint("dir1/test.txt", "6b24cc6a", "crc32")
	fp3 := testutil.CreateSparseFingerprint("hello.world", "a1b2c3d4", "crc32")
	memoryDatabase := dal.NewMemoryDatabase()
	memoryDatabase.AddFingerprint(fp1)
	memoryDatabase.AddFingerprint(fp2)
	memoryDatabase.AddFingerprint(fp3)
	testPath := testHelper.GetTestRootDirectory()
	verifier := NewVerifier(memoryDatabase, testPath)

	// Act.
	verifier.Verify(false, fpFilter)

	// Assert.
	if len(expectedCorruptFiles) > 0 {
		if !testHelper.HasStringItems(verifier.Report.CorruptFiles, expectedCorruptFiles...) {
			t.Errorf("File should be marked as corrupt: %s.", strings.Join(expectedCorruptFiles, ", "))
		}
	}
	if len(expectedMissingFiles) > 0 {
		if !testHelper.HasStringItems(verifier.Report.MissingFiles, expectedMissingFiles...) {
			t.Errorf("File should be marked as missing: %s.", strings.Join(expectedMissingFiles, ", "))
		}
	}
}
