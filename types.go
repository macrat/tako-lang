package main

import (
	"fmt"
	"strings"
)

type Number float64

func (n Number) String() string {
	return fmt.Sprint(float64(n))
}

func (n Number) Compute(ctx Context) (Expression, error) {
	return n, nil
}

func (n Number) Computable(ctx Context) bool {
	return false
}

type String string

func (s String) String() string {
	x := strings.Replace(string(s), "\\", "\\\\", -1)
	x = strings.Replace(x, "\n", "\\n", -1)
	x = strings.Replace(x, "\r", "\\r", -1)
	x = strings.Replace(x, "\t", "\\t", -1)
	x = strings.Replace(x, "'", "\\'", -1)
	return fmt.Sprintf("'%s'", x)
}

func (s String) Compute(ctx Context) (Expression, error) {
	return s, nil
}

func (s String) Computable(ctx Context) bool {
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
