package main

var (
	builtinContext = Context{
		values: map[string]Expression{
			"null": Null{},
			"true": Boolean(true),
			"false": Boolean(false),

			":+:": NewBuiltInFunction(func(ctx Context, args map[string]Expression) (Expression, error) {
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

			":-:": NewBuiltInFunction(func(ctx Context, args map[string]Expression) (Expression, error) {
				x, err := ctx.ComputeRecursive(args["x"])
				if err != nil {
					return nil, err
				}

				y, err := ctx.ComputeRecursive(args["y"])
				if err != nil {
					return nil, err
				}

				return Number(x.(Number) - y.(Number)), nil
			}, "x", "y"),

			":*:": NewBuiltInFunction(func(ctx Context, args map[string]Expression) (Expression, error) {
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

			":/:": NewBuiltInFunction(func(ctx Context, args map[string]Expression) (Expression, error) {
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

			"-:": NewBuiltInFunction(func(ctx Context, args map[string]Expression) (Expression, error) {
				x, err := ctx.ComputeRecursive(args["x"])
				if err != nil {
					return nil, err
				}

				return Number(-x.(Number)), nil
			}, "x"),

			"!:": NewBuiltInFunction(func(ctx Context, args map[string]Expression) (Expression, error) {
				x, err := ctx.ComputeRecursive(args["x"])
				if err != nil {
					return nil, err
				}

				return Boolean(!x.(Boolean)), nil
			}, "x"),

			":==:": NewBuiltInFunction(func(ctx Context, args map[string]Expression) (Expression, error) {
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

			":!=:": NewBuiltInFunction(func(ctx Context, args map[string]Expression) (Expression, error) {
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

			":<:": NewBuiltInFunction(func(ctx Context, args map[string]Expression) (Expression, error) {
				x, err := ctx.ComputeRecursive(args["x"])
				if err != nil {
					return nil, err
				}

				y, err := ctx.ComputeRecursive(args["y"])
				if err != nil {
					return nil, err
				}

				return Boolean(x.(Number) < y.(Number)), nil
			}, "x", "y"),

			":<=:": NewBuiltInFunction(func(ctx Context, args map[string]Expression) (Expression, error) {
				x, err := ctx.ComputeRecursive(args["x"])
				if err != nil {
					return nil, err
				}

				y, err := ctx.ComputeRecursive(args["y"])
				if err != nil {
					return nil, err
				}

				return Boolean(x.(Number) < y.(Number)), nil
			}, "x", "y"),

			":>:": NewBuiltInFunction(func(ctx Context, args map[string]Expression) (Expression, error) {
				x, err := ctx.ComputeRecursive(args["x"])
				if err != nil {
					return nil, err
				}

				y, err := ctx.ComputeRecursive(args["y"])
				if err != nil {
					return nil, err
				}

				return Boolean(x.(Number) > y.(Number)), nil
			}, "x", "y"),

			":>=:": NewBuiltInFunction(func(ctx Context, args map[string]Expression) (Expression, error) {
				x, err := ctx.ComputeRecursive(args["x"])
				if err != nil {
					return nil, err
				}

				y, err := ctx.ComputeRecursive(args["y"])
				if err != nil {
					return nil, err
				}

				return Boolean(x.(Number) >= y.(Number)), nil
			}, "x", "y"),

			":=:": NewBuiltInFunction(func(ctx Context, args map[string]Expression) (Expression, error) {
				value, err := ctx.ComputeRecursive(args["expression"])
				if err != nil {
					return nil, err
				}

				return value, ctx.Put(args["identifier"].(Identifier), value)
			}, "identifier", "expression"),

			"::=:": NewBuiltInFunction(func(ctx Context, args map[string]Expression) (Expression, error) {
				value, err := ctx.ComputeRecursive(args["expression"])
				if err != nil {
					return nil, err
				}

				return value, ctx.Define(args["identifier"].(Identifier), value)
			}, "identifier", "expression"),

			":.:": NewBuiltInFunction(func(ctx Context, args map[string]Expression) (Expression, error) {
				object, err := ctx.ComputeRecursive(args["object"])
				if err != nil {
					return nil, err
				}

				identifier := args["identifier"].(Identifier)

				if val, ok := object.(Object).Named[identifier.Key]; ok {
					return val, nil
				}

				return nil, NotDefinedError(identifier)
			}, "object", "identifier"),

			":[]:": NewBuiltInFunction(func(ctx Context, args map[string]Expression) (Expression, error) {
				object, err := ctx.ComputeRecursive(args["object"])
				if err != nil {
					return nil, err
				}

				index := int(args["index"].(Number))

				if 0 <= index && index < len(object.(Object).Indexed) {
					return object.(Object).Indexed[index], nil
				}

				max := len(object.(Object).Indexed) - 1
				if max < 0 {
					max = 0
				}

				return nil, OutOfBoundsError{
					max: max,
					got: index,
				}
			}, "object", "index"),

			":.=:": NewBuiltInFunction(func(ctx Context, args map[string]Expression) (Expression, error) {
				object, err := ctx.ComputeRecursive(args["object"])
				if err != nil {
					return nil, err
				}

				value, err := ctx.ComputeRecursive(args["value"])
				if err != nil {
					return nil, err
				}

				identifier := args["identifier"].(Identifier)

				if _, ok := object.(Object).Named[identifier.Key]; ok {
					object.(Object).Named[identifier.Key] = value
					return value, nil
				}

				return nil, NotDefinedError(identifier)
			}, "object", "identifier", "value"),

			":.:=:": NewBuiltInFunction(func(ctx Context, args map[string]Expression) (Expression, error) {
				object, err := ctx.ComputeRecursive(args["object"])
				if err != nil {
					return nil, err
				}

				value, err := ctx.ComputeRecursive(args["value"])
				if err != nil {
					return nil, err
				}

				identifier := args["identifier"].(Identifier)

				if _, ok := object.(Object).Named[identifier.Key]; ok {
					return nil, AlreadyDefinedError(identifier)
				}

				object.(Object).Named[identifier.Key] = value
				return value, nil
			}, "object", "identifier", "value"),

			":[]=:": NewBuiltInFunction(func(ctx Context, args map[string]Expression) (Expression, error) {
				object, err := ctx.ComputeRecursive(args["object"])
				if err != nil {
					return nil, err
				}

				value, err := ctx.ComputeRecursive(args["value"])
				if err != nil {
					return nil, err
				}

				index := int(args["index"].(Number))

				if 0 <= index && index < len(object.(Object).Indexed) {
					object.(Object).Indexed[index] = value
					return value, nil
				}

				max := len(object.(Object).Indexed) - 1
				if max < 0 {
					max = 0
				}

				return nil, OutOfBoundsError{
					max: max,
					got: index,
				}
			}, "object", "index", "value"),

			":[]:=:": NewBuiltInFunction(func(ctx Context, args map[string]Expression) (Expression, error) {
				_, err := ctx.ComputeRecursive(args["object"])
				if err != nil {
					return nil, err
				}

				_, err = ctx.ComputeRecursive(args["value"])
				if err != nil {
					return nil, err
				}

				return nil, TypeError{name: "index", excepts: []string{"string"}}
			}, "object", "index", "value"),
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
