package main

import "fmt"

type Number int64

func (n Number) String() string {
	return fmt.Sprint(int64(n))
}

func (n Number) Compute(ctx Context) (Expression, error) {
	return n, nil
}

func (n Number) Computable(ctx Context) bool {
	return false
}

type Boolean bool

func (b Boolean) String() string {
	return fmt.Sprint(bool(b))
}

func (b Boolean) Compute(ctx Context) (Expression, error) {
	return b, nil
}

func (b Boolean) Computable(ctx Context) bool {
	return false
}

type Null struct{}

func (n Null) String() string {
	return "null"
}

func (n Null) Compute(ctx Context) (Expression, error) {
	return n, nil
}

func (n Null) Computable(ctx Context) bool {
	return false
}
