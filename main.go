package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/Sheep42/Monkey-Lang/repl"
)

func main() {
	user, err := user.Current()

	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s! Welcome to Monkey!\n", user.Username)
	fmt.Printf("Feel free to type in commands\n\n")

	repl.Start(os.Stdin, os.Stdout)
}
