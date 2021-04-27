/*
 * Copyright (c) 2021.  -present, Broos Action, Inc. All rights reserved.
 *
 *  This source code is licensed under the MIT license
 *  found in the LICENSE file in the root directory of this source tree.
 */

package transformers

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/broosaction/gotext/tokenizers"

	"log"
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
	docIndex  map[string]int         // train document index in TermFreqs
	termFreqs []map[string]int       // term frequency for each train document
	termDocs  map[string]int         // documents number for each term in train data
	n         int                    // number of documents in train data
	stopWords map[string]interface{} // words to be filtered
	tokenizer tokenizers.Tokenizer         // tokenizer, space is used as default
}

// New new model with default
func NewTFIDF() *TFIDF {
	return &TFIDF{
		docIndex:  make(map[string]int),
		termFreqs: make([]map[string]int, 0),
		termDocs:  make(map[string]int),
		n:         0,
		tokenizer: &tokenizers.WordTokenizer{},
	}
}

// NewTokenizer new with specified tokenizer
// works well in GOLD
func NewTokenizer(tokenizer tokenizers.Tokenizer) *TFIDF {
	return &TFIDF{
		docIndex:  make(map[string]int),
		termFreqs: make([]map[string]int, 0),
		termDocs:  make(map[string]int),
		n:         0,
		tokenizer: tokenizer,
	}
}

func (f *TFIDF) initStopWords() {
	if f.stopWords == nil {
		f.stopWords = make(map[string]interface{})
	}

	lines, err := tokenizers.ReadLines("../data/stopword")
	if err != nil {
		log.Printf("init stop words with error: %s", err)
	}

	for _, w := range lines {
		f.stopWords[w] = nil
	}
}

// AddStopWords add stop words to be filtered
func (f *TFIDF) AddStopWords(words ...string) {
	if f.stopWords == nil {
		f.stopWords = make(map[string]interface{})
	}

	for _, word := range words {
		f.stopWords[word] = nil
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
		h := hash(doc)
		if f.docHashPos(h) >= 0 {
			return
		}

		termFreq := f.termFreq(doc)
		if len(termFreq) == 0 {
			return
		}

		f.docIndex[h] = f.n
		f.n++

		f.termFreqs = append(f.termFreqs, termFreq)

		for term := range termFreq {
			f.termDocs[term]++
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
		termFreq = f.termFreqs[docPos]
	}

	docTerms := 0
	for _, freq := range termFreq {
		docTerms += freq
	}
	for term, freq := range termFreq {
		weight[term] = tfidf(freq, docTerms, f.termDocs[term], f.n)
	}

	return weight
}

func (f *TFIDF) termFreq(doc string) (m map[string]int) {
	m = make(map[string]int)

	tokens := f.tokenizer.Tokenize(doc)
	if len(tokens) == 0 {
		return
	}

	for _, term := range tokens {
		if _, ok := f.stopWords[term]; ok {
			continue
		}

		m[term]++
	}

	return
}

func (f *TFIDF) docHashPos(hash string) int {
	if pos, ok := f.docIndex[hash]; ok {
		return pos
	}

	return -1
}

func (f *TFIDF) docPos(doc string) int {
	return f.docHashPos(hash(doc))
}

func hash(text string) string {
	h := md5.New()
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}

func tfidf(termFreq, docTerms, termDocs, N int) float64 {
	tf := float64(termFreq) / float64(docTerms)
	idf := math.Log(float64(1+N) / (1 + float64(termDocs)))
	return tf * idf
}
