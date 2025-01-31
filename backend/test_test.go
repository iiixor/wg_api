package main

import (
	"testing"
)

func TestEscapeMarkdownV2(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Базовые случаи
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "No special chars",
			input:    "Hello World",
			expected: "Hello World",
		},

		// Тесты для каждого специального символа
		{
			name:     "Underscore",
			input:    "_test_",
			expected: "\\_test\\_",
		},
		{
			name:     "Asterisk",
			input:    "*text*",
			expected: "\\*text\\*",
		},
		{
			name:     "Brackets",
			input:    "[link]",
			expected: "\\[link\\]",
		},
		{
			name:     "Parentheses",
			input:    "(text)",
			expected: "\\(text\\)",
		},
		{
			name:     "Tilde",
			input:    "~text~",
			expected: "\\~text\\~",
		},
		{
			name:     "Backtick",
			input:    "`code`",
			expected: "\\`code\\`",
		},
		{
			name:     "Greater than",
			input:    "> quote",
			expected: "\\> quote",
		},
		{
			name:     "Hash",
			input:    "# header",
			expected: "\\# header",
		},
		{
			name:     "Plus",
			input:    "+ item",
			expected: "\\+ item",
		},
		{
			name:     "Minus",
			input:    "- item",
			expected: "\\- item",
		},
		{
			name:     "Equal",
			input:    "a = b",
			expected: "a \\= b",
		},
		{
			name:     "Pipe",
			input:    "a|b",
			expected: "a\\|b",
		},
		{
			name:     "Braces",
			input:    "{text}",
			expected: "\\{text\\}",
		},
		{
			name:     "Dot",
			input:    "10.0.0.1",
			expected: "10\\.0\\.0\\.1",
		},
		{
			name:     "Exclamation",
			input:    "Hello!",
			expected: "Hello\\!",
		},
		{
			name:     "Real Name",
			input:    "2-Even_forgotten_my_name-2025-01-28-17-28-15",
			expected: "2\\-Even\\_forgotten\\_my\\_name\\-2025\\-01\\-28\\-17\\-28\\-15",
		},
		// Комбинированные случаи
		{
			name:     "Mixed characters",
			input:    "Hello_World* [test](link) ~1.2.3~",
			expected: "Hello\\_World\\* \\[test\\]\\(link\\) \\~1\\.2\\.3\\~",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := escapeMarkdownV2(tt.input)
			if result != tt.expected {
				t.Errorf(
					"Input: %q\nExpected: %q\nGot: %q",
					tt.input,
					tt.expected,
					result,
				)
			}
		})
	}
}

func BenchmarkEscapeMarkdownV2(b *testing.B) {
	testString := "Hello_World* [test](link) ~1.2.3~"
	for i := 0; i < b.N; i++ {
		escapeMarkdownV2(testString)
	}
}
