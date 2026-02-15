package repl

import (
	"bufio"
	"fmt"
	"interp/lexer"
	"interp/token"
	"io"
	"log"
	"os"
)

const PROMPT = ">>"

func Start(in io.Reader, out io.Writer) {
	file, err := os.OpenFile("output.txt", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
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
			s := fmt.Sprintf("{Type: %s, Literal: %s}\n", tok.Type, tok.Literal)
			file.WriteString(s)
		}
		fmt.Printf("\n")
	}
}
