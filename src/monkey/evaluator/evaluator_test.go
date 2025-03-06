package evaluator

import (
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)


func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input string
		expected int64
	}{
		{"return 5;", 5 },
		{"return 2 * 5;", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}


func TestWhileExpression(t *testing.T) {
	tests := []struct {
		input string
		expected interface{}
	}{
		{
			"let x = 0; while (x < 5) { let x = x + 1; }; x;",
			5, 
		},
		{
			"let a = 0; while (a < 10) { let a = a + 1; if (a == 5) { return a; }; }; 10;",
			5, 
		},
		{
			"let a = 0; let b = 0; while (a < 3) { let a = a + 1; let b = b + a; }; b;",
			6,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, int64(tt.expected.(int)))
	}
}


func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("dog")`, 3},
		{`len("jhin")`, 4},
		{`len(1)`, "argument to `len` not supported, got=INTEGER"},
		{`len("one", "two")`, "wrong number of arguments. got=2, want=1"},
		{`random()`, "wrong number of arguments. got=0, want=1"},
		{`random("abc")`, "argument to `random` not supported, got=STRING"},
		{`random(true)`, "argument to `random` not supported, got=BOOLEAN"},
		{`random(0)`, "argument to `random` must be a positive integer, got=0"},
		{`random(-6)`, "argument to `random` must be a positive integer, got=-6"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error) 
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
			}
		}
	}
}

func TestForLoopEvaluation(t *testing.T) {
	tests := []struct {
		input string
		expected int64
	}{ 
		{
			`
			let sum = 0;
			for (let i = 0; i < 5; i = i + 1) {
				sum = sum + i; 
			}
			sum; 
			`,
			10,
		},
		{
			`
			let sum = 0; 
			let i = 0; 
			for (; i < 5; i = i + 1) {
				sum = sum + i; 
			}
			sum;
			`,
			10, 
		},
		{
			`
			let sum = 0; 
			let i = 0;
			for (;;) {
				if (i >= 5) {
					break;
				}
			}
			sum;
			`,
			10, 
		},
		{
			`
			let sum = 0;
			for (let i = 0; i < 3; i = i + 1) {
				for (let j = 0; j < 2; j = j + 1) {
					sum = sum + (i * j);
				}
			}
			`,
			3,  
		},
	}

	for i, tt  := range tests {
		evaluated := testEval(tt.input)
		integer, ok := evaluated.(*object.Integer)

		if !ok {
			t.Errorf("test[%d] - object is not Integer. got=%T (%+v)", i, evaluated, evaluated)
			continue
		}

		if integer.Value != tt.expected {
			t.Errorf("test[%d] - wrong value. expected=%d, got=%d", i, tt.expected, integer.Value)
		}
	}
}

func TestBreakStatements(t *testing.T) { 
	tests := []struct {
		input string
		expected interface{}
	}{
	{
		`
		let result = 0; 
		for (int i = 0; i < 5; i = i + 1) {
			for (let j = 0; j < 5; j = j + 1) {
				result = i * 10 + j; 
				if (j == 3) {
					break; 
				}
			}
			if (i == 2) {
				break; 
			}
		}
		result;
		`,
		23,
	}, 
	{
		`
		let count = 0; 
		let result = 0; 
		while (count < 10) {
			count = count + 1; 
			result = count; 
			if (count == 5) { 
				break; 
			}
		}
		result; 
		`,
		5, 
		}
	}

	for i, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := evaluated.(*object.Integer)

		if !ok {
			t.Errorf("test[%d] - object is not Integer. got=%T (%+v)", i, evaluated, evaluated)
			continue
		}

		expected, ok := tt.expected.(int)
		if !ok {
			continue
		}

		if integer.Value != int64(expected) {
			t.Errorf("test[%d] - wrong value. expected=%d, got=%d", i, expected, integer.Value)
		}
	}
}

func TestForLoopScope(t *testing.T) {
	tests := []struct {
		input string 
		expected int64
	}{
		{
			`
			let result = 0; 
			for (let i = 0; i < 5; i = i + 1) {
				result = i; 
			}
			let i = 100; 
			result + 1;
			`,
			104,
		},
		{
			`
			let i = 100; 
			let result = 0; 
			for (let i = 0; i < 5; i = i + 1) { 
				result = i; 
			}
			result + i; 
			`,
			104,
		},
		{
			`
			let sum = 0; 
			for (let i = 0; i < 5; i = i + 1) {
				sum = sum + 1; 
			}
			sum; 
			`, 
			10,
		},
	}

	for i, tt := range tests { 
		evaluated := testEval(tt.input)

		integer, ok := evaluated.(*object.Integer)
		if !ok {
			t.Errorf("test[%d] - object is not Integer. got=%T (%+v)", i, evaluated, evaluated)
			continue
		}

		if integer.Value != tt.expected {
			t.Errorf("test[%d] - wrong value. expected=%d, got=%d", i, tt.expected, integer.Value)
		}
	}
}

func TestArrayMapMethod(t *testing.T) {
	tests := []struct {
		input string 
		expected interface{}
	}{
		{`[1, 2, 3].map(fn(x) { x * 4 })`, []int{2, 4, 6}},
		{`[1, 2, 3].map(fn(x) { x })`, []int{1, 2, 3}},
		{`[1, 2, 3].map(fn(x) { x * x})`, []int{1, 4, 9}},
		{`[].map(fn(x) { x })`, []int{}},
		{`["a", "b", "c"].map(fn(x) { x + "!" })`, []string{"a!", "b!", "c!"}},
		{`["hello", "world"].map(fn(x) { len(x) })`, []int{5, 5}},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case []int:
				testIntegerArray(t, evaluated, expected)
		case []string:
				testStringArray(t, evaluated, expected)
		}
	}
}


func TestArrayFilterMethod(t *testing.T) {
	tests := []struct {
		input string
		expected interface{}
	}{
		{`[1, 2, 3, 4].filter(fn(x) { x > 2 })`, []int{3, 4}},
		{`[1, 2, 3, 4].filter(fn(x) { x < 3 })`, []int{1, 2}},
		{`["a", "ab", "abc"].filter(fn(x) { len(x) > 1 })`, []string{"ab", "abc"}},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) { 
			case []int:
				testIntegerArray(t, evaluated, expected)
			case []string: 
				testStringArray(t, evaluated, expected)
		}
	}
}


func TestArrayReduceMethod(t *testing.T) { 
	tests := []struct {
		input string
		expected interface{}
	}{ 
		{`[1,2,3,4].reduce(fn(acc, x) { acc + x }, 0)`, 6},
		{`[].reduce(fn(acc, x) {acc + x}, 5)`, 5},
		{`[1, 2, 3].reduce(fn(acc, x) { acc * x }, 1)`, 6},
		{`["Josh", "Luke", "Nath"].reduce(fn(acc, name) {
			if (acc == "") {
				name
			} else {
				acc + ", " + name
			}
		}, "")`, "Josh, Luke, Nath"},
		{`["h","e","l","l","o"].reduce(fn(acc, x) { acc + x }, "")`, "hello"},
		{`["a", "b", "c"].reduce(fn(acc, x) { x + acc }, "")`, "cba"},
	}
	
	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
			case int:
				testIntegerObject(t, evaluated, int64(expected))
			case string:
				testStringObject(t, evaluated, expected)
		}
	}


}


func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input string
		expected interface{}
	}{ 
			{
				"[1, 2, 3][0];",
				1,
			},
			{
				"[1, 2, 3][1];",
				2,
			},
			{
				"[1, 2, 3][2];",
				3,
			},
			{
				"let myArray = [1,2,3]; myArray[2];",
				3,
			},
			{
				"let i = 0; [1][i];",
				1,
			},
			{
				"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
				6,
			},
			{
				"[1, 2, 3][3];",
				nil,
			}, 
			{
				"[1, 2, 3][-1];",
				nil,
			},

	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}

}

func TestArrayLiterals(t *testing.T) { 
	input := "[1, 2 * 2, 3 + 3]"

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}

	if len(result.Elements) != 3 { 
		t.Fatalf("array has wrong num of elements. got=%d", len(result.Elements))
	}

	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}


func TestLetStatements(t *testing.T) {
	tests := []struct {
		input string
		expected int64
	}{
		{"let a = 5; a;", 5}, 
		{"let a = 5 * 5; a;", 25}, 		
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)

	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) { 
	input := `"Hello" + " " + "World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)

	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}

}


