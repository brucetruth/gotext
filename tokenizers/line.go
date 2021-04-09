/*
 * Copyright (c) 2021.  -present, Broos Action, Inc. All rights reserved.
 *
 *  This source code is licensed under the MIT license
 *  found in the LICENSE file in the root directory of this source tree.
 */

package tokenizers

import (
	"bufio"
	stringUtils "goText/utils/strings"
	"io"
	"os"
	"strings"
)

/**
 * A very simple line splitter.
 *
 * [Author]: Bruce Mubangwa
 */

const (
	lINES = "/(?:\r\n|\n\r|\n|\r)/"
)

func Line(sentence string) []string {

	s := stringUtils.Cleanup(sentence)
	tokens := strings.Split(s, lINES)

	return tokens
}

func ReadLines(file string) ([]string, error) {
	return ReadSplitter(file, '\n')
}

func ReadSplitter(file string, splitter byte) (lines []string, err error) {
	fin, err := os.Open(file)
	if err != nil {
		return
	}

	r := bufio.NewReader(fin)
	for {
		line, err := r.ReadString(splitter)
		if err == io.EOF {
			break
		}
		line = strings.Replace(line, string(splitter), "", -1)
		lines = append(lines, line)
	}
	return
}
