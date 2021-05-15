/*
 * Copyright (c) 2021.  -present, Broos Action, Inc. All rights reserved.
 *
 *  This source code is licensed under the MIT license
 *  found in the LICENSE file in the root directory of this source tree.
 */

package types

import (
	"fmt"
	"github.com/broosaction/gotext/nlp/nlptools/en"
	"github.com/broosaction/gotext/nlp/pos"
	"github.com/broosaction/gotext/tokenizers"
	"strings"
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
		Text : 		text,
		PosTag : 	"",
		Stem : 		"",
		Vector : 	0,
		Meaning : 	"",
		Letters : 	[]string{},
		Synonyms : 	[]string{},
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

//Applies a stemmer to turn words like Cats to Cat
func (w *Word) ApplyStemmer() *Word{
	w.Stem = string(en.Stem([]byte(w.Text)))
	return w
}

func (w *Word) PrepareMeaning() *Word{

	return w
}

// Applies words with similer meaning
func (w *Word) PrepareSynonyms() *Word{

	return w
}

//explodes the words in an array of characters
func (w *Word) prepareLetters(){
  w.Letters = strings.Split(w.Text,"")
}

// Preloads all other needed fields
func (w *Word) Learn() *Word{
   w.AttachPOStag()
   w.ApplyStemmer()
   w.PrepareMeaning()
   w.PrepareSynonyms()
   w.prepareLetters()
  return w
}