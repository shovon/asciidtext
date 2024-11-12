// Package asciidtext provides functionality for reading and writing ASCII-delimited text files.
//
// The package uses ASCII control characters as delimiters: Unit Separator (0x1F) for fields
// within records, and Record Separator (0x1E) for separating records.
//
// Reading Files:
//
// Use Reader to read ASCII-delimited text files:
//
//	r := asciidtext.NewReader(input)
//	record, err := r.Read()        // read single record
//	records, err := r.ReadAll()    // read all records
//
// Writing Files:
//
// Use Writer to write ASCII-delimited text files:
//
//	w := asciidtext.NewWriter(output)
//	w.Write([]string{"field1", "field2"})    // write single record
//	w.WriteAll([][]string{...})              // write multiple records
//	w.Flush()                                // ensure all data is written
//
// The package is useful for handling structured data that needs to be stored or transmitted
// in a text format while preserving field and record boundaries using standard ASCII control
// characters.

package asciidtext
