package cmd

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"
)

func TestRenderReferenceTree(t *testing.T) {
	root := &cobra.Command{Use: "tool", Short: "A tool"}
	root.PersistentFlags().Bool("json", false, "print JSON")
	group := &cobra.Command{Use: "item", Short: "Work with items"}
	group.PersistentFlags().String("project", "", "the project")
	show := &cobra.Command{
		Use:     "show <id>",
		Short:   "Show one item",
		Aliases: []string{"get"},
		Long: `Print the item <id> with its
children, one per line.

  tool item show S-1
  tool item show S-1 --all

- first point
- second point`,
		Example: `  tool item show S-1`,
		Run:     func(*cobra.Command, []string) {},
	}
	show.Flags().Int("depth", 2, "how deep | to go")
	show.Flags().Bool("all", false, "every child")
	secret := &cobra.Command{Use: "secret", Short: "Hidden", Hidden: true, Run: func(*cobra.Command, []string) {}}
	group.AddCommand(show, secret)
	root.AddCommand(group)

	got := renderReference(root)
	for _, want := range []string{
		"| [item](#tool-item) | Work with items |",
		"## tool\n",
		"| `--json` | print JSON |",
		"### tool item\n",
		"- [show](#tool-item-show): Show one item\n",
		"#### tool item show\n\nShow one item.\n\n```text\ntool item show <id> [flags]\n```",
		"Print the item &lt;id&gt; with its children, one per line.\n",
		"```text\ntool item show S-1\ntool item show S-1 --all\n```",
		"- first point\n- second point\n",
		"Aliases: `get`.",
		"```bash\ntool item show S-1\n```",
		"| `--depth` int | how deep \\| to go (default `2`) |",
		"| `--all` | every child |",
		"Flags from `tool item`:\n\n| Flag | Meaning |\n|------|---------|\n| `--project` string | the project |",
		"completion bash",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("reference lacks %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "secret") || strings.Contains(got, "tool help") {
		t.Errorf("reference lists a hidden command:\n%s", got)
	}
	if strings.Contains(got, "\n\n\n") {
		t.Errorf("reference has a double blank line:\n%s", got)
	}
	if strings.Count(got, "`--json`") != 1 {
		t.Errorf("global flags should be listed once, at the root:\n%s", got)
	}
}

func TestReferenceCoversEveryCommand(t *testing.T) {
	root := newRootCmd(io.Discard, io.Discard)
	got := renderReference(root)
	var walk func(c *cobra.Command)
	walk = func(c *cobra.Command) {
		for _, sub := range availableCommands(c) {
			if !strings.Contains(got, "# "+sub.CommandPath()+"\n") {
				t.Errorf("no heading for %s", sub.CommandPath())
			}
			if sub.Short == "" {
				t.Errorf("%s has no short description", sub.CommandPath())
			}
			walk(sub)
		}
	}
	walk(root)
	if strings.Contains(got, "# flai reference\n") {
		t.Error("the hidden reference command documents itself")
	}
}

func TestReferenceWriteKeepsDateWhenUnchanged(t *testing.T) {
	path := filepath.Join(t.TempDir(), "flai-reference.md")
	day1 := time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)
	day2 := time.Date(2026, 9, 2, 12, 0, 0, 0, time.UTC)
	if _, errOut, code := runInAt(t, ".", day1, "reference", "--write", path); code != 0 {
		t.Fatalf("write: %s", errOut)
	}
	first, _ := os.ReadFile(path)
	if !strings.Contains(string(first), "updated: 2026-09-01\n") {
		t.Fatalf("front matter lacks the date:\n%.200s", first)
	}
	out, _, _ := runInAt(t, ".", day2, "reference", "--write", path)
	if again, _ := os.ReadFile(path); string(again) != string(first) || !strings.Contains(out, "up to date") {
		t.Fatalf("an unchanged reference was rewritten: %s", out)
	}
	_ = os.WriteFile(path, []byte(strings.Replace(string(first), "Print the kanban board", "stale", 1)), 0o644)
	runInAt(t, ".", day2, "reference", "--write", path)
	if again, _ := os.ReadFile(path); !strings.Contains(string(again), "updated: 2026-09-02\n") || strings.Contains(string(again), "stale") {
		t.Fatalf("a stale reference was not rewritten with the new date:\n%.200s", again)
	}
}

// The committed reference is what the help says now; make flai-reference
// regenerates it.
func TestReferenceIsCurrent(t *testing.T) {
	if testing.Short() {
		t.Skip("reads the monorepo's docs")
	}
	page, err := os.ReadFile("../../docs/users/flai-reference.md")
	if err != nil {
		t.Fatalf("read the reference: %v", err)
	}
	if referenceBody(string(page)) != renderReference(newRootCmd(io.Discard, io.Discard)) {
		t.Fatal("docs/users/flai-reference.md is stale: run make flai-reference")
	}
}
