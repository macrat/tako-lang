package main

import (
	"fmt"
	"os"
)

func Parse(file *os.File) Expression {
	l := NewLexer(file)
	l.Filename = file.Name()

	yyParse(l)

	return l.result
}

func main() {
	expr := Parse(os.Stdin)

	fmt.Println(expr)

	ctx := NewContext()
	if result, err := ctx.ComputeRecursive(expr); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(result)
	}
}
