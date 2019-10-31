package ctags

import (
	"bytes"
	"strings"
	"testing"
)

var ctagsTestFormat = `
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

func TestReadWrite(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			`ctags format example 1`,
			`asdf	sub.cc	/^asdf()$/;"	new_field:some\svalue	file:` + "\n",
			`asdf	sub.cc	/^asdf()$/;"	file:	new_field:some\\svalue` + "\n",
		},
		{
			`ctags format example 2`,
			`foo_t	sub.h	/^typedef foo_t$/;"	kind:t` + "\n",
			`foo_t	sub.h	/^typedef foo_t$/;"	kind:t` + "\n",
		},
		{
			`ctags format example 3`,
			`func3	sub.p	/^func3()$/;"	function:/func1/func2	file:` + "\n",
			`func3	sub.p	/^func3()$/;"	file:	function:/func1/func2` + "\n",
		},
		{
			`ctags format example 4`,
			`getflag	sub.c	/^getflag(arg)$/;"	kind:f	file:` + "\n",
			`getflag	sub.c	/^getflag(arg)$/;"	file:	kind:f` + "\n",
		},
		{
			`ctags format example 5`,
			`inc	sub.cc	/^inc()$/;"	file: class:PipeBuf` + "\n",
			`inc	sub.cc	/^inc()$/;"	file: class:PipeBuf` + "\n",
		},
		{
			`tagfield omitted "kind:" name 1`,
			`foo_t	sub.h	/^typedef foo_t$/;"	t` + "\n",
			`foo_t	sub.h	/^typedef foo_t$/;"	kind:t` + "\n",
		},
		{
			`tagfield omitted "kind:" name 2`,
			`getflag	sub.c	/^getflag(arg)$/;"	f	file:` + "\n",
			`getflag	sub.c	/^getflag(arg)$/;"	file:	kind:f` + "\n",
		},
		{
			`tagfield invalid name characters`,
			`foo	foo	1;"	invalid-name:value	validname:value` + "\n",
			`foo	foo	1;"	validname:value` + "\n",
		},
		{
			`tagfield empty name and value`,
			`foo	foo	1;"	:	foo:bar` + "\n",
			`foo	foo	1;"	foo:bar` + "\n",
		},
		{
			// tagaddress contains a backslash ("\\") followed by the character "t".
			// tagfields contains a newline ("\n"), backslash ("\\"), and tab ("\t")
			`escape characters`,
			`foo	foo	/^\\t$/;"	foo:bar\n\\\t` + "\n",
			`foo	foo	/^\\t$/;"	foo:bar\n\\\t` + "\n",
		},
	}

	for _, tc := range cases {
		var b bytes.Buffer

		r := NewReader(strings.NewReader(tc.input))
		w := NewWriter(&b)
		w.WriteAll(r.ReadAll())
		if actual := b.String(); actual != tc.expected {
			t.Errorf(ctagsTestFormat, tc.name, tc.input, tc.expected, actual)
		}
	}
}

func TestLine(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected int
	}{
		{
			`line from tagaddress`,
			`foo	foo	99;"`,
			99,
		},
		{
			`line from tagfield`,
			`foo	foo	/^foo$/;"	line:100`,
			100,
		},
		{
			`no line`,
			`foo	foo	/^foo$/;"`,
			-1,
		},
	}

	for _, tc := range cases {
		r := NewReader(strings.NewReader(tc.input))
		tags := r.ReadAll()
		if actual := tags[0].Line(); actual != tc.expected {
			t.Errorf(ctagsTestFormat, tc.name, tc.input, tc.expected, actual)
		}
	}
}

func TestKind(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			`kind from tagfield`,
			`foo	foo	/^foo$/;"	kind:v`,
			`v`,
		},
		{
			`kind from tagfield without "kind:" name`,
			`foo	foo	/^foo$/;"	v`,
			`v`,
		},
		{
			`no kind`,
			`foo	foo	/^foo$/;"`,
			``,
		},
	}

	for _, tc := range cases {
		r := NewReader(strings.NewReader(tc.input))
		tags := r.ReadAll()
		if actual := tags[0].Kind(); actual != tc.expected {
			t.Errorf(ctagsTestFormat, tc.name, tc.input, tc.expected, actual)
		}
	}
}
