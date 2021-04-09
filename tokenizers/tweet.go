/*
 * Copyright (c) 2021.  -present, Broos Action, Inc. All rights reserved.
 *
 *  This source code is licensed under the MIT license
 *  found in the LICENSE file in the root directory of this source tree.
 */

package tokenizers

import (
	"fmt"
	"regexp"
)

var (
	EMOTICONS = []string{
		"(?:",
		"[<>]?",
		"[:;=8]",
		"[\\-o*\"]?",
		"[)\\](\\[dDpP/:}{@|\"]",
		"|",
		"[)\\](\\[dDpP/:}{@|\"]",
		"[\\-o*\"]?",
		"[:;=8]",
		"[<>]?",
		"|",
		"<3",
		")",
	}

// URL pattern due to John Gruber, modified by Tom Winzig. See
// https://gist.github.com/winzig/8894715
	 URLS = []string{
		 "(?:",
		 "https?:",
		 "(?:",
		 "\\/{1,3}",
		 "|",
		 "[a-z0-9%]",
		 ")",
		 "|",
		 "[a-z0-9.\\-]+[.]",
		 "(?:[a-z]{2,13})",
		 "\\/",
		 ")",
		 "(?:",
		 "[^\\s()<>{}\\[\\]]+",
		 "|",
		 "\\([^\\s()]*?\\([^\\s()]+\\)[^\\s()]*?\\)",
		 "|",
		 "\\([^\\s]+?\\)",
		 ")+",
		 "(?:",
		 "\\([^\\s()]*?\\([^\\s()]+\\)[^\\s()]*?\\)",
		 "|",
		 "\\([^\\s]+?\\)",
		 "|",
		 "[^\\s`!()\\[\\]{};:\".,<>?«»“”‘’]",
		 ")",
		 "|",
		 "(?:",
		 // NOTE: Lookbehind doesn"t work in JavaScript. Skipping this for now.
		 // "(?<!@)",
		 "[a-z0-9]+",
		 "(?:[.\\-][a-z0-9]+)*",
		 "[.]",
		 "(?:[a-z]{2,13})",
		 "\\b",
		 "\\/?",
		 "(?!@)",
		 ")",
	 }

	 phone_numbers = []string{
		 "(?:",
		 "(?:",
		 "\\+?[01]",
		 "[\\-\\s.]*",
		 ")?",
		 "(?:",
		 "[\\(]?",
		 "\\d{3}",
		 "[\\-\\s.)]*",
		 ")?",
		 "\\d{3}",
		 "[\\-\\s.]*",
		 "\\d{4}",
		 ")",
	 }

	 type_words = []string{
		 // "(?:[^\\W\\d_](?:[^\\W\\d_]|[\"\\-_])+[^\\W\\d_])",
		 // "|",
		 "(?:[+\\-]?\"d+[,/.:-]\\d+[+\\-]?)",
		 "|",
		 "(?:[\\w_]+)",
		 "|",
		 "(?:\\.(?:\\s*\\.){1,})",
		 "|",
		 "(?:\\S)",
	 }

)



func Tweet(tweet string) []string{
// EMOTICONS :=  strings.Join(EMOTICONS,"")
 //URLS := strings.Join(URLS, "")
// phone_numbers :=  strings.Join(phone_numbers, "")
// type_words :=  strings.Join(type_words, "")
// html_tags := "<[^>\\\\s]+>"
//	ASCII_Arrows := "[\\\\-]+>|<[\\\\-]+"
	//Twitter_Usernames := "(?:@[\\\\w_]+)"
//	Twitter_hashtags := "(?:\\\\#+[\\\\w_]+[\\\\w\\'_\\\\-]*[\\\\w_]+)"
	Email_addresses := "[\\\\w.+-]+@[\\\\w-]+\\\\.(?:[\\\\w-]\\\\.?)+[\\\\w-]"
	HANG_RE := "/([^a-zA-Z0-9])\"1{3,}/g"


	r, _ := regexp.Compile(HANG_RE)

	safeText := r.ReplaceAllString(tweet,"$1$1$1")

	fmt.Println(safeText)
	rm, _ := regexp.Compile(Email_addresses)

	matches := rm.FindAllString(safeText, -1)

	return matches

}
