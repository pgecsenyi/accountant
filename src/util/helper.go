package util

import (
	"fmt"
	"log"
	"os"
)

// CheckErr Displays the given error message if an error has happened and interrupts execution.
func CheckErr(err error, message string) {

	if err != nil {
		if message != "" {
			fmt.Println(message)
		}
		panic(err)
	}
}

// CheckErrDontPanic Displays the given error message if an error has happened.
func CheckErrDontPanic(err error, message string) {

	if err != nil {
		if message == "" {
			log.Fatalln(err)
		} else {
			log.Fatalln(message)
		}
	}
}

// CheckIfFileExists Checks whether the given file exist or not.
func CheckIfFileExists(path string) bool {

	if stat, err := os.Stat(path); err == nil && !os.IsNotExist(err) && !stat.IsDir() {
		return true
	}

	return false
}

// Compare Compares 2 slices of bytes, returns true if they are equal, otherwise returns false.
func Compare(slice1 []byte, slice2 []byte) bool {

	if (slice1 == nil && slice2 != nil) || (slice1 != nil && slice2 == nil) {
		return false
	}
	if len(slice1) != len(slice2) {
		return false
	}

	for i := 0; i < len(slice1); i++ {
		if slice1[i] != slice2[i] {
			return false
		}
	}

	return true
}
