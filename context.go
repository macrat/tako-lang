package main

type Context map[Identifier]Expression

func NewContext() Context {
	return builtinContext.Copy()
}

func (c Context) ComputeRecursive(expr Expression) (Expression, error) {
	r := expr
	for {
		n, err := r.Compute(c)
		if err != nil {
			return nil, err
		}
		if n == r {
			return n, nil
		}
		r = n
	}
}

func (c Context) Copy() Context {
	result := make(map[Identifier]Expression)
	for k, v := range c {
		result[k] = v
	}
	return result
}
