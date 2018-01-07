package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/macrat/simplexer"
)

type Position struct {
	simplexer.Position

	Filename string
}

func (p Position) String() string {
	return fmt.Sprintf("%s:%d:%d", p.Filename, p.Line+1, p.Column+1)
}

type Lexer struct {
	lexer        *simplexer.Lexer
	result       Expression
	lastToken    *simplexer.Token
	lastPosition Position
	Filename     string
}

func NewLexer(reader io.Reader) *Lexer {
	l := simplexer.NewLexer(reader)

	l.Whitespace = regexp.MustCompile(`^[ \t]+`)
	l.TokenTypes = []simplexer.TokenType{
		simplexer.NewTokenType(NEWLINE, `^[\n\r]+`),
		simplexer.NewTokenType(NUMBER, `^[0-9]+`),
		simplexer.NewTokenType(COMPARE_OPERATOR, `^(==|!=|>=?|<=?)`),
		simplexer.NewTokenType(DEFINE_OPERATOR, `^(:=|=)`),
		simplexer.NewTokenType(CALCULATE_DEFINE_OPERATOR, `^(\+|-|\*|/)=`),
		simplexer.NewTokenType(FUNCTION_SEP, `^\){`),
		simplexer.NewTokenType(IF, `^if`),
		simplexer.NewTokenType(ELSE, `^else`),
		simplexer.NewTokenType(ELLIPSIS, `^\.{3}`),
		simplexer.NewTokenType(STRING, `^("((?:\\\\|\\"|[^"])*)"|'((?:\\\\|\\'|[^'])*)')`),
		simplexer.NewTokenType(IDENTIFIER, `^([a-zA-Z_][a-zA-Z0-9_]*|:[^ \t\n\r]:|[^ \t\n\r]:)`),
		simplexer.NewTokenType(0, `^.`),
	}

	return &Lexer{
		lexer: l,
	}
}

func (l *Lexer) Lex(lval *yySymType) int {
	token, err := l.lexer.Scan()
	if err != nil {
		if e, ok := err.(simplexer.UnknownTokenError); ok {
			fmt.Fprintln(os.Stderr, e.Error()+":")
			fmt.Fprintln(os.Stderr, l.lexer.GetLastLine())
			fmt.Fprintln(os.Stderr, strings.Repeat(" ", e.Position.Column)+strings.Repeat("^", len(e.Literal)))
		} else {
			l.Error(err.Error())
		}
		os.Exit(1)
	}
	if token == nil {
		return -1
	}

	tokenID := int(token.Type.ID)
	if tokenID == 0 {
		tokenID = int(token.Literal[0])
	}

	pos := Position{
		Position: token.Position,
		Filename: l.Filename,
	}

	lval.token = Token{
		Token:   tokenID,
		Literal: token.Literal,
		Pos:     pos,
	}

	switch tokenID {
	case CALCULATE_DEFINE_OPERATOR:
		lval.token.Literal = token.Submatches[0]
	case STRING:
		lval.token.Literal = regexp.MustCompile(`\\[nrt\\"']`).ReplaceAllStringFunc(token.Submatches[1]+token.Submatches[2], func(s string) string {
			switch s[1] {
			case 'n':
				return "\n"
			case 'r':
				return "\r"
			case 't':
				return "\t"
			case '\\':
				return "\\"
			case '"':
				return "\""
			case '\'':
				return "'"
			}
			return ""
		})
	}

	l.lastToken = token
	l.lastPosition = pos

	return tokenID
}

func (l *Lexer) Error(e string) {
	fmt.Fprintln(os.Stderr, e+":")
	fmt.Fprintln(os.Stderr, l.lexer.GetLastLine())
	fmt.Fprintln(os.Stderr, strings.Repeat(" ", l.lastToken.Position.Column)+strings.Repeat("^", len(l.lastToken.Literal)))
	os.Exit(1)
}
