/*
 * Copyright (c) 2021.  -present, Broos Action, Inc. All rights reserved.
 *
 *  This source code is licensed under the MIT license
 *  found in the LICENSE file in the root directory of this source tree.
 */

package tokenizers

import (
	stringUtils "github.com/broosaction/gotext/utils/strings"
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
// WordTokenizer is the primary interface for tokenizing words
type WordTokenizer struct{
	RemoveStopWords bool
}

func (wt WordTokenizer)Tokenize(text string) []string{
	safeText := stringUtils.Cleanup(text)
	words := strings.Fields(safeText)
	var tokens []string
	for _, w := range words {
		if wt.RemoveStopWords == true {
			if !stringUtils.IsStopword(w) {
				tokens = append(tokens, w)
			}
		}else{
			tokens = append(tokens, w)
		}
	}
	return tokens
}
