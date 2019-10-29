package mini_lisp

import (
	"fmt"
	"strings"
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

var consumers = []func([]rune, *int) TokenType{
	consumeWhitespace,
	consumeNewline,
	consumeComment,
	consumeExpressionOpen,
	consumeExpressionClose,
	consumeInt,
	consumeSymbol,
}

func Tokenize(input string) []Token {
	var tokens []Token

	pos := 0
	stream := []rune(input)

	currentLine := 1
	lineStartPos := 0

	for pos < len(input) {
		start := pos
		t := consumeNext(stream, &pos)

		if t == TokenTypeIgnore {
			continue
		}
		if t == TokenTypeInvalid {
			panic(fmt.Sprintf("invalid char '%c' at line %d:%d", stream[pos], currentLine, pos-lineStartPos+1))
		}

		str := input[start:pos]
		token := Token{Type: t, Str: str}
		tokens = append(tokens, token)

		if t == TokenTypeNewline {
			currentLine++
			lineStartPos = pos + 1
		}
	}

	return tokens
}

func consumeNext(stream []rune, pos *int) TokenType {
	start := *pos

	for _, c := range consumers {
		*pos = start
		t := c(stream, pos)
		if t != TokenTypeIgnore {
			return t
		}
	}

	// ignore any other control characters
	if unicode.IsControl(stream[*pos]) {
		return TokenTypeIgnore
	}

	return TokenTypeInvalid
}

func consumeWhitespace(stream []rune, pos *int) TokenType {
	startPos := *pos
	for *pos < len(stream) && stream[*pos] == ' ' {
		*pos++
	}
	if *pos != startPos {
		return TokenTypeWhitespace
	}
	return TokenTypeIgnore
}

func consumeNewline(stream []rune, pos *int) TokenType {
	if stream[*pos] == '\n' {
		*pos++
		return TokenTypeNewline
	}
	return TokenTypeIgnore
}

func consumeComment(stream []rune, pos *int) TokenType {
	if stream[*pos] != '#' {
		return TokenTypeIgnore
	}
	for *pos < len(stream) && stream[*pos] != '\n' {
		*pos++
	}
	return TokenTypeComment
}

func consumeExpressionOpen(stream []rune, pos *int) TokenType {
	if stream[*pos] != '(' {
		return TokenTypeIgnore
	}
	*pos++
	return TokenTypeExpOpen
}

func consumeExpressionClose(stream []rune, pos *int) TokenType {
	if stream[*pos] != ')' {
		return TokenTypeIgnore
	}
	*pos++
	return TokenTypeExpClose
}

func consumeInt(stream []rune, pos *int) TokenType {
	c := stream[*pos]

	if !unicode.IsDigit(c) {
		if c == '+' || c == '-' {
			// do not consume the token if the sign is not followed by a number
			*pos++
			if *pos >= len(stream) || !unicode.IsDigit(stream[*pos]) {
				return TokenTypeIgnore
			}
		} else {
			return TokenTypeIgnore
		}
	}

	for *pos < len(stream) && unicode.IsDigit(stream[*pos]) {
		*pos++
	}
	return TokenTypeInt
}

func consumeSymbol(stream []rune, pos *int) TokenType {
	const invalid = "()# "
	startPos := *pos
	for *pos < len(stream) && !strings.ContainsAny(string(stream[*pos]), invalid) {
		*pos++
	}
	if *pos != startPos {
		return TokenTypeSymbol
	}
	return TokenTypeIgnore
}
