package stringUtils

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
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

// generateUUID  generate simple uuid.
func generateUUID() string {
	b := make([]byte, 16)
	rand.Read(b)
	uuid := fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
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
