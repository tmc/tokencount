package openaitokenizer_test

import (
	"fmt"
	"log"

	"github.com/tmc/tokencount/openaitokenizer"
)

func Example() {
	enc, err := openaitokenizer.NewEncoder("o200k_base")
	if err != nil {
		log.Fatal(err)
	}

	count := enc.Count("Hello, world!")
	fmt.Printf("%d tokens\n", count)
	// Output: 4 tokens
}

func ExampleNewEncoder() {
	// Create an encoder for GPT-4o
	enc, err := openaitokenizer.NewEncoder("o200k_base")
	if err != nil {
		log.Fatal(err)
	}

	text := "The quick brown fox"
	fmt.Printf("%d tokens\n", enc.Count(text))
	// Output: 4 tokens
}

func ExampleEncoder_Count() {
	enc, err := openaitokenizer.NewEncoder("cl100k_base")
	if err != nil {
		log.Fatal(err)
	}

	// Count tokens for GPT-4/GPT-3.5-turbo
	count := enc.Count("Hello, world!")
	fmt.Printf("%d tokens\n", count)
	// Output: 4 tokens
}

func ExampleEncoder_Encode() {
	enc, err := openaitokenizer.NewEncoder("o200k_base")
	if err != nil {
		log.Fatal(err)
	}

	// Get token IDs
	tokens := enc.Encode("Hello")
	fmt.Printf("%d tokens: %v\n", len(tokens), tokens)
	// Output: 1 tokens: [13225]
}
