/*
 * Copyright (c) 2021.  -present, Broos Action, Inc. All rights reserved.
 *
 *  This source code is licensed under the MIT license
 *  found in the LICENSE file in the root directory of this source tree.
 */

package tokenizers

import (
	"fmt"
	stringUtils "github.com/broosaction/gotext/utils/strings"
	"math"
	"strings"
)

/**
 * N-gram
 *
 * N-grams are sequences of n-words of a given string. The N-gram tokenizer
 * outputs tokens of contiguous words ranging from *min* to *max* number of
 * words per token.
 *
 * @category    Machine Learning
 * @author      Bruce Mubangwa
 */
const (
	separator = " "
)

type NGramTokenizer struct{

	/**
	 * The minimum number of contiguous words to a single token.
	 *
	 * @var int
	 */
	Min int

	/**
	 * The maximum number of contiguous words to a single token.
	 *
	 * @var int
	 */
	Max int
}

/**
 * Tokenize a block of text.
 *
 */
func (ng NGramTokenizer) Tokenize(text string) []string {
	var nGrams []string


	words := strings.Fields(stringUtils.Cleanup(text))

	length := len(words)
	for index, word := range words {
		p := math.Min(float64(length-index), float64(ng.Max))

		for j := ng.Min; j <= int(p); j++{
			ngram := word
			for k := 1; k < j; k++ {
				ngram = string(ngram) + separator + string(words[index + k])

			}
			fmt.Println(ngram)
			nGrams = append(nGrams, ngram)
		}
	}


	return nGrams
}
