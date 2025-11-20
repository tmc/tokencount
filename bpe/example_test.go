package bpe_test

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/tmc/tokencount/bpe"
)

func Example() {
	enc, err := bpe.NewEncoder("o200k_base")
	if err != nil {
		log.Fatal(err)
	}

	count := enc.Count("Hello, world!")
	fmt.Printf("%d tokens\n", count)
	// Output: 4 tokens
}

func ExampleWriter() {
	w, err := bpe.NewWriter("anthropic")
	if err != nil {
		log.Fatal(err)
	}

	// Write data in chunks (simulating streaming)
	io.WriteString(w, "Hello ")
	io.WriteString(w, "world!")

	fmt.Printf("%d tokens\n", w.Count())
	// Output: 3 tokens
}

func ExampleWriter_streaming() {
	w, err := bpe.NewWriter("o200k_base")
	if err != nil {
		log.Fatal(err)
	}

	// Simulate streaming data from a reader
	reader := strings.NewReader("The quick brown fox jumps over the lazy dog")
	io.Copy(w, reader)

	fmt.Printf("%d tokens\n", w.Count())
	// Output: 9 tokens
}

func ExampleWriter_Reset() {
	w, err := bpe.NewWriter("anthropic")
	if err != nil {
		log.Fatal(err)
	}

	io.WriteString(w, "First message")
	fmt.Printf("First: %d tokens\n", w.Count())

	w.Reset()
	io.WriteString(w, "Second message")
	fmt.Printf("Second: %d tokens\n", w.Count())
	// Output:
	// First: 2 tokens
	// Second: 2 tokens
}

func ExampleCountReader() {
	reader := strings.NewReader("Hello, world!")
	count, err := bpe.CountReader(reader, "anthropic")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d tokens\n", count)
	// Output: 4 tokens
}

func ExampleCountReader_file() {
	// Count tokens from any io.Reader (file, network stream, etc.)
	data := "The quick brown fox jumps over the lazy dog"
	reader := strings.NewReader(data)

	count, err := bpe.CountReader(reader, "o200k_base")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d tokens\n", count)
	// Output: 9 tokens
}
