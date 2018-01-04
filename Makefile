tako: parser_gen.go *.go
	go build -o $@

parser_gen.go: parser.go.y
	go generate

.PHONY: clean
clean:
	- rm parser_gen.go y.output tako
