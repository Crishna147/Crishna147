package main

import (
	"fmt"
	"strings"
)

type Writer interface {
	Write(text string)
}

type ConsoleWriter struct{}

func (cw ConsoleWriter) Write(text string) {
	fmt.Println("Writing to console:", text)
}

type UpperCaseWriter struct {
	Writer Writer
}

func (uw UpperCaseWriter) Write(text string) {
	upperCaseText := strings.ToUpper(text)
	uw.Writer.Write(upperCaseText)
}

func main() {
	consoleWriter := ConsoleWriter{}
	upperCaseWriter := UpperCaseWriter{Writer: consoleWriter}

	upperCaseWriter.Write("Hello, World!")
}
---------------------------------------------------------------------------------
PS D:\go\sample1> go run interface.go   
Writing to console: HELLO, WORLD!


