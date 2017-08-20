package bll

import (
	"bll/testutil"
	"dal"
	"testing"
)

func TestVerifier(t *testing.T) {

	setupVerifierTests()

	t.Run("Verify", testVerifierVerify)
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

	// Arrange.
	fp1 := testutil.CreateSparseFingerprint("test.txt", "1d291cf2", "crc32")
	fp2 := testutil.CreateSparseFingerprint("dir1/test.txt", "6b24cc6a", "crc32")
	fp3 := testutil.CreateSparseFingerprint("hello.world", "a1b2c3d4", "crc32")
	memoryDatabase := dal.NewMemoryDatabase()
	memoryDatabase.AddFingerprint(fp1)
	memoryDatabase.AddFingerprint(fp2)
	memoryDatabase.AddFingerprint(fp3)
	testPath := testHelper.GetTestRootDirectory()
	verifier := NewVerifier(memoryDatabase, "calculation", testPath)

	// Act.
	verifier.Verify(false)

	// Assert.
	if !testHelper.HasStringItems(verifier.Report.CorruptFiles, "test.txt", "dir1/test.txt") {
		t.Error("File should be marked as corrupt: \"test.txt\", \"dir1/test.txt\".")
	}
	if !testHelper.HasStringItems(verifier.Report.MissingFiles, "hello.world") {
		t.Error("File should be marked as missing: \"hello.world\".")
	}
}

func testVerifierVerifyNamesOnly(t *testing.T) {

	// Arrange.
	fp1 := testutil.CreateSparseFingerprint("test.txt", "1d291cf2", "crc32")
	fp2 := testutil.CreateSparseFingerprint("hello.world", "a1b2c3d4", "crc32")
	memoryDatabase := dal.NewMemoryDatabase()
	memoryDatabase.AddFingerprint(fp1)
	memoryDatabase.AddFingerprint(fp2)
	testPath := testHelper.GetTestRootDirectory()
	verifier := NewVerifier(memoryDatabase, testPath, testPath)

	// Act.
	verifier.Verify(true)

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
