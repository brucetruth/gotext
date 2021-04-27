/*
 * Copyright (c) 2021.  -present, Broos Action, Inc. All rights reserved.
 *
 *  This source code is licensed under the MIT license
 *  found in the LICENSE file in the root directory of this source tree.
 */

package types

import (
	"github.com/broosaction/gotext/tokenizers"
	stringUtils "github.com/broosaction/gotext/utils/strings"
)

type Document struct {
	Text     	  string   `json:"text"`
	Summary  	  string   `json:"Summary"`
	Meaning 	  string   `json:"meaning"`
	Sentences     []Sentence   `json:"sentences"`
	Tokens   	  []string `json:"tokens"`
	WordFrequence FreqDist
	tokenizer 	  tokenizers.Tokenizer
}

func  NewDocument(text string) Document {
	return Document{
		Text 			: text,
		Summary	    	: "",
		Meaning     	: "",
		Sentences   	: []Sentence{},
		Tokens      	: []string{},
		WordFrequence   : FreqDist{},
		tokenizer   	: &tokenizers.SentenceTokenizer{},
	}
}

func (d *Document) PrepareSentences() {
	tokenizer := d.tokenizer
	for _, w := range tokenizer.Tokenize(d.Text) {

			sentence := NewSentence(w)
			sentence.Learn()
			d.Sentences = append(d.Sentences, sentence)

	}
}

func (d *Document) compute_word_freq() {
	tokenizer := tokenizers.DefaultTokenizer{}
	samples := map[string]int{}
	for _, w := range tokenizer.Tokenize(stringUtils.Cleanup(d.Text)) {

		wf, ok := samples[w]
		if !ok {
			wf = 0
		}
		wf++
		samples[w] = wf

	}

	d.WordFrequence.Samples = samples

}

func (d *Document) Learn() *Document{
	d.compute_word_freq()
	d.PrepareSentences()
	return d
}
