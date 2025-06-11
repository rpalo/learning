package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: jr FILENAME")
	}
	f, err := os.Open(os.Args[1])

	if err != nil {
		log.Fatal("Could not open input file.")
	}
	defer f.Close()

	tokens, err := Lex(f)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Got tokens: %v", tokens)

	obj, err := Parse(tokens)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(obj)
}
