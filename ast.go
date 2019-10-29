package mini_lisp

import (
	"fmt"
	"strconv"
	"strings"
)

type Equatable interface {
	Equals(other interface{}) bool
}

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

func (exp Expression) String() string {
	b := strings.Builder{}
	b.WriteString("(")

	for i := range exp.Values {
		if i > 0 {
			b.WriteString(" ")
		}
		b.WriteString(fmt.Sprintf("%s", exp.Values[i]))
	}

	b.WriteString(")")
	return b.String()
}

func (exp Expression) Equals(other interface{}) bool {
	o, ok := other.(*Expression)
	if !ok {
		return false
	}
	if len(exp.Values) != len(o.Values) {
		return false
	}
	for i := range exp.Values {
		e, ok := exp.Values[i].(Equatable)
		if !ok || !e.Equals(o.Values[i]) {
			return false
		}
	}
	return true
}

type Int struct {
	Value int
}

func (i Int) String() string {
	return strconv.Itoa(i.Value)
}

func (i Int) Equals(other interface{}) bool {
	o, ok := other.(*Int)
	if !ok {
		return false
	}
	return i.Value == o.Value
}

type Symbol struct {
	Value string
}

func (s Symbol) String() string {
	return s.Value
}

func (s Symbol) Equals(other interface{}) bool {
	o, ok := other.(*Symbol)
	if !ok {
		return false
	}
	return s.Value == o.Value
}

type True struct {
}

func (t True) String() string {
	return "true"
}

func (t True) Equals(other interface{}) bool {
	_, ok := other.(*True)
	return ok
}
