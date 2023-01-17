package acl

import (
	"bytes"
	"testing"
)

func TestTableMarshalText(t *testing.T) {
	table := Table{
		Header: []string{"Column 1", "Column 2", "Column 3"},
		Records: [][]string{
			[]string{"Row 1.1", "Row 1.2", "Row 1.3"},
			[]string{"Row 2.1", "Row 2.2", "Row 2.3"},
			[]string{"Row 3.1", "Row 3.2.4.5.6.7.8.9", "Row 3.3"},
			[]string{"Row 4.1", "Row 4.2", "Row 4.3"},
		},
	}

	expected := `> Column 1  Column 2             Column 3
> Row 1.1   Row 1.2              Row 1.3 
> Row 2.1   Row 2.2              Row 2.3 
> Row 3.1   Row 3.2.4.5.6.7.8.9  Row 3.3 
> Row 4.1   Row 4.2              Row 4.3 
`

	bytes := table.MarshalTextIndent("> ", "  ")

	if string(bytes) != expected {
		t.Errorf("MarshalIndent produced incorrect output:\n  expected:\n%v\n  got:\n%v", expected, string(bytes))
	}
}

func TestTableToCSV(t *testing.T) {
	table := Table{
		Header: []string{"Column 1", "Column 2", "Column 3"},
		Records: [][]string{
			[]string{"Row 1.1", "Row 1.2", "Row 1.3"},
			[]string{"Row 2.1", "Row 2.2", "Row 2.3"},
			[]string{"Row 3.1", "Row 3.2.4.5.6.7.8.9", "Row 3.3"},
			[]string{"Row 4.1", "Row 4.2", "Row 4.3"},
		},
	}

	expected := `Column 1	Column 2	Column 3
Row 1.1	Row 1.2	Row 1.3
Row 2.1	Row 2.2	Row 2.3
Row 3.1	Row 3.2.4.5.6.7.8.9	Row 3.3
Row 4.1	Row 4.2	Row 4.3
`

	var b bytes.Buffer

	err := table.ToTSV(&b)
	if err != nil {
		t.Fatalf("Unexpected error (%v)", err)
	}

	if b.String() != expected {
		t.Errorf("ToTSV produced incorrect output:\n  expected:\n%v\n  got:\n%v", expected, b.String())
		t.Errorf("ToTSV produced incorrect output:\n  expected:\n%v\n  got:\n%v", []byte(expected), b.Bytes())
	}
}
