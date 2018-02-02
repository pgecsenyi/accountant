package util

import (
	"container/list"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestTestHelper(t *testing.T) {

	setupTestHelperTests()

	t.Run("testCreateTestFileWithContent", testCreateTestFileWithContent)
	t.Run("CreateTestRootDirectory", testCreateTestRootDirectory)
	t.Run("CreateTestDirectory", testCreateTestDirectory)
	t.Run("GetTestDirectory", testGetTestDirectory)
	t.Run("HasFileInfoValues", testHasFileInfoValues)
	t.Run("HasStringItems", testHasStringItems)
	t.Run("HasStringValues", testHasStringValues)
	t.Run("RemoveTestDirectory", testRemoveTestDirectory)

	tearDownTestHelperTests()
}

func setupTestHelperTests() {

	testHelper.CreateTestRootDirectory()
	testHelper.CreateTestDirectory("dir1")
	testHelper.CreateTestFile("test.txt")
	testHelper.CreateTestFile("dir1/sample1.xml")
}

func testCreateTestFileWithContent(t *testing.T) {

	filename := "test1.txt"
	testHelper.CreateTestFileWithContent(filename, "Hello World!")

	content, err := ioutil.ReadFile(testHelper.GetTestPath(filename))
	if err == nil {
		lines := strings.Split(string(content), "\n")

		if len(lines) != 1 || lines[0] != "Hello World!" {
			t.Errorf("Unexpected content in file: \"%s\".", filename)
		}
	} else {
		t.Errorf("Cannot read the file: \"%s\": \"%s\".", filename, err)
	}
}

func testCreateTestRootDirectory(t *testing.T) {

	stat, err := os.Stat("testfiles")
	if err != nil {
		t.Error("Test root directory does not exist with the given name.")
	}
	if !stat.IsDir() {
		t.Error("Test root directory is not a directory.")
	}
	// if stat.Mode()&0777 == 0777 {
	// 	t.Errorf("The test root directory does not have the right permission setting: %d.", stat.Mode()&0777)
	// }
}

func testCreateTestDirectory(t *testing.T) {

	stat, err := os.Stat("testfiles/dir1")
	if err != nil {
		t.Error("Test directory does not exist with the given name.")
	}
	if !stat.IsDir() {
		t.Error("Test directory is not a directory.")
	}
	// if stat.Mode() != 0777 {
	// 	t.Error("The test directory does not have the right permission setting.")
	// }
}

func testGetTestDirectory(t *testing.T) {

	dirname := testHelper.GetTestDirectory("test2")
	if dirname != "testfiles/test2" {
		t.Error("Test directory name is invalid.")
	}
}

func testHasFileInfoValues(t *testing.T) {

	dirname := testHelper.GetTestDirectory("dir1")
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		t.Errorf("Cannot list directory: %s.", dirname)
	}

	if !testHelper.HasFileInfoValues(files, "test.txt", "sample1.xml") {
		t.Error("Should contain \"test.txt\" and \"sample1.xml\".")
	}
	if testHelper.HasFileInfoValues(files, "something.csv") {
		t.Error("Should not contain \"something.csv\".")
	}
}

func testHasStringItems(t *testing.T) {

	testList := list.New()
	testList.PushFront("a")
	testList.PushFront("b")

	if !testHelper.HasStringItems(testList, "a", "b") {
		t.Error("Should contain \"a\" and \"b\".")
	}
	if testHelper.HasStringItems(testList, "c") {
		t.Error("Should not contain \"c\".")
	}
}

func testHasStringValues(t *testing.T) {

	testSlice := make([]string, 2)
	testSlice[0] = "a"
	testSlice[1] = "b"

	if !testHelper.HasStringValues(testSlice, "a", "b") {
		t.Error("Should contain \"a\" and \"b\".")
	}
	if testHelper.HasStringValues(testSlice, "c") {
		t.Error("Should not contain \"c\".")
	}
}

func testRemoveTestDirectory(t *testing.T) {

	testHelper.RemoveTestDirectory("dir1")
	stat, err := os.Stat("dir1")

	if err == nil {
		t.Error("Unable to check if test directory \"dir1\" exists.")
	}
	if stat != nil && stat.IsDir() {
		t.Error("The test directory \"dir1\" should not exist.")
	}
}

func tearDownTestHelperTests() {

	testHelper.CleanUp()
}
