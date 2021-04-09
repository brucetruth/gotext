package stringUtils

import (
	"regexp"
	"strings"
)

// StringInSlice looks for a string in slice.
// It returns true or false and position of string in slice (false, -1 if not found).
func StringInSlice(element string, slice []string) (bool, int) {

	// for element in slice
	for index, value := range slice {
		if value == element {
			return true, index
		}
	}

	// return false, placeholder
	return false, -1

}


// cleanup remove none-alnum characters and lowercasize them
func Cleanup(sentence string) string {
	re := regexp.MustCompile("[^a-zA-Z 0-9]+")
	return re.ReplaceAllString(strings.ToLower(sentence), "")
}
