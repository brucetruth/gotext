/*
 * Copyright (c) 2021.  -present, Broos Action, Inc. All rights reserved.
 *
 *  This source code is licensed under the MIT license
 *  found in the LICENSE file in the root directory of this source tree.
 */

package nlptools

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/broosaction/gotext/tokenizers"
	stringUtils "github.com/broosaction/gotext/utils/strings"

	"math"
)

/*
TFIDF is a Term Frequency- Inverse
Document Frequency model that is created
from a trained NaiveBayes model (they
are very similar so you can just train
NaiveBayes and convert into TDIDF)

This is not a probabalistic model, necessarily,
and doesn't give classification. It can be
used to determine the 'importance' of a
word in a document, though, which is
useful in, say, keyword tagging.

Term frequency is basically just adjusted
frequency of a word within a document/sentence:
termFrequency(word, doc) = 0.5 * ( 0.5 * word.Count ) / max{ w.Count | w ∈ doc }

Inverse document frequency is basically how
little the term is mentioned within all of
your documents:
invDocumentFrequency(word, Docs) = log( len(Docs) ) - log( 1 + |{ d ∈ Docs | t ∈ d}| )

TFIDF is the multiplication of those two
functions, giving you a term that is larger
when the word is more important, and less when
the word is less important
*/
// TFIDF tfidf model
type TFIDF struct {
	// train document index in TermFreqs
	DocIndex  map[string]int
	// term frequency for each train document
	TermFreqs []map[string]int
	// documents number for each term in train data
	TermDocs  map[string]int
	// number of documents in train data
	N         int
	// words to be filtered
	StopWords map[string]struct{}
	// tokenizer, space is used as default
	Tokenizer string
}

// New new model with default
func NewTFIDF() *TFIDF {
	return &TFIDF{
		DocIndex:  make(map[string]int),
		TermFreqs: make([]map[string]int, 0),
		TermDocs:  make(map[string]int),
		N:         0,
		Tokenizer: tokenizers.WordTokenizer{}.GetName(),
	}
}

// NewTokenizer new with specified tokenizer
// works well in GOLD
func NewTokenizer(tokenizer tokenizers.Tokenizer) *TFIDF {
	return &TFIDF{
		DocIndex:  make(map[string]int),
		TermFreqs: make([]map[string]int, 0),
		TermDocs:  make(map[string]int),
		N:         0,
		Tokenizer: tokenizer.GetName(),
	}
}

func (f *TFIDF) initStopWords() {
	if f.StopWords == nil {
		f.StopWords = stringUtils.GetStopwords()
	}
}

// AddStopWords add stop words to be filtered
func (f *TFIDF) AddStopWords(words ...string) {
	if f.StopWords == nil {
		f.StopWords = stringUtils.GetStopwords()
	}

	for _, word := range words {
		f.StopWords[word] = struct{}{}
	}
}

// AddStopWordsFile add stop words file to be filtered, with one word a line
func (f *TFIDF) AddStopWordsFile(file string) (err error) {
	lines, err := tokenizers.ReadLines(file)
	if err != nil {
		return
	}

	f.AddStopWords(lines...)
	return
}

// AddDocs add train documents
func (f *TFIDF) AddDocs(docs ...string) {
	for _, doc := range docs {
		h := f.hash(doc)
		if f.docHashPos(h) >= 0 {
			return
		}

		termFreq := f.termFreq(doc)
		if len(termFreq) == 0 {
			return
		}

		f.DocIndex[h] = f.N
		f.N++

		f.TermFreqs = append(f.TermFreqs, termFreq)

		for term := range termFreq {
			f.TermDocs[term]++
		}
	}
}

// Cal calculate tf-idf weight for specified document
func (f *TFIDF) Cal(doc string) (weight map[string]float64) {
	weight = make(map[string]float64)

	var termFreq map[string]int

	docPos := f.docPos(doc)
	if docPos < 0 {
		termFreq = f.termFreq(doc)
	} else {
		termFreq = f.TermFreqs[docPos]
	}

	docTerms := 0
	for _, freq := range termFreq {
		docTerms += freq
	}
	for term, freq := range termFreq {
		weight[term] = f.tfidf(freq, docTerms, f.TermDocs[term], f.N)
	}

	return weight
}

func (f *TFIDF) termFreq(doc string) (m map[string]int) {
	m = make(map[string]int)

	tokens := tokenizers.GetTokenizer(f.Tokenizer).Tokenize(doc)
	if len(tokens) == 0 {
		return
	}

	for _, term := range tokens {
		if _, ok := f.StopWords[term]; ok {
			continue
		}

		m[term]++
	}

	return
}

func (f *TFIDF) docHashPos(hash string) int {
	if pos, ok := f.DocIndex[hash]; ok {
		return pos
	}

	return -1
}

func (f *TFIDF) docPos(doc string) int {
	return f.docHashPos(f.hash(doc))
}

func (f *TFIDF) hash(text string) string {
	h := md5.New()
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}

func (f *TFIDF) tfidf(termFreq, docTerms, termDocs, N int) float64 {
	tf := float64(termFreq) / float64(docTerms)
	idf := math.Log(float64(1+N) / (1 + float64(termDocs)))
	return tf * idf
}
