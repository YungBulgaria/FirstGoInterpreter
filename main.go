package main

import (
	"fmt"
    "os"
    "os/user"
	"donkey/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the Donkey programming language!\n", user.Username)
	fmt.Printf("This is my first interpreter so it may suck!\n")
	fmt.Printf("Feel free to type in commands\n")
	repl.Start(os.Stdin, os.Stdout)
}