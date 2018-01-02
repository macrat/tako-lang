package main

import (
	"text/scanner"
)

type Lexer struct {
	scanner.Scanner

	result Expression
}

func (l *Lexer) Lex(lval *yySymType) int {
	token := int(l.Scan())

	switch token {
	case scanner.Int:
		token = NUMBER
	case scanner.Ident:
		token = IDENTIFIER
	}

	lval.token = Token{
		Token:   token,
		Literal: l.TokenText(),
	}

	return token
}

func (l *Lexer) Error(e string) {
	panic(e)
}
