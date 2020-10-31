package utils

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

// InArray returns true if the string 's' is found in the array 'arr', otherwise false
func InArray(s string, arr []string) bool {
	for _, v := range arr {
		if v == s {
			return true
		}
	}
	return false
}

//GetCurrentDirectory Return the current path
func GetCurrentDirectory() string {
	ex, err := os.Getwd()
	if err != nil {
		log.Error(fmt.Errorf("Could get the path %v", err))
		return os.Getenv("PWD")
	}
	return ex
}

//ConvertToDate convert unix timestamp to dd-mm-YYYY hh:mm:ss
// func ConvertToDate(unixtime int64, location string) string {
// 	date

// return date
// }
