package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPrime(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	dir := "../internal/metrics/testdata/good"
	out, errOut, code := runIn(t, dir, "prime")
	if code != 0 {
		t.Fatalf("prime: %s", errOut)
	}
	want := "design/conventions/README.md\ndesign/conventions/session-start.md\ndesign/conventions/git.md\n"
	if out != want {
		t.Errorf("paths:\n%s", out)
	}
	out, _, _ = runIn(t, dir, "prime", "--cat")
	if !strings.Contains(out, "design/conventions/git.md\n===") || !strings.Contains(out, "- Commit at landing.") || strings.Index(out, "Session start") > strings.Index(out, "# Git") {
		t.Errorf("cat:\n%s", out)
	}
	out, _, _ = runIn(t, dir, "prime", "--json")
	var v struct {
		Files []struct {
			Name  string `json:"name"`
			Order int    `json:"order"`
		} `json:"files"`
	}
	if err := json.Unmarshal([]byte(out), &v); err != nil || len(v.Files) != 2 || v.Files[1].Order != 20 {
		t.Errorf("json: %v %s", err, out)
	}
	_, errOut, code = runIn(t, tempProject(t), "prime")
	if code == 0 || !strings.Contains(errOut, "no conventions folder") {
		t.Errorf("missing folder should fail: %d %s", code, errOut)
	}
}

