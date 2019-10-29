package mini_lisp

import (
	"fmt"
	"unicode"
)

type TokenType uint8

const (
	TokenTypeInvalid TokenType = iota
	TokenTypeIgnore

	TokenTypeWhitespace
	TokenTypeNewline
	TokenTypeComment

	TokenTypeExpOpen
	TokenTypeExpClose

	TokenTypeInt
	TokenTypeSymbol
)

func (t TokenType) String() string {
	return map[TokenType]string{
		TokenTypeInvalid: "TokenTypeInvalid",
		TokenTypeIgnore:  "TokenTypeIgnore",

		TokenTypeWhitespace: "TokenTypeWhitespace",
		TokenTypeNewline:    "TokenTypeNewline",
		TokenTypeComment:    "TokenTypeComment",

		TokenTypeExpOpen:  "TokenTypeExpOpen",
		TokenTypeExpClose: "TokenTypeExpClose",

		TokenTypeInt:    "TokenTypeInt",
		TokenTypeSymbol: "TokenTypeSymbol",
	}[t]
}

type Token struct {
	Type TokenType
	Str  string
}

func Tokenize(input string) []Token {
	var tokens []Token

	var lastCharTokenType TokenType
	var currentToken Token

	currentLine := 1
	lineStartPos := 0

	for i, c := range input {
		if c == '\n' {
			currentLine++
			lineStartPos = i + 1
		}

		charTokenType := getTokenTypeForChar(c, currentToken)
		end := isSingleToken(charTokenType)

		if charTokenType == TokenTypeIgnore {
			continue
		}
		if charTokenType == TokenTypeInvalid {
			panic(fmt.Sprintf("invalid char '%c' at line %d:%d", c, currentLine, i-lineStartPos+1))
		}

		if charTokenType != lastCharTokenType || end {
			if lastCharTokenType != TokenTypeInvalid {
				tokens = append(tokens, currentToken)
			}
			currentToken = Token{Type: charTokenType, Str: string(c)}
		} else {
			currentToken.Str += string(c)
		}

		lastCharTokenType = charTokenType
	}

	// add last pending token
	tokens = append(tokens, currentToken)

	return tokens
}

func getTokenTypeForChar(c rune, currentToken Token) TokenType {
	// allow comment to always interrupt a token
	if currentToken.Type != TokenTypeComment && c == '#' {
		return TokenTypeComment
	}

	// support multi-character tokens by checking the current state
	switch currentToken.Type {
	// continue comments until a newline
	case TokenTypeComment:
		if c != '\n' {
			return TokenTypeComment
		}

	// continue symbols until the expression is closed or a space is found
	case TokenTypeSymbol:
		if c != ')' && c != ' ' {
			return TokenTypeSymbol
		}

	// reduce more than one whitespace to one
	case TokenTypeWhitespace:
		if c == ' ' {
			return TokenTypeIgnore
		}
	}

	switch c {
	case ' ':
		return TokenTypeWhitespace
	case '#':
		return TokenTypeComment
	case '\n':
		return TokenTypeNewline
	case '(':
		return TokenTypeExpOpen
	case ')':
		return TokenTypeExpClose

	// special function names
	case '=', '+':
		return TokenTypeSymbol
	}

	if unicode.IsNumber(c) {
		return TokenTypeInt
	}
	if unicode.IsLetter(c) {
		return TokenTypeSymbol
	}

	// ignore any other control characters
	if unicode.IsControl(c) {
		return TokenTypeIgnore
	}

	return TokenTypeInvalid
}

func isSingleToken(t TokenType) bool {
	switch t {
	case TokenTypeNewline:
		return true
	case TokenTypeExpOpen, TokenTypeExpClose:
		return true
	}
	return false
}
