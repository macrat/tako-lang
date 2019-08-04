tako-lang
=========

	(X){


my experimental language.

might working. (probably not)

``` text
fizzbuzz := (n){
	loop := (i){
		println(if i % 15 == 0 {
			"fizzbuzz"
		} else if i % 3 == 0 {
			"fizz"
		} else if i % 5 == 0 {
			"buzz"
		} else {
			i
		})

		if i < n {
			loop(i + 1)
		}
	}
	
	loop(1)
}

fizzbuzz(30)
```

## feature
- everything is the expression and they have value
- dynamic typing
- object that like table of Lua
- no GC (plans; not yet)

## how to build
``` shell
$ go get https://github.com/macrat/tako-lang.git && cd tako-lang
$ go generate
$ go build
```
