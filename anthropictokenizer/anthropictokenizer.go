package anthropictokenizer

import (
	"bufio"
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/anthropic-ai/anthropic-sdk-go"
)

//go:embed anthropic_vocab.jsonl
var vocabFS embed.FS

type Tokenizer struct {
	client *anthropic.Client
	model  string
	vocab  map[string]bool
}

func NewTokenizer(apiKey string, model string) (*Tokenizer, error) {
	client, err := anthropic.NewClient(apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create Anthropic client: %w", err)
	}

	t := &Tokenizer{
		client: client,
		model:  model,
		vocab:  make(map[string]bool),
	}

	if err := t.loadVocab(); err != nil {
		return nil, fmt.Errorf("failed to load vocabulary: %w", err)
	}

	return t, nil
}

func (t *Tokenizer) loadVocab() error {
	file, err := vocabFS.Open("anthropic_vocab.jsonl")
	if err != nil {
		return fmt.Errorf("failed to open vocab file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var entry struct {
			Token string `json:"token"`
		}
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			return fmt.Errorf("failed to unmarshal vocab entry: %w", err)
		}
		t.vocab[entry.Token] = true
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error scanning vocab file: %w", err)
	}

	return nil
}

func (t *Tokenizer) TokenizeText(text string) ([]string, int, error) {
	ctx := context.Background()

	stream, err := t.client.Messages.Create(ctx, &anthropic.MessageCreateParams{
		Model: t.model,
		System: "Copy the text between <tocopy> markers. Include trailing spaces or breaklines. " +
			"Do not write anything else. One example \nInput: <tocopy>Example sentence.</tocopy>\nOutput: Example sentence.",
		Messages: []anthropic.Message{
			{
				Role:    anthropic.MessageRoleUser,
				Content: fmt.Sprintf("<tocopy>%s</tocopy>", text),
			},
		},
		MaxTokens: 1000,
		Stream:    true,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create message stream: %w", err)
	}
	defer stream.Close()

	var tokens []string
	var totalTokensUsage int

	for {
		event, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, 0, fmt.Errorf("error receiving stream event: %w", err)
		}

		if event.Type == anthropic.EventTypeContentBlockDelta {
			tokens = append(tokens, event.Delta.Text)
		}

		if event.Type == anthropic.EventTypeMessageDelta {
			totalTokensUsage = event.Usage.OutputTokens
		}
	}

	return tokens, totalTokensUsage, nil
}

func (t *Tokenizer) UpdateVocab(tokens []string) error {
	var buf bytes.Buffer
	for _, token := range tokens {
		if !t.vocab[token] {
			entry := struct {
				Token string `json:"token"`
			}{Token: token}
			if err := json.NewEncoder(&buf).Encode(entry); err != nil {
				return fmt.Errorf("failed to encode token entry: %w", err)
			}
			t.vocab[token] = true
		}
	}

	if buf.Len() > 0 {
		file, err := vocabFS.Open("anthropic_vocab.jsonl")
		if err != nil {
			return fmt.Errorf("failed to open vocab file: %w", err)
		}
		defer file.Close()

		if _, err := io.Copy(&buf, file); err != nil {
			return fmt.Errorf("failed to append to vocab file: %w", err)
		}

		if err := vocabFS.(*embed.FS).WriteFile("anthropic_vocab.jsonl", buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("failed to write updated vocab file: %w", err)
		}
	}

	return nil
}

func (t *Tokenizer) TokenizeFile(inputPath, outputPath string) error {
	// Implementation for tokenizing a file
	// This would involve reading the input JSONL file, tokenizing each entry,
	// and writing the results to the output JSONL file
	// For brevity, this implementation is omitted
	return nil
}