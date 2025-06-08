package main

import (
	"errors"
	"fmt"
)

var UnmatchedBracket = errors.New("Syntax error: unmatched square bracket")

type RJump struct {
	target int
}

type LJump struct {
	target int
}

type Add struct {
	amount int
}

type Right struct{}
type Left struct{}
type Input struct{}
type Output struct{}

type Opcode interface{}

func Compile(source string) ([]Opcode, error) {
	result := make([]Opcode, 0, len(source))

	for _, c := range source {
		switch c {
		case '+':
			result = append(result, &Add{1})
		case '-':
			result = append(result, &Add{-1})
		case '>':
			result = append(result, &Right{})
		case '<':
			result = append(result, &Left{})
		case ',':
			result = append(result, &Input{})
		case '.':
			result = append(result, &Output{})
		case '[':
			result = append(result, &RJump{-1})
		case ']':
			result = append(result, &LJump{-1})
		}
	}
	err := matchLoops(result)
	return result, err
}

func matchLoops(ops []Opcode) error {
	for i, op := range ops {
		if rjump, ok := op.(*RJump); ok {
			target, err := findMatchingLJump(ops, i)

			if err != nil {
				panic(i)
				return err
			}
			rjump.target = target
			ljump, _ := ops[target].(*LJump)
			ljump.target = i
		} else if ljump, ok := op.(*LJump); ok && ljump.target == -1 {
			// We should find all ljumps by running into their rjumps first
			return UnmatchedBracket
		}
	}
	return nil
}

func findMatchingLJump(ops []Opcode, start int) (int, error) {
	loopCounter := 0

	for i := start + 1; i < len(ops); i++ {
		if _, ok := ops[i].(*RJump); ok {
			loopCounter++
		} else if _, ok := ops[i].(*LJump); ok && loopCounter != 0 {
			loopCounter--
		} else if _, ok := ops[i].(*LJump); ok {
			return i, nil
		}
	}
	return 0, UnmatchedBracket
}

func PrintOps(ops []Opcode) {
	for i, op := range ops {
		fmt.Printf("%05d:\t%T%v\n", i, op, op)
	}
}
