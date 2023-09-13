/*
the REPL reads input, sends it to the interpreter for evaluation, prints the result/output of the interpreter and starts again. Read, Eval, Print, Loop.
*/

package repl

import (
	"bufio"
	"fmt"
	"io"
	"trash/lexer"
	"trash/parser"
)

const TRASH_ICON = `
 _                 _     
| |               | |			___^___
| |_ _ __ __ _ ___| |__	               |-------|
| __| '__/ _| / __| '_ \		| | | |
| |_| | | (_| \__ \ | | |		| | | |
 \__|_|  \__,_|___/_| |_|		|_____|
`

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	fmt.Print(TRASH_ICON)
	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.New(line)

		parser := parser.New(l)
		prog := parser.Parse()
		if len(parser.Errors()) != 0 {
			logErrors(out, parser.Errors())
			continue
		}
		io.WriteString(out, prog.String())
		io.WriteString(out, "\n")
	}
}

func logErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Oops errors :'( \n")
	io.WriteString(out, "Parser Errors: \n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
