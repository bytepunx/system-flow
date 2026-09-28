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
	if body != cat {
		t.Errorf("with every convention at [all] the pack is prime --cat:\n%s\n---\n%s", body, cat)
	}
	if !strings.Contains(head, fmt.Sprintf("size: %d bytes, %d", len(body), strings.Count(body, "\n"))) {
		t.Errorf("size:\n%s", head)
	}
	for _, id := range []string{"S-1", "S-0999"} {
		_, errOut, code = runIn(t, dir, "prime", "--story", id)
		if code == 0 || !strings.Contains(errOut, "flai prime --story "+id) {
			t.Errorf("%s should fail naming it: %d %s", id, code, errOut)
		}
	}
	if _, errOut, _ = runIn(t, dir, "prime", "--story", "S-1"); !strings.Contains(errOut, "S-001 is archived") {
		t.Errorf("archived: %s", errOut)
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
