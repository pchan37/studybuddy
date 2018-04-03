package utils

import (
	"log"
)

func LogFatalfOnError(message string, err error) {
	if err != nil {
		log.Fatalf(message, err)
	}
}

func LogFatalOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
