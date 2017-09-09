package util

import (
	"container/list"
	"os"
	"path"
	"testing"
)

// TestHelper Contains common operations needed for testing.
type TestHelper struct {
	testFolder string
}

// NewTestHelper Instantiates a new NewTestHelper object.
func NewTestHelper() TestHelper {

	return TestHelper{"testfiles"}
}

// AssertPanic Asserts if the given function panics.
func (th *TestHelper) AssertPanic(t *testing.T, f func()) {

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic.")
		}
	}()

	f()
}

// CleanUp Deletes all files in test folder as well the test folder itself.
func (th *TestHelper) CleanUp() {

	os.RemoveAll(th.testFolder)
}

// CreateTestDirectory Creates a test directory under the root path.
func (th *TestHelper) CreateTestDirectory(directoryPath string) {

	os.Mkdir(th.GetTestPath(directoryPath), 0777)
}

// CreateTestFile Creates a test file.
func (th *TestHelper) CreateTestFile(filePath string) {

	createdFile, _ := os.Create(th.GetTestPath(filePath))
	createdFile.Close()
}

// CreateTestFileWithContent Creates a test file with the given content.
func (th *TestHelper) CreateTestFileWithContent(filePath string, content string) {

	outputFile, _ := os.Create(th.GetTestPath(filePath))
	defer outputFile.Close()

	outputFile.WriteString(content)
}

// CreateTestRootDirectory Creates the root test folder.
func (th *TestHelper) CreateTestRootDirectory() {

	if !CheckIfDirectoryExists(th.testFolder) {
		os.Mkdir(th.testFolder, 0777)
	}
}

// GetTestDirectory Gets the path of a test folder.
func (th *TestHelper) GetTestDirectory(directory string) string {

	return path.Join(th.testFolder, directory)
}

// GetTestPath Gets the path for the given file under the test folder.
func (th *TestHelper) GetTestPath(filePath string) string {

	return path.Join(th.testFolder, filePath)
}

// GetTestRootDirectory Gets the path of the root test folder.
func (th *TestHelper) GetTestRootDirectory() string {

	return th.testFolder
}

// HasFileInfoValues Checks whether the given slice has all the os.FileInfo values provided.
func (th *TestHelper) HasFileInfoValues(targetSlice []os.FileInfo, values ...string) bool {

	for _, file := range targetSlice {
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

// HasStringItems Checks whether the given list has all the string values provided.
func (th *TestHelper) HasStringItems(targetList *list.List, expectedValues ...string) bool {

	for _, expectedValue := range expectedValues {
		match := false
		for element := targetList.Front(); element != nil; element = element.Next() {
			strElement := element.Value.(string)
			if strElement == expectedValue {
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

// HasStringValues Checks whether the given slice has all the string values provided.
func (th *TestHelper) HasStringValues(targetSlice []string, expectedValues ...string) bool {

	for _, expectedValue := range expectedValues {
		match := false
		for _, item := range targetSlice {
			if item == expectedValue {
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

// RemoveTestDirectory Deletes all files in the specified test folder.
func (th *TestHelper) RemoveTestDirectory(testDirectory string) {

	os.RemoveAll(th.GetTestPath(testDirectory))
}
