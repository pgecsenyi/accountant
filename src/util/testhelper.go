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

	return TestHelper{"test"}
}

func (th *TestHelper) AssertPanic(t *testing.T, f func()) {

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic.")
		}
	}()

	f()
}

func (th *TestHelper) CleanUp() {

	os.RemoveAll(th.testFolder)
}

func (th *TestHelper) CreateTestDirectory(directoryPath string) {

	os.Mkdir(th.GetTestPath(directoryPath), 0777)
}

func (th *TestHelper) CreateTestFile(filePath string) {

	createdFile, _ := os.Create(th.GetTestPath(filePath))
	createdFile.Close()
}

func (th *TestHelper) CreateTestRootDirectory() {

	if !CheckIfDirectoryExists(th.testFolder) {
		os.Mkdir(th.testFolder, 0777)
	}
}

func (th *TestHelper) CreateTestFileWithContent(filePath string, content string) {

	outputFile, _ := os.Create(th.GetTestPath(filePath))
	defer outputFile.Close()

	outputFile.WriteString(content)
}

func (th *TestHelper) GetTestPath(filePath string) string {

	return path.Join(th.testFolder, filePath)
}

func (th *TestHelper) GetTestRootDirectory() string {

	return th.testFolder
}

func (th *TestHelper) HasFileInfoValues(files []os.FileInfo, values ...string) bool {

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

func (th *TestHelper) HasStringItems(targetList *list.List, values ...string) bool {

	for element := targetList.Front(); element != nil; element = element.Next() {
		filename := element.Value.(string)
		match := false
		for _, value := range values {
			if filename == value {
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

func (th *TestHelper) HasStringValues(targetSlice []string, values ...string) bool {

	for _, item := range targetSlice {
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
