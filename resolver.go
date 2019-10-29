package mini_lisp

import "fmt"

type state struct {
	defines map[string]interface{}
}

func (s *state) assignDefine(name string, value interface{}) {
	if s.defines == nil {
		s.defines = make(map[string]interface{})
	}
	if _, exists := s.defines[name]; exists {
		panic(fmt.Sprintf("trying to reassign existing value '%s'", name))
	}
	s.defines[name] = value
}

func (s *state) resolveDefine(name string) interface{} {
	if d, exists := s.defines[name]; exists {
		return d
	}
	return nil
}

type inbuiltFunction func([]interface{}, *state) interface{}

// initialized in init because of cyclic initialization loop
var inbuiltFunctions map[string]inbuiltFunction

// TODO ?
func init() {
	inbuiltFunctions = map[string]inbuiltFunction{
		"+":      fnAdd,
		"=":      fnEquals,
		"print":  fnPrint,
		"list":   fnList,
		"first":  fnFirst,
		"rest":   fnRest,
		"append": fnAppend,
		"define": fnDefine,
		"if":     fnIf,
	}
}

func getInbuiltFunction(name string) inbuiltFunction {
	if fn, exists := inbuiltFunctions[name]; exists {
		return fn
	}
	return nil
}

func Resolve(ast *AST) {
	s := &state{}
	for _, e := range ast.Expressions {
		resolveExpression(e, s)
	}
}

func resolveExpression(exp *Expression, state *state) interface{} {
	if len(exp.Values) == 0 {
		return exp
	}

	if s, ok := exp.Values[0].(*Symbol); ok {
		if s.Value == "define" && exp.Root != nil {
			panic("function 'define' can only be used at the root level")
		}
		fn := getInbuiltFunction(s.Value)
		if fn != nil {
			return fn(exp.Values[1:], state)
		}
	}

	for i, v := range exp.Values {
		exp.Values[i] = autoResolve(v, state)
	}

	return exp
}

func resolveSymbol(sym *Symbol, state *state) interface{} {
	fn := getInbuiltFunction(sym.Value)
	if fn != nil {
		panic("functions are only allowed as the first symbol of an expression")
	}
	def := state.resolveDefine(sym.Value)
	if def == nil {
		panic(fmt.Sprintf("could not resolve symbol '%s'", sym.Value))
	}
	return def
}

func autoResolve(i interface{}, state *state) interface{} {
	switch i.(type) {
	case *Expression:
		return resolveExpression(i.(*Expression), state)
	case *Symbol:
		return resolveSymbol(i.(*Symbol), state)
	}
	return i
}

func fnAdd(args []interface{}, s *state) interface{} {
	if len(args) != 2 {
		panic("function '+' expects 2 parameters")
	}
	args[0] = autoResolve(args[0], s)
	args[1] = autoResolve(args[1], s)
	i1, ok1 := args[0].(*Int)
	i2, ok2 := args[1].(*Int)
	if !ok1 || !ok2 {
		panic("function '+' expects 2 int parameters")
	}
	return &Int{Value: i1.Value + i2.Value}
}

func fnEquals(args []interface{}, s *state) interface{} {
	if len(args) != 2 {
		panic("function '=' expects 2 parameters")
	}
	args[0] = autoResolve(args[0], s)
	args[1] = autoResolve(args[1], s)
	o1, ok1 := args[0].(Equatable)
	o2, ok2 := args[1].(Equatable)
	if !ok1 || !ok2 {
		panic("function '=' expects 2 equatable parameters")
	}
	if o1.Equals(o2) {
		return &True{}
	}
	// empty expression is equal to false
	return &Expression{}
}

func fnPrint(args []interface{}, s *state) interface{} {
	if len(args) != 1 {
		panic("function 'print' expects 1 parameter")
	}
	args[0] = autoResolve(args[0], s)
	fmt.Println(args[0])
	return args[0]
}

func fnList(args []interface{}, s *state) interface{} {
	return resolveExpression(&Expression{Values: args}, s)
}

func fnFirst(args []interface{}, s *state) interface{} {
	if len(args) == 0 {
		panic("function 'first' expects 1 parameter")
	}
	args[0] = autoResolve(args[0], s)
	list, ok := args[0].(*Expression)
	if !ok {
		panic("function 'first' expects a list as parameter 1")
	}
	if len(list.Values) == 0 {
		panic("function 'first' cannot be called on an empty list")
	}
	return list.Values[0]
}

func fnRest(args []interface{}, s *state) interface{} {
	if len(args) == 0 {
		panic("function 'rest' expects 1 parameter")
	}
	args[0] = autoResolve(args[0], s)
	list, ok := args[0].(*Expression)
	if !ok {
		panic("function 'rest' expects a list as parameter 1")
	}
	if len(list.Values) == 0 {
		panic("function 'rest' cannot be called on an empty list")
	}
	return &Expression{Values: list.Values[1:]}
}

func fnAppend(args []interface{}, s *state) interface{} {
	if len(args) != 2 {
		panic("function 'append' expects 2 parameters")
	}
	args[0] = autoResolve(args[0], s)
	base, ok := args[0].(*Expression)
	if !ok {
		panic("function 'append' expects a list as parameter 1")
	}
	return &Expression{Values: append(base.Values, autoResolve(args[1], s))}
}

func fnDefine(args []interface{}, s *state) interface{} {
	if len(args) != 2 {
		panic("function 'define' expects 2 parameters")
	}
	sym, ok := args[0].(*Symbol)
	if !ok {
		panic("function 'define' expects a symbol as parameter 1")
	}
	s.assignDefine(sym.Value, autoResolve(args[1], s))
	return nil
}

func fnIf(args []interface{}, s *state) interface{} {
	if len(args) != 3 {
		panic("function 'if' expects 3 parameters")
	}
	cond := autoResolve(args[0], s)
	// empty expression is equal to false
	if e, ok := cond.(*Expression); ok && len(e.Values) == 0 {
		return autoResolve(args[2], s)
	}
	return autoResolve(args[1], s)
}
