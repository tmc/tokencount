// Package bpe provides a unified interface for byte pair encoding tokenization.
//
// Supports both OpenAI (tiktoken) and Anthropic (Claude) tokenizers through
// a common Counter/Encoder interface. All tokenization is performed offline
// using embedded vocabularies.
//
// Basic usage:
//
//	enc, err := bpe.NewEncoder("o200k_base")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	count := enc.Count("Hello, world!")
//
// Streaming usage:
//
//	w, _ := bpe.NewWriter("anthropic")
//	io.Copy(w, reader)
//	count := w.Count()
//
// Supported encodings: anthropic, claude, o200k_base, cl100k_base, p50k_base, r50k_base.
package bpe
