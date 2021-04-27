/*
 * Copyright (c) 2021.  -present, Broos Action, Inc. All rights reserved.
 *
 *  This source code is licensed under the MIT license
 *  found in the LICENSE file in the root directory of this source tree.
 */

package classifiers

import (
	"fmt"
	"github.com/broosaction/gotext/tokenizers"
	"github.com/broosaction/gotext/utils/types"
	"math"
	"strings"
)


/**
 * Naive-Bayes Classifier
 *
 * This is a naive-bayes classifier that uses Laplace Smoothing.
 *
 * @category    Machine Learning
 * [Author]: Bruce Mubangwa
 */

/*
# `bayes`: A Naive-Bayes classifier for Go


`bayes` takes a document (piece of text), and tells you what category that document belongs to.

## What can I use this for?

You can use this for categorizing any text content into any arbitrary set of **categories**. For example:

- is an email **spam**, or **not spam** ?
- is a news article about **technology**, **politics**, or **sports** ?
- is a piece of text expressing **positive** emotions, or **negative** emotions?

## Usage

    classifier := classifiers.NewNaiveBayes()

	classifier.Learn("amazing, awesome movie!! Yeah!! Oh boy.", "positive")
	classifier.Learn("Sweet, this is incredibly, amazing, perfect, great!!", "positive")
	classifier.Learn("terrible, shitty thing. Damn. Sucks!!", "negative")

	fmt.Println(classifier.Classify("awesome, cool shitty thing"))
 */
type NaiveBayes struct {
	words           map[string]wordFrequency
	classes         map[string]Class
	vocabularySize 	int
	weigh 			weight
	tokenizer 		tokenizers.Tokenizer
}


/**
 * The possible class outcomes.
 */
type Class struct {
	name                    string
	counter                 int
	words                   map[string]types.Word
	Probability             int
	temp_tokenProbabilities float64
}

/**
 * The weight of each class as a proportion of the entire training set.
 *
 * @var array
 */
type weight struct {
	amount float64
	class  Class
}


// wordFrequency stores frequency of words. For example:
// wordFrequency{
//      word: "excellent"
//	counter: map[string]int{
//		"positive": 15
//		"negative": 0
//	}
// }
type wordFrequency struct {
	word    types.Word
	counter map[string]int
}


func NewNaiveBayes() *NaiveBayes {
	c := new(NaiveBayes)
	c.vocabularySize = 	0
	c.words = 			map[string]wordFrequency{}
	c.classes = 		map[string]Class{}
	c.weigh = 			weight{}
	c.tokenizer = 		&tokenizers.DefaultTokenizer{}
	return c
}

/**
 * Initialize each of our data structure entries for this new category or class
 *
 * @param  {String} categoryName
 */
func (nb *NaiveBayes) setClasses(name string) {
	wf, ok := nb.classes[name]
	if !ok {
		wf = Class{
			name:                    name,
			counter:                 0,
			words:                   map[string]types.Word{},
			Probability:             0,
			temp_tokenProbabilities: 0.0000,
		}
	}
	nb.classes[name] = wf

}




/**
 * train our naive-bayes classifier by telling it what `category`
 * the `text` corresponds to.
 */
func (nb *NaiveBayes) Learn(text, class string) {
	nb.setClasses(class)
	//normalize the text into a word array
	tokens := nb.tokenizer.Tokenize(text)

	for _, w := range tokens {
		nb.addWord(w, class)

	}
}

func (nb *NaiveBayes) LearnSentence(sentence types.Sentence, class string) {
	nb.setClasses(class)
	//normalize the Sentence into a word array
	sentence.PrepareWords()

	tokens := sentence.Words

	for _, w := range tokens {
		nb.addWord(w.Text, class)

	}
}

func (nb *NaiveBayes) LearnDocument(document types.Document, class string){
	nb.setClasses(class)
	document.PrepareSentences()
	sentences := document.Sentences
	for _, s := range sentences {
		words := s.Words
		for _, w := range words {
			nb.addWord(w.Text, class)
		}

	}
}

func (nb *NaiveBayes) addWord(word, class string) {
	word = strings.ToLower(word)
	wf, ok := nb.words[word]
	if !ok {
		wf = wordFrequency{word: types.NewWord(word), counter: map[string]int{}}
	}
	wf.counter[class]++
	nb.words[word] = wf
	nb.classes[class].words[word] = types.NewWord(word)
	nb.vocabularySize++

}

/**
 * Calculate probability that a `token` belongs to a `category` or class
 *
 */
func (nb *NaiveBayes) tokenProbability(token, category string) float64 {
	//how many times this word has occurred in documents mapped to this category
	wordFrequencyCount := nb.words[token].counter[category]
	//what is the count of all words that have ever been mapped to this category
	wordCount := len(nb.classes[category].words)

	//use laplace Add-1 Smoothing equation
	return (float64(wordFrequencyCount) + 1) / (float64(wordCount) + float64(nb.vocabularySize))
}

/**
 * Determine what category or class `text` belongs to.
 */
func (nb *NaiveBayes) Classify(text string) (string, float64) {
	text = strings.ToLower(text)
	//var maxProbability = -syscall.INFINITE
	var top_probability = 0.0000
	var chosenCategory string = ""
	tokens := nb.tokenizer.Tokenize(text)

	for _, class := range nb.classes {

		//start by calculating the overall probability of this category
		//=>  out of all documents we've ever looked at, how many were
		//    mapped to this category
		categoryProbability := len(class.words) / len(nb.words)
		class.Probability = categoryProbability

		//take the log to avoid underflow
		logProbability := math.Log(float64(categoryProbability))
		tokens_w := map[string]int{}
		for _, w := range tokens {

			tokens_w[w]++
		}
		for _, w := range tokens {

			frequencyInText := tokens_w[w]
			tokenProbability := nb.tokenProbability(w, class.name)
			fmt.Printf("token: %s \t\tcategory: %s \ttokenProbability: %f \n\n", w, class.name, tokenProbability)

			class.temp_tokenProbabilities += tokenProbability

			//determine the log of the P( w | c ) for this word
			logProbability += float64(frequencyInText) * math.Log(tokenProbability)
		}
		if class.temp_tokenProbabilities > top_probability {
			top_probability = class.temp_tokenProbabilities
			chosenCategory = class.name
		}
         nb.weigh = weight{
         	amount: class.temp_tokenProbabilities,
         	class: class,
		 }
		/*if (logProbability > maxProbability) {
			maxProbability = int(logProbability)
			chosenCategory = class.name
		} */

		//a reset is okay for online learning.
		class.temp_tokenProbabilities = 0.0000
	}
	return chosenCategory, top_probability
}