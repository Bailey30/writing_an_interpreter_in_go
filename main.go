package main

import (
	// "fmt"
	"monkey/lexer"
	"monkey/parser"
	// "monkey/repl"
	// "os"
	// "os/user"
)

func main() {
	// user, err := user.Current()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("Hello %s! This is the Monkey programming language!\n", user.Username)
	// fmt.Printf("Feel free to type commands\n")
	// repl.Start(os.Stdin, os.Stdout)

	input := `
    let x = 5;
    let y = 10;
    let foobar = 838383;
    `

	l := lexer.New(input)
	p := parser.New(l)

	p.ParseProgram()

}
