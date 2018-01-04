package main

import (
	"fmt"
	"strings"
)

var (
	builtinMethods = map[string]Expression{
		"length": NewBuiltInFunction(func(ctx Context, variables *Object, args map[string]Expression) (Expression, error) {
			self, err := ctx.ComputeRecursive(args["self"])
			if err != nil {
				return nil, err
			}

			return Number(len(self.(*Object).Indexed)), nil
		}, "", "self"),

		"size": NewBuiltInFunction(func(ctx Context, variables *Object, args map[string]Expression) (Expression, error) {
			self, err := ctx.ComputeRecursive(args["self"])
			if err != nil {
				return nil, err
			}

			return Number(len(self.(*Object).Indexed) + len(self.(*Object).Named)), nil
		}, "", "self"),

		"push": NewBuiltInFunction(func(ctx Context, variables *Object, args map[string]Expression) (Expression, error) {
			self, err := ctx.ComputeRecursive(args["self"])
			if err != nil {
				return nil, err
			}

			value, err := ctx.ComputeRecursive(args["value"])
			if err != nil {
				return nil, err
			}

			obj := self.(*Object)
			obj.Indexed = append(obj.Indexed, value)

			return args["self"], nil
		}, "", "self", "value"),

		"pop": NewBuiltInFunction(func(ctx Context, variables *Object, args map[string]Expression) (Expression, error) {
			self, err := ctx.ComputeRecursive(args["self"])
			if err != nil {
				return nil, err
			}

			obj := self.(*Object)
			obj.Indexed = obj.Indexed[:len(obj.Indexed) - 1]

			return args["self"], nil
		}, "", "self"),
	}
)

type Object struct {
	Indexed []Expression
	Named   map[string]Expression
}

func NewObject() *Object {
	return &Object {
		Indexed: []Expression{},
		Named: make(map[string]Expression),
	}
}

func (o *Object) String() string {
	ss := make([]string, len(o.Indexed))
	for i, e := range o.Indexed {
		ss[i] = fmt.Sprint(e)
	}
	for k, v := range o.Named {
		ss = append(ss, fmt.Sprintf("%s: %s", k, v))
	}
	return "[" + strings.Join(ss, ", ") + "]"
}

func (o *Object) Compute(ctx Context) (Expression, error) {
	result := &Object{
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

func (o *Object) Computable(ctx Context) bool {
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

func (o *Object) Get(key Expression) (Expression, error) {
	switch k := key.(type) {
	case Identifier:
		if r, ok := o.Named[k.Key]; ok {
			return r, nil
		} else if r, ok = builtinMethods[k.Key]; ok {
			return r, nil
		} else {
			return nil, NotDefinedError(k)
		}

	case Number:
		i := int(k)

		if 0 <= i && i < len(o.Indexed) {
			return o.Indexed[i], nil
		}

		max := len(o.Indexed) - 1
		if max < 0 {
			max = 0
		}

		return nil, OutOfBoundsError{
			max: max,
			got: i,
		}
	}

	return nil, TypeError{
		name: "index of object",
		excepts: []string{"identifier", "number"},
	}
}
