package parser

import (
	"testing"
	"monkey/ast"
	"monkey/lexer"
)


func TestLetStatements(t *testing.T) {
	input := `
	let x = 5;
	let y = 10;
	let foobar = 838383;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram() 
	if program == nil {
		t.Fatalf("ParseProgram returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements got=%d", len(program.Statements))
	}

	test := []struct {
		expectedIdentifier string
	}{
		{x},
		{y},
		{"foobar"},
	}

	for i, tt := range tests { 
		smts := program.Statements[i]
		
	}
	

}