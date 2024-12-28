package parser

import (
	"testing"
	"monkey/ast"
	"monkey/lexer"
)


func TestReturnStatements(t *testing.T) {
	input := `
	return 5;
	return 10;
	return 993322;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statments does not have 3 statements got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.Returnstatement got=%T", stmt)
			continue
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("not a return token got=%q", returnStmt.TokenLiteral())
		}
	}


}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if (len(program.Statements)) != 1 { 
		t.Fatalf("Expected 1 Statement got=%d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not an ExpressionStatement got=%T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("did not recieve a IntegerLiteral type got=%T", stmt.Expression)
	}

	if literal.Value != 5 { 
		t.Errorf("Literal value != 5 got=%d", literal.Value)
	}

	if literal.TokenLiteral() != "5" { 
		t.Errorf("literal.TokenLiteral() not 5. got=%s", literal.TokenLiteral())
	}


}


func TestIdentifierExpression(t *testing.T) { 
	input := "foobar;"

	l := lexer.New(input) // create a new lexer
	p := New(l) // create a parser
	program := p.ParseProgram() // call the parse program function to run parser operations
	checkParserErrors(t, p) // check for any errors 

	if len(program.Statements) != 1 {
		t.Fatalf("Program does not have 1 statment got=%d", len(program.Statements)) 
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if (!ok) {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier) 

	if !ok {
		t.Fatalf("Expression not a *ast.Identifier. got=%T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value not foobar got=%s", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not foobar got=%s", ident.TokenLiteral())
	}

	



}

func TestLetStatements(t *testing.T) {
	input := `
	let x = 5;
	let y = 10;
	let foobar = 838383;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram() 
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements got=%d", len(program.Statements))
	}

	

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests { 
		smts := program.Statements[i]
		if !testLetStatements(t, smts, tt.expectedIdentifier) {
			return
		}
	}
}

	func testLetStatements(t *testing.T, s ast.Statement, name string) bool {
		if s.TokenLiteral() != "let" {
			t.Errorf("s.TokenLiteral does not equal let got%T=", s)
			return false
		}
	

		letStmt, ok := s.(*ast.LetStatement) 
		if !ok {
			t.Errorf("s not *ast.LetStatement. got=%T", s)
			return false
		}
		

		if letStmt.Name.Value != name {
			t.Errorf("letStmt.Name.Value not '%s', got=%s", name, letStmt.Name.Value)
			return false
		}

		if letStmt.Name.TokenLiteral() != name {
			t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s", name, letStmt.Name.TokenLiteral())
		}

		return true
	}


	func checkParserErrors(t *testing.T, p *Parser) {
		errors := p.Errors()
		if len(errors) == 0 {
			return
		}

		t.Errorf("parser has %d errors", len(errors))
		for _, msg := range errors {
			t.Errorf("parser error: %q", msg)
		}
		t.FailNow()
	}
