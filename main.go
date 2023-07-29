package main

import (
	"fmt"
	"os"
	"os/user"
	"trash/repl"
)

const INTER_NAME = "Trash"

func main() {
	user, err := user.Current()
	if err != nil {
		fmt.Printf("Error : %s ", err.Error())
	}
	fmt.Printf("Hi %s!, Ever heard of %s ?\n", user.Username, INTER_NAME)
	fmt.Printf("Type something in %s\n", INTER_NAME)
	repl.Start(os.Stdin, os.Stdout)
}
