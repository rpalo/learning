package main

import (
	"errors"
	"regexp"
	"strings"
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

type Move struct {
	amount int
}
type Transfer struct {
	distance int
}
type FindEmpty struct {
	step int
}

type Input struct{}
type Output struct{}
type Clear struct {
	step bool
}
type Opcode any

func Compile(source string) ([]Opcode, error) {
	source = stripComments(source)
	source = replaceOptimizations(source)
	ops := make([]Opcode, 0, len(source))

	for i := 0; i < len(source); i++ {
		switch source[i] {
		case '+':
			count := consolidateRun(source, i)
			ops = append(ops, &Add{count})
			i += count - 1
		case '-':
			count := consolidateRun(source, i)
			ops = append(ops, &Add{-1 * count})
			i += count - 1
		case '>':
			count := consolidateRun(source, i)
			ops = append(ops, &Move{count})
			i += count - 1
		case '<':
			count := consolidateRun(source, i)
			ops = append(ops, &Move{-1 * count})
			i += count - 1
		case ',':
			ops = append(ops, &Input{})
		case '.':
			ops = append(ops, &Output{})
		case '[':
			ops = append(ops, &RJump{-1})
		case ']':
			ops = append(ops, &LJump{-1})
		case 'x':
			ops = append(ops, &Clear{false})
		case 'X':
			ops = append(ops, &Clear{true})
		}
	}
	err := matchLoops(ops)

	if err != nil {
		return nil, err
	}

	result := optimize(ops)
	err = matchLoops(result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func stripComments(source string) string {
	pattern := regexp.MustCompile(`[^\+\-\,\.\[\]\<\>]`)
	return pattern.ReplaceAllLiteralString(source, "")
}

func consolidateRun(source string, start int) int {
	count := 1
	for i := start + 1; i < len(source) && source[i] == source[i-1]; i++ {
		count++
	}
	return count
}

func replaceOptimizations(source string) string {
	source = strings.ReplaceAll(source, "[-]>", "X") // Clear cell and step
	source = strings.ReplaceAll(source, "[-]", "x")  // Clear cell
	return source
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

func optimize(ops []Opcode) []Opcode {
	result := make([]Opcode, 0, len(ops))

	for i := 0; i < len(ops); i++ {
		switch v := ops[i].(type) {
		case *RJump:
			if transfer := optimizeTransfer(ops, i); transfer != nil {
				result = append(result, transfer)
				i += 5
			} else if find := optimizeFindEmpty(ops, i); find != nil {
				result = append(result, find)
				i += 2
			} else {
				result = append(result, v)
			}
		default:
			result = append(result, v)
		}
	}
	return result
}

func optimizeTransfer(ops []Opcode, i int) *Transfer {
	rjump, _ := ops[i].(*RJump)
	if rjump.target != i+5 {
		return nil
	}
	sub, subOk := ops[i+1].(*Add)
	move, moveOk := ops[i+2].(*Move)
	add, addOk := ops[i+3].(*Add)
	back, backOk := ops[i+4].(*Move)

	if !subOk || !moveOk || !addOk || !backOk {
		return nil
	}

	if sub.amount == -1*add.amount && move.amount == -1*back.amount {
		return &Transfer{move.amount}
	}
	return nil
}

func optimizeFindEmpty(ops []Opcode, i int) *FindEmpty {
	rjump, _ := ops[i].(*RJump)
	if rjump.target != i+2 {
		return nil
	}

	move, ok := ops[i+1].(*Move)

	if !ok {
		return nil
	}

	return &FindEmpty{move.amount}
}
