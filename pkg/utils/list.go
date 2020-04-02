package utils

import (
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// IntArrayMatch tells whether array arr containing item s
func IntArrayMatch(arr []int, s int) bool { return IntArrayIndex(arr, s) >= 0 }

// IntArrayIndex tells whether array arr containing item s
func IntArrayIndex(arr []int, s int) int {
	for i, a := range arr {
		if a == s {
			return i
		}
	}
	return -1
}

// IntArrayUniqAppend returns
func IntArrayUniqAppend(needSort bool, arr []int, ss ...int) []int {
	newA := arr
	for _, s := range ss {
		if !IntArrayMatch(newA, s) {
			newA = append(newA, s)
		}
	}
	if needSort {
		sort.Ints(newA)
	}
	return newA
}

// IArrayUniqAppend returns
func IArrayUniqAppend(arr []interface{}, ss ...interface{}) []interface{} {
	newA := arr
	for _, s := range ss {
		if !IArrayMatch(newA, s) {
			newA = append(newA, s)
		}
	}
	return newA
}

// IArrayIndex tells whether array arr containing item s
func IArrayIndex(arr []interface{}, s interface{}) int {
	for i, a := range arr {
		if reflect.DeepEqual(a, s) {
			return i
		}
	}
	return -1
}

// IArrayMatch tells whether array arr containing item s
func IArrayMatch(arr []interface{}, s interface{}) bool { return IArrayIndex(arr, s) >= 0 }

// ArrayUniqAppend returns
func ArrayUniqAppend(needSort bool, removeBlank bool, arr []string, ss ...string) []string {
	// ret := IArrayUniqAppend(arr, ss...)
	// return ret.([]string)
	newA := arr
	for _, s := range ss {
		if removeBlank && strings.TrimSpace(s) == "" {
			continue
		}
		if !ArrayMatch(newA, s) {
			newA = append(newA, s)
		}
	}
	if needSort {
		sort.Strings(newA)
	}
	return newA
}

// ArrayMatch tells whether array arr containing item s
func ArrayMatch(arr []string, s string) bool { return ArrayIndex(arr, s) >= 0 }

// ArrayIndex tells whether array arr containing item s
func ArrayIndex(arr []string, s string) int {
	for i, a := range arr {
		// if (a == "" && s == "") || (a != "" && s != "") {
		if a == s {
			// FIXME: maybe some functions will be affeced
			// if matched, _ := regexp.Match(s, []byte(a)); matched || a == s {
			return i
		}
		// }
	}
	return -1
}

// ArrayItemCount tells how many times that item s duplicated in array arr
func ArrayItemCount(arr []string, s string) int {
	counter := make(map[string]int)
	for _, a := range arr {
		if assert.EqualValues(new(testing.T), a, s) {
			counter[s]++
		}
	}
	return counter[s]
}

// ArrayRemoveDupli returns a new string array which duplicated items has been removed
func ArrayRemoveDupli(arr []string, trimSpace bool) []string {
	var newA []string

	for _, a := range arr {
		if trimSpace {
			if strings.TrimSpace(a) != "" && !ArrayMatch(newA, a) {
				newA = append(newA, a)
			}
		} else if !ArrayMatch(newA, a) {
			newA = append(newA, a)
		}
	}
	return newA
}

// IntArrayRemoveDupli returns a new string array which duplicated items has been removed
func IntArrayRemoveDupli(arr []int) []int {
	var newA []int

	for _, a := range arr {
		if !IntArrayMatch(newA, a) {
			newA = append(newA, a)
		}
	}
	return newA
}
func UniqAppendString(arr []string, s ...string) []string {
	var list []string
	copy(arr, list)

	for _, single := range s {
		duplicated := false
		for _, a := range list {
			if a == single {
				duplicated = true
				break
			}
		}
		if !duplicated {
			list = append(list, single)
		}
	}
	return list
}
func UniqAppendInt(arr []int, s ...int) []int {
	var list []int
	copy(arr, list)

	for _, single := range s {
		duplicated := false
		for _, a := range list {
			if a == single {
				duplicated = true
				break
			}
		}
		if !duplicated {
			list = append(list, single)
		}
	}
	return list
}
