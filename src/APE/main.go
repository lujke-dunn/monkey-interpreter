package main 

import (
	"fmt"
	"os"
	"strings"
	"io/ioutil"
	"os/user"
	"APE/repl"
	"APE/lexer"
	"APE/parser"
	"APE/evaluator"
	"APE/object"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Printf("Hello %s! Welcome to the APE programming language! \n", user.Username)
		fmt.Printf("Type out commands\n")
		repl.Start(os.Stdin, os.Stdout)
	} else {
		filename := args[0]


		if !strings.HasSuffix(filename, ".ape") {
			fmt.Printf("Error: Expected file with .ape extension\n")
			return
		}

		content, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Printf("Error reading file: %s\n", err)
			return 
		}

		executeAPE(string(content))
	}
}

func executeAPE(input string) {
	env := object.NewEnvironment()
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		fmt.Print("error")
		return
	}
	
	evaluated := evaluator.Eval(program, env)

	if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
		fmt.Println(evaluated.Inspect())
	}
}

