package main

import (
	"fmt"
	"log"
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
