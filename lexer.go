package main

import (
	"fmt"
	"os"
	"text/scanner"
)

type Lexer struct {
	scanner.Scanner

	result Expression
}

func (l *Lexer) Lex(lval *yySymType) int {
	l.Whitespace = 1<<'\t' | 1<<'\r' | 1<<' '
	token := int(l.Scan())

	switch token {
	case scanner.Int:
		token = NUMBER
	case scanner.Ident:
		if l.TokenText() == "true" || l.TokenText() == "false" {
			token = BOOLEAN
		} else {
			token = IDENTIFIER
		}
	case '\n':
		token = NEWLINE
	case ':':
		if l.Peek() == '=' {
			token = DEFINE
			l.Next()
		}
	case '=':
		if l.Peek() == '=' {
			token = EQUALS
			l.Next()
		}
	case '!':
		if l.Peek() == '=' {
			token = NOT_EQUALS
			l.Next()
		}
	}

	lval.token = Token{
		Token:   token,
		Literal: l.TokenText(),
	}

	return token
}

func (l *Lexer) Error(e string) {
	fname := l.Position.Filename
	if fname == "" {
		fname = "unknown"
	}

	fmt.Fprintf(os.Stderr, "SyntaxError: %s:%d:%d\n", fname, l.Position.Line, l.Position.Column)
	os.Exit(1)
}
