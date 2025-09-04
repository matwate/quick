<div align="center">

```
██████╗ ██╗   ██╗██╗ ██████╗██╗  ██╗
██╔══██╗██║   ██║██║██╔════╝██║ ██╔╝
██████╔╝██║   ██║██║██║     █████╔╝
██╔══██╗██║   ██║██║██║     ██╔═██╗
██████╔╝╚██████╔╝██║╚██████╗██║  ██╗
╚═════╝  ╚═════╝ ╚═╝ ╚═════╝╚═╝  ╚═╝
```

**A command-line tool to simplify the testing of competitive programming problems.**

[![Go Version](https://img.shields.io/badge/Go-1.18+-blue.svg)](https://golang.org/dl/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

</div>

---

`quick` is a command-line tool designed to streamline the process of testing solutions for competitive programming problems. It automates running your code against multiple input/output files, saving you time and effort.

## 🚀 Installation

To get started, make sure you have Go installed on your system. Then, you can install `quick` with a single command:

```bash
go install
```

This will build the binary and place it in your `$GOPATH/bin` directory. If this directory is in your system's `PATH`, you'll be able to run `quick` from anywhere.

## 💡 Usage

`quick` is built around two main commands: `run` and `create`.

### ✅ Running Tests (`run`)

Use the `run` command to execute your solution against a set of test cases.

```bash
quick run <your_solution.py>
```

By default, `quick` scans the current directory for all `*.in` files and their corresponding `*.out` files, and runs your script against each pair.

#### Strict Mode

For more targeted testing, especially when you have multiple problems in the same directory, use the `--strict` or `-s` flag.

```bash
quick run --strict <your_solution.py>
```

In strict mode, the tool will only run on test files that are associated with your solution's name. For a script named `solution.py`, it will look for test files like `1solution.in`, `1solution.out`, `2solution.in`, `2solution.out`, etc.

### ✨ Creating Test Files (`create`)

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

## ⚙️ Configuration

You can customize the shell and the run command used by `quick` to execute your scripts.

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

## 🗺️ Roadmap

`quick` is actively being developed, and here are some of the features we're planning to add:

-   **🌐 Support for More Languages:** Add a `-lang` flag to specify other languages (e.g., C++, Java, Node.js) and their corresponding execution commands.
-   **✍️ Interactive Test Case Creation:** An interactive mode for the `create` command (`quick create -i`) that prompts the user to paste the input and expected output directly in the terminal.
-   **⏱️ Configurable Limits:** Allow users to configure the timeout and memory limits for test case execution.
-   **⚡ Parallel Test Execution:** Run test cases in parallel to speed up the verification process.

## 🙌 Contributing

We welcome contributions from the community! If you have an idea for a new feature or want to help improve `quick`, please check out our [Contributing Guide](CONTRIBUTING.md) (you'll need to create this file).

## 📄 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details (you'll need to create this file).