package main

import (
	"interp/repl"
	"os"
)

func main() {

	repl.Start(os.Stdin, os.Stdout)
}
