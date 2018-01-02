package main

type Context struct {
	parent *Context
	values map[Identifier]Expression
}

func NewContext() Context {
	return builtinContext.MakeScope()
}

func (c Context) Get(key Identifier) Expression {
	if v, ok := c.values[key]; ok {
		return v
	} else if c.parent != nil {
		return c.parent.Get(key)
	} else {
		return nil
	}
}

func (c Context) Put(key Identifier, value Expression) {
	for cur := &c; cur != nil; cur = cur.parent {
		if _, ok := cur.values[key]; ok {
			cur.values[key] = value
			return
		}
	}

	c.values[key] = value
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

func (c Context) MakeScope() Context {
	return Context {
		parent: &c,
		values: make(map[Identifier]Expression),
	}
}
