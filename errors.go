package main

import "fmt"

type NotDefinedError Identifier

func (e NotDefinedError) Error() string {
	i := Identifier(e)
	return fmt.Sprintf("%s: %s is not defined", i.Position(), i)
}

type AlreadyDefinedError Identifier

func (e AlreadyDefinedError) Error() string {
	i := Identifier(e)
	return fmt.Sprintf("%s: %s is already defined", i.Position(), i)
}

type NotFunctionError struct {
	value Expression
	pos   Position
}

func (e NotFunctionError) Error() string {
	return fmt.Sprintf("%s: %s is not function", e.pos, e.value)
}

type SyntaxError struct {
	pos     Position
	literal string
}

func (e SyntaxError) Error() string {
	return fmt.Sprintf("%s: syntax error near %#v", e.pos, e.literal)
}

type MissmatchArgumentError struct {
	excepted int
	got      int
	pos      Position
	name     string
}

func (e MissmatchArgumentError) Error() string {
	fname := e.name
	if fname == "" {
		fname = "[unnamed]"
	}
	return fmt.Sprintf("%s: %s excepted %d arguments but got %d arguments", e.pos, fname, e.excepted, e.got)
}
