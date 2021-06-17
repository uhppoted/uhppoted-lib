package acl

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
)

type Table struct {
	Header  []string
	Records [][]string
}

func (table *Table) ToTSV(f io.Writer) error {
	w := csv.NewWriter(f)
	w.Comma = '\t'

	w.Write(table.Header)
	for _, row := range table.Records {
		w.Write(row)
	}

	w.Flush()

	return nil
}

func (table *Table) MarshalText() []byte {
	return table.MarshalTextIndent("", " ")
}

func (table *Table) MarshalTextIndent(indent, gap string) []byte {
	var b bytes.Buffer

	widths := make([]int, len(table.Header))
	for i, field := range table.Header {
		if len(field) > widths[i] {
			widths[i] = len(field)
		}
	}

	for _, row := range table.Records {
		for i, field := range row {
			if len(field) > widths[i] {
				widths[i] = len(field)
			}
		}
	}

	for i := 1; i < len(widths); i++ {
		widths[i-1] += len(gap)
	}

	fmt.Fprintf(&b, "%s", indent)
	for i, field := range table.Header {
		fmt.Fprintf(&b, "%-*v", widths[i], field)
	}
	fmt.Fprintln(&b)

	for _, row := range table.Records {
		fmt.Fprintf(&b, "%s", indent)
		for i, field := range row {
			fmt.Fprintf(&b, "%-*v", widths[i], field)
		}
		fmt.Fprintln(&b)
	}

	return b.Bytes()
}
