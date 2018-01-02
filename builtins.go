package main

var (
	builtinContext = Context{
		Identifier("+"): NewBuiltInFunction(func(ctx Context, args map[Identifier]Expression) (Expression, error) {
			x, err := ctx.ComputeRecursive(args["__builtin_argument_x__"])
			if err != nil {
				return nil, err
			}

			y, err := ctx.ComputeRecursive(args["__builtin_argument_y__"])
			if err != nil {
				return nil, err
			}

			return Number(x.(Number) + y.(Number)), nil
		}, "__builtin_argument_x__", "__builtin_argument_y__"),

		Identifier("-"): NewBuiltInFunction(func(ctx Context, args map[Identifier]Expression) (Expression, error) {
			x, err := ctx.ComputeRecursive(args["__builtin_argument_x__"])
			if err != nil {
				return nil, err
			}

			y, err := ctx.ComputeRecursive(args["__builtin_argument_y__"])
			if err != nil {
				return nil, err
			}

			return Number(x.(Number) - y.(Number)), nil
		}, "__builtin_argument_x__", "__builtin_argument_y__"),

		Identifier("*"): NewBuiltInFunction(func(ctx Context, args map[Identifier]Expression) (Expression, error) {
			x, err := ctx.ComputeRecursive(args["__builtin_argument_x__"])
			if err != nil {
				return nil, err
			}

			y, err := ctx.ComputeRecursive(args["__builtin_argument_y__"])
			if err != nil {
				return nil, err
			}

			return Number(x.(Number) * y.(Number)), nil
		}, "__builtin_argument_x__", "__builtin_argument_y__"),

		Identifier("/"): NewBuiltInFunction(func(ctx Context, args map[Identifier]Expression) (Expression, error) {
			x, err := ctx.ComputeRecursive(args["__builtin_argument_x__"])
			if err != nil {
				return nil, err
			}

			y, err := ctx.ComputeRecursive(args["__builtin_argument_y__"])
			if err != nil {
				return nil, err
			}

			return Number(x.(Number) / y.(Number)), nil
		}, "__builtin_argument_x__", "__builtin_argument_y__"),

		Identifier(";"): NewBuiltInFunction(func(ctx Context, args map[Identifier]Expression) (Expression, error) {
			_, err := ctx.ComputeRecursive(args["__builtin_argument_x__"])
			if err != nil {
				return nil, err
			}

			result, err := ctx.ComputeRecursive(args["__builtin_argument_y__"])
			if err != nil {
				return nil, err
			}

			return result, nil
		}, "__builtin_argument_x__", "__builtin_argument_y__"),

		Identifier("="): NewBuiltInFunction(func(ctx Context, args map[Identifier]Expression) (Expression, error) {
			value, err := ctx.ComputeRecursive(args["__builtin_argument_expr__"])
			if err != nil {
				return nil, err
			}

			ctx[args["__builtin_argument_ident__"].(Identifier)] = value

			return value, nil
		}, "__builtin_argument_ident__", "__builtin_argument_expr__"),
	}
)

type BuiltInFunction struct {
	Arguments []Identifier
	Function  func(Context, map[Identifier]Expression) (Expression, error)
}

func NewBuiltInFunction(fun func(Context, map[Identifier]Expression) (Expression, error), arguments ...string) BuiltInFunction {
	args := make([]Identifier, len(arguments))
	for i, a := range arguments {
		args[i] = Identifier(a)
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

func (bf BuiltInFunction) Call(ctx Context, args map[Identifier]Expression) (Expression, error) {
	return bf.Function(ctx, args)
}

func (bf BuiltInFunction) GetArguments() []Identifier {
	return bf.Arguments
}
