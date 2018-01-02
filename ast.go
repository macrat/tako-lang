package main

import (
	"fmt"
)

type Expression interface {
	Compute(Context) (Expression, error)
}

type Token struct {
	Token   int
	Literal string
}

type Number int64

func (n Number) Compute(ctx Context) (Expression, error) {
	return n, nil
}

type Identifier string

func (i Identifier) Compute(ctx Context) (Expression, error) {
	if val, ok := ctx[i]; ok {
		return val, nil
	} else {
		return nil, fmt.Errorf("NameError: %s is not defined", i)
	}
}

type Function interface {
	Expression

	Call(Context, map[Identifier]Expression) (Expression, error)
	GetArguments() []Identifier
}

type FunctionDefine struct {
	Arguments  []Identifier
	Expression Expression
}

func (fd FunctionDefine) Compute(ctx Context) (Expression, error) {
	return fd, nil
}

func (fd FunctionDefine) GetArguments() []Identifier {
	return fd.Arguments
}

func (fd FunctionDefine) Call(ctx Context, args map[Identifier]Expression) (Expression, error) {
	return fd.Expression.Compute(ctx)
}

type FunctionCall struct {
	Function  Expression
	Arguments []Expression
}

func (fc FunctionCall) Compute(ctx Context) (Expression, error) {
	f, err := fc.Function.Compute(ctx)
	if err != nil {
		return nil, err
	}

	fd, ok := f.(Function)
	if !ok {
		return nil, fmt.Errorf("TypeError: %s is not function", fc.Function)
	}

	if len(fc.Arguments) != len(fd.GetArguments()) {
		return nil, fmt.Errorf("ArgumentsError: excepted %d arguments but got %d", len(fd.GetArguments()), len(fc.Arguments))
	}

	args := make(map[Identifier]Expression)
	newCtx := ctx.Copy()
	for i, x := range fd.GetArguments() {
		newCtx[x] = fc.Arguments[i]
		args[x] = fc.Arguments[i]
	}

	return fd.Call(newCtx, args)
}
