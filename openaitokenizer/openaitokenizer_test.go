package openaitokenizer

import "testing"

func TestEncoder(t *testing.T) {
	enc, err := NewEncoder("o200k_base")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		text string
		want int // approximate token count
	}{
		{"hello world", 2},
		{"The quick brown fox", 4},
		{"", 0},
	}

	for _, tt := range tests {
		got := enc.Count(tt.text)
		// Allow some variance in token count
		if got == 0 && tt.want != 0 {
			t.Errorf("Count(%q) = %d, want approximately %d", tt.text, got, tt.want)
		}
		if got > 0 && tt.want == 0 {
			t.Errorf("Count(%q) = %d, want %d", tt.text, got, tt.want)
		}
	}
}

func TestEncoderEncode(t *testing.T) {
	enc, err := NewEncoder("o200k_base")
	if err != nil {
		t.Fatal(err)
	}

	text := "hello"
	tokens := enc.Encode(text)
	if len(tokens) == 0 {
		t.Errorf("Encode(%q) returned no tokens", text)
	}
}

func TestAllEncodings(t *testing.T) {
	encodings := []string{"o200k_base", "cl100k_base", "p50k_base", "r50k_base"}

	for _, name := range encodings {
		t.Run(name, func(t *testing.T) {
			enc, err := NewEncoder(name)
			if err != nil {
				t.Fatalf("NewEncoder(%q) failed: %v", name, err)
			}

			text := "hello world"
			count := enc.Count(text)
			if count == 0 {
				t.Errorf("Count(%q) = 0, want > 0", text)
			}
		})
	}
}

func TestUnknownEncoding(t *testing.T) {
	_, err := NewEncoder("nonexistent")
	if err == nil {
		t.Error("NewEncoder(nonexistent) should return error")
	}
}
