package main

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"trash/repl"
)

const INTER_NAME = "Trash"

func main() {
	args := os.Args

	if len(args) == 2 {
		filePath := args[1]
		// Reading from a file
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			os.Exit(1)
		}
		defer file.Close()

		repl.StartWithFile(bufio.NewReader(file), os.Stdout)
	} else if len(args) == 1 {

		user, err := user.Current()
		if err != nil {
			fmt.Printf("Error : %s ", err.Error())
		}
		fmt.Printf("Hi %s!, Ever heard of %s ?\n", user.Username, INTER_NAME)
		fmt.Printf("Type something in %s\n", INTER_NAME)
		repl.Start(os.Stdin, os.Stdout)
	}
}
