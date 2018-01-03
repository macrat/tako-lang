package main

import (
	"fmt"
	"os"
	"io"
	"regexp"

	"github.com/macrat/simplexer"
)

type Lexer struct {
	lexer     *simplexer.Lexer
	result    Expression
	lastToken *simplexer.Token
}

func NewLexer(reader io.Reader) *Lexer {
	l := simplexer.NewLexer(reader)

	l.Whitespace = regexp.MustCompile(`^[ \t]+`)
	l.TokenTypes = []simplexer.TokenType{
		simplexer.NewTokenType(NEWLINE, `^[\n\r]+`),
		simplexer.NewTokenType(BOOLEAN, `^(true|false)`),
		simplexer.NewTokenType(NUMBER, `^[0-9]+`),
		simplexer.NewTokenType(COMPARE_OPERATOR, `^(==|!=)`),
		simplexer.NewTokenType(DEFINE_OPERATOR, `^(:=|=)`),
		simplexer.NewTokenType(IDENTIFIER, `^[a-zA-Z_][a-zA-Z0-9_]*`),
		simplexer.NewTokenType(0, `^.`),
	}

	return &Lexer {
		lexer: l,
	}
}

func (l *Lexer) Lex(lval *yySymType) int {
	token, err := l.lexer.Scan()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
	if token == nil {
		return -1
	}

	tokenID := int(token.Type.ID)
	if tokenID == 0 {
		tokenID = int(token.Literal[0])
	}

	lval.token = Token{
		Token:   tokenID,
		Literal: token.Literal,
	}

	l.lastToken = token

	return tokenID
}

func (l *Lexer) Error(e string) {
	fmt.Fprintf(os.Stderr, "SyntaxError:%d:%d: %#v\n", l.lexer.Position.Line, l.lexer.Position.Column, l.lastToken.Literal)
	os.Exit(1)
}
