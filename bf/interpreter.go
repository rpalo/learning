package main

import (
	"fmt"
	"log"
	"sort"
)

func EvalBf(source string, buffer_size int, debug bool) {
	i := 0
	d := 0
	loopCounter := 0
	buffer := make([]int, buffer_size)

	for i >= 0 && i < len(source) {
		if debug {
			fmt.Printf("%d: %c, %d: [%d]\n", i, source[i], d, buffer[d])
		}
		switch source[i] {
		case '>':
			d = (d + 1) % buffer_size
		case '<':
			d = (d - 1 + buffer_size) % buffer_size
		case '+':
			buffer[d]++
		case '-':
			buffer[d]--
		case '.':
			fmt.Printf("%c", buffer[d])
		case ',':
			_, err := fmt.Scanf("%c", &buffer[d])

			if err != nil {
				log.Fatal(err)
			}
		case '[':
			if buffer[d] == 0 {
				for i++; source[i] != ']' || loopCounter != 0; i++ {
					if source[i] == '[' {
						loopCounter++
					} else if source[i] == ']' {
						loopCounter--
					}
				}
			}
		case ']':
			if buffer[d] != 0 {
				for i--; source[i] != '[' || loopCounter != 0; i-- {
					if source[i] == '[' {
						loopCounter--
					} else if source[i] == ']' {
						loopCounter++
					}
				}
			}
		}
		i++
	}
}

func EvalBfOps(ops []Opcode) {
	i := 0
	d := 0
	buffer := make([]int, buffer_size)
	loopCount := make(map[int]int)

	for i >= 0 && i < len(ops) {
		if debug {
			fmt.Printf("%05d: %T%v, %d: [%d]\n", i, ops[i], ops[i], d, buffer[d])
		}
		switch v := ops[i].(type) {
		case *Move:
			d = (d + v.amount + buffer_size) % buffer_size
		case *Add:
			buffer[d] += v.amount
		case *Output:
			fmt.Printf(outputPattern, buffer[d])
		case *Input:
			_, err := fmt.Scanf("%c", &buffer[d])

			if err != nil {
				log.Fatal(err)
			}
		case *RJump:
			loopCount[i] += 1
			if buffer[d] == 0 {
				i = v.target
			}
		case *LJump:
			if buffer[d] != 0 {
				i = v.target
			}
		case *Clear:
			buffer[d] = 0
			if v.step {
				d++
			}
		case *Transfer:
			newInd := (d + v.distance + buffer_size) % buffer_size
			buffer[newInd] += buffer[d]
			buffer[d] = 0
		case *FindEmpty:
			for buffer[d] != 0 {
				d = (d + v.step + buffer_size) % buffer_size
			}
		default:
			panic(fmt.Sprintf("Unrecognized opcode %T\n", ops[i]))
		}
		i++
	}
	if loopcheck {
		printLoops(ops, loopCount)
	}
}

type KV struct {
	key   int
	value int
}

func printLoops(ops []Opcode, loops map[int]int) {
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
