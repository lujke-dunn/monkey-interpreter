package parser

import (
	"monkey/lexer"
	"monkey/token"
	"monkey/ast"
)


type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	p.nextToken() // read the next two tokens
	p.nextToken() 

	return p
}


func (p *Parser) nextToken() {
	 p.curToken = p.peekToken
	 p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program { 
	return nil
}