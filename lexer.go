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

var consumers = []func(*stream) TokenType{
	consumeWhitespace,
	consumeNewline,
	consumeComment,
	consumeExpressionOpen,
	consumeExpressionClose,
	consumeInt,
	consumeSymbol,
}

type stream struct {
	buffer   []rune
	startPos int
	pos      int
}

func (s *stream) Current() rune {
	if !s.Done() {
		return s.buffer[s.pos]
	}
	return rune(0)
}

func (s *stream) Read() rune {
	c := s.Current()
	s.pos++
	return c
}

func (s *stream) Reset() {
	s.pos = s.startPos
}

func (s *stream) HasAdvanced() bool {
	return s.pos > s.startPos
}

func (s *stream) ReadToken() string {
	if !s.HasAdvanced() {
		panic("pointer was not advanced before reading token")
	}
	token := s.buffer[s.startPos:s.pos]
	s.startPos = s.pos
	return string(token)
}

func (s *stream) Done() bool {
	return s.pos >= len(s.buffer)
}

func Tokenize(input string) []Token {
	var tokens []Token

	s := &stream{buffer: []rune(input)}

	currentLine := 1
	lineStartPos := 0

	for !s.Done() {
		t := consumeNext(s)

		if t == TokenTypeIgnore {
			continue
		}
		if t == TokenTypeInvalid {
			panic(fmt.Sprintf("invalid char '%c' at line %d:%d", s.Current(), currentLine, s.pos-lineStartPos+1))
		}

		token := Token{Type: t, Str: s.ReadToken()}
		tokens = append(tokens, token)

		if t == TokenTypeNewline {
			currentLine++
			lineStartPos = s.pos + 1
		}
	}

	return tokens
}

func consumeNext(s *stream) TokenType {
	for _, c := range consumers {
		s.Reset()
		t := c(s)
		if t != TokenTypeIgnore {
			return t
		}
	}

	// ignore any other control characters
	if unicode.IsControl(s.Read()) {
		return TokenTypeIgnore
	}

	return TokenTypeInvalid
}

func consumeWhitespace(s *stream) TokenType {
	for s.Current() == ' ' {
		s.Read()
	}
	if s.HasAdvanced() {
		return TokenTypeWhitespace
	}
	return TokenTypeIgnore
}

func consumeNewline(s *stream) TokenType {
	if s.Read() == '\n' {
		return TokenTypeNewline
	}
	return TokenTypeIgnore
}

func consumeComment(s *stream) TokenType {
	if s.Read() != '#' {
		return TokenTypeIgnore
	}
	for s.Current() != 0 && s.Current() != '\n' {
		s.Read()
	}
	return TokenTypeComment
}

func consumeExpressionOpen(s *stream) TokenType {
	if s.Read() == '(' {
		return TokenTypeExpOpen
	}
	return TokenTypeIgnore
}

func consumeExpressionClose(s *stream) TokenType {
	if s.Read() == ')' {
		return TokenTypeExpClose
	}
	return TokenTypeIgnore
}

func consumeInt(s *stream) TokenType {
	c := s.Read()

	if !unicode.IsDigit(c) {
		if c == '+' || c == '-' {
			// do not consume the token if the sign is not followed by a number
			if !unicode.IsDigit(s.Read()) {
				return TokenTypeIgnore
			}
		} else {
			return TokenTypeIgnore
		}
	}

	for unicode.IsDigit(s.Current()) {
		s.Read()
	}
	return TokenTypeInt
}

func consumeSymbol(s *stream) TokenType {
	const invalid = " ()#\n"
	for s.Current() != 0 && !strings.ContainsAny(string(s.Current()), invalid) {
		s.Read()
	}
	if s.HasAdvanced() {
		return TokenTypeSymbol
	}
	return TokenTypeIgnore
}
