//go:generate sh -c "curl -sL https://raw.githubusercontent.com/anthropics/anthropic-tokenizer-typescript/main/claude.json | gzip > claude.json.gz"

// Package anthropictokenizer implements token counting for Anthropic's Claude models.
//
// note: This tokenizer is for older Claude models (pre-Claude 3). For Claude 3+ models,
// this provides only a rough approximation. Use for estimation purposes only.
//
// The tokenizer uses byte-pair encoding (BPE) with Unicode NFKC normalization.
// All vocabulary data is embedded at compile time, requiring no runtime file I/O.
//
// Basic usage:
//
//	counter, err := anthropictokenizer.NewCounter()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	count := counter.Count("Hello, Claude!")
//
// The Counter is safe for concurrent use.
package anthropictokenizer
