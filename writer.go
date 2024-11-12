package asciidtext

import (
	"fmt"
	"io"
	"strings"
)

type Writer struct {
	writer io.Writer
	buffer []string
	err    error
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: w, buffer: []string{}}
}

func (w *Writer) Write(record []string) error {
	w.buffer = append(w.buffer, strings.Join(record, string(rune(unitSeparator))))
	return nil
}

func (w *Writer) WriteAll(records [][]string) error {
	for _, record := range records {
		err := w.Write(record)
		if err != nil {
			return err
		}
	}
	w.Flush()
	return w.err
}

func (w *Writer) Flush() {
	toWrite := []byte(strings.Join(w.buffer, string(rune(recordSeperator))))
	expectedLen := len(toWrite)
	writtenLength, err := w.writer.Write(toWrite)
	if err != nil {
		w.err = err
	}
	if writtenLength != expectedLen {
		w.err = fmt.Errorf("something went wrong. Expected to have written %d bytes but actually wrote %d", expectedLen, writtenLength)
	}
}

func (w *Writer) Error() error {
	return w.err
}
