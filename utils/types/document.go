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
	Tokenizer 	  string
}

func  NewDocument(text string) Document {
	tokenizer := tokenizers.SentenceTokenizer{}
	return Document{
		Text 			:    	text,
		Summary	    	: 		"",
		Meaning     	:     	"",
		Sentences   	:       []Sentence{},
		Tokens      	:      	[]string{},
		WordFrequence   :    	FreqDist{},
		Tokenizer   	:       tokenizer.GetName(),
	}
}

func (d *Document) PrepareSentences() {

	for _, w := range tokenizers.GetTokenizer(d.Tokenizer).Tokenize(d.Text) {

			sentence := NewSentence(w)
			sentence.Learn()
			d.Sentences = append(d.Sentences, sentence)

	}
}

func (d *Document) computeWordFreq() {
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
	d.computeWordFreq()
	d.PrepareSentences()
	return d
}
