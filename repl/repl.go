package repl

import (
	"bufio"
	"fmt"
	"interp/lexer"
	"interp/parser"
	"io"
	"log"
	"os"
)

func Start(in io.Reader, out io.Writer) {
	file, err := os.OpenFile("output.txt", os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(in)
	for {
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		fmt.Println(line)
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		file.WriteString(program.String())
		file.WriteString("\n")

	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Woops!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
