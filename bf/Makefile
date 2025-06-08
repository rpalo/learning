.PHONY: benchmark

sources := $(wildcard *.go)

bf: $(sources)
	go build

benchmark: bf
	time ./bf examples/mandelbrot.bf > /dev/null 2>&1
