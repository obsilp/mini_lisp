package mini_lisp

type AST struct {
	Expressions []*Expression
}

func (ast *AST) AddExpression() *Expression {
	e := &Expression{}
	ast.Expressions = append(ast.Expressions, e)
	return e
}

type Expression struct {
	Root   *Expression
	Values []interface{}
}

func (exp *Expression) AddSubExpression() *Expression {
	e := &Expression{Root: exp}
	exp.Add(e)
	return e
}

func (exp *Expression) Add(e interface{}) {
	exp.Values = append(exp.Values, e)
}

type Int struct {
	Value int
}

type Word struct {
	Value string
}
