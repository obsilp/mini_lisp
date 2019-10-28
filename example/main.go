package main

import "mini_lisp"

func main() {
	i := mini_lisp.NewInterpreter()
	i.Interpret("(if (= a 7) (twice a) (twice 2)) # evaluates to 4 (twice 2)")
}
