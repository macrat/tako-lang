package main

import (
	"fmt"
	"strings"
)

type Expression interface {
	Compute(Context) (Expression, error)
	Computable(Context) bool
}

type LocatedExpression interface {
	Expression

	Position() Position
}

type ExpressionList []Expression

func (el ExpressionList) String() string {
	ss := make([]string, len(el))
	for i, e := range el {
		ss[i] = fmt.Sprint(e)
	}
	return strings.Join(ss, "; ")
}

func (el ExpressionList) Compute(ctx Context) (exp Expression, err error) {
	if len(el) == 0 {
		return Null{}, nil
	}

	for _, e := range el {
		exp, err = e.Compute(ctx)
		if err != nil {
			return
		}
	}

	return
}

func (el ExpressionList) Computable(ctx Context) bool {
	if len(el) == 0 {
		return true
	}

	for _, e := range el {
		if e.Computable(ctx) {
			return true
		}
	}

	return false
}

type Object struct {
	Indexed []Expression
	Named   map[string]Expression
}

func (o Object) String() string {
	ss := make([]string, len(o.Indexed))
	for i, e := range o.Indexed {
		ss[i] = fmt.Sprint(e)
	}
	for k, v := range o.Named {
		ss = append(ss, fmt.Sprintf("%s: %s", k, v))
	}
	return "[" + strings.Join(ss, ", ") + "]"
}

func (o Object) Compute(ctx Context) (Expression, error) {
	result := Object{
		Indexed: make([]Expression, len(o.Indexed)),
		Named:   make(map[string]Expression),
	}

	for i, x := range o.Indexed {
		c, err := x.Compute(ctx)
		if err != nil {
			return nil, err
		}
		result.Indexed[i] = c
	}

	for k, v := range o.Named {
		c, err := v.Compute(ctx)
		if err != nil {
			return nil, err
		}
		result.Named[k] = c
	}

	return Expression(result), nil
}

func (o Object) Computable(ctx Context) bool {
	for _, e := range o.Indexed {
		if e.Computable(ctx) {
			return true
		}
	}

	for _, e := range o.Named {
		if e.Computable(ctx) {
			return true
		}
	}

	return false
}

type Token struct {
	Token   int
	Literal string
	Pos     Position
}

type Identifier struct {
	Key string
	Pos Position
}

func NewIdentifier(key string) Identifier {
	return Identifier{
		Key: key,
		Pos: Position{Filename: "builtin"},
	}
}

func (i Identifier) String() string {
	return i.Key
}

func (i Identifier) Compute(ctx Context) (Expression, error) {
	if val, err := ctx.Get(i); err != nil {
		return nil, err
	} else {
		return val, nil
	}
}

func (i Identifier) Computable(ctx Context) bool {
	return true
}

func (i Identifier) Position() Position {
	return i.Pos
}
