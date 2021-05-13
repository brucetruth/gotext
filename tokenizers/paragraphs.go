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
 * A very simple paragraph tokenizer.
 *
 * [Author]: Bruce Mubangwa
 */

type ParagraphTokenizer struct {}

const (
	regx = "/(?:\r\n|\n\r|\n|\r)/"
)

var ParagraphTokenizerName = "ParagraphTokenizer"

func (w *ParagraphTokenizer) Tokenize(text string) []string {

	s := stringUtils.Cleanup(text)
	words := strings.Split(s, regx)

	var tokens []string
	for _, w := range words {
		if !stringUtils.IsStopword(w) {
			tokens = append(tokens, w)
		}
	}
	return tokens
}

func (w *ParagraphTokenizer) GetName() string {
	return ParagraphTokenizerName
}