package repl

import (
	"bufio"
	"fmt"
	"interp/lexer"
	"interp/token"
	"io"
)

const PROMPT = ">>"

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

		for tok := l.NextToken(); tok.Type != token.LEX_EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
