package mini_lisp

import "fmt"

type state struct {
	defines map[string]interface{}
	funcs   map[string]*customFunc
}

func (s *state) clone() *state {
	clone := &state{}
	for k, v := range s.defines {
		clone.setDefine(k, v)
	}
	for k, v := range s.funcs {
		clone.setFunc(k, v)
	}
	return clone
}

func (s *state) setDefine(name string, value interface{}) {
	if _, exists := s.defines[name]; exists {
		panic(fmt.Sprintf("trying to reassign existing value '%s'", name))
	}
	s.setDefineUnsafe(name, value)
}

func (s *state) setDefineUnsafe(name string, value interface{}) {
	if s.defines == nil {
		s.defines = make(map[string]interface{})
	}
	s.defines[name] = value
}

func (s *state) resolveDefine(name string) interface{} {
	if d, exists := s.defines[name]; exists {
		return d
	}
	return nil
}

func (s *state) setFunc(name string, f *customFunc) {
	if s.funcs == nil {
		s.funcs = make(map[string]*customFunc)
	}
	if _, exists := s.funcs[name]; exists {
		panic(fmt.Sprintf("trying to reassign existing function '%s'", name))
	}
	s.funcs[name] = f
}

func (s *state) resolveFunc(name string) *customFunc {
	if f, exists := s.funcs[name]; exists {
		return f
	}
	return nil
}

type customFunc struct {
	name       string
	exp        *Expression
	paramNames []string
}

func (f *customFunc) invoke(args []interface{}, s *state) interface{} {
	if len(args) != len(f.paramNames) {
		panic(fmt.Sprintf("parameter count for function '%s' is not correct. expecting %d got %d", f.name, len(f.paramNames), len(args)))
	}
	stateClone := s.clone()
	for i, n := range f.paramNames {
		stateClone.setDefineUnsafe(n, autoResolve(args[i], s))
	}
	return autoResolve(f.exp, stateClone)
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
		cFn := state.resolveFunc(s.Value)
		if cFn != nil {
			return cFn.invoke(exp.Values[1:], state)
		}
	}

	values := make([]interface{}, len(exp.Values))
	for i, v := range exp.Values {
		values[i] = autoResolve(deepCopyExpression(v), state)
	}
	return &Expression{Values: values}
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

func deepCopyExpression(i interface{}) interface{} {
	if e, ok := i.(*Expression); ok {
		values := make([]interface{}, len(e.Values))
		for i, v := range e.Values {
			values[i] = deepCopyExpression(v)
		}
		return &Expression{Values: values}
	}
	return i
}

func fnAdd(args []interface{}, s *state) interface{} {
	if len(args) != 2 {
		panic("function '+' expects 2 parameters")
	}
	i1, ok1 := autoResolve(args[0], s).(*Int)
	i2, ok2 := autoResolve(args[1], s).(*Int)
	if !ok1 || !ok2 {
		panic("function '+' expects 2 int parameters")
	}
	return &Int{Value: i1.Value + i2.Value}
}

func fnEquals(args []interface{}, s *state) interface{} {
	if len(args) != 2 {
		panic("function '=' expects 2 parameters")
	}
	o1, ok1 := autoResolve(args[0], s).(Equatable)
	o2, ok2 := autoResolve(args[1], s).(Equatable)
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
	e := autoResolve(args[0], s)
	fmt.Println(e)
	return e
}

func fnList(args []interface{}, s *state) interface{} {
	return resolveExpression(&Expression{Values: args}, s)
}

func fnFirst(args []interface{}, s *state) interface{} {
	if len(args) == 0 {
		panic("function 'first' expects 1 parameter")
	}
	list, ok := autoResolve(args[0], s).(*Expression)
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
	list, ok := autoResolve(args[0], s).(*Expression)
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
	base, ok := autoResolve(args[0], s).(*Expression)
	if !ok {
		panic("function 'append' expects a list as parameter 1")
	}
	base = deepCopyExpression(base).(*Expression)

	add := autoResolve(args[1], s)
	if e, ok := add.(*Expression); ok {
		for _, v := range e.Values {
			base.Values = append(base.Values, v)
		}
	} else {
		base.Values = append(base.Values, add)
	}

	return &Expression{Values: base.Values}
}

func fnDefine(args []interface{}, s *state) interface{} {
	if len(args) == 2 {
		return defineSymbol(args, s)
	} else if len(args) == 4 {
		return defineFunction(args, s)
	}
	panic("function 'define' expects either 2 or 4 parameters")
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

func defineSymbol(args []interface{}, s *state) interface{} {
	sym, ok := args[0].(*Symbol)
	if !ok {
		panic("function 'define' expects a symbol as parameter 1")
	}
	s.setDefine(sym.Value, autoResolve(args[1], s))
	return nil
}

func defineFunction(args []interface{}, s *state) interface{} {
	if fmt.Sprint(args[1]) != "lambda" {
		panic("function 'define' with 4 parameter is only supported for lambda definitions")
	}

	name, ok := args[0].(*Symbol)
	if !ok {
		panic("function 'define' expects a symbol as parameter 1")
	}
	params, ok := args[2].(*Expression)
	if !ok {
		panic("function 'define' expects an expression as parameter 3")
	}
	def, ok := args[3].(*Expression)
	if !ok {
		panic("function 'define' expects an expression as parameter 4")
	}

	paramNames := make([]string, 0, len(def.Values))
	for _, v := range params.Values {
		sym, ok := v.(*Symbol)
		if !ok {
			panic("function 'define' only allows symbols as parameter names for lambda definitions")
		}
		paramNames = append(paramNames, sym.Value)
	}

	f := &customFunc{
		name:       name.Value,
		exp:        def,
		paramNames: paramNames,
	}
	s.setFunc(name.Value, f)

	return nil
}
