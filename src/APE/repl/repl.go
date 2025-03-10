package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"APE/evaluator"
	"APE/lexer"
	"APE/object"
	"APE/parser"
)

const PROMPT = "APE >>"
const CONT_PROMPT = "... "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	var inputBuffer strings.Builder
	inBlock := false

	for {
		if !inBlock {
			fmt.Fprintf(out, PROMPT)
		} else {
			fmt.Fprintf(out, CONT_PROMPT)
		}

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
	
		if inputBuffer.Len() > 0 {
			inputBuffer.WriteString("\n")
		}
		inputBuffer.WriteString(line)

		input := inputBuffer.String()
		braceCount := countBraces(input)

		if braceCount > 0 {
			inBlock = true
			continue
		}
		inBlock = false

		l := lexer.New(input)
		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			inputBuffer.Reset()
			continue
		}

		evaluated := evaluator.Eval(program, env)

		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}

		inputBuffer.Reset()
	}
}

func countBraces(input string) int { 
	count := 0
	for _, ch := range input {
		if ch == '{' {
			count++
		} else if ch == '}' {
			count-- 
		}
	}
	return count 
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "APE found an error!\n")
		io.WriteString(out, " \t"+msg+"\n")
	}
}
