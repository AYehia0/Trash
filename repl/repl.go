/*
the REPL reads input, sends it to the interpreter for evaluation, prints the result/output of the interpreter and starts again. Read, Eval, Print, Loop.
*/

package repl

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"trash/eval"
	"trash/lexer"
	"trash/object"
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
	env := object.NewEnv()

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

		evaluated := eval.Eval(prog, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func logErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Oops errors :'( \n")
	io.WriteString(out, "Parser Errors: \n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

func StartWithFile(input io.Reader, output io.Writer) {
	env := object.NewEnv()
	content, err := ioutil.ReadAll(input)
	if err != nil {
		fmt.Fprintln(output, "Error reading the file:", err)
		return
	}

	// Convert the content to a string
	codeBlock := string(content)

	// Parse and evaluate the code block
	l := lexer.New(codeBlock)
	p := parser.New(l)
	program := p.Parse()

	if len(p.Errors()) != 0 {
		logErrors(output, p.Errors())
	} else {
		evaluated := eval.Eval(program, env)
		if evaluated != nil {
			fmt.Fprintln(output, evaluated.Inspect())
		}
	}
}
func IsStartOfBlock(line string) bool {
	line = strings.TrimSpace(line)
	if len(line) == 0 {
		return false // Handle empty lines
	}
	lastChar := line[len(line)-1]
	return lastChar == '{'
}

// IsEndOfBlock checks if a line marks the end of a code block.
func IsEndOfBlock(line string) bool {
	line = strings.TrimSpace(line)
	if len(line) == 0 {
		return false // Handle empty lines
	}
	lastChar := line[len(line)-1]
	return lastChar == '}'
}
