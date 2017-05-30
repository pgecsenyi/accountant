package util

import "testing"
import "os"

var testHelper = NewTestHelper()

func Test_Util(t *testing.T) {

	setupUtilTests()
	t.Run("CheckErr", test_CheckErr)
	t.Run("CheckIfDirectoryExists", test_CheckIfDirectoryExists)
	t.Run("CheckIfFileExists", test_CheckIfFileExists)
	t.Run("ListDirectory", test_ListDirectory)
	t.Run("ListFilesRecursively", test_ListFilesRecursively)
	t.Run("NormalizePath", test_NormalizePath)
	t.Run("TrimPath", test_TrimPath)
	tearDownUtilTests()
}

func setupUtilTests() {

	testHelper.CreateTestRootDirectory()
	testHelper.CreateTestDirectory("dir1")
	testHelper.CreateTestDirectory("dir2")
	testHelper.CreateTestFile("test.txt")
	testHelper.CreateTestFile("dir1/sample1.xml")
	testHelper.CreateTestFile("dir1/sample2.png")
	testHelper.CreateTestFile("dir2/sample3.jpg")
	testHelper.CreateTestFile("dir2/sample4.go")
}

func test_CheckErr(t *testing.T) {

	testHelper.AssertPanic(t, func() {
		CheckErr(os.ErrExist, "")
	})
}

func test_CheckIfDirectoryExists(t *testing.T) {

	testPath1 := testHelper.GetTestPath("dir1")
	testPath2 := testHelper.GetTestPath("dir3")

	doesExist1 := CheckIfDirectoryExists(testPath1)
	doesExist2 := CheckIfDirectoryExists(testPath2)

	if !doesExist1 {
		t.Errorf("%s does exist, funcion indicate otherwise.", testPath1)
	}
	if doesExist2 {
		t.Errorf("%s does not exist, funcion indicate otherwise.", testPath2)
	}
}

func test_CheckIfFileExists(t *testing.T) {

	testPath1 := testHelper.GetTestPath("test.txt")
	testPath2 := testHelper.GetTestPath("test1.txt")

	doesExist1 := CheckIfFileExists(testPath1)
	doesExist2 := CheckIfFileExists(testPath2)

	if !doesExist1 {
		t.Errorf("%s does exist, funcion indicate otherwise.", testPath1)
	}
	if doesExist2 {
		t.Errorf("%s does not exist, funcion indicate otherwise.", testPath2)
	}
}

func test_ListDirectory(t *testing.T) {

	files := ListDirectory(testHelper.GetTestRootDirectory())

	if len(files) != 3 {
		t.Errorf("Wrong number of nodes listed.")
	}
	if !testHelper.HasFileInfoValues(files, "dir1", "dir2", "test.txt") {
		t.Errorf("Not all nodes are listed.")
	}
}

func test_ListFilesRecursively(t *testing.T) {

	files := ListFilesRecursively(testHelper.GetTestRootDirectory())

	if len(files) != 5 {
		t.Errorf("Wrong number of files listed.")
	}
	hasValues := testHelper.HasStringValues(
		files,
		"test.txt", "dir1/sample1.xml", "dir1/sample2.png", "dir2/sample3.jpg", "dir2/sample4.go")
	if !hasValues {
		t.Errorf("Not all files are listed.")
	}
}

func test_NormalizePath(t *testing.T) {

	path1 := ""
	path2 := "C:\\Temp/An interesting directory\\somefile.go"
	path3 := "/usr\\local/src"

	normalizedPath1 := NormalizePath(path1)
	normalizedPath2 := NormalizePath(path2)
	normalizedPath3 := NormalizePath(path3)

	if normalizedPath1 != "." {
		t.Errorf("\"%s\" is not the expected normalized result.", normalizedPath1)
	}
	if normalizedPath2 != "C:/Temp/An interesting directory/somefile.go" {
		t.Errorf("\"%s\" is not the expected normalized result.", normalizedPath2)
	}
	if normalizedPath3 != "/usr/local/src" {
		t.Errorf("\"%s\" is not the expected normalized result.", normalizedPath3)
	}
}

func test_TrimPath(t *testing.T) {

	path1 := ""
	path2 := "C:\\Temp/Test.xml"
	path3 := "C:\\Temp/An interesting directory\\somefile.go"
	path4 := "Files"

	trimmedPath1 := TrimPath(path1, "C:\\Temp")
	trimmedPath2 := TrimPath(path2, "C:/Temp\\Valami")
	trimmedPath3 := TrimPath(path3, "C:/Temp/")
	trimmedPath4 := TrimPath(path4, "Files/")

	if trimmedPath1 != "." {
		t.Errorf("\"%s\" is not the expected trimmed path.", trimmedPath1)
	}
	if trimmedPath2 != "C:/Temp/Test.xml" {
		t.Errorf("\"%s\" is not the expected trimmed path.", trimmedPath2)
	}
	if trimmedPath3 != "An interesting directory/somefile.go" {
		t.Errorf("\"%s\" is not the expected trimmed path.", trimmedPath3)
	}
	if trimmedPath4 != "" {
		t.Errorf("\"%s\" is not the expected trimmed path.", trimmedPath4)
	}
}

func tearDownUtilTests() {

	testHelper.CleanUp()
}
