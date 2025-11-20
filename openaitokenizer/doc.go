//go:generate sh -c "curl -s https://openaipublic.blob.core.windows.net/encodings/o200k_base.tiktoken -o o200k_base.tiktoken"
//go:generate sh -c "curl -s https://openaipublic.blob.core.windows.net/encodings/cl100k_base.tiktoken -o cl100k_base.tiktoken"
//go:generate sh -c "curl -s https://openaipublic.blob.core.windows.net/encodings/p50k_base.tiktoken -o p50k_base.tiktoken"
//go:generate sh -c "curl -s https://openaipublic.blob.core.windows.net/encodings/r50k_base.tiktoken -o r50k_base.tiktoken"
//go:generate sh -c "go run golang.org/x/exp/cmd/txtar@latest *.tiktoken > encodings.txtar"
//go:generate sh -c "gzip -f encodings.txtar"
//go:generate sh -c "rm *.tiktoken"

// Package openaitokenizer implements OpenAI tiktoken byte pair encoding.
//
// Supports o200k_base (GPT-4o), cl100k_base (GPT-4, GPT-3.5), p50k_base (Codex),
// and r50k_base (GPT-3). All vocabularies are embedded for offline use.
//
// Basic usage:
//
//	enc, err := openaitokenizer.NewEncoder("o200k_base")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	count := enc.Count("Hello, world!")
//	tokens := enc.Encode("Hello, world!")
package openaitokenizer
