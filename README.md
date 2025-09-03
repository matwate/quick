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

## Configuration

You can configure the shell and the run command used by `quick` to execute your scripts.

### Setting Configuration Values

To set a configuration value, use the `config set` command.

```bash
quick config set <key> <value>
```

For example, to change the shell to `zsh`:

```bash
quick config set shell zsh
```

### Windows Configuration

For Windows, you need to configure the shell to `cmd.exe` (or `powershell.exe`) and the run command to `python`.

```bash
quick config set shell cmd.exe
quick config set run_command python
```

You can also use `powershell.exe` as your shell:

```bash
quick config set shell powershell.exe
```

## Potential Upgrades

This tool provides a solid foundation, but there are many ways it could be extended. Here are a few ideas for future development:

*   **Support for More Languages:** The tool is currently hardcoded for Python 3. A great upgrade would be to add a `-lang` flag to specify other languages (e.g., C++, Java, Node.js) and their corresponding execution commands.

*   **Interactive Test Case Creation:** The `create` command could be enhanced with an interactive mode. For example, `quick create -i` could prompt the user to paste the input and expected output directly in the terminal, saving it to the files.

*   **Configurable Limits:** The 5-second timeout is currently fixed. This could be made configurable with a `--timeout` flag. Support for memory limit constraints could also be added.

*   **Parallel Test Execution:** To speed up the verification process, the test cases could be run in parallel instead of sequentially.
