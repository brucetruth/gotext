/*
 * Copyright (c) 2021.  -present, Broos Action, Inc. All rights reserved.
 *
 *  This source code is licensed under the MIT license
 *  found in the LICENSE file in the root directory of this source tree.
 */

package types

import (
	"github.com/broosaction/gotext/tokenizers"
)

type Sentence struct {
	Text     string   `json:"text"`
	Summary  string   `json:"Summary"`
	Meaning  string   `json:"meaning"`
	Words    []Word   `json:"words"`
	Tokens   []string  `json:"tokens"`
	Tokenizer string  `json:"tokenizer"`
}

func  NewSentence(text string) Sentence {
	tokenizer := tokenizers.DefaultTokenizer{}
	return Sentence{
		Text :      text,
		Summary :   "",
		Meaning :   "",
		Words :     []Word{},
		Tokens:     []string{},
		Tokenizer : tokenizer.GetName(),
	}
}

func (s *Sentence) ApplyTokenizer(tokenizer tokenizers.Tokenizer) {

	for _, w := range tokenizer.Tokenize(s.Text) {

			s.Tokens = append(s.Tokens, w)

	}
}

func (s *Sentence) PrepareWords() {

	for _, w := range tokenizers.GetTokenizer(s.Tokenizer).Tokenize(s.Text) {

			word := NewWord(w)
			word.Learn()
			s.Words = append(s.Words, word)

	}
}

func (s *Sentence) PrepareMeaning() *Sentence{
	return s
}

func (s *Sentence) PrepareSummary() *Sentence{

	return s
}

func (s *Sentence) Learn() *Sentence{
      s.PrepareWords()
      s.PrepareMeaning()
      s.PrepareSummary()
      s.ApplyTokenizer(tokenizers.GetTokenizer(s.Tokenizer))
	return s
}