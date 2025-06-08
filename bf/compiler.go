package main

// compiler.go contains all the functions for compiling bf code to opcodes
// and optimizing those opcodes for efficiency

import (
	"errors"
	"regexp"
	"strings"
)

var UnmatchedBracket = errors.New("Syntax error: unmatched square bracket")

// RJump tells the VM to jump "right" to the matching ']' if the current buffer
// value is 0.
type RJump struct {
	target int
}

// LJump tells the VM to jump "left" to the matching '[' if the current buffer
// value is not 0.
type LJump struct {
	target int
}

// Add increments the current buffer value by some amount.
type Add struct {
	amount int
}

// Move moves the buffer pointer some amount left or right.
type Move struct {
	amount int
}

// Transfer shifts all of the value from one buffer slot to another slot some
// distance away.
type Transfer struct {
	distance int
}

// Find empty skips forward some number of steps repeatedly until it finds an
// empty buffer slot
type FindEmpty struct {
	step int
}

// Input causes the interpreter to read a character of input from stdin into
// the current buffer slot
type Input struct{}

// Output causes the interpreter to write a character of output from the
// current buffer slot (using its ascii character value).
type Output struct{}

// Clear sets the current buffer slot to zero, and moves one buffer slot right
// if step is true.
type Clear struct {
	step bool
}
type Opcode any

// Compile compiles bf source to Opcodes, and optimizes them.
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

// stripComments removes any characters that are not canonical operation chars.
func stripComments(source string) string {
	pattern := regexp.MustCompile(`[^\+\-\,\.\[\]\<\>]`)
	return pattern.ReplaceAllLiteralString(source, "")
}

// consolidateRun counts how many of the same character are in a row starting
// at `start`
func consolidateRun(source string, start int) int {
	count := 1
	for i := start + 1; i < len(source) && source[i] == source[i-1]; i++ {
		count++
	}
	return count
}

// replaceOptimizations performs simple string replacement optimizations.
func replaceOptimizations(source string) string {
	source = strings.ReplaceAll(source, "[-]>", "X") // Clear cell and step
	source = strings.ReplaceAll(source, "[-]", "x")  // Clear cell
	return source
}

// matchLoops attempts to set the jump opcodes' `target` fields to point to
// their matching jumps and returns an error if there are unmatched brackets.
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

// findMatchingLJump finds the index of the matching jump op.
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

// optimize performs post-compilation optimizations on opcodes.  It returns
// a new slice of opcodes.  Note: currently these new opcodes' jump targets
// need re-matched again.
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

// optimizeTransfer finds the "transfer" idiom and replaces it with a transfer
// opcode.
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

// optimizeFindEmpty finds the "find empty" idiom and replaces it with a findempty
// opcode.
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
