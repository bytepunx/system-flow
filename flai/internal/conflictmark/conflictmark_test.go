package conflictmark

import (
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/execx"
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

// Branch reads only what the branch adds or changes since it left the base,
// numbering lines as the files hold them, and skips binary files (S-0276).
func TestBranch(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	root := t.TempDir()
	git := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = root
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	write := func(name, text string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(root, name), []byte(text), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	git("init", "-q", "-b", "main")
	git("config", "user.email", "t@t")
	git("config", "user.name", "t")
	write("main.md", ours+" HEAD\non main only\n"+theirs+" x\n")
	write("changed.md", "# Changed\n")
	write("gone.md", ours+" HEAD\n")
	git("add", "-A")
	git("commit", "-q", "-m", "init")
	git("checkout", "-q", "-b", "story/S-0001")
	write("added.md", "\n# Added\n"+ours+" HEAD\na\n"+divide+"\nb\n"+theirs+" topic\n")
	write("changed.md", "# Changed\n\n"+theirs+" topic\n")
	write("blob.bin", ours+" HEAD\n\x00\n")
	write("clean.md", "# Clean\n")
	git("rm", "-q", "gone.md")
	git("add", "-A")
	git("commit", "-q", "-m", "branch")

	got, err := Branch(execx.System{}, root, "main", "story/S-0001")
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"added.md:3", "added.md:7", "changed.md:3"}; !slices.Equal(got, want) {
		t.Errorf("Branch = %v, want %v", got, want)
	}

	git("checkout", "-q", "main")
	git("checkout", "-q", "-b", "story/S-0002")
	write("clean.md", "# Clean\n")
	git("add", "-A")
	git("commit", "-q", "-m", "clean")
	if got, err := Branch(execx.System{}, root, "main", "story/S-0002"); err != nil || len(got) != 0 {
		t.Errorf("a clean branch: %v, %v", got, err)
	}
}
