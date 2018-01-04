package main

import (
	"fmt"
	"io"
	"os"
	"regexp"

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
		simplexer.NewTokenType(COMPARE_OPERATOR, `^(==|!=|>|>=|<|<=)`),
		simplexer.NewTokenType(DEFINE_OPERATOR, `^(:=|=)`),
		simplexer.NewTokenType(CALCULATE_DEFINE_OPERATOR, `^(\+|-|\*|/)=`),
		simplexer.NewTokenType(FUNCTION_SEP, `^\){`),
		simplexer.NewTokenType(IF, `^if`),
		simplexer.NewTokenType(ELSE, `^else`),
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

	pos := Position{
		Position: l.lexer.Position,
		Filename: l.Filename,
	}

	lval.token = Token{
		Token:   tokenID,
		Literal: token.Literal,
		Pos:     pos,
	}

	if tokenID == CALCULATE_DEFINE_OPERATOR {
		lval.token.Literal = token.Submatches[0]
	}

	l.lastToken = token
	l.lastPosition = pos

	return tokenID
}

func (l *Lexer) Error(e string) {
	fmt.Fprintln(os.Stderr, (SyntaxError{pos: l.lastPosition, literal: l.lastToken.Literal}).Error())
	os.Exit(1)
}
