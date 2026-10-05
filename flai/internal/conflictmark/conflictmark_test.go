package conflictmark

import (
	"slices"
	"testing"
)

// The markers are built here rather than written out, so that no line of
// this file begins with one.
var (
	ours   = "<<<<<<" + "<"
	base   = "||||||" + "|"
	theirs = ">>>>>>" + ">"
	divide = "======" + "="
)

func TestLines(t *testing.T) {
	for _, tc := range []struct {
		name string
		text string
		want []int
	}{
		{"empty", "", nil},
		{"clean markdown", "# Title\n\nSome text.\n", nil},
		{"I-0066's conflict", "# Doc\n" + ours + " HEAD\nmine\n" + divide + "\ntheirs\n" + theirs + " f0f5443 (docs: [S-0201] x)\nafter\n", []int{2, 6}},
		{"diff3 base", ours + " HEAD\na\n" + base + " merged common ancestors\nb\n" + divide + "\nc\n" + theirs + " topic\n", []int{1, 3, 7}},
		{"bare markers", ours + "\n" + theirs + "\n", []int{1, 2}},
		{"CRLF line endings", "x\r\n" + ours + " HEAD\r\ny\r\n" + theirs + "\r\n", []int{2, 4}},
		{"divider alone is a setext underline", "Heading\n" + divide + "\n", nil},
		{"eight is not a marker", ours + "<\n" + theirs + ">\n" + base + "|\n", nil},
		{"six is not a marker", "<<<<<<\n>>>>>>\n", nil},
		{"text straight after is not a marker", ours + "HEAD\n" + theirs + "x\n", nil},
		{"indented is not a marker", " " + ours + " HEAD\n> " + theirs + " x\n", nil},
		{"tab after is not a marker", ours + "\tHEAD\n", nil},
		{"no final newline", "a\n" + theirs + " b", []int{2}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := Lines(tc.text); !slices.Equal(got, tc.want) {
				t.Errorf("Lines(%q) = %v, want %v", tc.text, got, tc.want)
			}
		})
	}
}

func TestIsMarker(t *testing.T) {
	for _, tc := range []struct {
		line string
		want bool
	}{
		{ours, true},
		{ours + " HEAD", true},
		{base + " base", true},
		{theirs + " abc1234 (msg)", true},
		{divide, false},
		{ours + "<", false},
		{"x" + ours, false},
		{"", false},
	} {
		if got := IsMarker(tc.line); got != tc.want {
			t.Errorf("IsMarker(%q) = %v, want %v", tc.line, got, tc.want)
		}
	}
}
