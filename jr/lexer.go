package main

import (
	"bufio"
	"fmt"
	"io"
)

type ErrLex struct {
	text string
}

func (e ErrLex) Error() string {
	return fmt.Sprintf("Unrecognized character: %s", e.text)
}

func Lex(reader io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanRunes)

	result := make([]string, 0)

	for {
		more := scanner.Scan()
		if !more {
			return result, nil
		}
		c := scanner.Text()

		switch c {
		case "{", "}":
			result = append(result, c)
		default:
			return nil, ErrLex{c}
		}
	}
}
