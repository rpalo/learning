package main

import (
	"fmt"
	"strings"
)

type ErrParse struct {
	reason   string
	expected string
	actual   string
}

func (e ErrParse) Error() string {
	return fmt.Sprintf("Error parsing: %s: %s, got %s", e.reason, e.expected, e.actual)
}

type JsonObject map[string]string

func (o JsonObject) String() string {
	result := []string{"{"}
	for key, value := range o {
		result = append(result, fmt.Sprintf("\t%s: %s,", key, value))
	}
	result = append(result, "}")
	return strings.Join(result, "\n")
}

func Parse(tokens []Token) (JsonObject, error) {
	result, tokens, err := parseObject(tokens)
	if err != nil {
		return nil, err
	}
	if len(tokens) > 0 {
		return nil, ErrParse{"Extra tokens after object.", "", tokens[0].value}
	}
	return result, nil
}

func parseObject(tokens []Token) (JsonObject, []Token, error) {
	if len(tokens) < 2 {
		return nil, nil, ErrParse{"Not enough tokens to parse object.", "", ""}
	}

	if _, err := expectRaw(tokens, "{"); err != nil {
		return nil, nil, err
	}

	result := make(map[string]string)

	if _, err := expectRaw(tokens[1:], "}"); err == nil {
		remaining, err := advance(tokens, 2)
		return result, remaining, err
	} else {
		tokens, _ = advance(tokens, 1)
	}

	for {
		if len(tokens) < 4 {
			return nil, nil, ErrParse{"Not enough tokens for key-value pair.", "", ""}
		}

		key, err := expectString(tokens)

		if err != nil {
			return nil, nil, err
		}

		if _, err := expectRaw(tokens[1:], ":"); err != nil {
			return nil, nil, err
		}

		value, err := expectString(tokens[2:])

		if err != nil {
			return nil, nil, err
		}
		result[key] = value

		if _, err := expectRaw(tokens[3:], "}"); err == nil {
			tokens = tokens[3:]
			break
		}

		if _, err := expectRaw(tokens[3:], ","); err != nil {
			return nil, nil, err
		}

		tokens, _ = advance(tokens, 4)
	}

	if _, err := expectRaw(tokens, "}"); err != nil {
		return nil, nil, err
	}
	remaining, err := advance(tokens, 1)
	return result, remaining, err
}

func advance(tokens []Token, count int) ([]Token, error) {
	if len(tokens) < count {
		return nil, ErrParse{"Not enough tokens to advance.", "", ""}
	}
	if len(tokens) == count {
		return nil, nil // No tokens left to advance
	}
	return tokens[count:], nil
}

func expectRaw(tokens []Token, value string) (string, error) {
	if len(tokens) == 0 {
		return "", ErrParse{"Missing expected token.", value, ""}
	}
	if tokens[0].kind != TokenRaw {
		return "", ErrParse{"Expected raw value.", value, tokens[0].value}
	}
	if tokens[0].value != value {
		return "", ErrParse{"Unexpected token.", value, tokens[0].value}
	}
	return tokens[0].value, nil
}

func expectString(tokens []Token) (string, error) {
	if len(tokens) == 0 {
		return "", ErrParse{"Missing expected string.", "string", ""}
	}
	if tokens[0].kind != TokenString {
		return "", ErrParse{"Expected string.", "string", tokens[0].value}
	}
	return tokens[0].value, nil
}