func TestPrimeStory(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	dir := "../internal/metrics/testdata/good"
	cat, _, _ := runIn(t, dir, "prime", "--cat")
	out, errOut, code := runIn(t, dir, "prime", "--story", "S-4")
	if code != 0 {
		t.Fatalf("prime --story: %s", errOut)
	}
	i := strings.Index(out, "contents, in bytes:\n")
	n := strings.Index(out[max(i, 0):], "\n\n")
	ok := i >= 0 && n >= 0
	head, body := out, ""
	if ok {
		head, body = out[:i+n+2], out[i+n+2:]
	}
	if !ok || !strings.HasPrefix(head, "S-004 context pack\n") || !strings.Contains(head, "story: S-004 Four\n") || !strings.Contains(head, "  all  every story\n") {
		t.Fatalf("header:\n%s", out)
	}
	if !strings.HasPrefix(body, cat) || !strings.HasSuffix(body, "\ndesign/system/overview.md: Overview (70 bytes; topics: all)\n") {
		t.Errorf("with every convention at [all] the pack is prime --cat, then the design:\n%s\n---\n%s", body, cat)
	}
	if !strings.Contains(head, fmt.Sprintf("size: %d bytes, %d lines", len(out), strings.Count(out, "\n"))) {
		t.Errorf("size:\n%s", head)
	}
	for _, id := range []string{"E-1", "S-0999"} {
		_, errOut, code = runIn(t, dir, "prime", "--story", id)
		if code == 0 || !strings.Contains(errOut, "flai prime --story "+id) {
			t.Errorf("%s should fail naming it: %d %s", id, code, errOut)
		}
	}
	if out, errOut, code = runIn(t, dir, "prime", "--story", "S-1"); code != 0 || !strings.HasPrefix(out, "S-001 context pack\n") {
		t.Errorf("an archived story gets a pack: %d %s", code, errOut)
	}

	root := tempProject(t)
	conv := filepath.Join(root, "design", "conventions")
	if err := os.MkdirAll(conv, 0o755); err != nil {
		t.Fatal(err)
	}
	write := func(name, s string) {
		if err := os.WriteFile(filepath.Join(conv, name), []byte(s), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("README.md", "# Conventions\n\n- [code-quality.md](code-quality.md)\n")
	write("code-quality.md", "---\ntitle: Code quality\nupdated: 2026-09-28\naudience: agent\norder: 60\nstatus: active\ntopics: [all]\n---\n\n# Code quality\n\n## Rules\n\n- Tests accompany the change.\n\n### Go <!-- topics: go -->\n\n- Standard library testing only.\n\n### Svelte <!-- topics: svelte -->\n\n- vitest.\n\n<!-- system-flow:end-of-baseline -->\n\n## Project additions <!-- topics: svelte -->\n\n- flaiover runs in Docker.\n")
	if _, errOut, code = runIn(t, root, "epic", "new", "Epic"); code != 0 {
		t.Fatal(errOut)
	}
	if _, errOut, code = runIn(t, root, "story", "new", "Slice", "--epic", "E-0001", "--topics", "go"); code != 0 {
		t.Fatal(errOut)
	}
	out, errOut, code = runIn(t, root, "prime", "--story", "S-0001")
	if code != 0 {
		t.Fatalf("prime --story: %s", errOut)
	}
	for _, want := range []string{"  go   own\n", "  left out    2 convention sections, listed at the end\n", "- Standard library testing only.", "<!-- system-flow:end-of-baseline -->\n\n## Project additions <!-- topics: svelte -->\n",
		"\nleft out\n========\n\n- code-quality.md § Rules › Svelte (svelte)\n- code-quality.md § Project additions (svelte)\n"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q:\n%s", want, out)
		}
	}
	for _, gone := range []string{"vitest", "Docker"} {
		if strings.Contains(out, gone) {
			t.Errorf("kept %q:\n%s", gone, out)
		}
	}
	out, _, _ = runIn(t, root, "prime", "--story", "S-0001", "--json")
	var v struct {
		Story  string `json:"story"`
		Topics []struct {
			Topic   string `json:"topic"`
			Sources []struct {
				Kind string `json:"kind"`
			} `json:"sources"`
		} `json:"topics"`
		Size        struct{ Bytes int } `json:"size"`
		Conventions []struct {
			Path    string                     `json:"path"`
			Kept    []struct{ Heading string } `json:"kept"`
			LeftOut []struct {
				Heading string   `json:"heading"`
				Topics  []string `json:"topics"`
			} `json:"left_out"`
		} `json:"conventions"`
		LeftOut []string `json:"left_out"`
	}
	if err := json.Unmarshal([]byte(out), &v); err != nil {
		t.Fatalf("json: %v %s", err, out)
	}
	if v.Story != "S-0001" || len(v.Topics) != 2 || v.Topics[0].Topic != "go" || v.Topics[0].Sources[0].Kind != "own" || v.Size.Bytes == 0 ||
		len(v.Conventions) != 1 || len(v.Conventions[0].Kept) != 3 || len(v.Conventions[0].LeftOut) != 2 ||
		v.Conventions[0].LeftOut[0].Heading != "Svelte" || v.Conventions[0].LeftOut[0].Topics[0] != "svelte" || len(v.LeftOut) != 2 {
		t.Errorf("json: %+v", v)
	}
}

func TestPrimeStoryDesign(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	write := func(rel, s string) {
		p := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(s), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("design/conventions/README.md", "# Conventions\n")
	write("design/system/cli.md", "---\ntitle: CLI\nupdated: 2026-09-28\nstatus: active\ntopics: [go]\n---\n\n# CLI\n\n## Commands\n\nAs ADR-0002 decides.\n\n## Board <!-- topics: dashboard -->\n\nThe board.\n")
	write("design/system/other.md", "---\ntitle: Other\nupdated: 2026-09-28\nstatus: active\n---\n\n# Other\n\nNothing.\n")
	adr := func(n, title, fm string) {
		write("design/adrs/"+n+".md", "---\nid: ADR-"+n[:4]+"\ntitle: "+title+"\nstatus: accepted\ndate: 2026-09-01\n"+fm+"---\n\n# ADR-"+n[:4]+" "+title+"\n\n## Decision\n\n"+title+".\n")
	}
	adr("0001-old", "Old", "superseded_by: [ADR-0003]\n")
	adr("0002-linked", "Linked", "")
	adr("0003-new", "New", "")
	if _, errOut, code := runIn(t, root, "epic", "new", "Epic"); code != 0 {
		t.Fatal(errOut)
	}
	if _, errOut, code := runIn(t, root, "story", "new", "Slice", "--epic", "E-0001", "--topics", "go"); code != 0 {
		t.Fatal(errOut)
	}
	story, err := filepath.Glob(filepath.Join(root, "wip", "kanban", "stories", "S-0001-*.md"))
	if err != nil || len(story) != 1 {
		t.Fatal(story, err)
	}
	data, err := os.ReadFile(story[0])
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(story[0], []byte(strings.Replace(string(data), "## Goal\n", "## Goal\n\nKeep to ADR-0001.\n", 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	out, errOut, code := runIn(t, root, "prime", "--story", "S-0001")
	if code != 0 {
		t.Fatalf("prime --story: %s", errOut)
	}
	for _, want := range []string{
		"size: ",
		"; budget 81920 bytes\n",
		"  named       design/adrs/0003-new.md\n",
		"  briefed     design/system/cli.md\n",
		"  briefed     1 ADR by its decision sentence, on one line\n",
		"  ranked      design/system/cli.md § Commands\n",
		"\ndesign/adrs/0003-new.md\n=======================\nreason: supersedes ADR-0001\n\n---\nid: ADR-0003\n",
		"\nbriefs\n======\n",
		"\ndesign/system/cli.md: CLI (",
		" bytes; topics: go)\n\nAs ADR-0002 decides.\n\n- Commands (selected, loaded)\n- Board\n",
		"\ndecisions\n=========\n",
		"- ADR-0002 Linked (design/adrs/0002-linked.md; linked from design/system/cli.md): Linked.\n",
		"\ndesign/system/cli.md § Commands\n",
		"reason: rank 1\n\n## Commands\n\nAs ADR-0002 decides.\n",
		"\ncatalog\n=======\n",
		"- design/system/other.md: Other\n",
		"- design/adrs/0001-old.md: Old (superseded by ADR-0003)\n",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "The board.") || strings.Contains(out, "over budget") {
		t.Errorf("printed a section whose topics miss, or is over budget:\n%s", out)
	}
	if !strings.Contains(out, fmt.Sprintf("size: %d bytes, %d lines, this header included", len(out), strings.Count(out, "\n"))) {
		t.Errorf("size is not the pack's (%d bytes):\n%s", len(out), out)
	}

	out, _, _ = runIn(t, root, "prime", "--story", "S-0001", "--json")
	var v struct {
		Size   struct{ Bytes int } `json:"size"`
		Budget int                 `json:"budget"`
		Items  []struct {
			Path    string   `json:"path"`
			ID      string   `json:"id"`
			Heading []string `json:"heading"`
			Step    string   `json:"step"`
			Reason  string   `json:"reason"`
			Size    int      `json:"size"`
			Whole   int      `json:"whole"`
			Text    string   `json:"text"`
		} `json:"items"`
		Catalog struct {
			NotLoaded []struct{ Path string } `json:"not_loaded"`
			InPart    []struct {
				Path    string `json:"path"`
				Outline []struct {
					Heading string `json:"heading"`
					Loaded  bool   `json:"loaded"`
				} `json:"outline"`
			} `json:"in_part"`
		} `json:"catalog"`
	}
	if err := json.Unmarshal([]byte(out), &v); err != nil {
		t.Fatalf("json: %v %s", err, out)
	}
	if len(v.Items) != 4 || v.Items[0].Path != "design/adrs/0003-new.md" || v.Items[0].Step != "named" || v.Items[0].Reason != "supersedes ADR-0001" ||
		v.Items[1].Path != "design/system/cli.md" || v.Items[1].Step != "briefed" || v.Items[1].Whole == 0 || v.Items[1].Size != len(v.Items[1].Text) ||
		v.Items[2].ID != "ADR-0002" || v.Items[2].Text != "Linked." || v.Items[3].Step != "ranked" || strings.Join(v.Items[3].Heading, "/") != "CLI/Commands" ||
		len(v.Catalog.NotLoaded) != 2 || len(v.Catalog.InPart) != 0 || v.Budget != 81920 {
		t.Errorf("json: %+v", v)
	}
	if v.Size.Bytes <= v.Items[0].Size+v.Items[1].Size+v.Items[2].Size {
		t.Errorf("the size leaves out the items: %d", v.Size.Bytes)
	}

	out, _, code = runIn(t, root, "prime", "--story", "S-0001", "--budget", "1")
	if code != 0 || !strings.Contains(out, "; budget 1 bytes\nover budget: the conventions alone exceed it") || strings.Contains(out, "reason:") {
		t.Errorf("--budget 1 (%d):\n%s", code, out)
	}
	mf := filepath.Join(root, "system-flow.yaml")
	data, err = os.ReadFile(mf)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(mf, append(data, []byte("prime:\n  budget: 2KB\n")...), 0o644); err != nil {
		t.Fatal(err)
	}
	if out, _, _ = runIn(t, root, "prime", "--story", "S-0001"); !strings.Contains(out, "; budget 2048 bytes\n") {
		t.Errorf("prime.budget not the default:\n%s", out)
	}
	if out, _, _ = runIn(t, root, "prime", "--story", "S-0001", "--budget", "100KB"); !strings.Contains(out, "; budget 102400 bytes\n") {
		t.Errorf("--budget does not win over prime.budget:\n%s", out)
	}
	for _, args := range [][]string{{"prime", "--story", "S-0001", "--budget", "lots"}, {"prime", "--budget", "80KB"}} {
		if _, errOut, code := runIn(t, root, args...); code == 0 || !strings.Contains(errOut, "budget") {
			t.Errorf("%v: code %d, %s", args, code, errOut)
		}
	}
}
