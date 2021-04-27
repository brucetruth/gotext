package stringUtils

import (
	"bufio"
	"io"
	"io/ioutil"
	"regexp"
	"strconv"
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

// determines if the string represents a number.
func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// cleanup remove none-alnum characters and lowercasize them
func Cleanup(sentence string) string {
	re := regexp.MustCompile("[^a-zA-Z 0-9]+")
	return re.ReplaceAllString(strings.ToLower(sentence), "")
}

// ReadFile returns the bytes of a file searched in the path and beyond it
func ReadFile(path string) (bytes []byte) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		bytes, err = ioutil.ReadFile("../" + path)
	}

	if err != nil {
		panic(err)
	}

	return bytes
}

func ReaderToString(r io.Reader,bufferSize int) chan string {
	tokenizer := bufio.NewScanner(r)
	tokenizer.Split(bufio.ScanWords)
	tokens := make(chan string, bufferSize)

	go func() {
		for tokenizer.Scan() {
			tokens <- tokenizer.Text()
		}
		close(tokens)
	}()

	return tokens
}