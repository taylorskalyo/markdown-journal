package journal

import (
	"bytes"
	"strings"
	"testing"

	"github.com/taylorskalyo/markdown-journal/ctags"
)

func TestParse(t *testing.T) {
	format := `
============= case %s ================
Markdown Input:
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
		filename string
		input    string
		expected string
	}{
		{
			`basics`,
			`2006-01-02.md`,
			`
# Foo

:bar:
			`,
			`
Foo	2006-01-02.md	2;"	kind:heading	line:2
bar	2006-01-02.md	4;"	kind:label	line:4
			`,
		},
		{
			`ignore labels in codefence`,
			`2006-01-02.md`,
			"# Foo\n```\n:bar:\n```",
			`
Foo	2006-01-02.md	1;"	kind:heading	line:1
			`,
		},
		{
			`capture labels inside headings`,
			`2006-01-02.md`,
			`
# Foo :bar:
			`,
			`
Foo bar	2006-01-02.md	2;"	kind:heading	line:2
bar	2006-01-02.md	2;"	kind:label	line:2
			`,
		},
	}

	for _, tc := range cases {
		var b bytes.Buffer

		p := NewEntryParser()
		lines, _ := p.parse(tc.filename, []byte(tc.input))
		w := ctags.NewWriter(&b)
		w.WriteAll(lines)

		actual := strings.TrimSpace(b.String())
		expected := strings.TrimSpace(tc.expected)
		if actual != expected {
			t.Errorf(format, tc.name, tc.input, expected, actual)
		}
	}
}
