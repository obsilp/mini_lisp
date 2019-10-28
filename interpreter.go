package mini_lisp

type Interpreter struct {
	resolver *Resolver
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		resolver: &Resolver{},
	}
}

func (i *Interpreter) Interpret(input string) {
	tokens := Tokenize(input)
	ast := Parse(tokens)
	i.resolver.Resolve(ast)
}
