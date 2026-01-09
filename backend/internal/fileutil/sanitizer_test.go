package fileutil

import (
	"testing"
)

func TestSanitizePath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "removes forward slashes",
			input:    "Book/Title",
			expected: "Book-Title",
		},
		{
			name:     "removes backslashes",
			input:    "Book\\Title",
			expected: "Book-Title",
		},
		{
			name:     "removes colons",
			input:    "Book: The Title",
			expected: "Book- The Title",
		},
		{
			name:     "removes asterisks",
			input:    "Book*Title",
			expected: "Book-Title",
		},
		{
			name:     "removes question marks",
			input:    "Book?Title",
			expected: "Book-Title",
		},
		{
			name:     "removes quotes",
			input:    "Book\"Title",
			expected: "Book-Title",
		},
		{
			name:     "removes angle brackets",
			input:    "Book<Title>",
			expected: "Book-Title-",
		},
		{
			name:     "removes pipes",
			input:    "Book|Title",
			expected: "Book-Title",
		},
		{
			name:     "trims leading whitespace",
			input:    "  Book Title",
			expected: "Book Title",
		},
		{
			name:     "trims trailing whitespace",
			input:    "Book Title  ",
			expected: "Book Title",
		},
		{
			name:     "collapses multiple spaces",
			input:    "Book   Title",
			expected: "Book Title",
		},
		{
			name:     "handles empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "handles whitespace only",
			input:    "   ",
			expected: "",
		},
		{
			name:     "handles unicode characters",
			input:    "Book Título über 日本",
			expected: "Book Título über 日本",
		},
		{
			name:     "prevents path traversal with dots",
			input:    "../../../etc/passwd",
			expected: "..-..-..-etc-passwd",
		},
		{
			name:     "handles multiple invalid characters",
			input:    "Book: The/Title\\Part*2?<Test>|End",
			expected: "Book- The-Title-Part-2--Test--End",
		},
		{
			name:     "preserves valid characters",
			input:    "Book Title - Part 2 (2024)",
			expected: "Book Title - Part 2 (2024)",
		},
		{
			name:     "handles newlines and tabs",
			input:    "Book\nTitle\tPart",
			expected: "Book Title Part",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizePath(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizePath(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "prevents path traversal in filename",
			input:    "../secret.txt",
			expected: "..-secret.txt",
		},
		{
			name:     "removes invalid characters from filename",
			input:    "file:name*.txt",
			expected: "file-name-.txt",
		},
		{
			name:     "preserves file extension",
			input:    "book title.mp3",
			expected: "book title.mp3",
		},
		{
			name:     "handles filename with multiple dots",
			input:    "book.title.part.1.m4b",
			expected: "book.title.part.1.m4b",
		},
		{
			name:     "handles empty filename",
			input:    "",
			expected: "",
		},
		{
			name:     "handles unicode in filename",
			input:    "日本語.txt",
			expected: "日本語.txt",
		},
		{
			name:     "prevents directory separator injection",
			input:    "file/name\\test",
			expected: "file-name-test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeFilename(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeFilename(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// Benchmark to ensure performance is acceptable
func BenchmarkSanitizePath(b *testing.B) {
	input := "Book: The/Title\\Part*2?<Test>|End"
	for i := 0; i < b.N; i++ {
		SanitizePath(input)
	}
}

func BenchmarkSanitizeFilename(b *testing.B) {
	input := "../file:name*.txt"
	for i := 0; i < b.N; i++ {
		SanitizeFilename(input)
	}
}