func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"

	evaluated := testEval(input)

	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("Object is not function got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v", fn.Parameters)
	}

	expectedBody := "(x + 2)"

	if fn.Body.String() != expectedBody {
		t.Fatalf("Body is not %q. got=%q", expectedBody, fn.Body.String())

	}
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input string
		expected int64
	}{
		{"5", 5}, 
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60}, 
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected) 
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 > 3", false},
		{"3 < 1", false},
		{"3 > 2", true},
		{"5 < 6", true}, 
		{"1 == 1", true},
		{"1 != 1", false}, 
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true}, 
		{"false == false", true}, 
		{"true == false", false}, 
		{"true != false", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false}, 
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true}, 
	}
	

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		input string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10}, 
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)

		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}


func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input string
		expected bool
	}{
		{"!true", false},
		{"!!true", true},
		{"!false", true},
		{"!!false", false},
		{"!5", false}, 
		{"!!5", true}, 
	}

	for _, tt := range tests {
		evaluator := testEval(tt.input)
		testBooleanObject(t, evaluator, tt.expected)
	}
}
 

func testEval(input string) object.Object { 
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()
	return Eval(program, env)
}


func testStringObject(t *testing.T, obj object.Object, expected string) bool {
	result, ok := obj.(*object.String)
	if !ok {
		t.Errorf("object is not String. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("String has wrong value. expected=%q, got=%q", expected, result.Value)
		return false
	}

	return true
}



func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)

	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}

	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)

	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
		return false
	}

	return true
}

