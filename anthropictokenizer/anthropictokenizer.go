package anthropictokenizer

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	_ "embed"

	"golang.org/x/text/unicode/norm"
)

//go:embed claude.json.gz
var configDataGZ []byte

// A Counter counts tokens in text using Claude's tokenization scheme.
type Counter struct {
	vocab   map[string]int
	pattern *regexp.Regexp
	special map[string]int
}

type config struct {
	PatternStr     string         `json:"pat_str"`
	SpecialTokens  map[string]int `json:"special_tokens"`
	ExplicitNVocab int            `json:"explicit_n_vocab"`
	BPERanks       string         `json:"bpe_ranks"`
}

// NewCounter creates a new token counter.
// The returned Counter is safe for concurrent use.
func NewCounter() (*Counter, error) {
	// Decompress the embedded config
	gr, err := gzip.NewReader(bytes.NewReader(configDataGZ))
	if err != nil {
		return nil, fmt.Errorf("decompress config: %w", err)
	}
	defer gr.Close()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(gr); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg config
	if err := json.Unmarshal(buf.Bytes(), &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	vocab, err := parseRanks(cfg.BPERanks)
	if err != nil {
		return nil, fmt.Errorf("parse ranks: %w", err)
	}

	// Simplified pattern compatible with standard regexp
	// Original pattern: 's|'t|'re|'ve|'m|'ll|'d| ?\p{L}+| ?\p{N}+| ?[^\s\p{L}\p{N}]+|\s+(?!\S)|\s+
	// Simplified: removes negative lookahead (?!\S) and uses standard character classes
	pattern := regexp.MustCompile(
		`'[sStTdDmM]|'[rR][eE]|'[vV][eE]|'[lL][lL]|[^\r\n\p{L}\p{N}]?\p{L}+|\p{N}{1,3}| ?[^\s\p{L}\p{N}]+[\r\n]*|\s*[\r\n]+|\s+`)

	return &Counter{
		vocab:   vocab,
		pattern: pattern,
		special: cfg.SpecialTokens,
	}, nil
}

// Count returns the number of tokens in text.
// Text is normalized using Unicode NFKC normalization before tokenization.
func (c *Counter) Count(text string) int {
	return len(c.Encode(text))
}

// Encode returns the token IDs for the given text.
// Text is normalized using Unicode NFKC normalization before tokenization.
func (c *Counter) Encode(text string) []int {
	// Apply NFKC normalization
	text = norm.NFKC.String(text)

	var tokens []int
	pos := 0

	for pos < len(text) {
		// Check for special tokens
		matched := false
		for tok, id := range c.special {
			if strings.HasPrefix(text[pos:], tok) {
				tokens = append(tokens, id)
				pos += len(tok)
				matched = true
				break
			}
		}
		if matched {
			continue
		}

		// Match pattern
		loc := c.pattern.FindStringIndex(text[pos:])
		if loc == nil {
			break
		}

		chunk := text[pos : pos+loc[1]]
		tokens = append(tokens, c.encodeChunk([]byte(chunk))...)
		pos += len(chunk)
	}

	return tokens
}

// encodeChunk applies BPE to a single chunk of bytes.
func (c *Counter) encodeChunk(chunk []byte) []int {
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
			if rank, ok := c.vocab[pair]; ok && rank < bestRank {
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
		if rank, ok := c.vocab[string(part)]; ok {
			tokens = append(tokens, rank)
		}
	}

	return tokens
}

// parseRanks parses the BPE merge ranks from the space-separated format in claude.json.
// Format: space-separated base64-encoded tokens in rank order.
func parseRanks(data string) (map[string]int, error) {
	fields := strings.Fields(data)
	ranks := make(map[string]int, len(fields))

	for rank, b64Token := range fields {
		tokenBytes, err := base64.StdEncoding.DecodeString(b64Token)
		if err != nil {
			// Skip non-base64 tokens (like the first token "!" which is literal)
			ranks[b64Token] = rank
			continue
		}

		ranks[string(tokenBytes)] = rank
	}

	return ranks, nil
}
