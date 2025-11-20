package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/tmc/tokencount/bpe"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "tokencount: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	encoding := flag.String("encoding", "anthropic", "Encoding to use (anthropic, o200k_base, cl100k_base, p50k_base, r50k_base)")
	verbose := flag.Bool("verbose", false, "Verbose output")
	flag.Parse()

	enc, err := bpe.NewEncoder(*encoding)
	if err != nil {
		return fmt.Errorf("failed to get encoding: %w", err)
	}

	files := flag.Args()
	if len(files) == 0 {
		files = []string{"-"} // Use stdin if no files specified
	}

	for _, file := range files {
		if err := processFile(file, enc, *verbose); err != nil {
			return err
		}
	}

	return nil
}

func processFile(filename string, enc bpe.Counter, verbose bool) error {
	var reader io.Reader
	if filename == "-" {
		reader = os.Stdin
	} else {
		file, err := os.Open(filename)
		if err != nil {
			return fmt.Errorf("failed to open file %s: %w", filename, err)
		}
		defer file.Close()
		reader = file
	}

	content, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}

	totalTokens := enc.Count(string(content))

	if verbose {
		fmt.Printf("Tokens in %s: %d\n", filename, totalTokens)
	}
	fmt.Printf("\t%d %s\n", totalTokens, filename)
	return nil
}
