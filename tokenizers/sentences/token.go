/*
 * Copyright (c) 2021.  -present, Broos Action, Inc. All rights reserved.
 *
 *  This source code is licensed under the MIT license
 *  found in the LICENSE file in the root directory of this source tree.
 */

package sentences

import (
	"fmt"
	"regexp"
)

// TokenGrouper two adjacent tokens together.
type TokenGrouper interface {
	Group([]*Token) [][2]*Token
}

// DefaultTokenGrouper is the default implementation of TokenGrouper
type DefaultTokenGrouper struct{}

// Group is the primary logic for implementing TokenGrouper
func (p *DefaultTokenGrouper) Group(tokens []*Token) [][2]*Token {
	if len(tokens) == 0 {
		return nil
	}

	pairTokens := make([][2]*Token, 0, len(tokens))

	prevToken := tokens[0]
	for _, tok := range tokens {
		if prevToken == tok {
			continue
		}
		pairTokens = append(pairTokens, [2]*Token{prevToken, tok})
		prevToken = tok
	}
	pairTokens = append(pairTokens, [2]*Token{prevToken, nil})

	return pairTokens
}

// Token stores a token of text with annotations produced during sentence boundary detection.
type Token struct {
	Tok         string
	Position    int
	SentBreak   bool
	ParaStart   bool
	LineStart   bool
	Abbr        bool
	periodFinal bool
	ReEllipsis  *regexp.Regexp
	ReNumeric   *regexp.Regexp
	ReInitial   *regexp.Regexp
	ReAlpha     *regexp.Regexp
}

var reEllipsis = regexp.MustCompile(`\.\.+$`)
var reNumeric = regexp.MustCompile(`-?[\.,]?\d[\d,\.-]*\.?$`)
var reInitial = regexp.MustCompile(`^[A-Za-z]\.$`)
var reAlpha = regexp.MustCompile(`^[A-Za-z]+$`)

// NewToken is the default implementation of the Token struct
func NewToken(token string) *Token {
	tok := Token{
		Tok:        token,
		ReEllipsis: reEllipsis,
		ReNumeric:  reNumeric,
		ReInitial:  reInitial,
		ReAlpha:    reAlpha,
	}

	return &tok
}

// String is the string representation of Token
func (p *Token) String() string {
	return fmt.Sprintf("<Token Tok: %q, SentBreak: %t, Abbr: %t, Position: %d>", p.Tok, p.SentBreak, p.Abbr, p.Position)
}
