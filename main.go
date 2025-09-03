package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/pmezard/go-difflib/difflib"
	"github.com/urfave/cli/v3"
)

const (

	// Emojis
	EmojiPassed  = "✅"
	EmojiFailed  = "❌"
	EmojiTimeout = "⌛"
)

func run(c *cli.Command) {
	pyFilename := c.Args().Get(0)
	strict := c.Bool("strict")

	Ins_Outs := make(map[string]string)

	if strict {
		// Strict mode: look for files matching *{py_filename}.in and *{py_filename}.out
		baseFilename := strings.TrimSuffix(pyFilename, ".py")
		inputs, err := os.ReadDir(".")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Found matching input files (strict mode):")
		for _, file := range inputs {
			if !file.IsDir() && strings.Contains(file.Name(), baseFilename) &&
				strings.HasSuffix(file.Name(), ".in") {
				inFilename := file.Name()
				outFilename := strings.TrimSuffix(inFilename, ".in") + ".out"

				// Check if the corresponding .out file exists
				if _, err := os.Stat(outFilename); err == nil {
					fmt.Println(" -", inFilename)
					Ins_Outs[inFilename] = outFilename
				}
			}
		}
	} else {
		// Scan *.in files and *.out files
		inputs, err := os.ReadDir(".")
		ins, outs := make([]string, 0), make([]string, 0)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Found input files:")
		for _, file := range inputs {
			if !file.IsDir() && len(file.Name()) > 3 && file.Name()[len(file.Name())-3:] == ".in" {
				fmt.Println(" -", file.Name())
				ins = append(ins, file.Name())
			}
		}
		fmt.Println("Found output files:")
		for _, file := range inputs {
			if !file.IsDir() && len(file.Name()) > 4 && file.Name()[len(file.Name())-4:] == ".out" {
				fmt.Println(" -", file.Name())
				outs = append(outs, file.Name())
			}
		}

		// Match each *.in file with its corresponding *.out file
		for _, in := range ins {
			out := in[:len(in)-3] + ".out"

			for _, o := range outs {
				if o == out {
					Ins_Outs[in] = o
					break
				}
			}
		}
	}
	// Run a python subprocess with the input specified in the *.in file and capture the output

	passed := true
	failed := []string{}
	for in, out := range Ins_Outs {
		fmt.Printf("Running test case: %s -> %s\n ", in, out)

		// Read input from the *.in file
		inputData, err := os.ReadFile(in)
		if err != nil {
			log.Fatalf("Failed to read input file %s: %v", in, err)
		}

		// Execute the Python script with the input inputData
		cmd := fmt.Sprintf("python3 %s", pyFilename)

		// Create a context with a timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		command := exec.CommandContext(ctx, "zsh", "-c", cmd)
		command.Stdin = bytes.NewReader(inputData)

		var outBuffer, errBuffer bytes.Buffer
		command.Stdout = &outBuffer
		command.Stderr = &errBuffer

		startTime := time.Now()
		err = command.Run()
		duration := time.Since(startTime)

		if ctx.Err() == context.DeadlineExceeded {
			color.RGB(0xf9, 0xe2, 0xaf).
				Printf("%s Test case %s timed out! (%s)\n", EmojiTimeout, in, duration)
			continue
		}

		if err != nil {
			log.Fatalf("Failed to execute command: %v Stderr: %s", err, errBuffer.String())
		}

		// Read expected output from the *.out file
		expectedOutput, err := os.ReadFile(out)
		if err != nil {
			log.Fatalf("Failed to read output file %s: %v", out, err)
		}

		// Compare the actual output with the expected output
		actualOutput := outBuffer.Bytes()
		if bytes.Equal(bytes.TrimSpace(actualOutput), bytes.TrimSpace(expectedOutput)) {
			color.RGB(0xa6, 0xe3, 0xa1).
				Printf("%s Test case %s passed! (%s)\n", EmojiPassed, in, duration)
			passed = passed && true
		} else {
			color.RGB(0xf3, 0x8b, 0xa8).Printf("%s Test case %s failed! (%s)", EmojiFailed, in, duration)

			diff := difflib.UnifiedDiff{
				A:        difflib.SplitLines(string(expectedOutput)),
				B:        difflib.SplitLines(string(actualOutput)),
				FromFile: "Expected",
				ToFile:   "Got",
				Context:  3,
			}
			diffStr, _ := difflib.GetUnifiedDiffString(diff)
			fmt.Println(colorizeDiff(diffStr))
			passed = false
			failed = append(failed, in)
		}

	}
	if passed {
		color.RGB(0xa6, 0xe3, 0xa1).Println("All test cases passed! 🎉")
	} else {
		for _, f := range failed {
			color.RGB(0xf3, 0x8b, 0xa8).Printf(" - %s\n", f)
		}
		color.RGB(0xf3, 0x8b, 0xa8).Printf("%d test cases failed.\n", len(failed))
	}
}

