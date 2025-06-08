# BF (Brainf*ck) Toolkit

[https://codingchallenges.fyi/challenges/challenge-brainfuck]

This is a Brainf*ck interpreter written in Go, with some other learning
utilities added in.

## Usage

```shell
# Compile to bytecode and output the bytecode
bf compile example.bf

# Compile to bytecode and run the bytecode
bf run example.bf

# Run an interactive repl (sort of) (this fell by the wayside, needs attention)
bf repl
```

Additionally, the following env vars can be set to modify the execution:

- `BF_BUFFER_SIZE`: Modifies the size of the memory array for bf
- `BF_DEBUG`: Outputs more info about operation including what opcode each step
  plus buffer values, etc.
- `BF_NUMBERS`: If set, output memory will be output as numbers instead of their
  char code (useful for debugging)
- `BF_LOOPCHECK`: After running a program, will output each encountered loop
  sorted by number of iterations run, as a way of tracking down possibly useful
  optimizations

## Design

The compiler is set up to operate in the following steps:

1. Strip out comment characters (i.e. non-code chars).
2. Perform a simpler, string-pattern-based optimization pass, converting
   defined, frequently encountered patterns into other operation characters.
3. Iterate through the modified source-code, converting characters (including
   the new expanded chars) into opcodes. a. Condense repeated chars for some
   opcodes e.g. `+++++` becomes `Add{5}`.
4. Link up matching loop brackets for instant jumps (and error check).
5. Perform another optimization pass of trickier optimizations that require
   opcode form to find the patterns. These are defined as their own functions.
   a. Doing these optimizations botches the loop linking, so we need to re-link
   the loops after. b. In theory, we could get to a point where we need to fix
   the loops after every optimization (or use pointers), but for now, the
   optimizations don't overlap, so we can just fix them all once we're done with
   this pass.

After this, the interpreter runs through the ops in a pretty naive way, as you
would expect. We need to ensure that any new opcodes created as optimizations
get handled in the interpreter too. It could possibly be a fancier VM struct
object or something, but there aren't that many state variables, so using a
single function with local stack state variables seems like it makes the most
sense.

## Speed

Throughout this build, one of the driving goals was to reduce the speed of
operation where possible. When I first started with my initial pass, evaluating
from the source string direction (that vestigial code still lives in the
interpreter for reference), was slow, with the mandelbrot example taking around
40s to complete. After migrating to opcodes and a VM and linking loops, we
picked up 5 seconds or so. Condensing runs of the same op was a huge speedup, at
least 10s. Subsequent optimizations like the `Clear`, `FindEmpty`, and
`Transfer` ops got us down finally to about 6s, around 15% of the original
runtime.

Running the loopcheck shows that there are more idioms present that we could
optimize for, but I'm sleepy, so I'm going to leave it alone for now.

## Todo

I probably won't get to these, but I'm at least acknowledging that the tasks
exist:

- Fix up REPL, and keep the buffer state between commands.
- Tests, currently I just tested manually on known bf files with lots of
  debugging, but formal unit tests would probably be wise.
- Further compilation i.e. actual compiling of bf files to assembly/binaries.

## Contributions/Comments

This isn't really _for_ anything, but I'm always happy to look at comments/ideas
to improve, especially as I learn and practive my Golang skills more.
