package journal

import (
	"bytes"
	"strings"
	"testing"

	"github.com/taylorskalyo/markdown-journal/ctags"
)

var format = `
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
			`example 1`,
			`
02 Monday	diary/2006-01-02.md	1;"	kind:heading	line:1
03 Tuesday	diary/2006-01-03.md	1;"	kind:heading	line:1
recipe	diary/2006-01-03.md	5;"	kind:label	line:5
05 Sunday	diary/2006-02-05.md	1;"	kind:heading	line:1
30 Friday	diary/2007-11-30.md	1;"	kind:heading	line:1
			`,
			`
# 2007

## November
* [30 Fri](diary/2007-11-30.md) - 30 Friday

# 2006

## February
* [05 Sun](diary/2006-02-05.md) - 05 Sunday

## January
* [03 Tue](diary/2006-01-03.md) - 03 Tuesday
* [02 Mon](diary/2006-01-02.md) - 02 Monday
			`,
		},
		{
			`no heading`,
			`
recipe	diary/2006-01-04.md	5;"	kind:label	line:5
			`,
			`
# 2006

## January
* [04 Wed](diary/2006-01-04.md)
			`,
		},
		{
			`heading and label on same line`,
			`
Tantanmen :recipe:	diary/2006-01-03.md	1;"	kind:heading	line:1
recipe	diary/2006-01-03.md	1;"	kind:label	line:1
			`,
			`
# 2006

## January
* [03 Tue](diary/2006-01-03.md) - Tantanmen :recipe:
			`,
		},
	}

	for _, tc := range cases {
		var b bytes.Buffer

		r := ctags.NewReader(strings.NewReader(tc.input))
		j := NewJournal(r.ReadAll())
		j.WriteTimeline(&b)
		actual := strings.TrimSpace(b.String())
		expected := strings.TrimSpace(tc.expected)
		if actual != expected {
			t.Errorf(format, tc.name, tc.input, expected, actual)
		}
	}
}
