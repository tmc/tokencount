package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkoukk/tiktoken-go"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "tokencount: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	text, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	encoding := "cl100k_base"

	tke, err := tiktoken.GetEncoding(encoding)
	if err != nil {
		return err
	}

	// encode
	token := tke.Encode(string(text), nil, nil)

	// num_tokens
	fmt.Println(len(token))
	return nil
}
