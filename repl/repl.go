//repl/repl.go

package repl

import (
	"bufio" // for reading input
	"fmt"
	"io"
	"monkey/lexer"
	"monkey/parser"
)

const PROMPT = ">> "
const WELCOME_ASCII = `
 ### ###  ##   ##  ###  ##   ##  ##   ##   #  ##  ##
 #######  ##   ##  #### ##   ## ##    ## #    ##  ##
 #######  ##   ##  ## ####   ####     ####     ####
 ## # ##  ##   ##  ##  ###   ## ##    ## #      ##
 ##   ##  ##   ##  ##   ##   ##  ##   ##   #    ##
 ##   ##   #####   ##   ##  ###  ##  #######   ####


`

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
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
		io.WriteString(out, program.String()+"\n")
	}
}

func printParseErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
