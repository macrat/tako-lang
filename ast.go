package main

import (
	"fmt"
)

type Expression interface {
	Compute(Context) (Expression, error)
	Computable(Context) bool
}

type ExpressionList []Expression

func (el ExpressionList) Compute(ctx Context) (exp Expression, err error) {
	for _, e := range el {
		exp, err = e.Compute(ctx)
		if err != nil {
			return
		}
	}
	return
}

func (el ExpressionList) Computable(ctx Context) bool {
	for _, e := range el {
		if e.Computable(ctx) {
			return true
		}
	}
	return false
}

type Token struct {
	Token   int
	Literal string
}

type Number int64

func (n Number) Compute(ctx Context) (Expression, error) {
	return n, nil
}

func (n Number) Computable(ctx Context) bool {
	return false
}

type Identifier string

func (i Identifier) Compute(ctx Context) (Expression, error) {
	if val := ctx.Get(i); val != nil {
		return val, nil
	} else {
		return nil, fmt.Errorf("NameError: %s is not defined", i)
	}
}

func (i Identifier) Computable(ctx Context) bool {
	return true
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

func (fd FunctionDefine) Computable(ctx Context) bool {
	return false
}

func (fd FunctionDefine) GetArguments() []Identifier {
	return fd.Arguments
}

func (fd FunctionDefine) Call(ctx Context, args map[Identifier]Expression) (Expression, error) {
	newCtx := ctx.MakeScope()
	for k, v := range args {
		newCtx.Put(k, v)
	}
	return fd.Expression.Compute(newCtx)
}

type FunctionCall struct {
	Function  Expression
	Arguments []Expression
}

func (fc FunctionCall) GetFunction(ctx Context) (Function, error) {
	raw, err := ctx.ComputeRecursive(fc.Function)
	if err != nil {
		return nil, err
	}

	f, ok := raw.(Function)
	if !ok {
		return nil, fmt.Errorf("TypeError: %s is not function", fc.Function)
	}

	return f, nil
}

func (fc FunctionCall) Compute(ctx Context) (Expression, error) {
	f, err := fc.GetFunction(ctx)
	if err != nil {
		return nil, err
	}

	if len(fc.Arguments) != len(f.GetArguments()) {
		return nil, fmt.Errorf("ArgumentsError: excepted %d arguments but got %d", len(f.GetArguments()), len(fc.Arguments))
	}

	args := make(map[Identifier]Expression)
	for i, x := range f.GetArguments() {
		args[x] = fc.Arguments[i]
	}

	return f.Call(ctx, args)
}

func (fc FunctionCall) Computable(ctx Context) bool {
	return true
}
