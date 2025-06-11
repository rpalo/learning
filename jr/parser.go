package main

import (
	"fmt"
	"strconv"
	"strings"
)

type ErrParse struct {
	reason   string
	expected string
	actual   string
}

func (e ErrParse) Error() string {
	return fmt.Sprintf("Error parsing: %s: Wanted '%s', got '%s'", e.reason, e.expected, e.actual)
}

type JsonAny interface {
	String() string
}

type JsonObject map[JsonString]JsonAny

func (o JsonObject) String() string {
	result := make([]string, 0)
	for key, value := range o {
		result = append(result, fmt.Sprintf("%s: %s", key, value))
	}
	return "{" + strings.Join(result, ", ") + "}"
}

type JsonArray []JsonAny

func (a JsonArray) String() string {
	result := make([]string, len(a))
	for i, value := range a {
		result[i] = value.String()
	}
	return "[" + strings.Join(result, ", ") + "]"
}

type JsonString string

func (s JsonString) String() string {
	return string(s)
}

type JsonNumber float64

func (n JsonNumber) String() string {
	return strconv.FormatFloat(float64(n), 'f', -1, 64)
}

type JsonBool bool

func (b JsonBool) String() string {
	if b {
		return "true"
	}
	return "false"
}

type JsonNull struct{}

func (n JsonNull) String() string {
	return "null"
}

func Parse(tokens []Token) (JsonAny, error) {
	if len(tokens) == 0 {
		return nil, ErrParse{"No tokens to parse.", "", ""}
	}
	result, tokens, err := parseAny(tokens)
	if err != nil {
		return nil, err
	}
	if len(tokens) > 0 {
		return nil, ErrParse{"Extra tokens after parsing.", "", tokens[0].value}
	}
	return result, nil
}

func parseAny(tokens []Token) (JsonAny, []Token, error) {
	if len(tokens) == 0 {
		return nil, nil, ErrParse{"No tokens to parse.", "", ""}
	}
	switch tokens[0].kind {
	case TokenString:
		result, err := parseString(tokens)
		if err != nil {
			return nil, nil, err
		}
		tokens, err = advance(tokens, 1)
		return result, tokens, err
	case TokenNumber:
		result, err := parseNumber(tokens)
		if err != nil {
			return nil, nil, err
		}
		tokens, err = advance(tokens, 1)
		return result, tokens, err
	case TokenKeyword:
		switch tokens[0].value {
		case "true":
			tokens, err := advance(tokens, 1)
			return JsonBool(true), tokens, err
		case "false":
			tokens, err := advance(tokens, 1)
			return JsonBool(false), tokens, err
		case "null":
			tokens, err := advance(tokens, 1)
			return JsonNull{}, tokens, err
		default:
			return nil, nil, ErrParse{"Unknown keyword.", "", tokens[0].value}
		}
	case TokenRaw:
		switch tokens[0].value {
		case "{":
			return parseObject(tokens)
		case "[":
			return parseList(tokens)
		default:
			return nil, nil, ErrParse{"Unexpected raw token.", "", tokens[0].value}
		}
	default:
		return nil, nil, ErrParse{"Unexpected token kind.", "", tokenTypeNames[tokens[0].kind]}
	}
}

func parseObject(tokens []Token) (JsonObject, []Token, error) {
	if len(tokens) < 2 {
		return nil, nil, ErrParse{"Not enough tokens to parse object.", "", ""}
	}

	if _, err := expectRaw(tokens, "{"); err != nil {
		return nil, nil, err
	}

	result := make(map[JsonString]JsonAny)

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

		key, err := parseString(tokens)
		if err != nil {
			return nil, nil, err
		}

		if _, err := expectRaw(tokens[1:], ":"); err != nil {
			return nil, nil, err
		}

		var value JsonAny
		value, tokens, err = parseAny(tokens[2:])
		if err != nil {
			return nil, nil, err
		}
		result[key] = value
		if _, err := expectRaw(tokens, "}"); err == nil {
			break
		}

		if _, err := expectRaw(tokens, ","); err != nil {
			return nil, nil, err
		}

		tokens = tokens[1:]
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

func parseString(tokens []Token) (JsonString, error) {
	if len(tokens) == 0 {
		return "", ErrParse{"Missing expected string.", "string", ""}
	}
	if tokens[0].kind != TokenString {
		return "", ErrParse{"Expected string.", "string", tokens[0].value}
	}
	return JsonString(tokens[0].value), nil
}

func parseNumber(tokens []Token) (JsonNumber, error) {
	if len(tokens) == 0 {
		return 0, ErrParse{"Missing expected number.", "number", ""}
	}
	if tokens[0].kind != TokenNumber {
		return 0, ErrParse{"Expected number.", "number", tokens[0].value}
	}
	value, err := strconv.ParseFloat(tokens[0].value, 64)
	if err != nil {
		return 0, ErrParse{"Invalid number format.", "number", tokens[0].value}
	}
	return JsonNumber(value), nil
}

func parseList(tokens []Token) (JsonArray, []Token, error) {
	if len(tokens) < 2 {
		return nil, nil, ErrParse{"Not enough tokens to parse list.", "", ""}
	}

	if _, err := expectRaw(tokens, "["); err != nil {
		return nil, nil, err
	}

	result := make([]JsonAny, 0)

	if _, err := expectRaw(tokens[1:], "]"); err == nil {
		remaining, err := advance(tokens, 2)
		return result, remaining, err
	} else {
		tokens, _ = advance(tokens, 1)
	}

	for {
		value, remainingTokens, err := parseAny(tokens)
		if err != nil {
			return nil, nil, err
		}
		result = append(result, value)
		tokens = remainingTokens

		if len(tokens) == 0 {
			return nil, nil, ErrParse{"Unexpected end of tokens while parsing list.", "", ""}
		}
		if tokens[0].value == "]" {
			break
		}

		if _, err := expectRaw(tokens, ","); err != nil {
			return nil, nil, err
		}
		tokens = tokens[1:]
	}

	if _, err := expectRaw(tokens, "]"); err != nil {
		return nil, nil, err
	}
	remaining, err := advance(tokens, 1)
	return result, remaining, err
}
