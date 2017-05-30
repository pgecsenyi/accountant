package util

import (
	"container/list"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// CheckIfDirectoryExists Checks whether the given directory exist or not.
func CheckIfDirectoryExists(path string) bool {

	if stat, err := os.Stat(path); err == nil && stat.IsDir() {
		return true
	}

	return false
}

// CheckIfFileExists Checks whether the given file exist or not.
func CheckIfFileExists(path string) bool {

	if stat, err := os.Stat(path); err == nil && !os.IsNotExist(err) && !stat.IsDir() {
		return true
	}

	return false
}

// ListDirectory Lists the given directory (only the first level of the hierarchy).
func ListDirectory(path string) []os.FileInfo {

	files, err := ioutil.ReadDir(path)
	CheckErr(err, fmt.Sprintf("Cannot list files in directory %s.", path))

	return files
}

// ListFilesRecursively Lists the given directory recursively. Returns a single path list that does not contain
// directories.
func ListFilesRecursively(p string) []string {

	resultList := listDirectoryRecursively(p)
	result := make([]string, resultList.Len())

	counter := 0
	for element := resultList.Front(); element != nil; element = element.Next() {
		result[counter] = element.Value.(string)
		counter++
	}

	return result
}

// NormalizePath Normalizes the given path (e.g. replaces each '\' delimiter with a '/').
func NormalizePath(p string) string {

	p = path.Clean(p)
	p = strings.Replace(p, "\\", "/", -1)

	return p
}

// TrimPath Removes the given basePath from fullPath if appropriate.
func TrimPath(fullPath string, basePath string) string {

	normalizedFullPath := NormalizePath(fullPath)
	normalizedBasePath := NormalizePath(basePath)

	if len(normalizedFullPath) < len(normalizedBasePath) {
		return normalizedFullPath
	}

	nFullPathLen := len(normalizedFullPath)
	nBasePathLen := len(normalizedBasePath)

	if normalizedFullPath[:nBasePathLen] == normalizedBasePath {
		if nBasePathLen+1 > nFullPathLen {
			return fullPath[nBasePathLen:]
		}
		return normalizedFullPath[nBasePathLen+1:]
	}

	return normalizedFullPath
}

func listDirectoryRecursively(p string) *list.List {

	result := list.New()

	files, err := ioutil.ReadDir(p)
	CheckErr(err, "Cannot list files in directory "+p+".")

	for _, file := range files {
		if file.IsDir() {
			subDirName := file.Name()
			if subDirName != "." && subDirName != ".." {
				fullSubDirPath := path.Join(p, subDirName)
				subFiles := listDirectoryRecursively(fullSubDirPath)
				mergePathLists(subFiles, result, subDirName)
			}
		} else {
			result.PushFront(file.Name())
		}
	}

	return result
}

func mergePathLists(source *list.List, target *list.List, prefix string) {

	for element := source.Front(); element != nil; element = element.Next() {
		file := element.Value.(string)
		path := path.Join(prefix, file)
		target.PushFront(path)
	}
}
