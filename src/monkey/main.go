package main 

import (
	"fmt"
	"os"
	"os/user"
	"monkey/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the monkey programming language! \n", user.Username)
	fmt.Printf("Type out commands\n")
	repl.Start(os.Stdin, os.Stdout)
}