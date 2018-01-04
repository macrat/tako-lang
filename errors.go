package main

import (
	"fmt"
	"strings"
)

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

type OutOfBoundsError struct {
	max int
	got int
}

func (e OutOfBoundsError) Error() string {
	return fmt.Sprintf("index %d is out of bounds (must be between 0 and %d)", e.got, e.max)
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

type ConditionTypeError struct {
	pos Position
}

func (e ConditionTypeError) Error() string {
	return fmt.Sprintf("%s: condition value must be boolean value", e.pos)
}

type TypeError struct {
	name    string
	excepts []string
	pos     Position
}

func (e TypeError) Error() string {
	except := ""
	if len(e.excepts) == 1 {
		except = e.excepts[0]
	} else {
		except = strings.Join(e.excepts[:len(e.excepts)-1], ", ") + " or " + e.excepts[len(e.excepts)-1]
	}
	return fmt.Sprintf("%s: %s must be %s", e.pos, e.name, except)
}
