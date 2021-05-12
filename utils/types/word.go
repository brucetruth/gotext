/*
 * Copyright (c) 2021.  -present, Broos Action, Inc. All rights reserved.
 *
 *  This source code is licensed under the MIT license
 *  found in the LICENSE file in the root directory of this source tree.
 */

package types

import (
	"fmt"
	"github.com/broosaction/gotext/nlp/pos"
	"github.com/broosaction/gotext/tokenizers"
	"github.com/broosaction/gotext/transformers"
)

type Word struct {
	Text     string   `json:"text"`
	PosTag   string   `json:"pos"`
	Stem     string   `json:"stem"`
	Vector   float64  `json:"vector"`
	Meaning  string   `json:"meaning"`
	Letters  []string `json:"letters"`
	Synonyms []string `json:"synonyms"`
}


func  NewWord(text string) Word {
	return Word{
		Text : text,
		PosTag : "",
		Stem : "",
		Vector : 0,
		Meaning : "",
		Letters : []string{},
		Synonyms : []string{},
	}
}

func (w *Word) AttachPOStag() {
	//w.PosTag = pos.Tag(w.Text)[0][1]
	if m, e := pos.LoadModel("./data/models/pos/en"); e != nil {
		fmt.Println(e)
		pos.BuildModel(w.Text)
	} else {
		var tagger *pos.Tagger = pos.NewTagger(m)
		var tokens []string = tokenizers.WordTokenizer{}.Tokenize(w.Text)
		var tagged [][]string = tagger.Tag(tokens)

		for _, token := range tagged {
			//fmt.Printf("%q ", token)
			w.PosTag = token[1]
		}
	}
}

func (w *Word) ApplyStemmer() *Word{
	w.Stem = string(transformers.Stem([]byte(w.Text)))
	return w
}

func (w *Word) PrepareMeaning() *Word{

	return w
}

func (w *Word) PrepareSynonyms() *Word{

	return w
}

func (w *Word) prepareLetters(){

}

func (w *Word) Learn() *Word{
   w.AttachPOStag()
   w.ApplyStemmer()
   w.PrepareMeaning()
   w.PrepareSynonyms()
   w.prepareLetters()
  return w
}