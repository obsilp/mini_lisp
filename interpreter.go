package mini_lisp

import "io/ioutil"

func Interpret(input string) {
	tokens := Tokenize(input)
	ast := Parse(tokens)
	Resolve(ast)
}

func InterpretFile(path string) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	Interpret(string(buf))
}
