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
	TokenString  TokenType = 1
	TokenRaw     TokenType = 2
	TokenNumber  TokenType = 3
	TokenKeyword TokenType = 4
)

var tokenTypeNames = map[TokenType]string{
	TokenString:  "String",
	TokenRaw:     "Raw",
	TokenNumber:  "Number",
	TokenKeyword: "Keyword",
}

type Token struct {
	kind  TokenType
	value string
}

func Lex(reader io.Reader) ([]Token, error) {
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanRunes)

	result := make([]Token, 0)
	scanner.Scan()

	for {
		var more bool
		var err error
		switch c := scanner.Text(); c {
		case "{", "}", ":", ",", "[", "]":
			result = append(result, Token{kind: TokenRaw, value: c})
			more = scanner.Scan()
		case " ", "\n", "\t":
			more = scanner.Scan()
		case "\"":
			var text string
			text, more, err = lexString(scanner)

			if err != nil {
				return nil, err
			}
			result = append(result, Token{kind: TokenString, value: text})
		case "-", ".", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			var value string
			value, more, err = lexNumber(scanner)

			if err != nil {
				return nil, err
			}
			result = append(result, Token{kind: TokenNumber, value: value})
		case "t":
			err := scanKeyword(scanner, "rue")

			if err != nil {
				return nil, err
			}
			result = append(result, Token{kind: TokenKeyword, value: "true"})
			more = scanner.Scan()
		case "f":
			err := scanKeyword(scanner, "alse")
			if err != nil {
				return nil, err
			}
			result = append(result, Token{kind: TokenKeyword, value: "false"})
			more = scanner.Scan()
		case "n":
			err := scanKeyword(scanner, "ull")
			if err != nil {
				return nil, err
			}
			result = append(result, Token{kind: TokenKeyword, value: "null"})
			more = scanner.Scan()
		default:
			return nil, ErrLex{c}
		}
		if !more {
			break
		}
	}
	return result, nil
}

func lexString(scanner *bufio.Scanner) (string, bool, error) {
	chars := make([]string, 0)

	for scanner.Scan() {
		if scanner.Text() == "\"" && (len(chars) == 0 || chars[len(chars)-1] != "\\") {
			more := scanner.Scan() // consume the closing quote
			return strings.Join(chars, ""), more, nil
		}
		chars = append(chars, scanner.Text())
	}
	return "", false, ErrLex{"Unterminated string"}
}

func lexNumber(scanner *bufio.Scanner) (string, bool, error) {
	chars := []string{scanner.Text()}
	for scanner.Scan() {
		c := scanner.Text()
		if !strings.Contains("0123456789.-", c) {
			return strings.Join(chars, ""), true, nil
		}
		chars = append(chars, c)
	}
	return strings.Join(chars, ""), false, nil
}

func scanKeyword(scanner *bufio.Scanner, keyword string) error {
	for _, c := range keyword {
		if !scanner.Scan() || scanner.Text() != string(c) {
			return ErrLex{"Expected keyword: " + keyword}
		}
	}
	return nil
}
