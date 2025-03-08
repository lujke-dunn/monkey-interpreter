package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"testing"
)


func TestAssignmentExpression(t *testing.T) { 
	tests := []struct {
		input string
		expectedIdentifier string
		expectedValue interface{}
	}{
		{"x = 5;", "x", 5},
		{"y = true;", "y", true},
		{"foobar = y;", "foobar", "y"},
		{"dog = dog + 5;", "dog", "dog + 5"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		assignment, ok := stmt.Expression.(*ast.TestAssignmentExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.AssignmentExpression. got=%T", stmt.Expression)
		}

		if assignment.Name.Value != tt.expectedIdentifier {
			t.Errorf("assignment.Name.Value not '%s' got=%s", tt.expectedIdentifier, assignment.Name.Value)
		}
	}
	
	valueStr := assignment.Value.String()
	var expectedValueStr string 

	switch v := tt.expectedValue.(type) {
	case int: 
		expectedValueStr = fmt.Sprintf("%d", v)
	case bool: 
		expectedValueStr = fmt.Sprintf("%t", v)
	case string:
		expectedValueStr = v
	}

	if !strings.Contain(valueStr, expectedValueStr) {
		t.Errorf("assignment.Value.String() not containing '%s'. got=%s",  expectedValueStr, valueStr)
	}


}

func TestForLoopParsing(t *testing.T) {
	tests := []struct {
		input string
		expectedInit string
		expectedCondition string
		expectedUpdate string
		expectedBody string
	}{
		{
			`for (let i = 0; i < 10; i = i + 1) { x = x + i; }`,
			"let i = 0;",
			"(i < 10)",
			"(i = (i + 1))",
			"(x = (x + i))",
		},		{
			`for (; i < 10; i = i + 1) { x = x + i; }`,
			"",
			"(i < 10)",
			"(i = (i + 1))",
			"(x = (x + i))",
		},
		{
			`for (let i = 0;; i = i + 1) { x = x + i; }`,
			"let i = 0;",
			"",
			"(i = (i + 1))",
			"(x = (x + i))",
		},
		{
			`for (let i = 0; i < 10;) { x = x + i; }`,
			"let i = 0;",
			"(i < 10)",
			"",
			"(x = (x + i))",
		},
		{
			`for (;;) { x = x + 1; }`,
			"",
			"",
			"",
			"(x = (x + 1))",
		},
		{
			`for (i = 0; i < 10; i = i + 1) { x = x + i; }`,
			"(i = 0)",
			"(i < 10)",
			"(i = (i + 1))",
			"(x = (x + i))",
		},
	}

	for i, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("test[%d] - program.Statements does not contain 1 statement. got=%d", i, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("test[%d] - stmt.Expression is not ast.ForExpression. got=%T", i, stmt.Expression)
		}
		
		forLoop, ok := stmt.Expression.(*ast.ForExpression)
		if !ok {
			t.Fatalf("test[%d] - stmt.Expression is not ast.ForExpression. got=%T", i, stmt.Expression)
		}

		if forLoop.Init != nil {
			if forLoop.Init.String() != tt.expectedInit { 
				t.Errorf("test[%d] - init wrong. expected=%q, got=%q", i, tt.expectedInit, forLoop.Init.String())
			}
		} else if tt.expectedInit != "" {
			t.Errorf("test[%d] - init is nil, expected=%q", i, tt.expectedInit)
		}

		if forLoop.Condition != nil {
			if forLoop.Condition.String() != tt.expectedCondition {
				t.Errorf("test[%d] - condition wrong. expected=%q, got=%q", i, tt.expectedCondition, forLoop.Condition.String())
			} 
		} else if tt.expectedCondition != "" {
			t.Errorf("test[%d] - init is nil, expected=%q", i, tt.expectedCondition)
		}

		if forLoop.Update != nil {
			if forLoop.Update.String() != tt.expectedUpdate {
				t.Errorf("test[%d] - update wrong. expected=%q got=%q", i , tt.expectedUpdate, forLoop.Update.String())
			} 
		} else if tt.expectedUpdate != "" {
			t.Errorf("test[%d] - update is nil, expected=%q", i, tt.expectedUpdate)
		}

		if len(forLoop.Body.Statements) != 1 { 
			t.Fatalf("test[%d] - forLoop.Body.Statements does not contain 1 statement. got=%d", i, len(forLoop.Body.Statements))
		}

		bodyStmt, ok := forLoop.Body.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("test[%d] - forLoop.Body.Statements[0] is not ast.ExpressionStatement. got=%T", i, forLoop.Body.Statements[0])
		}

		if bodyStmt.Expression.String() != tt.expectedBody {
			t.Errorf("test[%d]  - body wrong. expected=%q, got=%q", i, tt.expectedBody, bodyStmt.Expression.String())
		}
	}
}

func TestDotOperatorParsing(t *testing.T) {
	input := `[1, 2, 3].map(fn(x) {x * 2});`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statment got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatment. got=%T", program.Statements[0])
	}

	methodCall, ok := stmt.Expression.(*ast.MethodCallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not an ast.MethodCallExpression. got=%T", stmt.Expression)
	}

	if methodCall.Method != "map" {
		t.Fatalf("methodCall.Method not 'map', got=%q", methodCall.Method)
	}

	array, ok := methodCall.Object.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("methodCall.Object is not an ast.ArrayLiteral. got=%T", methodCall.Object)
	}

	if len(array.Elements) != 3 {
		t.Fatalf("array has wrong number of elements. got=%d", len(array.Elements))
	}

	if len(methodCall.Arguments) != 1 {
		t.Fatalf("wrong number of arguments. got=%d", len(methodCall.Arguments))
	}

	function, ok := methodCall.Arguments[0].(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("argument is not an ast.FunctionLiteral. got=%T", methodCall.Arguments[0])
	}

	if len(function.Parameters) != 1 {
		t.Fatalf("function literal has wrong parameters. got=%d", len(function.Parameters))
	}

	if function.Parameters[0].Value != "x" { 
		t.Fatalf("parameter is not x. got=%q", function.Parameters[0].Value)
	}
}

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


