package main

import "fmt"

func PrintOps(ops []Opcode) {
	for i, op := range ops {
		fmt.Printf("%05d:\t%T%v\n", i, op, op)
	}
}

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
