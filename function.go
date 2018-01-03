package main

type Function interface {
	Expression

	Call(Context, map[Identifier]Expression) (Expression, error)
	GetArguments() []Identifier
}

type FunctionDefine struct {
	Arguments  []Identifier
	Expression Expression
	Pos        Position
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

func (fd FunctionDefine) Call(ctx Context, args map[Identifier]Expression) (Expression, error) {
	newCtx := ctx.MakeScope()
	for k, v := range args {
		if err := newCtx.Define(k, v); err != nil {
			return nil, err
		}
	}
	return fd.Expression.Compute(newCtx)
}

type FunctionCall struct {
	Function  Expression
	Arguments []Expression
	Pos       Position
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

	if len(fc.Arguments) != len(f.GetArguments()) {
		err := MissmatchArgumentError{
			excepted: len(f.GetArguments()),
			got: len(fc.Arguments),
			pos: fc.Position(),
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

	return f.Call(ctx, args)
}

func (fc FunctionCall) Computable(ctx Context) bool {
	return true
}

func (fc FunctionCall) Position() Position {
	return fc.Pos
}
