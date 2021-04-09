/*
 * Copyright (c) 2021.  -present, Broos Action, Inc. All rights reserved.
 *
 *  This source code is licensed under the MIT license
 *  found in the LICENSE file in the root directory of this source tree.
 */

package tokenizers

import (
	stringUtils "goText/utils/strings"
	"strings"
)
/**
 * A very simple paragraph tokenizer.
 *
 * [Author]: Bruce Mubangwa
 */
const (
	par_regx = "/(?:\r\n|\n\r|\n|\r)/"
)


func Paragraphs(text string) []string {

	s := stringUtils.Cleanup(text)
	tokens := strings.Split(s, par_regx)

	return tokens
}
