/*
 * Copyright (c) 2021.  -present, Broos Action, Inc. All rights reserved.
 *
 *  This source code is licensed under the MIT license
 *  found in the LICENSE file in the root directory of this source tree.
 */

package tokenizers

import (
	stringUtils "github.com/broosaction/gotext/utils/strings"
	"regexp"
	"strings"
)
type DefaultTokenizer struct {
	RemoveStopWords bool

}

/**
 * Given an input string, tokenize it into an array of word tokens.
 * This is the default tokenization function used if user does not provide one in `options`.
 *
 * @param  {String} text
 * @return {Array}
 * [Author]: Bruce Mubangwa
 */
func (d *DefaultTokenizer) Tokenize(text string) []string {
	//remove punctuation from text - remove anything that isn't a word char or a space
	rgxPunctuation := "[^(a-zA-ZA-Яa-я0-9_)+\\s]"


	r, _ := regexp.Compile(rgxPunctuation)

	safeText := r.ReplaceAllString(text, " ")
	words := strings.Fields(safeText)
	var tokens []string
	for _, w := range words {
		if d.RemoveStopWords == true {
			if !stringUtils.IsStopword(w) {
				tokens = append(tokens, w)
			}
		}else{
			tokens = append(tokens, w)
		}

	}
	return tokens
}
