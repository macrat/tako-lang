package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kingpin"
)

var (
	source = kingpin.Arg("source", "source file.").ExistingFile()
	debug  = kingpin.Flag("debug", "show debug messages.").Bool()
)

func Parse(file *os.File) Expression {
	l := NewLexer(file)
	l.Filename = file.Name()

	yyParse(l)

	return l.result
}

func main() {
	kingpin.Parse()

	file := os.Stdin
	if *source != "" {
		var err error
		file, err = os.Open(*source)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	}

	expr := Parse(file)

	if *debug {
		fmt.Println(expr)
	}

	ctx := NewContext()
	if _, err := ctx.ComputeRecursive(expr); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}
