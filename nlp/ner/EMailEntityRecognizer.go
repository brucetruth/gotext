/*
 * Copyright (c) 2021.  -present, Broos Action, Inc. All rights reserved.
 *
 *  This source code is licensed under the MIT license
 *  found in the LICENSE file in the root directory of this source tree.
 */

package ner

import (
	"regexp"
)

/**
 * Can recognize e-mail addresses in a sentence.
 * [Author]: Bruce Mubangwa
 */

var (
	regex   = "(([^<>()\\[\\]\\\\.,;:\\s@\"]+(\\.[^<>()\\[\\]\\\\.,;:\\s@\"]+)*)|(\".+\"))@((\\[[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}])|(([a-zA-Z\\-0-9]+\\.)+[a-zA-Z]{2,}))"
    ismatch = false
	emails  []string
)

/**
 * Returns all found entities in a sentence. Returned entities value is <tt>string</tt>.
 */
func EMailEntityRecognizer(text string) (bool, []string){

	r, _ := regexp.Compile(regex)
	ismatch = r.MatchString(text)
	if ismatch{
		emails = r.FindAllString(text,-1)

	}
	//todo
	//return the entity index
	return ismatch, emails
}