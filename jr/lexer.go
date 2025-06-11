package main

// lexer.go contains the lexer implementation, separating the input string into tokens.

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

// ErrLex represents an error encountered during lexing.
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

// Token is a single individual unit of the input text.
type Token struct {
	kind  TokenType
	value string
}

// Lex returns a slice of tokens parsed from the input reader.
func Lex(reader io.Reader) ([]Token, error) {
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanRunes)

	result := make([]Token, 0)
	scanner.Scan()
	var more bool = true
	var err error

	for more {
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
		case "0", "-", ".", "1", "2", "3", "4", "5", "6", "7", "8", "9":
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

var UnicodeEscapePattern = regexp.MustCompile("[0-9a-fA-F]")

// lexString reads a string from the scanner, handling escape sequences and unicode escapes.
//
// It returns the string, a boolean indicating if there are more characters to read, and an error if any occurred.
// Some edge cases required by the JSON specification:
// - Raw newlines and tabs are not allowed and must be escaped.
// - Unicode escape sequences must be 4 hex digits that evaluate to a valid unicode character.
// - The only allowed escape sequences are \", \\, \/, \b, \f, \n, \r, \t, and \uXXXX where XXXX is a 4-digit hexadecimal number.
func lexString(scanner *bufio.Scanner) (string, bool, error) {
	chars := make([]string, 0)

	for scanner.Scan() {
		if scanner.Text() == "\\" {
			// Handle escape sequences
			if !scanner.Scan() {
				return "", false, ErrLex{"Unterminated string with escape sequence"}
			}
			if !strings.Contains("\"\\/bfnrtu", scanner.Text()) {
				return "", false, ErrLex{"Invalid escape sequence: \\" + scanner.Text()}
			}
			if scanner.Text() == "u" {
				c, err := lexUnicodeEscape(scanner)
				if err != nil {
					return "", false, err
				}
				chars = append(chars, c)
				continue
			}
			chars = append(chars, scanner.Text())
			continue
		}
		if strings.Contains("\t\n", scanner.Text()) {
			return "", false, ErrLex{"Character not allowed in string.  Use escape sequences for newlines and tabs."}
		}
		if scanner.Text() == "\"" {
			more := scanner.Scan() // consume the closing quote
			return strings.Join(chars, ""), more, nil
		}
		chars = append(chars, scanner.Text())
	}
	return "", false, ErrLex{"Unterminated string"}
}

// lexUnicodeEscape reads and expects 4 hexadecimal digits that form a valid unicode escape sequence, and either returns the
// successful decoded character or an error.
func lexUnicodeEscape(scanner *bufio.Scanner) (string, error) {
	chars := make([]string, 4)
	for i := 0; i < 4; i++ {
		if !scanner.Scan() {
			return "", ErrLex{"Incomplete unicode escape sequence"}
		}
		c := scanner.Text()
		if !UnicodeEscapePattern.MatchString(c) {
			return "", ErrLex{"Invalid unicode escape sequence: " + c}
		}
		chars[i] = c
	}
	result, err := strconv.Unquote(`'\u` + strings.Join(chars, "") + `'`)
	if err != nil {
		return "", ErrLex{"Invalid unicode escape sequence: " + strings.Join(chars, "") + err.Error()}
	}
	return result, nil
}

// lexNumber reads a number from the scanner, handling integers, floats, and scientific notation.
//
// Leading zeros area not allowed on multi-digit integers.
func lexNumber(scanner *bufio.Scanner) (string, bool, error) {
	chars := []string{scanner.Text()}
	for scanner.Scan() {
		c := scanner.Text()
		if !strings.Contains("0123456789.-+eE", c) {
			num := strings.Join(chars, "")
			err := checkInt(num)
			if err != nil {
				return "", false, err
			}

			return num, true, nil
		}
		chars = append(chars, c)
	}
	num := strings.Join(chars, "")
	err := checkInt(num)
	if err != nil {
		return "", false, err
	}

	return num, true, nil
}

// checkInt checks if the given string represents a valid integer.  Leading zeros are not allowed on multi-digit integers.
func checkInt(value string) error {
	if value[0] == '0' && !strings.Contains(value, ".") && len(value) > 1 {
		return ErrLex{"Invalid number: leading zero in " + value}
	}
	return nil
}

// scans the input for a specific keyword, which is a sequence of characters that must be matched exactly.  It throws an error if a match is not found.
func scanKeyword(scanner *bufio.Scanner, keyword string) error {
	for _, c := range keyword {
		if !scanner.Scan() || scanner.Text() != string(c) {
			return ErrLex{"Expected keyword: " + keyword}
		}
	}
	return nil
}
