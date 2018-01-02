package main

type Context map[Identifier]Expression

func NewContext() Context {
	return builtinContext.Copy()
}

func (c Context) ComputeRecursive(expr Expression) (result Expression, err error) {
	r := expr
	for r.Computable(c) {
		r, err = r.Compute(c)
		if err != nil {
			return nil, err
		}
	}
	return r, nil
}

func (c Context) Copy() Context {
	result := make(map[Identifier]Expression)
	for k, v := range c {
		result[k] = v
	}
	return result
}
