package main

import "fmt"

type ErrParse struct {
	reason string
	token  string
}

func (e ErrParse) Error() string {
	return fmt.Sprintf("Error parsing: %s", e.token)
}

type JsonObject struct{}

func (o *JsonObject) String() string {
	return "{}"
}

func Parse(tokens []string) (*JsonObject, error) {
	return parseObject(tokens)
}

func parseObject(tokens []string) (*JsonObject, error) {
	if len(tokens) == 0 {
		return nil, ErrParse{"Empty file.", ""}
	}
	if tokens[0] != "{" {
		return nil, ErrParse{"Unexpected token", tokens[0]}
	}
	if len(tokens) == 1 {
		return nil, ErrParse{"Unclosed object.", ""}
	}
	if tokens[1] != "}" {
		return nil, ErrParse{"Unexpected token", tokens[1]}
	}
	if len(tokens) > 2 {
		return nil, ErrParse{"Extra token", tokens[2]}
	}
	return &JsonObject{}, nil
}
