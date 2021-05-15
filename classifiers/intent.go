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
	"errors"
	"fmt"
	"github.com/broosaction/gotext/tokenizers"
	"github.com/broosaction/gotext/utils/persist"
	"io"
	"log"
	"sync"
)

type IntentClassifier struct {
	Feat2cat  map[string]map[string]int
	CatCount  map[string]int
	Mu        sync.RWMutex
	Tokenizer string
}

var(
// ErrNotClassified indicates that a document could not be classified
 ErrNotClassified = errors.New("unable to classify document")
)

// New initializes a new naive Classifier using the standard tokenizer
func NewIntentClassifier() *IntentClassifier {
	tokenizer := tokenizers.DefaultTokenizer{}
	c := &IntentClassifier{
		Feat2cat:  make(map[string]map[string]int),
		CatCount:  make(map[string]int),
		Tokenizer: tokenizer.GetName(),
	}
	return c
}

func (c *IntentClassifier) getMeta() (string, string) {
	return "IntentClassifier", "01"
}

// Train provides supervisory training to the classifier
func (c *IntentClassifier) Train(r string, category string) error {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	for _, feature := range tokenizers.GetTokenizer(c.Tokenizer).Tokenize(r) {
		c.addFeature(feature, category)
	}

	c.addCategory(category)
	return nil
}



// Classify attempts to classify a document. If the document cannot be classified
// (eg. because the classifier has not been trained), an error is returned.
func (c *IntentClassifier) Classify(r string) (string, error) {
	max := 0.0
	var err error
	classification := ""
	probabilities := make(map[string]float64)

	c.Mu.RLock()
	defer c.Mu.RUnlock()

	for _, category := range c.categories() {
		probabilities[category] = c.probability(r, category)
		if probabilities[category] > max {
			max = probabilities[category]
			classification = category
		}
	}

	if classification == "" {
		return "", ErrNotClassified
	}
	return classification, err
}



func (c *IntentClassifier) addFeature(feature string, category string) {
	if _, ok := c.Feat2cat[feature]; !ok {
		c.Feat2cat[feature] = make(map[string]int)
	}
	c.Feat2cat[feature][category]++
}

func (c *IntentClassifier) featureCount(feature string, category string) float64 {
	if _, ok := c.Feat2cat[feature]; ok {
		return float64(c.Feat2cat[feature][category])
	}
	return 0.0
}

func (c *IntentClassifier) addCategory(category string) {
	c.CatCount[category]++
}

func (c *IntentClassifier) categoryCount(category string) float64 {
	if _, ok := c.CatCount[category]; ok {
		return float64(c.CatCount[category])
	}
	return 0.0
}

func (c *IntentClassifier) count() int {
	sum := 0
	for _, value := range c.CatCount {
		sum += value
	}
	return sum
}

func (c *IntentClassifier) categories() []string {
	var keys []string
	for k := range c.CatCount {
		keys = append(keys, k)
	}
	return keys
}

func (c *IntentClassifier) featureProbability(feature string, category string) float64 {
	if c.categoryCount(category) == 0 {
		return 0.0
	}
	return c.featureCount(feature, category) / c.categoryCount(category)
}

func (c *IntentClassifier) weightedProbability(feature string, category string) float64 {
	return c.variableWeightedProbability(feature, category, 1.0, 0.5)
}

func (c *IntentClassifier) variableWeightedProbability(feature string, category string, weight float64, assumedProb float64) float64 {
	sum := 0.0
	probability := c.featureProbability(feature, category)
	for _, category := range c.categories() {
		sum += c.featureCount(feature, category)
	}
	return ((weight * assumedProb) + (sum * probability)) / (weight + sum)
}

func (c *IntentClassifier) probability(r string, category string) float64 {
	categoryProbability := c.categoryCount(category) / float64(c.count())
	docProbability := c.docProbability(r, category)
	return docProbability * categoryProbability
}

func (c *IntentClassifier) docProbability(r string, category string) float64 {
	probability := 1.0
	for _,feature := range tokenizers.GetTokenizer(c.Tokenizer).Tokenize(r) {
		probability *= c.weightedProbability(feature, category)
	}
	return probability
}

func asReader(text string) io.Reader {
	return bytes.NewBufferString(text)
}


//save to a file
func (c *IntentClassifier) Save(file string) error {

	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)

	err := encoder.Encode(c)
	if err != nil {
		return fmt.Errorf("error encoding model: %s", err)
	}

	name, version := c.getMeta()
	persist.Save(file, persist.Modeldata{
		Data: buf.Bytes(),
		Name: name,
		Version: version,
	})
	return nil
}

// Load from the output file.
func (c *IntentClassifier) Load(filePath string) error {
	log.Printf("Loading Classifier from %s...", filePath)
	meta := persist.Load(filePath)
	//get the classifier current meta data
	name, version := c.getMeta()
	if meta.Name != name {
		return fmt.Errorf("This file doesn't contain a KNearestNeighbors classifier")
	}
	if meta.Version != version {
		return fmt.Errorf("Can't understand this file format")
	}

	decoder := gob.NewDecoder(bytes.NewBuffer(meta.Data))
	err := decoder.Decode(&c)
	if err != nil {
		return  fmt.Errorf("error decoding RNN checkpoint file: %s", err)
	}

	checkpointFile = filePath
	return nil
}
