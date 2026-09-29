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
	head, body, ok := strings.Cut(out, " lines below this header\n\n")
	if !ok || !strings.HasPrefix(head, "S-004 context pack\n") || !strings.Contains(head, "story: S-004 Four\n") || !strings.Contains(head, "  all  every story\n") {
		t.Fatalf("header:\n%s", out)
	}
	if !strings.HasPrefix(body, cat) || !strings.HasSuffix(body, "\ndesign/system/overview.md\n=========================\nreason: topics: all\n\n---\ntitle: Overview\nupdated: 2026-08-01\ntopics: [all]\n---\n\n# Overview\n") {
		t.Errorf("with every convention at [all] the pack is prime --cat, then the design:\n%s\n---\n%s", body, cat)
	}
	if !strings.Contains(head, fmt.Sprintf("size: %d bytes, %d", len(body), strings.Count(body, "\n"))) {
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
	for _, want := range []string{"  go   own\n", "left out: 2 sections, listed at the end\n", "- Standard library testing only.", "<!-- system-flow:end-of-baseline -->\n\n## Project additions <!-- topics: svelte -->\n",
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
		"design: 3 items, ",
		"catalog: 2 documents not loaded, 1 loaded in part\n",
		"\ndesign/system/cli.md § CLI\n",
		"reason: topics: go\n\n# CLI\n\n## Commands\n\nAs ADR-0002 decides.\n\n\n",
		"\ndesign/adrs/0003-new.md\n=======================\nreason: supersedes ADR-0001\n\n---\nid: ADR-0003\n",
		"\ndesign/adrs/0002-linked.md\n==========================\nreason: linked from design/system/cli.md § Commands\n",
		"\ncatalog\n=======\n",
		"- design/system/other.md: Other\n",
		"- design/adrs/0001-old.md: Old (superseded by ADR-0003)\n",
		"- design/system/cli.md: CLI\n  - CLI (loaded)\n    - Commands (loaded)\n    - Board\n",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "The board.") {
		t.Errorf("printed a section whose topics miss:\n%s", out)
	}
	head, body, _ := strings.Cut(out, " lines below this header\n\n")
	if !strings.Contains(head, fmt.Sprintf("size: %d bytes, %d", len(body), strings.Count(body, "\n"))) {
		t.Errorf("size:\n%s", head)
	}

	out, _, _ = runIn(t, root, "prime", "--story", "S-0001", "--json")
	var v struct {
		Size  struct{ Bytes int } `json:"size"`
		Items []struct {
			Path    string   `json:"path"`
			Heading []string `json:"heading"`
			Reason  string   `json:"reason"`
			Size    int      `json:"size"`
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
	if len(v.Items) != 3 || v.Items[0].Path != "design/system/cli.md" || strings.Join(v.Items[0].Heading, "/") != "CLI" || v.Items[0].Reason != "topics: go" ||
		v.Items[0].Size != len(v.Items[0].Text) || v.Items[1].Reason != "supersedes ADR-0001" || v.Items[2].Path != "design/adrs/0002-linked.md" ||
		len(v.Catalog.NotLoaded) != 2 || len(v.Catalog.InPart) != 1 || len(v.Catalog.InPart[0].Outline) != 3 || v.Catalog.InPart[0].Outline[2].Loaded {
		t.Errorf("json: %+v", v)
	}
	if v.Size.Bytes <= v.Items[0].Size+v.Items[1].Size+v.Items[2].Size {
		t.Errorf("the size leaves out the items: %d", v.Size.Bytes)
	}
}
