package main

import "mini_lisp"

func main() {
	mini_lisp.Interpret("(if (= a 7) (twice a) (twice 2)) # evaluates to 4 (twice 2)")
}
