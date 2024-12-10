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

func (l *Lexer) peekChar() byte { // looks ahead by one to check the character ahead. 
	if l.readPosition >= len(l.input) { // since readposition is +1 char of position we are looking ahead in order to  
		return 0
	} else {
		return l.input[l.readPosition]
	}
}


func (l *Lexer) readIdentifier() string { // finds the lexers position and then increments until white space returning where it started and finished 
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool { // returns all character from a to Z, and also _ as a special case for functions identifing valid characters. 
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) consumeWhitespace() { // consumes a whitespace if it is found. 
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readNumber() string { // checks the digit of the currently selected character by the lexer
	position := l.position 
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isDigit(ch byte) bool { 
	return '0' <= ch && ch <= '9'
}


func (l *Lexer) NextToken() token.Token { // identifies and returns the next token in the input
	var tok token.Token

	l.consumeWhitespace()

	switch l.ch {
		case '=': 
			if l.peekChar() == '=' {
				ch := l.ch
				l.readChar()
				literal := string(ch) + string(l.ch)
				tok = token.Token{Type: token.EQ, Literal: literal}
			} else {
				tok = newToken(token.ASSIGN, l.ch)
			}
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
		case '|': 
			if l.peekChar() == '|' {
				ch := l.ch
				l.readChar()
				literal := string(ch) + string(l.ch)
				tok = token.Token{Type: token.OR, Literal: literal}
			} else {
				tok = newToken(token.EOF, l.ch)
			}
		case '-':
			tok = newToken(token.MINUS, l.ch)
		case '!':
			if l.peekChar() == '=' {
				ch := l.ch 
				l.readChar()
				literal := string(ch) + string(l.ch)
				tok = token.Token{Type: token.NOT_EQ, Literal: literal}
			} else {
				tok = newToken(token.BANG, l.ch)
			}
		case '/':
			tok = newToken(token.SLASH, l.ch)
		case '*':
			tok = newToken(token.ASTERISK, l.ch)
		case '<':
			tok = newToken(token.LT, l.ch)
		case '>':
			tok = newToken(token.GT, l.ch) 
		case 0:
			tok.Literal = ""
			tok.Type = token.EOF
		default: // if it is not one of the above identifiers read the identifier and return its token else return illegal token
			if isLetter(l.ch) {
				tok.Literal = l.readIdentifier()
				tok.Type = token.LookupIdent(tok.Literal)
				return tok 
			} else if isDigit(l.ch) {
				tok.Type = token.INT
				tok.Literal = l.readNumber()
				return tok
			} else {
				tok = newToken(token.ILLEGAL, l.ch)
			}
	}

	l.readChar()
	return tok
}