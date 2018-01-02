tako: parser.go *.go
	go build -o $@

parser.go: parser.go.y
	goyacc -o $@ $^

.PHONY: clean
clean:
	- rm parser.go y.output tako
