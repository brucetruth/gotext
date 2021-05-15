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
	"github.com/broosaction/gotext/utils/persist"
	"log"
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

	// The training samples that make up the neighborhood of the problem space.
	Samples [][]float64


	//The memoized labels of the training set.
	Labels  []string
}

//SortedDistance
type SortedDistance struct {
	Idx  []int
	Dist []float64

	Cur int
}

var(
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

func (k *KNearestNeighbors) getMeta() (string, string) {
	return "KNearestNeighbors", "01"
}
/**
 * Store the sample and outcome arrays. No other work to be done as this is
 * a lazy learning algorithm.
 *
 */
func (k *KNearestNeighbors) LearnBatch(train [][]float64, label[]string){
	k.Samples = train
	k.Labels = label
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

		sortedDistance = NewSortedDistance(len(k.Samples))
		for i, _train := range k.Samples {
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
			_label := k.Labels[idx]
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
		Idx:  make([]int, size),
		Dist: make([]float64, size),
		Cur:  0,
	}
	return s
}

//Put
func (s *SortedDistance) Put(idx int, dist float64) {
	s.Idx[s.Cur] = idx
	s.Dist[s.Cur] = dist

	s.Cur++
}

//Len
func (s *SortedDistance) Len() int {
	return len(s.Dist)
}

//Less return true if [i] < [j].
func (s *SortedDistance) Less(i, j int) bool {
	if s.Dist[i] < s.Dist[j] {
		return true
	}
	return false
}

//Swap
func (s *SortedDistance) Swap(i, j int) {
	_tempIdx := s.Idx[i]
	_tempDistance := s.Dist[i]

	s.Idx[i] = s.Idx[j]
	s.Dist[i] = s.Dist[j]

	s.Idx[j] = _tempIdx
	s.Dist[j] = _tempDistance
}

//SelectTopKIdx return the index of the train
func (s *SortedDistance) SelectTopKIdx(k int) []int {
	return s.Idx[0:k]
}

//GetIdx
func (s *SortedDistance) GetIdx() []int {
	return s.Idx
}

//save to a file
func (k *KNearestNeighbors) Save(file string) error {

	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)

	err := encoder.Encode(k)
	if err != nil {
		return fmt.Errorf("error encoding model: %s", err)
	}

	name, version := k.getMeta()
	persist.Save(file, persist.Modeldata{
		Data: buf.Bytes(),
		Name: name,
		Version: version,
	})
	return nil
}

// Load from the output file.
func (k *KNearestNeighbors) Load(filePath string) error {
	log.Printf("Loading Classifier from %s...", filePath)
	meta := persist.Load(filePath)
	//get the classifier current meta data
	name, version := k.getMeta()
	if meta.Name != name {
		return fmt.Errorf("This file doesn't contain a KNearestNeighbors classifier")
	}
	if meta.Version != version {
		return fmt.Errorf("Can't understand this file format")
	}

	decoder := gob.NewDecoder(bytes.NewBuffer(meta.Data))
	err := decoder.Decode(&k)
	if err != nil {
		return  fmt.Errorf("error decoding RNN checkpoint file: %s", err)
	}

	checkpointFile = filePath
	return nil
}