func testIntegerArray(t *testing.T, obj object.Object, expected []int) bool {
	array, ok := obj.(*object.Array)

	if !ok {
		t.Errorf("object is not Array. got=%T, (%+v)", obj, obj)
		return false
	}

	if len(array.Elements) != len(expected) {
		t.Errorf("wrong array length. expected=%d, got=%d", len(expected), len(array.Elements))
		return false
	}

	for i, expectedElement := range expected {
		 if !testIntegerObject(t, array.Elements[i], int64(expectedElement)) {
			 return false
		 }
	}

	return true
}


// i love repeating myself might turn this into a switch statement inside this later as a refactor

func testStringArray(t *testing.T, obj object.Object, expected []string) bool {
	array, ok := obj.(*object.Array)
	if !ok {
		t.Errorf("Object not an Array. got=%T (%+v)", obj, obj)
		return false
	}

	if len(array.Elements) != len(expected) {
		t.Errorf("wrong array length. expected=%d, got=%d", len(expected), len(array.Elements))
		return false
	}

	for i, expectedElement := range expected {
		if !testStringObject(t, array.Elements[i], expectedElement) {
			return false
		}
	}

	return true
}



func TestHashLiterals(t *testing.T)  {
	input := `let two = "two";
	{
		"one": 10 - 9,
		two: 1 + 1,
		"thr" + "ee": 6 / 2, 
		4: 4,
		true: 5,
		false: 6
	}`

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("Eval didn't return Hash. got=%T (%+v)", evaluated, evaluated)
	}

	expected := map[object.HashKey]int64 {
		(&object.String{Value: "one"}).HashKey(): 1,
		(&object.String{Value: "two"}).HashKey(): 2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey(): 4,
		TRUE.HashKey(): 5,
		FALSE.HashKey(): 6,
	}

	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong num of pairs. got=%d", len(result.Pairs))
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key in Pairs")
		}

		testIntegerObject(t, pair.Value, expectedValue)
	}

}


func TestHashIndexExpressions(t *testing.T) { 
	tests := []struct {
		input string
		expected interface{}
	}{
	{
		`{"foo": 5}["foo"]`,
		5, 
	},
	{
		`{"foo": 5}["bar"]`,
		nil,
	},
	{
		`{}["bar"]`,
		nil,
	},
	{
		`{5: 5}[5]`,
		5,
	},
	{
		`{true: 5}[true]`,
		5,
	},
	{	
		`{false: 5}[false]`, 
		5,
	},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