func colorizeDiff(diff string) string {
	var coloredDiff strings.Builder
	for _, line := range strings.Split(diff, "\n") {
		switch {
		case strings.HasPrefix(line, "+"):
			coloredDiff.WriteString(color.RGB(0xa6, 0xe3, 0xa1).Sprint(line) + "\n")
		case strings.HasPrefix(line, "-"):
			coloredDiff.WriteString(color.RGB(0xf3, 0x8b, 0xa8).Sprint(line) + "\n")
		case strings.HasPrefix(line, "@@"):
			coloredDiff.WriteString(color.RGB(0xf9, 0xe2, 0xaf).Sprint(line) + "\n")
		default:
			coloredDiff.WriteString(line + "\n")
		}
	}
	return coloredDiff.String()
}

func createFiles(c *cli.Command) error {
	numFilesStr := c.Args().Get(0)
	if numFilesStr == "" {
		log.Println("Error: Number of files is required.")
		return fmt.Errorf("number of files is required")
	}
	numFiles, err := strconv.Atoi(numFilesStr)
	if err != nil {
		log.Printf("Error: Invalid number provided: %s", numFilesStr)
		return err
	}

	problemName := c.Args().Get(1)

	for i := 1; i <= numFiles; i++ {
		var baseFilename string
		if problemName == "" {
			baseFilename = fmt.Sprintf("%d", i)
		} else {
			problemName = strings.TrimSuffix(problemName, ".py")
			baseFilename = fmt.Sprintf("%d%s", i, problemName)
		}
		inFilename := baseFilename + ".in"
		outFilename := baseFilename + ".out"

		// Handle .in file
		if _, err := os.Stat(inFilename); os.IsNotExist(err) {
			if err := os.WriteFile(inFilename, []byte{}, 0o644); err != nil {
				color.New(color.FgRed).Printf("❌ Failed to create %s: %v\n", inFilename, err)
			} else {
				color.New(color.FgGreen).Printf("✅ Created %s\n", inFilename)
			}
		} else {
			color.New(color.FgYellow).Printf("⚠️ File %s already exists.\n", inFilename)
		}

		// Handle .out file
		if _, err := os.Stat(outFilename); os.IsNotExist(err) {
			if err := os.WriteFile(outFilename, []byte{}, 0o644); err != nil {
				color.New(color.FgRed).Printf("❌ Failed to create %s: %v\n", outFilename, err)
			} else {
				color.New(color.FgGreen).Printf("✅ Created %s\n", outFilename)
			}
		} else {
			color.New(color.FgYellow).Printf("⚠️ File %s already exists.\n", outFilename)
		}
	}
	return nil
}

func main() {
	cmd := &cli.Command{
		Name:  "quick",
		Usage: "A tool to simplify testing competitive programming problems.",
		Commands: []*cli.Command{
			{
				Name:      "run",
				Usage:     "Run tests for a python script.",
				ArgsUsage: "<python_script>",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "strict",
						Aliases: []string{"s"},
						Usage:   "only run the .in and .out files in the format *{name of the python file without the .py}.in(out)",
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					if !c.Args().Present() {
						log.Println("Error: Python script path is required.")
						return fmt.Errorf("python script path is required")
					}
				run(c)
					return nil
				},
			},
			{
				Name:      "create",
				Usage:     "Create n pairs of .in/.out files.",
				ArgsUsage: "<number_of_files> [problem_name]",
				Action: func(ctx context.Context, c *cli.Command) error {
					return createFiles(c)
				},
			},
		},
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
