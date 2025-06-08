package main

// output.go contains helpers for outputting debugging info

import (
	"fmt"
	"sort"
)

// PrintOps prints opcodes in a basic way.  Sort of a dissassembler for bf syntax.
func PrintOps(ops []Opcode) {
	for i, op := range ops {
		fmt.Printf("%05d:\t%T%v\n", i, op, op)
	}
}

// PrintOpsCompact converts opcodes back into a processed almost-bf syntax for
// quick checks.
func PrintOpsCompact(ops []Opcode) {
	for _, op := range ops {
		switch v := op.(type) {
		case *Add:
			fmt.Printf("%d%c", v.amount, '+')
		case *Move:
			fmt.Printf("%d%c", v.amount, '>')
		case *Input:
			fmt.Print(",")
		case *Output:
			fmt.Print(".")
		case *RJump:
			fmt.Print("[")
		case *LJump:
			fmt.Print("]")
		case *Clear:
			if v.step {
				fmt.Print("X")
			} else {
				fmt.Print("x")
			}
		case *Transfer:
			fmt.Printf("%dT", v.distance)
		case *FindEmpty:
			fmt.Printf("%dF", v.step)
		default:
			panic(fmt.Sprintf("Unrecognized op %T\n", op))
		}
	}
}

// KV is used to sort loop count maps
type KV struct {
	key   int
	value int
}

// PrintLoops prints the encountered loops in order of increasing number of
// iterations run.
func PrintLoops(ops []Opcode, loops map[int]int) {
	items := make([]KV, 0, len(loops))
	for i, count := range loops {
		items = append(items, KV{i, count})
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].value < items[j].value
	})
	for _, pair := range items {
		start, _ := ops[pair.key].(*RJump)
		fmt.Printf("%d: ", pair.value)
		PrintOpsCompact(ops[pair.key : start.target+1])
		fmt.Print("\n")
	}
}
