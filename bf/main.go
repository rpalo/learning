package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strconv"
)

const USAGE = `
usage: bf [-h|FILENAME]

Evaluate brainf*ck.  If FILENAME is provided, run that.  Otherwise run a REPL.
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
	if slices.Contains(os.Args, "-h") {
		fmt.Print(USAGE)
		return
	}

	if len(os.Args) == 2 {
		contents, err := os.ReadFile(os.Args[1])

		if err != nil {
			log.Fatal(err)
		}

		ops, err := Compile(string(contents))

		if err != nil {
			log.Fatal(err)
		}

		if debug {
			PrintOps(ops)
		}
		EvalBfOps(ops, buffer_size, debug)
	} else {
		repl()
	}
}

func repl() {
	reader := bufio.NewReader(os.Stdin)

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
