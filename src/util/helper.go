package util

import (
	"fmt"
	"log"
)

func CheckErr(err error, message string) {

	if err != nil {
		if message != "" {
			fmt.Println(message)
		}
		panic(err)
	}
}

func CheckErrDontPanic(err error, message string) {

	if err != nil {
		if message == "" {
			log.Fatalln(err)
		} else {
			log.Fatalln(message)
		}
	}
}
