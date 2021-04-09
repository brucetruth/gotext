/*
 * Copyright (c) 2021.  -present, Broos Action, Inc. All rights reserved.
 *
 *  This source code is licensed under the MIT license
 *  found in the LICENSE file in the root directory of this source tree.
 */

package tokenizers

import (
	stringUtils "goText/utils/strings"
	"strings"
)

type Whitespace struct{
	/**
	 * The whitespace character that delimits each token.
	 *
	 * @var string
	 */
	delimiter string
}

/**
 * Whitespace
 *
 * Separate each token by a user-specified delimiter such as a single
 * whitespace.
 * tokenize create an array of words from a sentence
 * [Author]: Bruce Mubangwa
**/

func (wp Whitespace) Whitespace(sentence string) []string {

	if wp.delimiter == "" {
		wp.delimiter = " "
	}
	s := stringUtils.Cleanup(sentence)
	words := strings.Split(s, wp.delimiter)
	var tokens []string
	for _, w := range words {
		if !stringUtils.IsStopword(w) {
			tokens = append(tokens, w)
		}
	}
	return tokens
}

