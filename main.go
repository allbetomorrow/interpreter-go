package main

import (
	"interp/lexer"
	"interp/parser"
	"io"
	"log"
	"os"
)

func main() {
	file, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatal(err)
	}

	file_write, err := os.OpenFile("output.txt", os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file_write.Close()

	l := lexer.New(string(file))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		for _, msg := range p.Errors() {
			io.WriteString(os.Stdout, "\t"+msg+"\n")
		}
		return
	}

	file_write.WriteString(program.String())
	file_write.WriteString("\n")

}
