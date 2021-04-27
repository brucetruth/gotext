/*
 * Copyright (c) 2021.  -present, Broos Action, Inc. All rights reserved.
 *
 *  This source code is licensed under the MIT license
 *  found in the LICENSE file in the root directory of this source tree.
 */

package stringUtils

import "regexp"

var(
	//like N.D.A
	periodAcronym = "([A-Z]\\.)+[A-Z]?,?$"
	//like F.
	oneLetterAcronym = "^[A-Z]\\.,?$"
	//like NDA
	noPeriodAcronym = "[A-Z]{2,}('s|,)?$"
	//like c.e.o
	lowerCaseAcronym = "([a-z]\\.)+[a-z]\\.?$"
	acronym  []string
	ismatch = false

)

func IsAcronym(text string) bool {
	r, _ := regexp.Compile(periodAcronym+"|"+oneLetterAcronym+"|"+noPeriodAcronym+"|"+lowerCaseAcronym)
	ismatch = r.MatchString(text)
	if ismatch{
		return true

	}

	return ismatch
}
