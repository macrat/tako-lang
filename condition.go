package main

import "fmt"

type Condition struct {
	Condition Expression
	Then      Expression
	Else      Expression
	Pos       Position
}

func (c Condition) String() string {
	els := "(){}"
	if c.Else != nil {
		els = fmt.Sprint(c.Else)
	}
	return fmt.Sprintf("if(%s, %s, %s)", c.Condition, c.Then, els)
}

func (c Condition) Compute(ctx Context) (Expression, error) {
	cond, err := ctx.ComputeRecursive(c.Condition)
	if err != nil {
		return nil, err
	}

	var expr Expression
	if b, ok := cond.(Boolean); !ok {
		return nil, ConditionTypeError{pos: c.Position()}
	} else if b {
		expr = c.Then
	} else {
		expr = c.Else
	}

	if expr == nil {
		return Null{}, nil
	}

	return ctx.ComputeRecursive(expr)
}

func (c Condition) Computable(ctx Context) bool {
	return true
}

func (c Condition) Position() Position {
	return c.Pos
}
