/*
 * Copyright (c) 2021.  -present, Broos Action, Inc. All rights reserved.
 *
 *  This source code is licensed under the MIT license
 *  found in the LICENSE file in the root directory of this source tree.
 */

package classifiers

import (
	"errors"
	"fmt"
	"math"
)

type DistanceMethodType uint8

const (
	//Euler distance method
	DMT_EulerMethod DistanceMethodType = iota
	// Huffman
	DMT_HuffmanMethod
)

/**
 * K Nearest Neighbors
 *
 * A distance-based algorithm that locates the K nearest neighbors from the
 * training set and uses a majority vote to classify the unknown sample. K
 * Nearest Neighbors is considered a lazy learning Estimator because it does all
 * of its computation at prediction time.
 *
 * @category    Machine Learning
 * @author      Bruce Mubangwa

  **usage
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
	for i := 1; i < 12; i++ {

		knn := classifiers.NewKNearestNeighbors(i, dm, w)
		knn.LearnBatch(train, labels)
		res := knn.Classify(test)
		if res[0] != "Setosa" {
			fmt.Printf("k = %d failed", i)
			fmt.Println()
		}
	}
 */
type KNearestNeighbors struct {
	/**
	 * The number of neighbors to consider when making a prediction.
	 *
	 * @var int
	 */
	K             int

	/**
	 * The distance function to use when computing the distances.
	 *
	 */
	DistanceMethod DistanceMethodType

	/**
	 * Should we use the inverse distances as confidence scores when
	 * making predictions?
	 * Weight define the weight vector for multi-dimension data
	 */
	Weight         []float64
}

//SortedDistance
type SortedDistance struct {
	idx  []int
	dist []float64

	cur int
}

var(
	/**
	 * The training samples that make up the neighborhood of the problem space.
	 *
	 */
	samples [][]float64

	/**
	 * The memoized labels of the training set.
	 */
	labels  []string

	errNotEqualDataLength = errors.New("The data length is not equal.")
)


func NewKNearestNeighbors(k int, dm DistanceMethodType, w []float64) *KNearestNeighbors {
	knn := &KNearestNeighbors{
		K:              k,
		DistanceMethod: dm,
		Weight:         w,
	}
	return knn
}

/**
 * Store the sample and outcome arrays. No other work to be done as this is
 * a lazy learning algorithm.
 *
 */
func (k *KNearestNeighbors) LearnBatch(train [][]float64, label[]string){
	samples = train
	labels = label
}


/**
 * Make predictions from a dataset.
 * Classify, train is an nxp matrix, where n denotes the n training data,
 * and p represents the number of attributes. The test is an mxp matrix.
 */
func (k *KNearestNeighbors) Classify(test [][]float64) []string {

	// computing the distance between one testing data and all the training data
	// fmt.Println(test)
	// fmt.Println(len(test))
	var sortedDistance *SortedDistance
	result := make([]string, len(test))
	for j, _test := range test {

		sortedDistance = NewSortedDistance(len(samples))
		for i, _train := range samples {
			var _dist float64
			switch k.DistanceMethod {
			case DMT_EulerMethod:
				_dist = EulerDistance(_test, _train, k.Weight)
				break
			case DMT_HuffmanMethod:
				_dist = HuffmanDistance(_test, _train, k.Weight)
				break
			default:
				_dist = EulerDistance(_test, _train, k.Weight)
			}
			sortedDistance.Put(i, _dist)
		}
		topKIdx := sortedDistance.SelectTopKIdx(k.K)
		//statistic the frequent of every labels
		freqLabels := make(map[string]int, k.K)
		var maxFreq int = 0
		for _, idx := range topKIdx {
			_label := labels[idx]
			v, ok := freqLabels[_label]
			if ok {
				freqLabels[_label] = v + 1
			} else {
				freqLabels[_label] = 1
			}

			if maxFreq < v+1 {
				maxFreq = v + 1
			}

		}
		//get the label with maximum frequent
		_out := make([]string, 0)
		for k, v := range freqLabels {
			if v == maxFreq {
				_out = append(_out, k)
			}
		}
		if len(_out) == 1 {
			result[j] = _out[0]
		} else {
			fmt.Printf("the %d test data has %d classification results.", j, len(_out))
			 fmt.Println()
			_str := ""
			for _, v := range _out {
				_str += v + "#"
			}
			result[j] = _str[0:len(_str)]
		}

	}
	return result
}

//
//EluerDistance
func EulerDistance(X, Y, W []float64) float64 {

	_len := len(X)
	if _len != len(Y) || _len != len(W) {
		fmt.Println(errNotEqualDataLength)
	}

	var x, y, w float64

	var sum float64
	for i := 0; i < _len; i++ {
		x = X[i]
		y = Y[i]
		w = W[i]
		sum += w * math.Pow(x-y, 2)
	}
	return math.Sqrt(sum)
}

//HuffmanDistance
func HuffmanDistance(X, Y, W []float64) float64 {

	_len := len(X)
	if _len != len(Y) || _len != len(W) {
		fmt.Println(errNotEqualDataLength)
	}

	var x, y, w float64

	var sum float64
	for i := 0; i < _len; i++ {
		x = X[i]
		y = Y[i]
		w = W[i]
		sum += w * math.Abs(x-y)
	}

	return math.Sqrt(sum)
}


//NewSortedDistance initial the SortedDistance with the size
func NewSortedDistance(size int) *SortedDistance {
	fmt.Printf("size: %d", size)
	fmt.Println()
	s := &SortedDistance{
		idx:  make([]int, size),
		dist: make([]float64, size),
		cur:  0,
	}
	return s
}

//Put
func (s *SortedDistance) Put(idx int, dist float64) {
	s.idx[s.cur] = idx
	s.dist[s.cur] = dist

	s.cur++
}

//Len
func (s *SortedDistance) Len() int {
	return len(s.dist)
}

//Less return true if [i] < [j].
func (s *SortedDistance) Less(i, j int) bool {
	if s.dist[i] < s.dist[j] {
		return true
	}
	return false
}

//Swap
func (s *SortedDistance) Swap(i, j int) {
	_tempIdx := s.idx[i]
	_tempDistance := s.dist[i]

	s.idx[i] = s.idx[j]
	s.dist[i] = s.dist[j]

	s.idx[j] = _tempIdx
	s.dist[j] = _tempDistance
}

//SelectTopKIdx return the index of the train
func (s *SortedDistance) SelectTopKIdx(k int) []int {
	return s.idx[0:k]
}

//GetIdx
func (s *SortedDistance) GetIdx() []int {
	return s.idx
}
