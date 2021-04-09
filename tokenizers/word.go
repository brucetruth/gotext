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

/**
 * Word
 *
 * This tokenizer matches words with 1 or more characters.
 *
 * @category    Machine Learning
 * [Author]: Bruce Mubangwa
 */
func Word(sentence string) []string {

	s := stringUtils.Cleanup(sentence)
	words := strings.Fields(s)
	var tokens []string
	for _, w := range words {
		if !stringUtils.IsStopword(w) {
			tokens = append(tokens, w)
		}
	}
	return tokens
}
