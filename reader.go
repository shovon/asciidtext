package asciidtext

import (
	"io"
	"strings"
)

// Reader reads records from an ASCII-delimited text file.
//
// Each record is a slice of fields, separated by the ASCII Unit Separator
// (0x1F). Records are separated by the ASCII Record Separator (0x1E).
//
// A Reader may be created with NewReader.
type Reader struct {
	reader io.Reader
	offset int64
}

// NewReader creates a new Reader that reads from r.
//
// The returned Reader will read records from r, where each record is a slice of
// fields separated by ASCII Unit Separator (0x1F) and records are separated by
// ASCII Record Separator (0x1E).
func NewReader(r io.Reader) *Reader {
	return &Reader{reader: r}
}

// Read reads one record (a slice of fields).
//
// The record is a slice of strings, where each string is a field from the
// input. Fields are separated by ASCII Unit Separator (0x1F) and records are
// separated by ASCII Record Separator (0x1E).
//
// Read returns io.EOF when there is no more data to read. If an error occurs
// during reading, that error will be returned.
func (r *Reader) Read() ([]string, error) {
	b := make([]byte, 1)
	var builder strings.Builder
	for {
		n, err := r.reader.Read(b)
		r.offset++
		if err != nil {
			if err == io.EOF {
				r.offset--
				if builder.Len() > 0 {
					record := strings.Split(builder.String(), string(rune(unitSeparator)))
					return record, nil
				}
				return nil, err
			}
			return nil, err
		}
		if n == 0 {
			continue
		}
		if b[0] == recordSeperator {
			record := strings.Split(builder.String(), string(rune(unitSeparator)))
			return record, nil
		}
		builder.WriteByte(b[0])
	}
}

// InputOffset returns the current offset in the input stream.
//
// The offset represents the number of bytes read from the underlying reader.
// This can be useful for tracking position or resuming reading from a specific
// point.
func (r Reader) InputOffset() int64 {
	return r.offset
}

func (r Reader) ReadAll() ([][]string, error) {
	b, err := io.ReadAll(r.reader)
	if err != nil {
		return nil, err
	}
	str := string(b)

	lines := strings.Split(str, string(rune(recordSeperator)))

	records := [][]string{}

	for _, line := range lines {
		record := strings.Split(line, string(rune(unitSeparator)))
		records = append(records, record)
	}

	return records, nil
}
