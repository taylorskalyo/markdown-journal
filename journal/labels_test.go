package journal

import (
	"bytes"
	"strings"
	"testing"

	"github.com/taylorskalyo/markdown-journal/ctags"
)

func TestWriteLabels(t *testing.T) {
	format := `
============= case %s ================
Ctags Input:
-----------
%v
Expected Output:
----------
%v
Actual Output:
----------
%v
`

	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			`basics`,
			`
03 Tuesday	diary/2006-01-03.md	1;"	kind:title	line:1
recipe	diary/2006-01-03.md	5;"	heading:03 Tuesday	kind:label	line:5
30 Friday	diary/2007-11-30.md	1;"	kind:title	line:1
recipe	diary/2007-11-30.md	3;"	heading:30 Friday	kind:label	line:3
groceries	diary/2007-11-30.md	14;"	heading:Ingredients	kind:label	line:14
groceries	diary/2007-11-30.md	16;"	heading:Ingredients	kind:label	line:16
			`,
			`
# groceries
* [Ingredients](diary/2007-11-30.md:16)
* [Ingredients](diary/2007-11-30.md:14)

# recipe
* [30 Friday](diary/2007-11-30.md:3)
* [03 Tuesday](diary/2006-01-03.md:5)
			`,
		},
		{
			`no heading`,
			`
recipe	diary/2006-01-04.md	5;"	kind:label	line:5
			`,
			`
# recipe
* [diary/2006-01-04.md:5](diary/2006-01-04.md:5)
			`,
		},
	}

	for _, tc := range cases {
		var b bytes.Buffer

		r := ctags.NewReader(strings.NewReader(tc.input))
		j := NewJournal(r.ReadAll())
		j.WriteLabels(&b)
		actual := strings.TrimSpace(b.String())
		expected := strings.TrimSpace(tc.expected)
		if actual != expected {
			t.Errorf(format, tc.name, tc.input, expected, actual)
		}
	}
}
