package anthropictokenizer

import "testing"

func TestCount(t *testing.T) {
	counter, err := NewCounter()
	if err != nil {
		t.Fatalf("NewCounter: %v", err)
	}

	tests := []struct {
		name  string
		input string
		want  int
	}{
		{
			name:  "simple text",
			input: "hello world!",
			want:  3,
		},
		{
			name:  "unicode normalization - trademark",
			input: "™",
			want:  1,
		},
		{
			name:  "unicode normalization - greek",
			input: "ϰ",
			want:  1,
		},
		{
			name:  "special token - EOT",
			input: "<EOT>",
			want:  1,
		},
		{
			name:  "empty string",
			input: "",
			want:  0,
		},
		{
			name:  "contractions",
			input: "I'm, you're, they've, we'll, it's",
			want:  14,
		},
		{
			name:  "whitespace handling",
			input: "spaces   between   words",
			want:  5,
		},
		{
			name:  "trailing whitespace",
			input: "text with trailing  ",
			want:  4,
		},
		{
			name:  "numbers",
			input: "The year 2024 has 365 days",
			want:  9,
		},
		{
			name:  "mixed punctuation",
			input: "Hello, world! How are you?",
			want:  8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := counter.Count(tt.input)
			if got != tt.want {
				t.Errorf("Count(%q) = %d; want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestEncode(t *testing.T) {
	counter, err := NewCounter()
	if err != nil {
		t.Fatalf("NewCounter: %v", err)
	}

	tests := []struct {
		name string
		input string
		wantLen int
	}{
		{
			name: "simple text",
			input: "hello world!",
			wantLen: 3,
		},
		{
			name: "special token",
			input: "<EOT>",
			wantLen: 1,
		},
		{
			name: "empty string",
			input: "",
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := counter.Encode(tt.input)
			if len(tokens) != tt.wantLen {
				t.Errorf("Encode(%q) returned %d tokens; want %d", tt.input, len(tokens), tt.wantLen)
			}
			// Verify all token IDs are valid
			for i, id := range tokens {
				if id < 0 {
					t.Errorf("Encode(%q) token[%d] = %d; want non-negative", tt.input, i, id)
				}
			}
		})
	}
}

func BenchmarkCount(b *testing.B) {
	counter, err := NewCounter()
	if err != nil {
		b.Fatalf("NewCounter: %v", err)
	}

	texts := []string{
		"hello world!",
		"The quick brown fox jumps over the lazy dog.",
		"This is a longer text that contains multiple sentences. It should give us a better sense of tokenization performance.",
	}

	for _, text := range texts {
		name := text
		if len(name) > 30 {
			name = name[:30] + "..."
		}
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = counter.Count(text)
			}
		})
	}
}
