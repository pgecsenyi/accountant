package util

import (
	"container/list"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// ListDirectory Lists the given directory (only the first level of the hierarchy).
func ListDirectory(p string) []os.FileInfo {

	files, err := ioutil.ReadDir(p)
	CheckErr(err, "Cannot list files in directory "+p+".")

	return files
}

// ListDirectoryRecursively Lists the given directory recursively. Returns a single path list that does not contain
// directories.
func ListDirectoryRecursively(p string) []string {

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
