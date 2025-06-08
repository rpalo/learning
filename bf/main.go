package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

const USAGE = `
usage: bf COMMAND [ARGS...]

commands:
	compile FILENAME: compile the bf file at FILENAME and output the ops.
	run FILENAME: compile the bf file at FILENAME and evaluate
	repl: Initiate an interactive repl
`

var buffer_size = 30000
var debug = false

func init() {
	if val := os.Getenv("BF_BUFFER_SIZE"); val != "" {
		size, err := strconv.Atoi(val)

		if err != nil {
			log.Fatalf("Env var BF_BUFFER_SIZE is not an integer: %s", val)
		}
		buffer_size = size
	}
	if os.Getenv("BF_DEBUG") != "" {
		debug = true
	}
}

func main() {
	if len(os.Args) == 1 || os.Args[1] == "-h" {
		fmt.Print(USAGE)
		return
	}

	command := os.Args[1]
	var ops []Opcode

	if len(os.Args) == 3 {
		bytes_, err := os.ReadFile(os.Args[2])

		if err != nil {
			log.Fatal(err)
		}
		contents := string(bytes_)

		ops, err = Compile(contents)

		if err != nil {
			log.Fatal(err)
		}
	}

	switch command {
	case "compile":
		PrintOps(ops)
	case "run":
		EvalBfOps(ops, buffer_size, debug)
	case "repl":
		repl()
	default:
		fmt.Print(USAGE)
	}
}

func repl() {
	reader := bufio.NewReader(os.Stdin)
	// TODO: make repl remember buffer between lines

	for {
		fmt.Print("bf> ")
		line, err := reader.ReadString('\n')

		if errors.Is(err, io.EOF) {
			return
		} else if err != nil {
			log.Fatal(err)
		}

		line = line[:len(line)-1]

		if len(line) == 0 {
			return
		}
		ops, err := Compile(line)
		if err != nil {
			fmt.Println(err)
		} else {
			EvalBfOps(ops, buffer_size, debug)
		}
	}
}
