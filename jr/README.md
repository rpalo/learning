# JR (JSON Parser)

Get it?  Like, Ryan's JQ, but also `jr` because it's just, like, a little bit better than `jq`.  Just kidding.

## Usage

```shell
# Parse the input JSON file and print it back out
jr example.json
```

Additionally, the following env vars can be set to modify the execution:

- `JR_MAX_DEPTH`: Sets the maximum nesting depth allowed before erroring

## Running the Tests

There are two test suites:

1. The test suite that I put together as sanity checks during development, found in examples.
2. The test suite provided by JSON.org.

You can run both by building `jr` and then running the test script:

```shell
$ go build
$ ./simple_test
..................................................
Done.  All tests passed.
```

## Todo

As-is this is really more of a JSON validator.  I could probably restructure the project to make it an importable package.

## Contributions/Comments

This isn't really _for_ anything, but I'm always happy to look at comments/ideas
to improve, especially as I learn and practive my Golang skills more.
