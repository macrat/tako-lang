package main

var (
	builtinContext = Context{
		values: map[string]Expression{
			"+": NewBuiltInFunction(func(ctx Context, args map[string]Expression) (Expression, error) {
				x, err := ctx.ComputeRecursive(args["x"])
				if err != nil {
					return nil, err
				}

				y, err := ctx.ComputeRecursive(args["y"])
				if err != nil {
					return nil, err
				}

				return Number(x.(Number) + y.(Number)), nil
			}, "x", "y"),

			"-": NewBuiltInFunction(func(ctx Context, args map[string]Expression) (Expression, error) {
				x, err := ctx.ComputeRecursive(args["x"])
				if err != nil {
					return nil, err
				}

				y, err := ctx.ComputeRecursive(args["y"])
				if err != nil {
					return nil, err
				}

				return Number(x.(Number) - y.(Number)), nil
			}, "y", "x"),

			"*": NewBuiltInFunction(func(ctx Context, args map[string]Expression) (Expression, error) {
				x, err := ctx.ComputeRecursive(args["x"])
				if err != nil {
					return nil, err
				}

				y, err := ctx.ComputeRecursive(args["y"])
				if err != nil {
					return nil, err
				}

				return Number(x.(Number) * y.(Number)), nil
			}, "x", "y"),

			"/": NewBuiltInFunction(func(ctx Context, args map[string]Expression) (Expression, error) {
				x, err := ctx.ComputeRecursive(args["x"])
				if err != nil {
					return nil, err
				}

				y, err := ctx.ComputeRecursive(args["y"])
				if err != nil {
					return nil, err
				}

				return Number(x.(Number) / y.(Number)), nil
			}, "x", "y"),

			"==": NewBuiltInFunction(func(ctx Context, args map[string]Expression) (Expression, error) {
				x, err := ctx.ComputeRecursive(args["x"])
				if err != nil {
					return nil, err
				}

				y, err := ctx.ComputeRecursive(args["y"])
				if err != nil {
					return nil, err
				}

				return Boolean(x == y), nil
			}, "x", "y"),

			"!=": NewBuiltInFunction(func(ctx Context, args map[string]Expression) (Expression, error) {
				x, err := ctx.ComputeRecursive(args["x"])
				if err != nil {
					return nil, err
				}

				y, err := ctx.ComputeRecursive(args["y"])
				if err != nil {
					return nil, err
				}

				return Boolean(x != y), nil
			}, "x", "y"),

			"=": NewBuiltInFunction(func(ctx Context, args map[string]Expression) (Expression, error) {
				value, err := ctx.ComputeRecursive(args["expression"])
				if err != nil {
					return nil, err
				}

				return value, ctx.Put(args["identifier"].(Identifier), value)
			}, "identifier", "expression"),

			":=": NewBuiltInFunction(func(ctx Context, args map[string]Expression) (Expression, error) {
				value, err := ctx.ComputeRecursive(args["expression"])
				if err != nil {
					return nil, err
				}

				return value, ctx.Define(args["identifier"].(Identifier), value)
			}, "identifier", "expression"),
		},
	}
)

type BuiltInFunction struct {
	Arguments []Identifier
	Function  func(Context, map[string]Expression) (Expression, error)
}

func NewBuiltInFunction(fun func(Context, map[string]Expression) (Expression, error), arguments ...string) BuiltInFunction {
	args := make([]Identifier, len(arguments))
	for i, a := range arguments {
		args[i] = NewIdentifier("__builtin_functions_argument_" + a + "__")
	}

	return BuiltInFunction{
		Arguments: args,
		Function:  fun,
	}
}

func (bf BuiltInFunction) Compute(ctx Context) (Expression, error) {
	return bf, nil
}

func (bf BuiltInFunction) Computable(ctx Context) bool {
	return false
}

func (bf BuiltInFunction) Position() Position {
	return Position{Filename: "builtin"}
}

func (bf BuiltInFunction) Call(ctx Context, args map[Identifier]Expression) (Expression, error) {
	as := make(map[string]Expression)
	for _, a := range bf.Arguments {
		as[a.Key[len("__builtin_functions_argument_"):len(a.Key)-len("__")]] = args[a]
	}

	return bf.Function(ctx, as)
}

func (bf BuiltInFunction) GetArguments() []Identifier {
	return bf.Arguments
}
