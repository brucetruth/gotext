/*
 * Copyright (c) 2021.  -present, Broos Action, Inc. All rights reserved.
 *
 *  This source code is licensed under the MIT license
 *  found in the LICENSE file in the root directory of this source tree.
 */

package tokenizers

import (
	"github.com/broosaction/gotext/utils/xstrings"
	"log"
	"regexp"
	"strings"
)

type BertWordPieceTokenizer struct {
	RemoveStopWords bool

	ClsToken    	string
	MaskToken   	string
	PadToken  	 	string
	SepToken 		string
	UnkToken 		string

	TokenPosition   tokenPosition
	AffixMaxLength 	int
	words 			map[string]int
	extras 			map[string]int
	affixes 		map[string]int
}

type tokenPosition struct {
	ClsToken    	int
	MaskToken   	int
	PadToken  	 	int
	SepToken 		int
	UnkToken 		int
}

type wordToken struct {
	text 	string
	start 	int
	typ 	string
}


var BertWordPieceTokenizerName = "BertWordPieceTokenizer"

func (b *BertWordPieceTokenizer) init() {
	b.ClsToken 	= "[CLS]"
	b.MaskToken = "[MASK]"
	b.PadToken 	= "[PAD]"
	b.SepToken 	= "[SEP]"
	b.UnkToken 	= "[UNK]"
	b.TokenPosition = tokenPosition{
		ClsToken: 101,
		MaskToken: 103,
		PadToken: 0,
		SepToken: 102,
		UnkToken: 100,
	}



}

func (b *BertWordPieceTokenizer) Tokenizer(text string) {
	b.init()
}


func (b *BertWordPieceTokenizer) loadDictionary(inputVocabContent string){
	vocabContent := inputVocabContent

	vocContent := strings.Split(vocabContent,"/\\r?\\n/")

    b.AffixMaxLength = 0
	i := 0
	for  i = 0; i < len(vocContent); i += 1 {
		word := vocContent[i]
		b.words[word] = i
		if strings.HasPrefix(word,"##"){
			affix := xstrings.Slice(word,2,0)
			if len(affix) > b.AffixMaxLength {
				b.AffixMaxLength = len(affix)
			}
			b.affixes[affix] = i
		}
	}

}

func (b *BertWordPieceTokenizer) createToken(text string, start int, srcType string) wordToken {
	typ := srcType
	if typ == ""{
		if strings.Contains(text,"\r") || strings.Contains(text,"\n") || strings.Contains(text," ") || strings.Contains(text,"\t"){
			typ = "space"
		} else{
			typ = "separator"
		}
	}
	return wordToken{
		text:  text,
		start: start + len(text) - 1,
		typ:   typ,
	}

}

//ooh chill. lol
func (b *BertWordPieceTokenizer) splitSentence(text string) []wordToken {
	rNormalizer, _ := regexp.Compile("[\\u0300-\\u036f]")
	normalized := text
	if rNormalizer.MatchString(text){
		normalized =  rNormalizer.ReplaceAllString(text,"")
	}

	var result []wordToken
	rx,_ := regexp.Compile("\\W+")
	
	lastEnd := 0
	found := rx.FindAllString(normalized, -1)
	foundCounter := 0
	for _, w := range found{
		log.Println(w)
		chars := strings.Split(found[0],"")
		i := 0
		for i = 0;i<len(chars);i+=1{
			token := b.createToken(chars[i],foundCounter + i,"")
			if token.start > lastEnd {
				wordtoken := b.createToken(xstrings.Slice(text,token.start,lastEnd), lastEnd, "word")

				result = append(result,wordtoken)
			}
			result = append(result,token)
		}
		foundCounter +=1
	}
	if lastEnd < len(text){
		result = append(result,b.createToken(xstrings.Slice(text,len(text),lastEnd), lastEnd, "word"))
	}
	return result
}

func (b *BertWordPieceTokenizer) getBestAffix(word string){

}

func (b *BertWordPieceTokenizer) GetName() string {
	return BertWordPieceTokenizerName
}
