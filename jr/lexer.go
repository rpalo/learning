package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type ErrLex struct {
	text string
}

func (e ErrLex) Error() string {
	return fmt.Sprintf("Unrecognized character: %s", e.text)
}

type TokenType int

const (
	TokenString TokenType = 1
	TokenRaw              = 2
)

type Token struct {
	kind  TokenType
	value string
}

func Lex(reader io.Reader) ([]Token, error) {
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanRunes)

	result := make([]Token, 0)

	for {
		more := scanner.Scan()
		if !more {
			return result, nil
		}
		c := scanner.Text()

		switch c {
		case "{", "}", ":", ",":
			result = append(result, Token{kind: TokenRaw, value: c})
		case " ", "\n", "\t":
			// ignore whitespace (note could be more robust)
		case "\"":
			text, err := lexString(scanner)

			if err != nil {
				return nil, err
			}
			result = append(result, Token{kind: TokenString, value: text})
		default:
			return nil, ErrLex{c}
		}
	}
}

func lexString(scanner *bufio.Scanner) (string, error) {
	chars := make([]string, 0)

	for scanner.Scan() {
		if scanner.Text() == "\"" && (len(chars) == 0 || chars[len(chars)-1] != "\\") {
			return strings.Join(chars, ""), nil
		}
		chars = append(chars, scanner.Text())
	}
	return "", ErrLex{"Unterminated string"}
}
