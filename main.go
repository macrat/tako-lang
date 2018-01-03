package main

import (
	"fmt"
	"io"
	"os"
)

func Parse(file io.Reader) Expression {
	l := NewLexer(file)

	yyParse(l)

	return l.result
}

func main() {
	expr := Parse(os.Stdin)

	fmt.Println(expr)

	ctx := NewContext()
	if result, err := expr.Compute(ctx); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(result)
	}
}
