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
	tokenizer tokenizers.Tokenizer
}

func  NewSentence(text string) Sentence {
	return Sentence{
		Text : text,
		Summary : "",
		Meaning : "",
		Words : []Word{},
		Tokens: []string{},
		tokenizer : &tokenizers.DefaultTokenizer{},
	}
}

func (s *Sentence) ApplyTokenizer(tokenizer tokenizers.Tokenizer) {
	s.tokenizer = tokenizer
	for _, w := range tokenizer.Tokenize(s.Text) {

			s.Tokens = append(s.Tokens, w)

	}
}

func (s *Sentence) PrepareWords() {
	tokenizer := s.tokenizer
	for _, w := range tokenizer.Tokenize(s.Text) {

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
      s.ApplyTokenizer(s.tokenizer)
	return s
}