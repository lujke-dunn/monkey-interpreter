package ast

import (
	"bytes"
	"monkey/token"
)

type LetStatement struct { 
	Token token.Token
	Name *Identifier
	Value Expression
}


func (ls *LetStatement) statementNode() {}

func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

func (ls *LetStatement) String() string { // This is writing and retrieving the total let statement expression here 
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + "  ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")
	
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

type ReturnStatement struct {
	Token token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}

func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

func (rs *ReturnStatement) String() string { // makes the return statement have a human readable form 
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")
	
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()

}


type ExpressionStatement struct {
	Token token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}


type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode() { }

func (b *Boolean) TokenLiteral() string { return b.Token.Literal}

func (b *Boolean) String() string { return b.Token.Literal }



type IntegerLiteral struct { 
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}

func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal}

func (il *IntegerLiteral) String() string {return il.Token.Literal} // used for more human readable code in the test cases

type InfixExpression struct {
	Token token.Token
	Left Expression
	Operator string
	Right Expression	
}
 
func (ie *InfixExpression) expressionNode() {}

func (ie *InfixExpression) TokenLiteral() string { 
    return ie.Token.Literal 
}

func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}


type PrefixExpression struct {
	Token token.Token
	Operator string
	Right Expression
}

func (pe *PrefixExpression) expressionNode() {}

func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal}

func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) String() string {
	return i.Value
}

// TokenLiteral returns the literal value of the token associated with the ExpressionStatement.
// This is useful for debugging and logging purposes, as it provides the exact string representation
// of the token.


func (i *Identifier) expressionNode() {}

func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
 
type Node interface {
	TokenLiteral() string
	String() string // used for debugging and logging purposes
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

// TokenLiteral returns the literal value of the token associated with the Program. 
// This is useful for debugging and logging purposes, as it provides the exact string representation
func (p *Program) tokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}