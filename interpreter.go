package mini_lisp

type Interpreter struct {
	parser   *Parser
	resolver *Resolver
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		parser:   &Parser{},
		resolver: &Resolver{},
	}
}

func (i *Interpreter) Interpret(input string) {
	tokens := Tokenize(input)
	ast := i.parser.Parse(tokens)
	i.resolver.Resolve(ast)
}
