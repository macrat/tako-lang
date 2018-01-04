package main

import (
	"fmt"
	"strings"
)

type Function interface {
	Expression

	Call(Context, map[Identifier]Expression, *Object) (Expression, error)
	GetArguments() []Identifier
	GetVariableArgument() *Identifier
}

type FunctionDefine struct {
	Arguments        []Identifier
	VariableArgument *Identifier
	Expression       Expression
	Pos              Position
}

func (fd FunctionDefine) String() string {
	args := make([]string, len(fd.Arguments))
	for i, a := range fd.Arguments {
		args[i] = a.String()
	}
	if v := fd.GetVariableArgument(); v != nil {
		return fmt.Sprintf("(%s...){%s}", strings.Join(append(args, v.Key), ", "), fd.Expression)
	} else {
		return fmt.Sprintf("(%s){%s}", strings.Join(args, ", "), fd.Expression)
	}
}

func (fd FunctionDefine) Compute(ctx Context) (Expression, error) {
	return fd, nil
}

func (fd FunctionDefine) Computable(ctx Context) bool {
	return false
}

func (fd FunctionDefine) Position() Position {
	return fd.Pos
}

func (fd FunctionDefine) GetArguments() []Identifier {
	return fd.Arguments
}

func (fd FunctionDefine) GetVariableArgument() *Identifier {
	return fd.VariableArgument
}

func (fd FunctionDefine) Call(ctx Context, args map[Identifier]Expression, variables *Object) (Expression, error) {
	newCtx := ctx.MakeScope()
	for k, v := range args {
		v_, err := ctx.ComputeRecursive(v)
		if err != nil {
			return nil, err
		}

		if err := newCtx.Define(k, v_); err != nil {
			return nil, err
		}
	}

	if vi := fd.GetVariableArgument(); vi != nil {
		vo, err := ctx.ComputeRecursive(variables)
		if err != nil {
			return nil, err
		}

		newCtx.Define(*vi, vo)
	}

	return fd.Expression.Compute(newCtx)
}

type FunctionCall struct {
	Function  Expression
	Arguments []Expression
	Pos       Position
}

func (fc FunctionCall) String() string {
	args := make([]string, len(fc.Arguments))
	for i, a := range fc.Arguments {
		args[i] = fmt.Sprint(a)
	}
	return fmt.Sprintf("%s(%s)", fc.Function, strings.Join(args, ", "))
}

func (fc FunctionCall) GetFunction(ctx Context) (Function, error) {
	raw, err := ctx.ComputeRecursive(fc.Function)
	if err != nil {
		return nil, err
	}

	f, ok := raw.(Function)
	if !ok {
		return nil, NotFunctionError{value: fc.Function, pos: fc.Pos}
	}

	return f, nil
}

func (fc FunctionCall) Compute(ctx Context) (Expression, error) {
	f, err := fc.GetFunction(ctx)
	if err != nil {
		return nil, err
	}

	va := f.GetVariableArgument()

	if (va == nil && len(fc.Arguments) != len(f.GetArguments())) || (va != nil && len(fc.Arguments) < len(f.GetArguments())) {
		err := MissmatchArgumentError{
			excepted: len(f.GetArguments()),
			got:      len(fc.Arguments),
			pos:      fc.Position(),
		}
		if ident, ok := fc.Function.(Identifier); ok {
			err.name = ident.String()
		}
		return nil, err
	}

	args := make(map[Identifier]Expression)
	for i, x := range f.GetArguments() {
		args[x] = fc.Arguments[i]
	}

	var obj *Object
	if va != nil {
		obj = NewObject()
		for _, x := range fc.Arguments[len(f.GetArguments()):] {
			obj.Indexed = append(obj.Indexed, x)
		}
	}

	return f.Call(ctx, args, obj)
}

func (fc FunctionCall) Computable(ctx Context) bool {
	return true
}

func (fc FunctionCall) Position() Position {
	return fc.Pos
}
