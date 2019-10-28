package main

import "mini_lisp"

func main() {
	i := mini_lisp.NewInterpreter()
	i.Interpret("(print # 1) ())( )(99()() alskdföl\n    ((test1324 4539asdf +test +))")
}
