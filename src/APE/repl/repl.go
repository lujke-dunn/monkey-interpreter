package repl

import (
	"bufio"
	"fmt"
	"io"

=======
	"strings"
	"APE/evaluator"
	"APE/lexer"
	"APE/object"
	"APE/parser"
>>>>>>> 02d6d68 (changed name to APE and added the ability to run files):src/APE/repl/repl.go
)

const PROMPT = ">>"

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)

		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}


func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "Monkey find error!\n")
		io.WriteString(out, " \t"+msg+"\n")
	}
}
