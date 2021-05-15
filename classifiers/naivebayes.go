/*
 * Copyright (c) 2021.  -present, Broos Action, Inc. All rights reserved.
 *
 *  This source code is licensed under the MIT license
 *  found in the LICENSE file in the root directory of this source tree.
 */

package classifiers

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/broosaction/gotext/tokenizers"
	"github.com/broosaction/gotext/utils/persist"
	"github.com/broosaction/gotext/utils/types"
	"go.uber.org/zap"
	"log"
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
	tokenizer 		string//tokenizers.Tokenizer
}


/**
 * The possible class outcomes.
 */
type Class struct {
	Name                    string
	Counter                 int
	Words                   map[string]types.Word
	Probability             int
	Temp_tokenProbabilities float64
}

/**
 * The weight of each class as a proportion of the entire training set.
 *
 * @var array
 */
type weight struct {
	Amount float64
	Class  Class
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
	Word    types.Word
	Counter map[string]int
}

var (
	checkpointFile    string
)

//changing this will make past saved models not work



func NewNaiveBayes() *NaiveBayes {
	c := new(NaiveBayes)
	c.vocabularySize = 	0
	c.words = 			map[string]wordFrequency{}
	c.classes = 		map[string]Class{}
	c.weigh = 			weight{}
	tokenizer := tokenizers.DefaultTokenizer{}
	c.tokenizer = 		tokenizer.GetName()
	return c
}


func (nb *NaiveBayes) getMeta() (string, string) {
	return "NaiveBayes", "01"
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
			Name:                    name,
			Counter:                 0,
			Words:                   map[string]types.Word{},
			Probability:             0,
			Temp_tokenProbabilities: 0.0000,
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

	tokens := tokenizers.GetTokenizer(nb.tokenizer).Tokenize(text)// nb.tokenizer.Tokenize(text)

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
		wf = wordFrequency{Word: types.NewWord(word), Counter: map[string]int{}}
	}
	wf.Counter[class]++
	nb.words[word] = wf
	nb.classes[class].Words[word] = types.NewWord(word)
	nb.vocabularySize++

}

/**
 * Calculate probability that a `token` belongs to a `category` or class
 *
 */
func (nb *NaiveBayes) tokenProbability(token, category string) float64 {
	//how many times this word has occurred in documents mapped to this category
	wordFrequencyCount := nb.words[token].Counter[category]
	//what is the count of all words that have ever been mapped to this category
	wordCount := len(nb.classes[category].Words)

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
	tokens :=  tokenizers.GetTokenizer(nb.tokenizer).Tokenize(text)//nb.tokenizer.Tokenize(text)

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	for _, class := range nb.classes {

		//start by calculating the overall probability of this category
		//=>  out of all documents we've ever looked at, how many were
		//    mapped to this category
		categoryProbability := len(class.Words) / len(nb.words)
		class.Probability = categoryProbability

		//take the log to avoid underflow
		logProbability := math.Log(float64(categoryProbability))
		tokens_w := map[string]int{}
		for _, w := range tokens {

			tokens_w[w]++
		}
		for _, w := range tokens {

			frequencyInText := tokens_w[w]
			tokenProbability := nb.tokenProbability(w, class.Name)


			logger.Info("processed",
				// Structured context as strongly typed Field values.
				zap.String("token", w),
				zap.String("category", class.Name),
				zap.Float64("token Probability", tokenProbability),
			)

			class.Temp_tokenProbabilities += tokenProbability

			//determine the log of the P( w | c ) for this word
			logProbability += float64(frequencyInText) * math.Log(tokenProbability)
		}
		if class.Temp_tokenProbabilities > top_probability {
			top_probability = class.Temp_tokenProbabilities
			chosenCategory = class.Name
		}
         nb.weigh = weight{
         	Amount: class.Temp_tokenProbabilities,
         	Class: class,
		 }
		/*if (logProbability > maxProbability) {
			maxProbability = int(logProbability)
			chosenCategory = class.name
		} */

		//a reset is okay for online learning.
		class.Temp_tokenProbabilities = 0.0000
	}
	return chosenCategory, top_probability
}




// GobEncode implements GobEncoder. This is necessary because RNN contains several unexported fields.
// It would be easier to simply export them by changing to uppercase, but for comparison purposes,
// I wanted to keep the field names the same between Go and the original Python code.
func (nb *NaiveBayes) GobEncode() ([]byte, error) {
	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)

	var err error

	encode := func(data interface{}) {
		// no-op if we've already seen an err
		if err == nil {
			err = encoder.Encode(data)
		}
	}
	encode(nb.words)
	encode(nb.classes)
	encode(nb.vocabularySize)
	encode(nb.weigh)
	encode(nb.tokenizer)

	return b.Bytes(), err
}

func (nb *NaiveBayes) Save(file string) error {

	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)

	err := encoder.Encode(nb)
	if err != nil {
		return fmt.Errorf("error encoding model: %s", err)
	}

	name, version := nb.getMeta()
	persist.Save(file, persist.Modeldata{
		Data: buf.Bytes(),
		Name: name,
		Version: version,
	})
  return nil
}

// Load from the output file.
func (nb *NaiveBayes) Load(filePath string) error {
	log.Printf("Loading Classifier from %s...", filePath)
	meta := persist.Load(filePath)
	//get the classifier current meta data
	name, version := nb.getMeta()
	if meta.Name != name {
		return fmt.Errorf("This file doesn't contain a Naive-Bayes classifier")
	}
	if meta.Version != version {
		return fmt.Errorf("Can't understand this file format")
	}

	decoder := gob.NewDecoder(bytes.NewBuffer(meta.Data))
	err := decoder.Decode(&nb)
	if err != nil {
		return  fmt.Errorf("error decoding RNN checkpoint file: %s", err)
	}

	checkpointFile = filePath
   return nil
}

// GobDecode implements GoDecoder.
func (nb *NaiveBayes) GobDecode(data []byte) error {
	b := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(b)

	var err error

	decode := func(data interface{}) {
		// no-op if we've already seen an err
		if err == nil {
			err = decoder.Decode(data)
		}
	}
	decode(&nb.words)
	decode(&nb.classes)
	decode(&nb.vocabularySize)
	decode(&nb.weigh)
	decode(&nb.tokenizer)
	return err
}