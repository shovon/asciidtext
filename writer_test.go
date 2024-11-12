package asciidtext_test

import (
	"bytes"
	"testing"

	"github.com/shovon/asciidtext"
)

func TestWriter_Write(t *testing.T) {
	tests := []struct {
		name     string
		records  [][]string
		expected string
	}{
		{
			name:     "single record",
			records:  [][]string{{"hello", "world"}},
			expected: "hello\x1Fworld",
		},
		{
			name:     "multiple records",
			records:  [][]string{{"hello", "world"}, {"foo", "bar"}},
			expected: "hello\x1Fworld\x1Efoo\x1Fbar",
		},
		{
			name:     "empty record",
			records:  [][]string{{}},
			expected: "",
		},
		{
			name:     "record with empty fields",
			records:  [][]string{{"", ""}},
			expected: "\x1F",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			w := asciidtext.NewWriter(buf)

			for _, record := range tt.records {
				err := w.Write(record)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
			w.Flush()

			if got := buf.String(); got != tt.expected {
				t.Errorf("Writer.Write() = %q, want %q", got, tt.expected)
			}

			if err := w.Error(); err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestWriter_WriteAll(t *testing.T) {
	tests := []struct {
		name     string
		records  [][]string
		expected string
	}{
		{
			name:     "multiple records",
			records:  [][]string{{"hello", "world"}, {"foo", "bar"}},
			expected: "hello\x1Fworld\x1Efoo\x1Fbar",
		},
		{
			name:     "single record",
			records:  [][]string{{"test", "data"}},
			expected: "test\x1Fdata",
		},
		{
			name:     "empty records",
			records:  [][]string{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			w := asciidtext.NewWriter(buf)

			err := w.WriteAll(tt.records)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if got := buf.String(); got != tt.expected {
				t.Errorf("Writer.WriteAll() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestWriter_Error(t *testing.T) {
	buf := &bytes.Buffer{}
	w := asciidtext.NewWriter(buf)

	if err := w.Error(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Write some valid data
	err := w.Write([]string{"test"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	w.Flush()

	if err := w.Error(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
