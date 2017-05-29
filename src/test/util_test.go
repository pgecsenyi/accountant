package test

import (
	"os"
	"path"
	"testing"
	"util"
)

var testFolder = "test"

func Test_Filesystem(t *testing.T) {

	setupFilesystemTests()
	t.Run("CheckIfDirectoryExists", test_CheckIfDirectoryExists)
	t.Run("CheckIfFileExists", test_CheckIfFileExists)
	t.Run("ListDirectory", test_ListDirectory)
	t.Run("ListFilesRecursively", test_ListFilesRecursively)
	t.Run("NormalizePath", test_NormalizePath)
	t.Run("TrimPath", test_TrimPath)
	tearDownFilesystemTests()
}

func setupFilesystemTests() {

	if !checkIfDirectoryExists(testFolder) {
		os.Mkdir(testFolder, 0777)
	}

	createTestDirectory("dir1")
	createTestDirectory("dir2")
	createTestFile("test.txt")
	createTestFile("dir1/sample1.xml")
	createTestFile("dir1/sample2.png")
	createTestFile("dir2/sample3.jpg")
	createTestFile("dir2/sample4.go")
}

func checkIfDirectoryExists(path string) bool {

	if stat, err := os.Stat(path); err == nil && stat.IsDir() {
		return true
	}

	return false
}

func createTestDirectory(directoryPath string) {

	os.Mkdir(path.Join(testFolder, directoryPath), 0777)
}

func createTestFile(filePath string) {

	createdFile, _ := os.Create(path.Join(testFolder, filePath))
	createdFile.Close()
}

func test_CheckIfDirectoryExists(t *testing.T) {

	testPath1 := path.Join(testFolder, "dir1")
	testPath2 := path.Join(testFolder, "dir3")

	doesExist1 := util.CheckIfDirectoryExists(testPath1)
	doesExist2 := util.CheckIfDirectoryExists(testPath2)

	if !doesExist1 {
		t.Errorf("%s does exist, funcion indicate otherwise.", testPath1)
	}
	if doesExist2 {
		t.Errorf("%s does not exist, funcion indicate otherwise.", testPath2)
	}
}

func test_CheckIfFileExists(t *testing.T) {

	testPath1 := path.Join(testFolder, "test.txt")
	testPath2 := path.Join(testFolder, "test1.txt")

	doesExist1 := util.CheckIfFileExists(testPath1)
	doesExist2 := util.CheckIfFileExists(testPath2)

	if !doesExist1 {
		t.Errorf("%s does exist, funcion indicate otherwise.", testPath1)
	}
	if doesExist2 {
		t.Errorf("%s does not exist, funcion indicate otherwise.", testPath2)
	}
}

func test_ListDirectory(t *testing.T) {

	files := util.ListDirectory(testFolder)

	if len(files) != 3 {
		t.Errorf("Wrong number of nodes listed.")
	}
	if !hasFileInfoValues(files, "dir1", "dir2", "test.txt") {
		t.Errorf("Not all nodes are listed.")
	}
}

func test_ListFilesRecursively(t *testing.T) {

	files := util.ListFilesRecursively(testFolder)

	if len(files) != 5 {
		t.Errorf("Wrong number of files listed.")
	}
	hasValues := hasStringValues(
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

	normalizedPath1 := util.NormalizePath(path1)
	normalizedPath2 := util.NormalizePath(path2)
	normalizedPath3 := util.NormalizePath(path3)

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

	trimmedPath1 := util.TrimPath(path1, "C:\\Temp")
	trimmedPath2 := util.TrimPath(path2, "C:/Temp\\Valami")
	trimmedPath3 := util.TrimPath(path3, "C:/Temp/")
	trimmedPath4 := util.TrimPath(path4, "Files/")

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

func tearDownFilesystemTests() {

	os.RemoveAll(testFolder)
}

func hasFileInfoValues(files []os.FileInfo, values ...string) bool {

	for _, file := range files {
		match := false
		for _, value := range values {
			if file.Name() == value {
				match = true
				break
			}
		}
		if !match {
			return false
		}
	}

	return true
}

func hasStringValues(list []string, values ...string) bool {

	for _, item := range list {
		match := false
		for _, value := range values {
			if item == value {
				match = true
				break
			}
		}
		if !match {
			return false
		}
	}

	return true
}
