package util

import "log"

// CheckErr Displays the given error message if an error has happened and interrupts execution.
func CheckErr(err error, message string) {

	if err != nil {
		if message != "" {
			log.Println(message)
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
