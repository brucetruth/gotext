/*
 * Copyright (c) 2021.  -present, Broos Action, Inc. All rights reserved.
 *
 *  This source code is licensed under the MIT license
 *  found in the LICENSE file in the root directory of this source tree.
 */

package main

import (
	"fmt"
	"github.com/broosaction/gotext/classifiers"
	"github.com/broosaction/gotext/nlp/ner"
	"github.com/broosaction/gotext/transformers"
	stringUtils "github.com/broosaction/gotext/utils/strings"
	"github.com/broosaction/gotext/utils/types"
)

func main() {

	classifier := classifiers.NewNaiveBayes()

	classifier.Learn("amazing, awesome movie!! Yeah!! Oh boy.", "positive")
	classifier.Learn("Sweet, this is incredibly, amazing, perfect, great!!", "positive")
	classifier.Learn("terrible, shitty thing. Damn. Sucks!! it is bad, hell no i hate it", "negative")

	fmt.Println(classifier.Classify("Damn, james this is a great thing, yeah "))

	//wordlab
	word := types.NewWord("agreed")
	word.ApplyStemmer()
	word.Text = "cats"
	word.Learn()
	fmt.Println(word.Stem)

	train := [][]float64{
		[]float64{5.3, 3.7},
		[]float64{5.1, 3.8},
		[]float64{7.2, 3},
		[]float64{5.4, 3.4},
		[]float64{5.1, 3.3},
		[]float64{5.4, 3.9},
		[]float64{7.4, 2.8},
		[]float64{6.1, 2.8},
		[]float64{7.3, 2.9},
		[]float64{6, 2.7},
		[]float64{5.8, 2.8},
		[]float64{6.3, 2.3},
		[]float64{5.1, 2.5},
		[]float64{6.3, 2.5},
		[]float64{5.5, 2.4},
	}

	labels := []string{
		"Setosa",
		"Setosa",
		"Virginica",
		"Setosa",
		"Setosa",
		"Setosa",
		"Virginica",
		"Versicolor",
		"Virginica",
		"Versicolor",
		"Virginica",
		"Versicolor",
		"Versicolor",
		"Versicolor",
		"Versicolor",
	}
	dm := classifiers.DMT_EulerMethod
	w := []float64{0.5, 0.5}

	var test [][]float64

	test = [][]float64{
		[]float64{5.2, 3.1},
	}


		knn := classifiers.NewKNearestNeighbors(len(labels), dm, w)
		knn.LearnBatch(train, labels)
		res := knn.Classify(test)

			fmt.Println(res)



	doc :=  types.NewDocument(`A perennial also-ran, Stallings won his seat when longtime lawmaker David Holmes
    died 11 days after the filing deadline. Suddenly, Stallings was a shoo-in, not
    the long shot. In short order, the Legislature attempted to pass a law allowing
    former U.S. Rep. Carolyn Cheeks Kilpatrick to file; Stallings challenged the
    law in court and won. Kilpatrick mounted a write-in campaign, but Stallings won.`)
	doc.Learn()
	fmt.Println(doc.Sentences[1].Text)
	fmt.Println(len(doc.Sentences))
	fmt.Println(ner.EMailEntityRecognizer("hi my email adress is bruce@broos.link and my other is broosaction@gmail.com"))
     fmt.Println(stringUtils.IsAcronym("NDA"))
	fmt.Println(transformers.ToPlural("ox"))

	intent := classifiers.NewIntentClassifier()
	intent.Train("amazing, awesome movie!! Yeah!! Oh boy.", "helli")
	intent.Train("Sweet, this is incredibly, amazing, perfect, great!!", "positive")

	fmt.Println(intent.Classify(" perfect is movie was"))
}
