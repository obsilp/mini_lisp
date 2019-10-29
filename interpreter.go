package mini_lisp

func Interpret(input string) {
	tokens := Tokenize(input)
	ast := Parse(tokens)
	Resolve(ast)
}
