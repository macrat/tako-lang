package main

type Expression interface {
	Compute(Context) (Expression, error)
	Computable(Context) bool
}

type LocatedExpression interface {
	Expression

	Position() Position
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
	Pos     Position
}

type Identifier struct {
	Key string
	Pos Position
}

func NewIdentifier(key string) Identifier {
	return Identifier{
		Key: key,
		Pos: Position{Filename: "builtin"},
	}
}

func (i Identifier) String() string {
	return i.Key
}

func (i Identifier) Compute(ctx Context) (Expression, error) {
	if val, err := ctx.Get(i); err != nil {
		return nil, err
	} else {
		return val, nil
	}
}

func (i Identifier) Computable(ctx Context) bool {
	return true
}

func (i Identifier) Position() Position {
	return i.Pos
}
