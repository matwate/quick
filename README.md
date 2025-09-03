# Quick

A command-line tool to simplify the verification of LeetCode-style coding problems by automating the input/output process.

## Installation

To make the `quick` command available system-wide, you can use `go install`:

```bash
go install
```

This will build the binary and place it in your `$GOPATH/bin` directory. If this directory is in your system's `PATH`, you'll be able to run `quick` from anywhere.

## Usage

This tool is split into two main commands: `run` and `create`.

### Running Tests (`run`)

To run your Python solution against a set of test cases, use the `run` command.

```bash
quick run <your_solution.py>
```

By default, the tool will scan the directory for all `*.in` files and their corresponding `*.out` files, and run your script against each pair.

#### Strict Mode

For more targeted testing, especially when you have multiple problems in the same directory, use the `--strict` or `-s` flag.

```bash
quick run --strict <your_solution.py>
```

In strict mode, the tool will only run on test files that are associated with your solution's name. For a script named `solution.py`, it will look for test files like `1solution.in`, `1solution.out`, `2solution.in`, `2solution.out`, etc.

### Creating Test Files (`create`)

To quickly generate empty `.in` and `.out` file pairs, use the `create` command.

**Usage:** `quick create <number_of_pairs> [problem_name]`

-   `<number_of_pairs>`: The number of pairs to create.
-   `[problem_name]`: An optional name to associate the files with a specific problem.

**Examples:**

To create 3 generic test file pairs:
```bash
quick create 3
```
This will create `1.in`, `1.out`, `2.in`, `2.out`, `3.in`, `3.out`.

To create 2 test file pairs for a problem named `solution`:
```bash
quick create 2 solution
```
This will create `1solution.in`, `1solution.out`, `2solution.in`, `2solution.out`. These files will be automatically picked up when you run `quick run --strict solution.py`.