package main

import "mini_lisp"

func main() {
	mini_lisp.Interpret("(define a 5)(print (= (+ a 1) (+ a 1)))")
}
