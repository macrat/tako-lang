package main

type Context struct {
	parent *Context
	values map[string]Expression
}

func NewContext() Context {
	return builtinContext.MakeScope()
}

func (c Context) GetByString(key string) (Expression, error) {
	return c.Get(NewIdentifier(key))
}

func (c Context) Get(key Identifier) (Expression, error) {
	if v, ok := c.values[key.Key]; ok {
		return v, nil
	}

	for cur := &c; cur != nil; cur = cur.parent {
		if val, ok := cur.values[key.Key]; ok {
			return val, nil
		}
	}

	return nil, NotDefinedError(key)
}

func (c Context) Put(key Identifier, value Expression) error {
	for cur := &c; cur != nil; cur = cur.parent {
		if _, ok := cur.values[key.Key]; ok {
			cur.values[key.Key] = value
			return nil
		}
	}

	return NotDefinedError(key)
}

func (c Context) Define(key Identifier, value Expression) error {
	if _, ok := c.values[key.Key]; ok {
		return AlreadyDefinedError(key)
	}

	c.values[key.Key] = value

	return nil
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
	return Context{
		parent: &c,
		values: make(map[string]Expression),
	}
}
