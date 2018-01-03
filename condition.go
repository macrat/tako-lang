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

	var fun Expression
	if b, ok := cond.(Boolean); !ok {
		return nil, ConditionTypeError{pos: c.Position()}
	} else if b {
		fun = c.Then
	} else {
		fun = c.Else
	}

	if fun == nil {
		return Null{}, nil
	}

	return ctx.ComputeRecursive(FunctionCall{
		Function:  fun,
		Arguments: []Expression{},
		Pos:       c.Position(),
	})
}

func (c Condition) Computable(ctx Context) bool {
	return true
}

func (c Condition) Position() Position {
	return c.Pos
}
