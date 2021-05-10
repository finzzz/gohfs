package logger

import (
	"log"
)

func LogErr(funcName string, err error) (bool) {
	if err != nil {
		log.Printf("ERROR at %s: %s\n",funcName, err.Error())
		return true
	}

	return false
}