package acl

import (
	"fmt"
	"io"

	"github.com/uhppoted/uhppote-core/uhppote"
)

func MakeFlatFile(acl ACL, devices []uhppote.Device, f io.Writer) error {
	t, err := MakeTable(acl, devices)
	if err != nil {
		return err
	}

	formats := make([]string, len(t.Header))

	for i, h := range t.Header {
		width := len(h)
		for _, r := range t.Records {
			if len(r[i]) > width {
				width = len(r[i])
			}
		}

		formats[i] = fmt.Sprintf("%%-%ds", width)
	}

	separator := ""
	for i, h := range t.Header {
		fmt.Fprintf(f, "%s", separator)
		fmt.Fprintf(f, formats[i], h)
		separator = "  "
	}
	fmt.Fprintln(f)

	for _, r := range t.Records {
		separator := ""
		for i, v := range r {
			fmt.Fprintf(f, "%s", separator)
			fmt.Fprintf(f, formats[i], v)
			separator = "  "
		}
		fmt.Fprintln(f)
	}

	return nil
}
