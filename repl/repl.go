/*
the REPL reads input, sends it to the interpreter for evaluation, prints the result/output of the interpreter and starts again. Read, Eval, Print, Loop.
*/

package repl

import (
	"bufio"
	"fmt"
	"io"
	"trash/lexer"
	"trash/token"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.New(line)
		for t := l.NextToken(); t.Type != token.EOF; t = l.NextToken() {
			fmt.Printf("%+v\n", t)
		}
	}
}
