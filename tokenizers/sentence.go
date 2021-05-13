/*
 * Copyright (c) 2021.  -present, Broos Action, Inc. All rights reserved.
 *
 *  This source code is licensed under the MIT license
 *  found in the LICENSE file in the root directory of this source tree.
 */

package tokenizers

import (
	"github.com/broosaction/gotext/tokenizers/sentences/english"
)

type SentenceTokenizer struct {
	RemoveStopWords bool
}

var SentenceTokenizerName = "SentenceTokenizer"

func (s *SentenceTokenizer) Tokenize(text string) []string {


	tokenizer, err := english.NewSentenceTokenizer(nil)
	if err != nil {
		panic(err)
	}

	sentences := tokenizer.Tokenize(text)
	var tokens []string
	for _, s := range sentences {
		tokens = append(tokens, s.Text)
	}
  return tokens
}

func (s *SentenceTokenizer) GetName() string {
	return SentenceTokenizerName
}