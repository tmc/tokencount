package anthropictokenizer_test

import (
	"fmt"
	"log"

	"github.com/tmc/tokencount/anthropictokenizer"
)

func ExampleCounter_Count() {
	counter, err := anthropictokenizer.NewCounter()
	if err != nil {
		log.Fatal(err)
	}

	n := counter.Count("hello world!")
	fmt.Printf("%d tokens\n", n)
	// Output: 3 tokens
}

func ExampleCounter_Count_unicode() {
	counter, err := anthropictokenizer.NewCounter()
	if err != nil {
		log.Fatal(err)
	}

	// Unicode characters are normalized before counting
	n := counter.Count("™")
	fmt.Printf("%d token\n", n)
	// Output: 1 token
}

func ExampleCounter_Encode() {
	counter, err := anthropictokenizer.NewCounter()
	if err != nil {
		log.Fatal(err)
	}

	tokens := counter.Encode("hello world!")
	fmt.Printf("%d tokens: %v\n", len(tokens), tokens)
	// Output: 3 tokens: [9378 2250 2]
}

func ExampleNewCounter() {
	// Create once at startup
	counter, err := anthropictokenizer.NewCounter()
	if err != nil {
		log.Fatal(err)
	}

	// Reuse for all requests
	texts := []string{
		"First message",
		"Second message",
		"Third message",
	}

	for _, text := range texts {
		n := counter.Count(text)
		fmt.Printf("%d tokens\n", n)
	}
	// Output:
	// 2 tokens
	// 2 tokens
	// 2 tokens
}
