// Package openaitokenizer implements OpenAI's tiktoken BPE tokenization.
//
// This package provides offline tokenization for OpenAI models using
// embedded vocabulary files. All tokenization is performed locally
// without any network requests.
package openaitokenizer

import (
	"bufio"
	"bytes"
	"compress/gzip"
	_ "embed"
	"encoding/base64"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/tools/txtar"
)

//go:embed encodings.txtar.gz
var vocabData []byte

// An Encoder tokenizes text using byte pair encoding.
type Encoder struct {
	vocab   map[string]int
	pattern *regexp.Regexp
}

// NewEncoder returns a new encoder for the named encoding.
// Supported encodings: o200k_base, cl100k_base, p50k_base, r50k_base.
func NewEncoder(name string) (*Encoder, error) {
	// Decompress the gzipped txtar archive
	gr, err := gzip.NewReader(bytes.NewReader(vocabData))
	if err != nil {
		return nil, fmt.Errorf("failed to decompress encodings: %w", err)
	}
	defer gr.Close()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(gr); err != nil {
		return nil, fmt.Errorf("failed to read encodings: %w", err)
	}

	// Parse the txtar archive
	ar := txtar.Parse(buf.Bytes())

	// Find the requested encoding
	var data []byte
	for _, f := range ar.Files {
		if f.Name == name+".tiktoken" {
			data = f.Data
			break
		}
	}
	if data == nil {
		return nil, fmt.Errorf("unknown encoding %q", name)
	}

	e := &Encoder{
		vocab: make(map[string]int),
	}

	// Parse tiktoken format: base64token rank
	s := bufio.NewScanner(bytes.NewReader(data))
	for s.Scan() {
		line := s.Text()
		fields := strings.Fields(line)
		if len(fields) != 2 {
			continue
		}
		tokenBytes, err := base64.StdEncoding.DecodeString(fields[0])
		if err != nil {
			continue
		}
		rank, err := strconv.Atoi(fields[1])
		if err != nil {
			continue
		}
		e.vocab[string(tokenBytes)] = rank
	}
	if err := s.Err(); err != nil {
		return nil, err
	}

	// Pattern splits on word boundaries, whitespace, and punctuation
	e.pattern = regexp.MustCompile(
		`'[sStTdDmM]|'[rR][eE]|'[vV][eE]|'[lL][lL]|[^\r\n\p{L}\p{N}]?\p{L}+|\p{N}{1,3}| ?[^\s\p{L}\p{N}]+[\r\n]*|\s*[\r\n]+|\s+`)

	return e, nil
}

// Encode returns the token IDs for the given text.
func (e *Encoder) Encode(text string) []int {
	var tokens []int

	// Split text into chunks using the pattern
	for _, chunk := range e.pattern.FindAllString(text, -1) {
		// Apply BPE to each chunk
		tokens = append(tokens, e.encodeChunk([]byte(chunk))...)
	}

	return tokens
}

// encodeChunk applies BPE to a single chunk of bytes.
func (e *Encoder) encodeChunk(chunk []byte) []int {
	if len(chunk) == 0 {
		return nil
	}

	// Start with individual bytes
	parts := make([][]byte, len(chunk))
	for i := range chunk {
		parts[i] = chunk[i : i+1]
	}

	// Repeatedly merge the best pair until no more merges possible
	for len(parts) > 1 {
		bestIdx := -1
		bestRank := int(^uint(0) >> 1) // max int

		// Find the pair with the lowest rank in vocab
		for i := 0; i < len(parts)-1; i++ {
			pair := string(append(parts[i], parts[i+1]...))
			if rank, ok := e.vocab[pair]; ok && rank < bestRank {
				bestRank = rank
				bestIdx = i
			}
		}

		if bestIdx == -1 {
			break // No more pairs to merge
		}

		// Merge the best pair
		merged := append(parts[bestIdx], parts[bestIdx+1]...)
		parts = append(parts[:bestIdx], append([][]byte{merged}, parts[bestIdx+2:]...)...)
	}

	// Convert parts to token IDs
	tokens := make([]int, 0, len(parts))
	for _, part := range parts {
		if rank, ok := e.vocab[string(part)]; ok {
			tokens = append(tokens, rank)
		}
	}

	return tokens
}

// Count returns the number of tokens in the text.
func (e *Encoder) Count(text string) int {
	return len(e.Encode(text))
}
