/*
 * Copyright (c) 2021.  -present, Broos Action, Inc. All rights reserved.
 *
 *  This source code is licensed under the MIT license
 *  found in the LICENSE file in the root directory of this source tree.
 */

package tokenizers

type Tokenizer interface {
	// Compute the output value.
	Tokenize(string) []string

    GetName() string
}

func GetTokenizer(name string) Tokenizer{

	switch name {
	case DefaultTokenizerName:
		return &DefaultTokenizer{}
	case LineTokenizerName:
		return &LineTokenizer{}
	case NGramTokenizerName:
		return &NGramTokenizer{}
	case ParagraphTokenizerName:
		return &ParagraphTokenizer{}
	case SentenceTokenizerName:
		return &SentenceTokenizer{}
	case WordTokenizerName:
		return &WordTokenizer{}
	case WhitespaceTokenizerName:
		return &WhitespaceTokenizer{}
	}

	return nil
}