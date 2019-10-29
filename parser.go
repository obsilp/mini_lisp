package mini_lisp

import (
	"fmt"
	"strconv"
)

type parserState uint8

const (
	parserStateEmpty parserState = iota
	parserStateExpression
	parserStateParam
	parserStateMore
)

func (s parserState) String() string {
	return map[parserState]string{
		parserStateEmpty:      "parserStateEmpty",
		parserStateExpression: "parserStateExpression",
		parserStateParam:      "parserStateParam",
		parserStateMore:       "parserStateMore",
	}[s]
}

var parserTransitions = map[parserState]map[TokenType]parserState{
	parserStateEmpty: {
		TokenTypeExpOpen: parserStateExpression,
	},
	parserStateExpression: {
		TokenTypeExpClose: parserStateEmpty,
		TokenTypeExpOpen:  parserStateParam, // recursion happens here
		TokenTypeInt:      parserStateParam,
		TokenTypeSymbol:   parserStateParam,
	},
	parserStateParam: {
		TokenTypeExpClose:   parserStateEmpty,
		TokenTypeWhitespace: parserStateMore,
	},
	parserStateMore: {
		TokenTypeExpOpen: parserStateParam, // recursion happens here
		TokenTypeInt:     parserStateParam,
		TokenTypeSymbol:  parserStateParam,
	},
}

func Parse(tokens []Token) *AST {
	var stateStack []parserState
	var state parserState

	ast := &AST{}
	var currentExpression *Expression

	for _, t := range tokens {
		// ignore all newlines and comments
		if t.Type == TokenTypeNewline || t.Type == TokenTypeComment {
			continue
		}
		// ignore whitespaces outside any statements
		if currentExpression == nil && t.Type == TokenTypeWhitespace {
			continue
		}

		if _, ok := parserTransitions[state][t.Type]; !ok {
			panic(fmt.Sprintf("invalid transition %s in state %s", t.Type, state))
		}

		// handle state transitions
		oldState := state
		state = parserTransitions[state][t.Type]
		if t.Type == TokenTypeExpOpen {
			// the first transition need a manual override to ensure that the empty
			// state is committed as the root of the stack
			if oldState == parserStateEmpty {
				state = parserStateEmpty
			}

			stateStack = append(stateStack, state)

			// this step allows for recursion to take place
			state = parserStateExpression

			if currentExpression == nil {
				currentExpression = ast.AddExpression()
			} else {
				currentExpression = currentExpression.AddSubExpression()
			}
		}
		if state == parserStateEmpty {
			if len(stateStack) == 0 || currentExpression == nil {
				panic("non-equal amount of statement openings and closings")
			}

			state = stateStack[len(stateStack)-1]
			stateStack = stateStack[:len(stateStack)-1]

			currentExpression = currentExpression.Root
		}

		// parse everything that is not an expression
		if currentExpression != nil {
			switch t.Type {
			case TokenTypeInt:
				i, _ := strconv.Atoi(t.Str)
				currentExpression.Add(&Int{Value: i})
			case TokenTypeSymbol:
				currentExpression.Add(&Symbol{Value: t.Str})
			}
		}
	}

	return ast
}
