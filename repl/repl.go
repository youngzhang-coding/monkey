//repl/repl.go

package repl

import (
	"bufio" // for reading input
	"fmt"
	"io"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
)

const PROMPT = ">> "
const WELCOME_ASCII = `
 ##   ##   #####   ##   ##   ##  ###  ####### ##  ##
 ### ###  ##   ##  ###  ##   ##  ##   ##      ##  ##
 #######  ##   ##  #### ##   ## ##    ##      ##  ##
 #######  ##   ##  ## ####   ####     ######   ####
 ## # ##  ##   ##  ##  ###   ## ##    ##        ##
 ##   ##  ##   ##  ##   ##   ##  ##   ##        ##
 ##   ##   #####   ##   ##  ###  ##   #######  ####


`

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	environment := object.NewEnvironment()
	fmt.Fprintf(out, WELCOME_ASCII)
	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParseErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, environment)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect()+"\n")
		}
	}
}

func printParseErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
