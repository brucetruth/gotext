/*
 * Copyright (c) 2021.  -present, Broos Action, Inc. All rights reserved.
 *
 *  This source code is licensed under the MIT license
 *  found in the LICENSE file in the root directory of this source tree.
 */

package transformers

import "syscall"

/**
 * Word Count Vectorizer
 *
 * In machine learning, word counts are often used to represent natural language
 * as numerical vectors. The Word Count Vectorizer builds a vocabulary using
 * hash tables from the training samples during fitting and transforms an array
 * of strings (text blobs) into sparse feature vectors. Each feature column
 * represents a word from the vocabulary and the value denotes the number of times
 * that word appears in a given sample.
 *
 * @category    Machine Learning
 * @author      Bruce Mubangwa
 */

type WordCountVectorizer struct {
	/**
	 * The maximum size of the vocabulary.
	 */
	maxVocabulary         		int

	/**
	 * The minimum number of documents a word must appear in to be added to
	 * the vocabulary.
	 */
	minDocumentFrequency         int

	/**
	 * The vocabulary of the fitted training set per categorical feature
	 * column.
	 */
	vocabulary					  []string
}


func NewWordCountVectorizer() *WordCountVectorizer {
	return &WordCountVectorizer{
		maxVocabulary:  syscall.INFINITE,
		minDocumentFrequency: 1,
		vocabulary:  []string{},
	}
}

/**
 * Return the size of the vocabulary in words.
 *
 */
func (wcv *WordCountVectorizer) size() int{
	return len(wcv.vocabulary)
}

//todo