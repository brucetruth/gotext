/*
 * Copyright (c) 2021.  -present, Broos Action, Inc. All rights reserved.
 *
 *  This source code is licensed under the MIT license
 *  found in the LICENSE file in the root directory of this source tree.
 */

package classifiers

import (
	"bytes"
	"errors"
	"github.com/broosaction/gotext/tokenizers"
	"io"
	"sync"
)

type IntentClassifier struct {
	feat2cat  map[string]map[string]int
	catCount  map[string]int
	mu        sync.RWMutex
	tokenizer tokenizers.Tokenizer
}

var(
// ErrNotClassified indicates that a document could not be classified
 ErrNotClassified = errors.New("unable to classify document")
)

// New initializes a new naive Classifier using the standard tokenizer
func NewIntentClassifier() *IntentClassifier {
	c := &IntentClassifier{
		feat2cat:  make(map[string]map[string]int),
		catCount:  make(map[string]int),
		tokenizer: &tokenizers.DefaultTokenizer{},
	}
	return c
}

// Train provides supervisory training to the classifier
func (c *IntentClassifier) Train(r string, category string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, feature := range c.tokenizer.Tokenize(r) {
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

	c.mu.RLock()
	defer c.mu.RUnlock()

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
	if _, ok := c.feat2cat[feature]; !ok {
		c.feat2cat[feature] = make(map[string]int)
	}
	c.feat2cat[feature][category]++
}

func (c *IntentClassifier) featureCount(feature string, category string) float64 {
	if _, ok := c.feat2cat[feature]; ok {
		return float64(c.feat2cat[feature][category])
	}
	return 0.0
}

func (c *IntentClassifier) addCategory(category string) {
	c.catCount[category]++
}

func (c *IntentClassifier) categoryCount(category string) float64 {
	if _, ok := c.catCount[category]; ok {
		return float64(c.catCount[category])
	}
	return 0.0
}

func (c *IntentClassifier) count() int {
	sum := 0
	for _, value := range c.catCount {
		sum += value
	}
	return sum
}

func (c *IntentClassifier) categories() []string {
	var keys []string
	for k := range c.catCount {
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
	for _,feature := range c.tokenizer.Tokenize(r) {
		probability *= c.weightedProbability(feature, category)
	}
	return probability
}

func asReader(text string) io.Reader {
	return bytes.NewBufferString(text)
}
