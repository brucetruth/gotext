/*
 * Copyright (c) 2021.  -present, Broos Action, Inc. All rights reserved.
 *
 *  This source code is licensed under the MIT license
 *  found in the LICENSE file in the root directory of this source tree.
 */

package main

import (
	"fmt"
	"goText/classifiers"
)

func main() {

	classifier := classifiers.NewNaiveBayes()

	classifier.Learn("amazing, awesome movie!! Yeah!! Oh boy.", "positive")
	classifier.Learn("Sweet, this is incredibly, amazing, perfect, great!!", "positive")
	classifier.Learn("terrible, shitty thing. Damn. Sucks!!", "negative")

	fmt.Println(classifier.Classify("awesome, cool shitty thing"))

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


}
