package asciidtext

import (
	"strings"
	"testing"
)

func TestReader_Read(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected [][]string
	}{
		{
			name:     "single record",
			input:    "hello\x1Fworld",
			expected: [][]string{{"hello", "world"}},
		},
		{
			name:     "multiple records",
			input:    "hello\x1Fworld\x1Efoo\x1Fbar",
			expected: [][]string{{"hello", "world"}, {"foo", "bar"}},
		},
		{
			name:     "empty record",
			input:    "",
			expected: [][]string{},
		},
		{
			name:     "record with empty fields",
			input:    "\x1F",
			expected: [][]string{{"", ""}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewReader(strings.NewReader(tt.input))
			var records [][]string

			for {
				record, err := r.Read()
				if err != nil {
					break
				}
				records = append(records, record)
			}

			if len(records) != len(tt.expected) {
				t.Errorf("Reader.Read() got %d records, want %d", len(records), len(tt.expected))
				return
			}

			for i := range records {
				if len(records[i]) != len(tt.expected[i]) {
					t.Errorf("Record %d: got %d fields, want %d", i, len(records[i]), len(tt.expected[i]))
					continue
				}
				for j := range records[i] {
					if records[i][j] != tt.expected[i][j] {
						t.Errorf("Record %d field %d: got %q, want %q", i, j, records[i][j], tt.expected[i][j])
					}
				}
			}
		})
	}
}

func TestReader_ReadAll(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected [][]string
	}{
		{
			name:     "multiple records",
			input:    "hello\x1Fworld\x1Efoo\x1Fbar",
			expected: [][]string{{"hello", "world"}, {"foo", "bar"}},
		},
		{
			name:     "single record",
			input:    "test\x1Fdata",
			expected: [][]string{{"test", "data"}},
		},
		{
			name:     "empty input",
			input:    "",
			expected: [][]string{{""}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewReader(strings.NewReader(tt.input))
			records, err := r.ReadAll()
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if len(records) != len(tt.expected) {
				t.Errorf("Reader.ReadAll() got %d records, want %d", len(records), len(tt.expected))
				return
			}

			for i := range records {
				if len(records[i]) != len(tt.expected[i]) {
					t.Errorf("Record %d: got %d fields, want %d", i, len(records[i]), len(tt.expected[i]))
					continue
				}
				for j := range records[i] {
					if records[i][j] != tt.expected[i][j] {
						t.Errorf("Record %d field %d: got %q, want %q", i, j, records[i][j], tt.expected[i][j])
					}
				}
			}
		})
	}
}

func TestReader_InputOffset(t *testing.T) {
	input := "hello\x1Fworld\x1Efoo\x1Fbar\x1Ebaz\x1Fqux"
	r := NewReader(strings.NewReader(input))

	expectedOffsets := []int64{0, 12, 20, 27}
	var actualOffsets []int64

	actualOffsets = append(actualOffsets, r.InputOffset())

	record, _ := r.Read()
	if len(record) != 2 || record[0] != "hello" || record[1] != "world" {
		t.Errorf("unexpected first record: %v", record)
	}
	actualOffsets = append(actualOffsets, r.InputOffset())

	record, _ = r.Read()
	if len(record) != 2 || record[0] != "foo" || record[1] != "bar" {
		t.Errorf("unexpected second record: %v", record)
	}
	actualOffsets = append(actualOffsets, r.InputOffset())

	record, _ = r.Read()
	if len(record) != 2 || record[0] != "baz" || record[1] != "qux" {
		t.Errorf("unexpected second record: %v", record)
	}
	actualOffsets = append(actualOffsets, r.InputOffset())

	for i, expected := range expectedOffsets {
		if actualOffsets[i] != expected {
			t.Errorf("at step %d: got offset %d, want %d", i, actualOffsets[i], expected)
		}
	}
}
