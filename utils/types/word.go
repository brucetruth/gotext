/*
 * Copyright (c) 2021.  -present, Broos Action, Inc. All rights reserved.
 *
 *  This source code is licensed under the MIT license
 *  found in the LICENSE file in the root directory of this source tree.
 */

package types

import (
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
	//todo
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

func (w *Word) Learn() *Word{
   w.AttachPOStag()
   w.ApplyStemmer()
   w.PrepareMeaning()
   w.PrepareSynonyms()

  return w
}