tako-lang
=========

	(X){


my experimental language.

might working. (probably not)

``` text
mod := (x, y){
	if x >= y {
		mod(x - y, y)
	} else {
		x
	}
}

fizzbuzz := (n){
	loop := (i){
		println(if mod(i, 15) == 0 {
			"fizzbuzz"
		} else if mod(i, 3) == 0 {
			"fizz"
		} else if mod(i, 5) == 0 {
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
