package main

type Number int64

func (n Number) Compute(ctx Context) (Expression, error) {
	return n, nil
}

func (n Number) Computable(ctx Context) bool {
	return false
}

type Boolean bool

func (b Boolean) Compute(ctx Context) (Expression, error) {
	return b, nil
}

func (b Boolean) Computable(ctx Context) bool {
	return false
}
