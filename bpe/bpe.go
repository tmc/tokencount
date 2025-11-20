package bpe

import (
	"fmt"
	"io"

	"github.com/tmc/tokencount/anthropictokenizer"
	"github.com/tmc/tokencount/openaitokenizer"
)

// Counter provides a unified interface for token counting.
type Counter interface {
	Count(text string) int
}

// Encoder extends Counter with token ID encoding capabilities.
type Encoder interface {
	Counter
	Encode(text string) []int
}

// A Writer counts tokens as data is written to it.
// Create a Writer using NewWriter; the zero value is not usable.
type Writer struct {
	enc Counter
	buf []byte
	n   int
}

// NewWriter returns a Writer that counts tokens using the named encoding.
func NewWriter(encoding string) (*Writer, error) {
	enc, err := NewEncoder(encoding)
	if err != nil {
		return nil, err
	}
	return &Writer{enc: enc}, nil
}

// Write writes p to the token counter.
// It always returns len(p), nil.
func (w *Writer) Write(p []byte) (n int, err error) {
	w.buf = append(w.buf, p...)
	return len(p), nil
}

// Count returns the number of tokens written so far.
func (w *Writer) Count() int {
	if len(w.buf) == 0 {
		return w.n
	}
	return w.enc.Count(string(w.buf))
}

// Reset resets the Writer to be empty.
func (w *Writer) Reset() {
	w.buf = w.buf[:0]
	w.n = 0
}

// CountReader counts tokens from an io.Reader.
// It reads the entire contents into memory.
func CountReader(r io.Reader, encoding string) (int, error) {
	w, err := NewWriter(encoding)
	if err != nil {
		return 0, err
	}
	if _, err := io.Copy(w, r); err != nil {
		return 0, err
	}
	return w.Count(), nil
}

// NewEncoder returns an encoder for the named tokenizer.
//
// Supported encodings:
//   - "anthropic" or "claude": Anthropic's Claude tokenizer
//   - "o200k_base": OpenAI GPT-4o and newer (default)
//   - "cl100k_base": OpenAI GPT-4, GPT-3.5-turbo
//   - "p50k_base": OpenAI Codex models
//   - "r50k_base": OpenAI GPT-3 models
func NewEncoder(name string) (Encoder, error) {
	switch name {
	case "anthropic", "claude":
		counter, err := anthropictokenizer.NewCounter()
		if err != nil {
			return nil, fmt.Errorf("failed to create anthropic tokenizer: %w", err)
		}
		return &anthropicEncoder{counter}, nil
	case "o200k_base", "cl100k_base", "p50k_base", "r50k_base", "":
		if name == "" {
			name = "o200k_base"
		}
		return openaitokenizer.NewEncoder(name)
	default:
		return nil, fmt.Errorf("unknown encoding %q", name)
	}
}

// NewCounter returns a counter for the named tokenizer.
// This is a convenience wrapper around NewEncoder for when you only need counting.
func NewCounter(name string) (Counter, error) {
	return NewEncoder(name)
}

// anthropicEncoder wraps anthropictokenizer.Counter to implement Encoder.
type anthropicEncoder struct {
	*anthropictokenizer.Counter
}

// Encode is not supported for Anthropic tokenizer.
// Use Count() instead for token counting.
func (e *anthropicEncoder) Encode(text string) []int {
	// Anthropic tokenizer doesn't expose token IDs, only counts
	panic("Encode not supported for Anthropic tokenizer, use Count() instead")
}
