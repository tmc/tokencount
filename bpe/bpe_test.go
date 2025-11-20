package bpe

import (
	"io"
	"strings"
	"testing"
)

func TestNewEncoder(t *testing.T) {
	tests := []struct {
		name    string
		encoder string
		want    bool
	}{
		{"anthropic", "anthropic", true},
		{"claude", "claude", true},
		{"o200k_base", "o200k_base", true},
		{"cl100k_base", "cl100k_base", true},
		{"p50k_base", "p50k_base", true},
		{"r50k_base", "r50k_base", true},
		{"default", "", true},
		{"unknown", "nonexistent", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc, err := NewEncoder(tt.encoder)
			if tt.want {
				if err != nil {
					t.Errorf("NewEncoder(%q) unexpected error: %v", tt.encoder, err)
				}
				if enc == nil {
					t.Errorf("NewEncoder(%q) returned nil encoder", tt.encoder)
				}
			} else {
				if err == nil {
					t.Errorf("NewEncoder(%q) expected error, got nil", tt.encoder)
				}
			}
		})
	}
}

func TestCount(t *testing.T) {
	tests := []struct {
		encoder string
		text    string
	}{
		{"anthropic", "hello world"},
		{"o200k_base", "hello world"},
		{"cl100k_base", "hello world"},
	}

	for _, tt := range tests {
		t.Run(tt.encoder, func(t *testing.T) {
			enc, err := NewEncoder(tt.encoder)
			if err != nil {
				t.Fatal(err)
			}

			count := enc.Count(tt.text)
			if count == 0 {
				t.Errorf("Count(%q) = 0, want > 0", tt.text)
			}
		})
	}
}

func TestEncode(t *testing.T) {
	enc, err := NewEncoder("o200k_base")
	if err != nil {
		t.Fatal(err)
	}

	tokens := enc.Encode("hello")
	if len(tokens) == 0 {
		t.Error("Encode returned no tokens")
	}
}

func TestWriter(t *testing.T) {
	w, err := NewWriter("anthropic")
	if err != nil {
		t.Fatal(err)
	}

	// Write in chunks
	io.WriteString(w, "hello ")
	io.WriteString(w, "world")

	count := w.Count()
	if count != 2 {
		t.Errorf("Count() = %d, want 2", count)
	}

	// Reset and write again
	w.Reset()
	io.WriteString(w, "test")
	count = w.Count()
	if count != 1 {
		t.Errorf("After reset, Count() = %d, want 1", count)
	}
}

func TestWriterCopy(t *testing.T) {
	w, err := NewWriter("o200k_base")
	if err != nil {
		t.Fatal(err)
	}

	// Use io.Copy to write data
	n, err := io.Copy(w, strings.NewReader("The quick brown fox"))
	if err != nil {
		t.Fatal(err)
	}
	if n != 19 {
		t.Errorf("io.Copy wrote %d bytes, want 19", n)
	}

	count := w.Count()
	if count == 0 {
		t.Error("Count() = 0 after io.Copy")
	}
}

func TestCountReader(t *testing.T) {
	tests := []struct {
		encoding string
		text     string
	}{
		{"anthropic", "hello world"},
		{"o200k_base", "The quick brown fox"},
		{"cl100k_base", "test input"},
	}

	for _, tt := range tests {
		t.Run(tt.encoding, func(t *testing.T) {
			r := strings.NewReader(tt.text)
			count, err := CountReader(r, tt.encoding)
			if err != nil {
				t.Fatal(err)
			}
			if count == 0 {
				t.Errorf("CountReader() = 0, want > 0")
			}
		})
	}
}
