package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	repl()

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
		fmt.Println(line)
	}
}
