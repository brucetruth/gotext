/*
 * Copyright (c) 2021.  -present, Broos Action, Inc. All rights reserved.
 *
 *  This source code is licensed under the MIT license
 *  found in the LICENSE file in the root directory of this source tree.
 */

package transformers

import (
	"regexp"
)

/** Turn a singular noun into a plural
 * assume the given string is singular
 * @author      Bruce Mubangwa
 */

type suffix struct {
	Tag       string
	Patterns  string
	Responses string
}

var (
	suffixes = map[string][]suffix{}
	addE = "(x|ch|sh|s|z)$"
	irregulars map[string]string
)

//plural english rules
func iniSuffixes(){
	/** patterns for turning "bus" to "buses"*/
	registerSuffixes("en",[]suffix{
		{
			Tag: "a",
			Patterns:
				"(antenn|formul|nebul|vertebr|vit)a$",
			Responses:
				"$1ae",
		},
		{
			Tag: "a",
			Patterns:
			"([ti])a$",
			Responses:
			"$1a",
		},
		{
			Tag: "e",
			Patterns:
			"(kn|l|w)ife$",
			Responses:
			"${1}ives",
		},
		{
			Tag: "e",
			Patterns:
			"(hive)$",
			Responses:
			"${1}s",
		},
		{
			Tag: "e",
			Patterns:
			"([m|l])ouse$",
			Responses:
			"${1}ice",
		},
		{
			Tag: "e",
			Patterns:
			"([m|l])ice$",
			Responses:
			"${1}ice",
		},

		{
			Tag: "f",
			Patterns:
			"^(dwar|handkerchie|hoo|scar|whar)f$",
			Responses:
			"${1}ves",
		},
		{
			Tag: "f",
			Patterns:
			"^((?:ca|e|ha|(?:our|them|your)?se|she|wo)l|lea|loa|shea|thie)f$",
			Responses:
			"${1}ves",
		},
		{
			Tag: "i",
			Patterns:
			"(octop|vir)i$",
			Responses:
			"${1}i",
		},
		{
			Tag: "m",
			Patterns:
			"([ti])um$",
			Responses:
			"${1}a",
		},
		{
			Tag: "n",
			Patterns:
			"^(oxen)$",
			Responses:
			"${1}",
		},
		{
			Tag: "o",
			Patterns:
			"(al|ad|at|er|et|ed|ad)o$",
			Responses:
			"${1}oes",
		},
		{
			Tag: "s",
			Patterns:
			"(ax|test)is$",
			Responses:
			"${1}es",
		},
		{
			Tag: "s",
			Patterns:
			"(alias|status)$",
			Responses:
			"${1}es",
		},
		{
			Tag: "s",
			Patterns:
			"sis$",
			Responses:
			"ses",
		},
		{
			Tag: "s",
			Patterns:
			"(sis)$",
			Responses:
			"ses",
		},

		{
			Tag: "n",
			Patterns:
			"(wo)man$",
			Responses:
			"${1}men",
		},
		{
			Tag: "s",
			Patterns:
			"^(?!talis|.*hu)(.*)man$",
			Responses:
			"${1}men",
		},
		{
			Tag: "s",
			Patterns:
			"(octop|vir|radi|nucle|fung|cact|stimul)us$",
			Responses:
			"${1}i",
		},
		{
			Tag: "x",
			Patterns:
			"(matr|vert|ind|cort)(ix|ex)$",
			Responses:
			"${1}ices",
		},
		{
			Tag: "x",
			Patterns:
			"^(ox)$",
			Responses:
			"${1}en",
		},
		{
			Tag: "y",
			Patterns:
			"([^aeiouy]|qu)y$",
			Responses:
			"${1}ies",
		},
		{
			Tag: "z",
			Patterns:
			"(quiz)$",
			Responses:
			"${1}zes",
		},
	})

	irregulars= map[string]string{
		"addendum": "addenda",
		"alga": "algae",
		"alumna": "alumnae",
		"alumnus": "alumni",
		"analysis": "analyses",
		"antenna": "antennae",
		"appendix": "appendices",
		"avocado": "avocados",
		"axis": "axes",
		"bacillus": "bacilli",
		"barracks": "barracks",
		"beau": "beaux",
		"bus": "buses",
		"cactus": "cacti",
		"chateau": "chateaux",
		"child": "children",
		"circus": "circuses",
		"clothes": "clothes",
		"corpus": "corpora",
		"criterion": "criteria",
		"curriculum": "curricula",
		"database": "databases",
		"deer": "deer",
		"diagnosis": "diagnoses",
		"echo": "echoes",
		"embargo": "embargoes",
		"epoch": "epochs",
		"foot": "feet",
		"formula": "formulae",
		"fungus": "fungi",
		"genus": "genera",
		"goose": "geese",
		"halo": "halos",
		"hippopotamus": "hippopotami",
		"index": "indices",
		"larva": "larvae",
		"leaf": "leaves",
		"libretto": "libretti",
		"loaf": "loaves",
		"man": "men",
		"matrix": "matrices",
		"memorandum": "memoranda",
		"modulus": "moduli",
		"mosquito": "mosquitoes",
		"mouse": "mice",
		// move: "moves",
		"nebula": "nebulae",
		"nucleus": "nuclei",
		"octopus": "octopi",
		"opus": "opera",
		"ovum": "ova",
		"ox": "oxen",
		"parenthesis": "parentheses",
		"person": "people",
		"phenomenon": "phenomena",
		"prognosis": "prognoses",
		"quiz": "quizzes",
		"radius": "radii",
		"referendum": "referenda",
		"rodeo": "rodeos",
		"sex": "sexes",
		"shoe": "shoes",
		"sombrero": "sombreros",
		"stimulus": "stimuli",
		"stomach": "stomachs",
		"syllabus": "syllabi",
		"synopsis": "synopses",
		"tableau": "tableaux",
		"thesis": "theses",
		"thief": "thieves",
		"tooth": "teeth",
		"tornado": "tornados",
		"tuxedo": "tuxedos",
		"vertebra": "vertebrae",
		// virus: "viri",
		// zero: "zeros",
	}
}

// RegisterModules registers an array of modules into the map
func registerSuffixes(locale string, _suffix []suffix) {
	suffixes[locale] = append(suffixes[locale], _suffix...)
}

func trySuffix(str string) string{
	c := str[len(str) - 1]
	for _, _suffix := range suffixes["en"] {
		if string(c) == _suffix.Tag {
			r, _ := regexp.Compile(_suffix.Patterns)
			if r.MatchString(str){
             return r.ReplaceAllString(str,_suffix.Responses)
			}
		}
	}
	return ""
}

func ToPlural(text string) string{
   iniSuffixes()
	// check irregulars list
	if irregulars[text] != "" {
		return irregulars[text]
	}
	//we have some rules to try-out
	plural := trySuffix(text)
	if plural != "" {
		return plural
	}

	//like 'church'
	r, _ := regexp.Compile(addE)
	if r.MatchString(text){
		return text + "es"
	}

	// ¯\_(ツ)_/¯
	return text + "s"
}
