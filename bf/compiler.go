package main

import "fmt"

type Opcode uint8

const (
	Inc Opcode = iota
	Dec
	Right
	Left
	RJump
	LJump
	Input
	Output
)

func Compile(source string) []Opcode {
	result := make([]Opcode, 0, len(source))

	mapping := map[rune]Opcode{
		'+': Inc,
		'-': Dec,
		'>': Right,
		'<': Left,
		'[': RJump,
		']': LJump,
		',': Input,
		'.': Output,
	}

	for _, c := range source {
		if op, ok := mapping[c]; ok {
			result = append(result, op)
		}
	}
	return result
}

func PrintOps(ops []Opcode) {

	mapping := map[Opcode]string{
		Inc:    "INC",
		Dec:    "DEC",
		Right:  "RIGHT",
		Left:   "LEFT",
		RJump:  "RJUMP",
		LJump:  "LJUMP",
		Input:  "INPUT",
		Output: "OUTPUT",
	}

	for i, op := range ops {
		fmt.Printf("%05d:\t%s\n", i, mapping[op])
	}
}
