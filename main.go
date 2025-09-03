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

	"github.com/pmezard/go-difflib/difflib"
	"github.com/urfave/cli/v3"

	"github.com/matwate/quick/config"
	renderer "github.com/matwate/quick/renderer"
)

const (
	// Emojis
	EmojiPassed  = "✅"
	EmojiFailed  = "❌"
	EmojiTimeout = "⌛"
)

func run(c *cli.Command, cfg *config.Config) {
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

		renderer.Render("Found matching input files (strict mode):")
		for _, file := range inputs {
			if !file.IsDir() && strings.Contains(file.Name(), baseFilename) &&
				strings.HasSuffix(file.Name(), ".in") {
				inFilename := file.Name()
				outFilename := strings.TrimSuffix(inFilename, ".in") + ".out"

				// Check if the corresponding .out file exists
				if _, err := os.Stat(outFilename); err == nil {
					renderer.Render(fmt.Sprintf(" - %s", inFilename))
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
		renderer.Render("Found input files:")
		for _, file := range inputs {
			if !file.IsDir() && len(file.Name()) > 3 && file.Name()[len(file.Name())-3:] == ".in" {
				renderer.Render(fmt.Sprintf(" - %s", file.Name()))
				ins = append(ins, file.Name())
			}
		}
		renderer.Render("Found output files:")
		for _, file := range inputs {
			if !file.IsDir() && len(file.Name()) > 4 && file.Name()[len(file.Name())-4:] == ".out" {
				renderer.Render(fmt.Sprintf(" - %s", file.Name()))
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

	if len(Ins_Outs) == 0 {
		renderer.RenderYellow("No test cases found.")
		return
	}

	passed := true
	failed := []string{}
	for in, out := range Ins_Outs {
		renderer.Render(
			fmt.Sprintf("Running test case: %s -> %s\n ", in, out),
		) // Note: Added a space after \n to match original formatting

		// Read input from the *.in file
		inputData, err := os.ReadFile(in)
		if err != nil {
			log.Fatalf("Failed to read input file %s: %v", in, err)
		}

		// Execute the Python script with the input inputData
		cmd := fmt.Sprintf("%s %s", cfg.RunCommand, pyFilename)

		// Create a context with a timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		command := exec.CommandContext(ctx, cfg.Shell, "-c", cmd)
		command.Stdin = bytes.NewReader(inputData)

		var outBuffer, errBuffer bytes.Buffer
		command.Stdout = &outBuffer
		command.Stderr = &errBuffer

		startTime := time.Now()
		err = command.Run()
		duration := time.Since(startTime)

		if ctx.Err() == context.DeadlineExceeded {
			renderer.RenderYellow(
				fmt.Sprintf("%s Test case %s timed out! (%s)", EmojiTimeout, in, duration),
			)
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
			renderer.RenderGreen(
				fmt.Sprintf("%s Test case %s passed! (%s)", EmojiPassed, in, duration),
			)
			passed = passed && true
		} else {
			renderer.RenderRed(fmt.Sprintf("%s Test case %s failed! (%s)", EmojiFailed, in, duration))

			diff := difflib.UnifiedDiff{
				A:        difflib.SplitLines(string(expectedOutput)),
				B:        difflib.SplitLines(string(actualOutput)),
				FromFile: "Expected",
				ToFile:   "Got",
				Context:  3,
			}
			diffStr, _ := difflib.GetUnifiedDiffString(diff)
			colorizeDiff(diffStr)
			passed = false
			failed = append(failed, in)
		}

	}
	if passed {
		renderer.RenderGreen("All test cases passed! 🎉")
	} else {
		for _, f := range failed {
			renderer.RenderRed(fmt.Sprintf(" - %s", f))
		}
		renderer.RenderRed(fmt.Sprintf("%d test cases failed.", len(failed)))
	}
}

func colorizeDiff(diff string) {
	for _, line := range strings.Split(diff, "\n") {
		switch {
		case strings.HasPrefix(line, "+"):
			renderer.RenderGreen(line)
		case strings.HasPrefix(line, "-"):
			renderer.RenderRed(line)
		case strings.HasPrefix(line, "@@"):
			renderer.RenderYellow(line)
		default:
			renderer.Render(line)
		}
	}
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
				renderer.RenderRed(fmt.Sprintf("❌ Failed to create %s: %v", inFilename, err))
			} else {
				renderer.RenderGreen(fmt.Sprintf("✅ Created %s", inFilename))
			}
		} else {
			renderer.RenderYellow(fmt.Sprintf("⚠️ File %s already exists.", inFilename))
		}

		// Handle .out file
		if _, err := os.Stat(outFilename); os.IsNotExist(err) {
			if err := os.WriteFile(outFilename, []byte{}, 0o644); err != nil {
				renderer.RenderRed(fmt.Sprintf("❌ Failed to create %s: %v", outFilename, err))
			} else {
				renderer.RenderGreen(fmt.Sprintf("✅ Created %s", outFilename))
			}
		} else {
			renderer.RenderYellow(fmt.Sprintf("⚠️ File %s already exists.", outFilename))
		}
	}
	return nil
}

func main() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

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
					run(c, cfg)
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
			{
				Name:  "config",
				Usage: "Manage configuration.",
				Commands: []*cli.Command{
					{
						Name:  "get",
						Usage: "Get a configuration value.",
						Action: func(ctx context.Context, c *cli.Command) error {
							key := c.Args().Get(0)
							switch key {
							case "shell":
								renderer.Render(cfg.Shell)
							case "run_command":
								renderer.Render(cfg.RunCommand)
							default:
								renderer.Render("Unknown configuration key.")
							}
							return nil
						},
					},
					{
						Name:  "set",
						Usage: "Set a configuration value.",
						Action: func(ctx context.Context, c *cli.Command) error {
							key := c.Args().Get(0)
							value := c.Args().Get(1)
							switch key {
							case "shell":
								cfg.Shell = value
							case "run_command":
								cfg.RunCommand = value
							default:
								renderer.Render("Unknown configuration key.")
							}
							err := config.SaveConfig("config.yaml", cfg)
							if err != nil {
								log.Fatal(err)
							}
							return nil
						},
					},
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
