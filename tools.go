package main

import "strings"

func stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func checkIfStringExists(strs, str string) bool {
	elements := strings.Split(strs, ",")
	for _, element := range elements {
		if element == str {
			return true
		}
	}
	return false
}