func TestParsingHashLiteralsStringKeys(t *testing.T) { 
	input := `{"one": 1, "two": 2, "three": 3}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}

	expected := map[string]int64 {
		"one": 1,
		"two": 2, 
		"three": 3, 
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)

		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
		}

		expectedValue := expected[literal.String()]

		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingEmptyHashLiteral(t *testing.T) {
	input := "{}"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)

	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}
	
	if len(hash.Pairs) != 0 { 
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}
}



func TestParsingArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not ast.ArrayLiteral. got=%T", stmt.Expression)
	}

	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) not 3. got=%d", len(array.Elements))
	}

	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 2)
	testInfixExpression(t, array.Elements[2], 3, "+", 3)
}

func TestStringLiteral(t *testing.T) {
	input := `"hello world";`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p) 

	stmt := program.Statements[0].(*ast.ExpressionStatement) 
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. got=%T", stmt.Expression)
	}

	if literal.Value != "hello world" {
		t.Errorf("literal.Value not %q. got=%q", "hello world", literal.Value)
	}
}


func TestParsinngIndexExpression(t *testing.T) { 
	input := "myArray[1 + 1]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp not *ast.IndexExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, indexExp.Left, "myArray") {
		return 
	}

	if !testInfixExpression(t, indexExp.Index, 1, "+", 1) {
		return
	}
}



func TestCallExpression(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5)"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p) 

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}


	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, exp.Function, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}


func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if (len(program.Statements)) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, program.Statements[0])
	}


	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}


	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n", len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statements was not nil. got=%+v", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) { 
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if (len(program.Statements)) != 1 {
		t.Fatalf("Expected 1 statement got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements is not *ast.Expression statement got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 { 
		t.Errorf("consequence is not 1 statements. got=%d\n", len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not a expression statement got=%T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("got more then one alternative statement expected 1 got=%d\n", len(exp.Alternative.Statements))
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("statement is not an ast.ExpressionStatement. got=%T", exp.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
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


func TestBooleanLiteralExpression(t *testing.T) {
    tests := []struct {
        input string
        expected bool
    }{
        {"true;", true},
        {"false;", false},
    }

    for _, tt := range tests {
        l := lexer.New(tt.input)
        p := New(l)
        program := p.ParseProgram()
        checkParserErrors(t, p)

        if len(program.Statements) != 1 {
            t.Fatalf("program has not enough statements. got=%d",
                len(program.Statements))
        }

        stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
        if !ok {
            t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
                program.Statements[0])
        }

        boolean, ok := stmt.Expression.(*ast.Boolean)
        if !ok {
            t.Fatalf("exp not *ast.Boolean. got=%T", stmt.Expression)
        }
        if boolean.Value != tt.expected {
            t.Errorf("boolean.Value not %t. got=%t", tt.expected,
                boolean.Value)
        }
    }
}


func TestParsingInfixExpressions(t *testing.T) {
	infixTest := []struct {
		input	   string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},	
		{"true == true", true, "==", true}, 
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false}, 
	}

	for _, tt := range infixTest {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
	

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}	

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statments[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

 	if !testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue) {
		return 
	}
}
}


func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input 	 string
		expected string 
	} {
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		}, 
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		}, 
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4", 
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		}, 
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)", 
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)", 
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{	
			"-(5 + 5)",
			"(-(5 + 5))",
		}, 
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}

	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests :=  []struct {
		input    	 string
		operator 	 string
		Value 	     interface{}
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 { 
			t.Fatalf("program.Statements does not contain 1 statement got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statement[0] is not ast.ExpressionStatment got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression) 
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", program.Statements[0])
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not %s got=%s", tt.operator, exp.Operator)
		}
		// if !testIntegerLiteral(t, exp.Right, tt.Value) {
		//	return
		//}
	}
}


func testIdentifier(t *testing.T, exp ast.Expression, value string) bool { 
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}
	
	if ident.TokenLiteral() != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	} 

	return true
}


func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) { 
		case int: 
			return testIntegerLiteral(t, exp, int64(v))
		case int64:
			return testIntegerLiteral(t, exp, v)
		case string: 
			return testIdentifier(t, exp, v) 
		case bool: 
			return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}

	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral() not %t. got=%s", value, bo.TokenLiteral())
		return false
	}

	return true
}


func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}, ) bool {
	opExp, ok := exp.(*ast.InfixExpression) 
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false 
	}

	if opExp.Operator != operator {
		t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) { 
		return false 
	}
	
	return true
}





func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral is not %d got=%s", value, integ.TokenLiteral())
		return false
	}
	return true
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
	tests := []struct {
		input 	string
		expectedIdentifier string
		expectedValue 	   interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input) 
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testLetStatements(t, stmt, tt.expectedIdentifier) {
			return
		}
		
		val := stmt.(*ast.LetStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
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
