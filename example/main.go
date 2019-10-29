package main

import "mini_lisp"

func main() {
	mini_lisp.Interpret("(print (if (= 4 4) (+ 1 1) (+ 2)))")
}
