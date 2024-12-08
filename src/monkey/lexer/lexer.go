package lexer

import (
	"monkey/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func newToken(tokenType token.TokenType, ch byte) token.Token { // initializes a token called from NextToken and returns the type and the associated character 
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) readChar() { // check current character and increments the counter to check next character 
	if l.readPosition >= len(l.input) { // check if its the last character in the input
		l.ch = 0 // indicates EOF or end of input
	} else {
		l.ch = l.input[l.readPosition] // read char
	}
	l.position = l.readPosition
	l.readPosition += 1 // increment to next character
}

func (l *Lexer) NextToken() token.Token { // identifies and returns the next token in the input
	var tok token.Token

	switch l.ch {
		case '=': 
			tok = newToken(token.ASSIGN, l.ch)
		case '+':
			tok = newToken(token.PLUS, l.ch)
		case '(':
			tok = newToken(token.LPAREN, l.ch)
		case ')': 
			tok = newToken(token.RPAREN, l.ch)
		case '{':
			tok = newToken(token.LBRACE, l.ch)
		case '}':
			tok = newToken(token.RBRACE, l.ch)
		case ',': 
			tok = newToken(token.COMMA, l.ch)
		case ';':
			tok = newToken(token.SEMICOLON, l.ch)
		case 0:
			tok.Literal = ""
			tok.Type = token.EOF
	}

	l.readChar()
	return tok
}

