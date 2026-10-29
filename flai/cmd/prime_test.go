package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/mcpserver"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
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

// ADR-0104: --part prints a part of the pack as JSON, the text the MCP
// prime tool returns for it; it is refused with --cat and past the last part.
func TestPrimePart(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	dir := "../internal/metrics/testdata/good"
	out, errOut, code := runIn(t, dir, "prime", "--story", "S-4", "--part", "1")
	if code != 0 {
		t.Fatalf("prime --part 1: %s", errOut)
	}
	var v struct {
		Story string `json:"story"`
		Part  int    `json:"part"`
		Parts int    `json:"parts"`
	}
	if err := json.Unmarshal([]byte(out), &v); err != nil || v.Story != "S-004" || v.Part != 1 || v.Parts != 1 {
		t.Errorf("json: %v %+v %s", err, v, out)
	}

	repo, err := workitem.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ct, st := mcp.NewInMemoryTransports()
	if _, err := mcpserver.New(mcpserver.Options{Repo: repo, Agent: "claude", Version: "test", MaxWait: time.Second}).Connect(ctx, st, nil); err != nil {
		t.Fatal(err)
	}
	cs, err := mcp.NewClient(&mcp.Implementation{Name: "test-agent", Version: "0"}, nil).Connect(ctx, ct, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = cs.Close() }()
	res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "prime", Arguments: map[string]any{"story": "S-4", "part": 1}})
	if err != nil || res.IsError || len(res.Content) != 1 {
		t.Fatalf("prime tool: %v %+v", err, res)
	}
	if text := res.Content[0].(*mcp.TextContent).Text; out != text+"\n" {
		t.Errorf("--part 1 is not the tool's text:\n%s\n---\n%s", out, text)
	}

	for _, c := range []struct {
		args []string
		want string
	}{
		{[]string{"prime", "--story", "S-4", "--part", "1", "--cat"}, "--cat prints the conventions"},
		{[]string{"prime", "--cat", "--part", "1"}, "--cat prints the conventions"},
		{[]string{"prime", "--story", "S-4", "--part", "2"}, "part 2 of a pack in 1 part; give part 1"},
		{[]string{"prime", "--story", "S-4", "--part", "0"}, "parts count from 1"},
		{[]string{"prime", "--part", "1"}, "give --story"},
	} {
		if _, errOut, code := runIn(t, dir, c.args...); code == 0 || !strings.Contains(errOut, c.want) {
			t.Errorf("%v: %d %q, want %q", c.args, code, errOut, c.want)
		}
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

// S-0175, ADR-0059, ADR-0068: --role primes a sub-agent with the conventions
// its role reads, the story's goal and criteria, and briefs, in half the
// budget; the story's own pack has the conventions the story's agent reads,
// and --cat every convention.
// ADR-0207, S-0415 criteria 1 and 3: flai prime --role takes a role the
// project defines in a role file, with the item flag its file takes, in
// parts; it refuses an item kind the file does not take, names the defined
// roles for a role no file defines, and refuses a role file that does not
// parse with its finding.
func TestPrimeForARoleTheProjectDefines(t *testing.T) {
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
	write("design/conventions/README.md", "# Conventions\n\n- [review.md](review.md)\n- [stream.md](stream.md)\n")
	write("design/conventions/review.md", "---\ntitle: Review\nupdated: 2026-10-28\naudience: agent\norder: 30\nstatus: active\nroles: [review]\n---\n\n# Review\n\n- Review rule.\n")
	write("design/conventions/stream.md", "---\ntitle: Stream\nupdated: 2026-10-28\naudience: agent\norder: 20\nstatus: active\nroles: [story]\n---\n\n# Stream\n\n- Sync at every task.\n")
	role := "design/workflows/roles/reviewer.md"
	write(role, "---\ntitle: Reviewer\nupdated: 2026-10-28\nstatus: active\nname: reviewer\nrole: review\nitem: [epic]\ntopics: [review]\n---\n\nYou are the reviewer.\n")
	if _, errOut, code := runIn(t, root, "epic", "new", "Epic"); code != 0 {
		t.Fatal(errOut)
	}
	if _, errOut, code := runIn(t, root, "story", "new", "Slice", "--epic", "E-0001"); code != 0 {
		t.Fatal(errOut)
	}
	out, errOut, code := runIn(t, root, "prime", "--role", "review", "--epic", "E-0001", "--part", "1")
	if code != 0 {
		t.Fatal(errOut)
	}
	var v struct {
		Role, Item, Story string
		Part, Parts       int
		Conventions       []struct{ Path string }
	}
	if err := json.Unmarshal([]byte(out), &v); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	if v.Role != "review" || v.Item != "E-0001" || v.Story != "" || v.Part != 1 || v.Parts < 1 || len(v.Conventions) != 1 || v.Conventions[0].Path != "design/conventions/review.md" {
		t.Errorf("pack: %+v", v)
	}
	if out, _, _ = runIn(t, root, "prime", "--role", "review", "--epic", "E-0001", "--budget", "1KB"); !strings.Contains(out, "; budget 1024 bytes\n") {
		t.Errorf("--budget 1KB:\n%s", out)
	}
	for _, c := range []struct {
		args []string
		want string
	}{
		{[]string{"prime", "--role", "review", "--story", "S-0001"}, "the reviewer takes an epic, not a story; give --epic E-nnnn"},
		{[]string{"prime", "--role", "review"}, "give --epic E-nnnn"},
		{[]string{"prime", "--role", "write"}, "give a strategic role: plan, orchestrate, analyze, curate, review"},
		{[]string{"prime", "--role", "write", "--story", "S-0001"}, "a strategic agent's plan, orchestrate, analyze, curate, review"},
	} {
		if _, errOut, code := runIn(t, root, c.args...); code == 0 || !strings.Contains(errOut, c.want) {
			t.Errorf("%v: code %d, %s", c.args, code, errOut)
		}
	}
	write(role, "---\ntitle: Reviewer\nrole: [\n---\n")
	for _, r := range []string{"review", "orchestrate"} {
		if _, errOut, code := runIn(t, root, "prime", "--role", r); code == 0 || !strings.Contains(errOut, "role file "+role) {
			t.Errorf("--role %s with a role file that does not parse: code %d, %s", r, code, errOut)
		}
	}
	if _, errOut, code := runIn(t, root, "prime", "--role", "explore", "--story", "S-0001"); code != 0 {
		t.Errorf("--role explore with a role file that does not parse: %s", errOut)
	}
}

func TestPrimeStoryRole(t *testing.T) {
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
	conv := func(name, title string, order int, roles string) {
		fm := fmt.Sprintf("---\ntitle: %s\nupdated: 2026-10-01\naudience: agent\norder: %d\nstatus: active\ntopics: [all]\n%s---\n\n# %s\n\n## Rules\n\n- %s rule.\n\n### Svelte <!-- topics: svelte -->\n\n- %s svelte rule.\n\n<!-- system-flow:end-of-baseline -->\n\n## Project additions\n", title, order, roles, title, title, title)
		write("design/conventions/"+name, fm)
	}
	write("design/conventions/README.md", "# Conventions\n\n- [session-start.md](session-start.md)\n- [stream.md](stream.md)\n- [code-quality.md](code-quality.md)\n- [safety.md](safety.md)\n")
	conv("session-start.md", "Session start", 10, "")
	conv("stream.md", "Stream", 20, "roles: [story]\n")
	conv("code-quality.md", "Code quality", 60, "roles: [verify]\n")
	conv("safety.md", "Safety", 80, "roles: [explore, verify]\n")
	write("design/system/cli.md", "---\ntitle: CLI\nupdated: 2026-09-28\nstatus: active\ntopics: [go]\n---\n\n# CLI\n\nThe command line.\n\n## Commands\n\nAs ADR-0002 decides.\n")
	write("design/system/named.md", "---\ntitle: Named\nupdated: 2026-09-28\nstatus: active\n---\n\n# Named\n\nWhat the story links.\n\n## Part\n\nThe body of the part.\n")
	for n, title := range map[string]string{"0001-named": "Named decision", "0002-linked": "Linked decision"} {
		write("design/adrs/"+n+".md", "---\nid: ADR-"+n[:4]+"\ntitle: "+title+"\nstatus: accepted\ndate: 2026-09-01\n---\n\n# ADR-"+n[:4]+" "+title+"\n\n## Decision\n\n"+title+" is made.\n")
	}
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
	body := strings.Replace(string(data), "## Goal\n", "## Goal\n\nKeep to ADR-0001 and [named](../../../design/system/named.md).\n", 1)
	body = strings.Replace(body, "## Notes\n", "## Notes\n\nA note the role pack leaves out.\n", 1)
	if err := os.WriteFile(story[0], []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	out, errOut, code := runIn(t, root, "prime", "--story", "S-0001")
	if code != 0 {
		t.Fatalf("prime --story: %s", errOut)
	}
	for _, want := range []string{"design/conventions/README.md\n", "- Session start rule.", "- Stream rule."} {
		if !strings.Contains(out, want) {
			t.Errorf("story: missing %q:\n%s", want, out)
		}
	}
	for _, gone := range []string{"Code quality rule", "Safety rule"} {
		if strings.Contains(out, gone) {
			t.Errorf("story: printed %q:\n%s", gone, out)
		}
	}
	out, _, _ = runIn(t, root, "prime", "--cat")
	for _, want := range []string{"- Session start rule.", "- Stream rule.", "- Code quality rule.", "- Safety rule."} {
		if !strings.Contains(out, want) {
			t.Errorf("--cat: missing %q:\n%s", want, out)
		}
	}

	out, errOut, code = runIn(t, root, "prime", "--story", "S-0001", "--role", "explore")
	if code != 0 {
		t.Fatalf("prime --role explore: %s", errOut)
	}
	for _, want := range []string{
		"S-0001 context pack for the explore role\n",
		"role: explore, a sub-agent of the story's agent (ADR-0059)",
		"; budget 40960 bytes\n",
		"design/conventions/session-start.md\n",
		"- Session start rule.",
		"design/conventions/safety.md\n",
		"- Safety rule.",
		"\nstory\n=====\n\n## Goal\n\nKeep to ADR-0001",
		"## Acceptance criteria\n",
		"\ndesign/system/named.md: Named (",
		"bytes; linked from S-0001)\n\nWhat the story links.\n\n- Part\n",
		"\ndesign/system/cli.md: CLI (",
		"- ADR-0001 Named decision (design/adrs/0001-named.md; linked from S-0001): Named decision is made.\n",
		"- ADR-0002 Linked decision (design/adrs/0002-linked.md; linked from design/system/cli.md",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("explore: missing %q:\n%s", want, out)
		}
	}
	for _, gone := range []string{"Stream rule", "Code quality rule", "Safety svelte rule", "The body of the part", "A note the role pack", "open issues", "\ncatalog\n", "reason:"} {
		if strings.Contains(out, gone) {
			t.Errorf("explore: printed %q:\n%s", gone, out)
		}
	}
	if !strings.Contains(out, fmt.Sprintf("size: %d bytes, %d lines, this header included", len(out), strings.Count(out, "\n"))) {
		t.Errorf("size is not the pack's (%d bytes):\n%s", len(out), out)
	}

	out, _, _ = runIn(t, root, "prime", "--story", "S-0001", "--role", "verify", "--json")
	var v struct {
		Role        string `json:"role"`
		Goal        string `json:"goal"`
		Budget      int    `json:"budget"`
		Conventions []struct {
			Path string `json:"path"`
		} `json:"conventions"`
		Items []struct {
			Path string `json:"path"`
			Step string `json:"step"`
		} `json:"items"`
		LeftOut int `json:"briefs_left_out"`
	}
	if err := json.Unmarshal([]byte(out), &v); err != nil {
		t.Fatalf("json: %v %s", err, out)
	}
	if v.Role != "verify" || !strings.HasPrefix(v.Goal, "## Goal\n") || v.Budget != 40960 || len(v.Conventions) != 3 || v.Conventions[0].Path != "design/conventions/session-start.md" ||
		v.Conventions[1].Path != "design/conventions/code-quality.md" || v.Conventions[2].Path != "design/conventions/safety.md" || len(v.Items) != 4 || v.LeftOut != 0 {
		t.Errorf("verify json: %+v", v)
	}
	for _, it := range v.Items {
		if it.Step != "briefed" {
			t.Errorf("a role pack loaded %s (%s)", it.Path, it.Step)
		}
	}

	out, _, _ = runIn(t, root, "prime", "--story", "S-0001", "--role", "explore", "--budget", "2KB")
	if len(out) > 2048 {
		t.Errorf("--budget 2KB printed %d bytes", len(out))
	}
	if !strings.Contains(out, "; budget 2048 bytes\n") || !strings.Contains(out, "briefs with no room in the budget, not printed\n") || strings.Contains(out, "over budget") {
		t.Errorf("--budget 2KB:\n%s", out)
	}
	if out, _, _ = runIn(t, root, "prime", "--story", "S-0001", "--role", "verify", "--budget", "1"); !strings.Contains(out, "over budget: the role's conventions and the story's goal alone exceed it") || strings.Contains(out, "\nbriefs\n") {
		t.Errorf("--budget 1:\n%s", out)
	}
	for _, c := range []struct {
		args []string
		want string
	}{
		{[]string{"prime", "--story", "S-0001", "--role", "write"}, "no such role; a sub-agent's role is explore or verify"},
		{[]string{"prime", "--role", "explore"}, "give --story too"},
	} {
		if _, errOut, code := runIn(t, root, c.args...); code == 0 || !strings.Contains(errOut, c.want) {
			t.Errorf("%v: code %d, %s", c.args, code, errOut)
		}
	}
}

// ADR-0152, S-0386: without a story, prime lists and prints the
// design-system convention only when the project has a design system's
// guidance.
func TestPrimeCatHasTheDesignSystemConventionOnlyWithGuidance(t *testing.T) {
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
	conv := func(name, title string, order int, roles string) {
		write("design/conventions/"+name, fmt.Sprintf("---\ntitle: %s\nupdated: 2026-10-14\naudience: agent\norder: %d\nstatus: active\n%s---\n\n# %s\n\n- %s rule.\n\n<!-- system-flow:end-of-baseline -->\n\n## Project additions\n", title, order, roles, title, title))
	}
	write("design/conventions/README.md", "# Conventions\n\n- [safety.md](safety.md)\n- [design-system.md](design-system.md)\n")
	conv("safety.md", "Safety", 10, "")
	conv("design-system.md", "Design system", 150, "roles: [story, verify]\n")
	for _, guidance := range []bool{false, true} {
		if guidance {
			write("design/system/design-system.md", "---\ntitle: Design system\nupdated: 2026-10-14\nstatus: active\ntopics: [all]\n---\n\n# Design system\n\nWhen to use each token.\n")
		}
		list, errOut, code := runIn(t, root, "prime")
		if code != 0 {
			t.Fatalf("prime: %s", errOut)
		}
		cat, _, _ := runIn(t, root, "prime", "--cat")
		js, _, _ := runIn(t, root, "prime", "--json")
		for what, out := range map[string]string{"list": list, "--cat": cat, "--json": js} {
			if !strings.Contains(out, "safety.md") || strings.Contains(out, "design/conventions/design-system.md") != guidance {
				t.Errorf("guidance %v, %s:\n%s", guidance, what, out)
			}
		}
		if strings.Contains(cat, "- Design system rule.") != guidance {
			t.Errorf("guidance %v, --cat:\n%s", guidance, cat)
		}
	}
}

// primeDigestProject is a project whose digest holds a decision on cli,
// sharing words with its story S-0001 on cli, and two answers: D-0002, to
// the question TH-0001 asks again, and D-0003, newer, sharing fewer of its
// words. TH-0001 awaits the operator.
func primeDigestProject(t *testing.T) string {
	t.Helper()
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
	write("design/conventions/README.md", "# Conventions\n\n- [safety.md](safety.md)\n")
	write("design/conventions/safety.md", "---\ntitle: Safety\nupdated: 2026-10-23\naudience: agent\norder: 10\nstatus: active\n---\n\n# Safety\n\n- Safety rule.\n\n<!-- system-flow:end-of-baseline -->\n\n## Project additions\n")
	write("design/system/cli.md", "---\ntitle: CLI\nupdated: 2026-10-23\nstatus: active\ntopics: [cli]\n---\n\n# CLI\n\nThe command line.\n")
	entry := func(id, slug, kind, title, topics, date, body string) {
		write("design/digest/"+id+"-"+slug+".md", "---\nid: "+id+"\ntitle: "+title+"\nkind: "+kind+"\nsources: [TH-0099]\ntopics: ["+topics+"]\npaths: []\ndate: "+date+"\nstatus: active\nabsorbed_by: \"\"\nasked_again: []\ncited: []\n---\n\n# "+id+" "+title+"\n\n"+body+"\n")
	}
	entry("D-0001", "loads", "decision", "Prime loads digest entries", "cli", "2026-10-01", "The pack loads digest entries ranked against the story.")
	entry("D-0002", "split", "answer", "Split a story whose forecast runs over five days", "orchestration", "2026-10-02", "The operator answered: split a story whose forecast runs over five days into two stories.")
	entry("D-0003", "forecasts", "answer", "Forecasts count days", "orchestration", "2026-10-03", "Each forecast counts days.")
	if _, errOut, code := runIn(t, root, "epic", "new", "Epic"); code != 0 {
		t.Fatal(errOut)
	}
	if _, errOut, code := runIn(t, root, "story", "new", "Prime loads the digest", "--epic", "E-0001", "--topics", "cli"); code != 0 {
		t.Fatal(errOut)
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := threads.New(repo, threads.NewOptions{Title: "Split a story over five days?", On: "design/system/cli.md", Author: "agent-S-0001", Text: "Should the planner split a story whose forecast runs over five days?", Now: time.Date(2026, 10, 23, 9, 0, 0, 0, time.UTC)}); err != nil {
		t.Fatal(err)
	}
	return root
}

// packItem is an item of a context pack as --json prints it.
type packItem struct {
	Path   string `json:"path"`
	ID     string `json:"id"`
	Step   string `json:"step"`
	Reason string `json:"reason"`
	Size   int    `json:"size"`
	Text   string `json:"text"`
}

// packDigest is the digest items of a pack --json printed, in order.
func packDigest(t *testing.T, out string) []packItem {
	t.Helper()
	var v struct {
		Items []packItem `json:"items"`
	}
	if err := json.Unmarshal([]byte(out), &v); err != nil {
		t.Fatalf("json: %v %s", err, out)
	}
	var got []packItem
	for _, it := range v.Items {
		if it.Step == "digest" {
			got = append(got, it)
		}
	}
	return got
}

// ADR-0194, S-0430 criteria 2 and 3: flai prime --story, as text and
// --json, loads the digest entry its topics select, marked digest with its
// reason and size; flai prime --role orchestrate loads the answer to the
// question a thread awaiting the operator asks again first for it, with
// the reason thread TH-nnnn; rank 1, ahead of a newer answer sharing fewer
// of its words.
func TestPrimeLoadsTheDigest(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := primeDigestProject(t)
	one := "# D-0001 Prime loads digest entries\n\nThe pack loads digest entries ranked against the story.\n"
	out, errOut, code := runIn(t, root, "prime", "--story", "S-0001")
	if code != 0 {
		t.Fatalf("prime --story: %s", errOut)
	}
	for _, want := range []string{"1 digest entry with its body\n", "\ndigest\n======\n", "\ndesign/digest/D-0001-loads.md\n=============================\nreason: topics: cli; rank 1\n\n" + one} {
		if !strings.Contains(out, want) {
			t.Errorf("story pack lacks %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "D-0002") || strings.Contains(out, "over budget") {
		t.Errorf("story pack loads an answer on orchestration, or is over budget:\n%s", out)
	}
	out, _, _ = runIn(t, root, "prime", "--story", "S-0001", "--json")
	if got := packDigest(t, out); len(got) != 1 || got[0].ID != "D-0001" || got[0].Path != "design/digest/D-0001-loads.md" || got[0].Reason != "topics: cli; rank 1" || got[0].Text != one || got[0].Size != len(one) {
		t.Errorf("story pack digest items: %+v", got)
	}

	two := "# D-0002 Split a story whose forecast runs over five days\n\nThe operator answered: split a story whose forecast runs over five days into two stories.\n"
	out, errOut, code = runIn(t, root, "prime", "--role", "orchestrate")
	if code != 0 {
		t.Fatalf("prime --role orchestrate: %s", errOut)
	}
	for _, want := range []string{"\ndigest\n======\n", "\ndesign/digest/D-0002-split.md\n=============================\nreason: thread TH-0001; rank 1\n\n" + two} {
		if !strings.Contains(out, want) {
			t.Errorf("orchestrator's pack lacks %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "D-0001") || strings.Contains(out, "over budget") {
		t.Errorf("orchestrator's pack loads a decision, or is over budget:\n%s", out)
	}
	out, _, _ = runIn(t, root, "prime", "--role", "orchestrate", "--json")
	got := packDigest(t, out)
	if len(got) == 0 || got[0].ID != "D-0002" || got[0].Path != "design/digest/D-0002-split.md" || got[0].Reason != "thread TH-0001; rank 1" || got[0].Text != two || got[0].Size != len(two) {
		t.Fatalf("orchestrator's digest items: %+v", got)
	}
	for _, it := range got[1:] {
		if it.ID != "D-0003" || it.Reason != "thread TH-0001; rank 2" {
			t.Errorf("after D-0002: %+v", it)
		}
	}
}

// primeReasonsProject is a project whose orchestrator's convention tags its
// sections by reason (ADR-0217 §2), with a design document a general pack
// briefs, an epic E-0001 with stories S-0001 and S-0002, a stored flai
// verify result for S-0001, TH-0001 on S-0001 from its agent, and the
// planner's open TH-0002 on E-0001.
func primeReasonsProject(t *testing.T) string {
	t.Helper()
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
	write("design/conventions/README.md", "# Conventions\n\n- [orchestration.md](orchestration.md)\n")
	write("design/conventions/orchestration.md", "---\ntitle: Orchestration\nupdated: 2026-10-29\naudience: agent\norder: 50\nstatus: active\nroles: [orchestrate]\n---\n\n# Orchestration\n\n"+
		"## Every run\n\n- Every run rule.\n\n"+
		"## Reviews <!-- topics: all, review -->\n\n- Review rule.\n\n"+
		"## Threads <!-- topics: all, thread -->\n\n- Thread rule.\n\n"+
		"## Analysis <!-- topics: all, analysis -->\n\n- Analysis rule.\n\n"+
		"<!-- system-flow:end-of-baseline -->\n\n## Project additions\n")
	write("design/system/orchestration.md", "---\ntitle: Orchestration design\nupdated: 2026-10-29\nstatus: active\ntopics: [orchestration]\n---\n\n# Orchestration design\n\nHow the orchestrator runs.\n")
	write(".flai-cache/verify/S-0001.json", "{\"story\": \"S-0001\", \"passed\": true}\n")
	if _, errOut, code := runIn(t, root, "epic", "new", "Epic"); code != 0 {
		t.Fatal(errOut)
	}
	for _, title := range []string{"Slice", "Second slice"} {
		if _, errOut, code := runIn(t, root, "story", "new", title, "--epic", "E-0001"); code != 0 {
			t.Fatal(errOut)
		}
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 10, 29, 9, 0, 0, 0, time.UTC)
	for _, th := range []threads.NewOptions{
		{Title: "Which way round?", On: "S-0001", Author: "agent-S-0001", Text: "Which way round should the slice go?", Now: now},
		{Title: "Split the epic?", On: "E-0001", Author: "planner-E-0001", Text: "Should the epic split in two?", Now: now},
	} {
		if _, err := threads.New(repo, th); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

// reasonPack is a pack primed for reasons as --json prints it.
type reasonPack struct {
	Role        string   `json:"role"`
	Reasons     []string `json:"reasons"`
	Part        int      `json:"part"`
	Parts       int      `json:"parts"`
	OpenIssues  int      `json:"open_issues"`
	Conventions []struct {
		Path    string                     `json:"path"`
		Kept    []struct{ Heading string } `json:"kept"`
		LeftOut []struct{ Heading string } `json:"left_out"`
	} `json:"conventions"`
	Items   []packItem `json:"items"`
	Catalog struct {
		NotLoaded []struct{ Path string } `json:"not_loaded"`
		InPart    []struct{ Path string } `json:"in_part"`
	} `json:"catalog"`
}

// ADR-0217 §1, §3, §5, S-0537 criterion 1: flai prime --role orchestrate
// --reason primes the orchestrator for each kind of reason, loading what it
// needs, each with step reason and the reason, and keeping the convention
// sections its kind names; it loads no brief and no catalog, as --json, as
// text, and in parts. The general pack is printed without --reason as
// before.
func TestPrimeForTheOrchestratorsReasons(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := primeReasonsProject(t)
	story := func(id string) string {
		m, err := filepath.Glob(filepath.Join(root, "wip", "kanban", "stories", id+"-*.md"))
		if err != nil || len(m) != 1 {
			t.Fatal(m, err)
		}
		return "wip/kanban/stories/" + filepath.Base(m[0])
	}
	epic, err := filepath.Glob(filepath.Join(root, "wip", "kanban", "epics", "E-0001-*.md"))
	if err != nil || len(epic) != 1 {
		t.Fatal(epic, err)
	}
	thread := func(id string) string {
		m, err := filepath.Glob(filepath.Join(root, "wip", "threads", id+"-*.md"))
		if err != nil || len(m) != 1 {
			t.Fatalf("%s: %v %v", id, m, err)
		}
		return "wip/threads/" + filepath.Base(m[0])
	}
	for _, c := range []struct {
		reason, canon string
		paths         []string // what the reason loads, in order
		kept          string   // the tagged section it keeps beside Every run
	}{
		{"thread:th-1", "thread:TH-0001", []string{thread("TH-0001"), story("S-0001")}, "Threads"},
		{"review:S-1", "review:S-0001", []string{story("S-0001"), ".flai-cache/verify/S-0001.json"}, "Reviews"},
		{"planned:E-0001", "planned:E-0001", []string{"wip/kanban/epics/" + filepath.Base(epic[0]), story("S-0001"), story("S-0002"), thread("TH-0002")}, ""},
		{"draft:S-0002", "draft:S-0002", []string{story("S-0002")}, ""},
		{"release", "release", []string{"release_evaluate"}, ""},
		{"ready", "ready", []string{"board", "promote_candidates"}, ""},
	} {
		out, errOut, code := runIn(t, root, "prime", "--role", "orchestrate", "--reason", c.reason, "--json")
		if code != 0 {
			t.Fatalf("--reason %s: %s", c.reason, errOut)
		}
		var v reasonPack
		if err := json.Unmarshal([]byte(out), &v); err != nil {
			t.Fatalf("--reason %s: %v %s", c.reason, err, out)
		}
		if v.Role != "orchestrate" || len(v.Reasons) != 1 || v.Reasons[0] != c.canon || v.OpenIssues != 0 || len(v.Catalog.NotLoaded)+len(v.Catalog.InPart) != 0 {
			t.Errorf("--reason %s: role %q, reasons %v, open issues %d, catalog %+v", c.reason, v.Role, v.Reasons, v.OpenIssues, v.Catalog)
		}
		var paths []string
		for _, it := range v.Items {
			if it.Step == "digest" {
				continue
			}
			if it.Step != "reason" || it.Reason != c.canon {
				t.Errorf("--reason %s loaded %s as %s, %q", c.reason, it.Path, it.Step, it.Reason)
			}
			paths = append(paths, it.Path)
		}
		if strings.Join(paths, " ") != strings.Join(c.paths, " ") {
			t.Errorf("--reason %s loaded %v; want %v", c.reason, paths, c.paths)
		}
		want := "Orchestration Every run"
		if c.kept != "" {
			want += " " + c.kept
		}
		if len(v.Conventions) != 1 {
			t.Fatalf("--reason %s: conventions %+v", c.reason, v.Conventions)
		}
		var kept []string
		for _, k := range v.Conventions[0].Kept {
			if k.Heading != "Project additions" && k.Heading != "" {
				kept = append(kept, k.Heading)
			}
		}
		if strings.Join(kept, " ") != want {
			t.Errorf("--reason %s kept %v; want %s", c.reason, kept, want)
		}
	}

	both := []string{"prime", "--role", "orchestrate", "--reason", "thread:TH-0001", "--reason", "review:S-0001"}
	out, errOut, code := runIn(t, root, both...)
	if code != 0 {
		t.Fatalf("%v: %s", both, errOut)
	}
	for _, want := range []string{"primed for thread:TH-0001, review:S-0001 (ADR-0217)", "- Thread rule.", "- Review rule.", "\n.flai-cache/verify/S-0001.json\n", "Which way round should the slice go?"} {
		if !strings.Contains(out, want) {
			t.Errorf("text lacks %q:\n%s", want, out)
		}
	}
	for _, gone := range []string{"- Analysis rule.", "\nbriefs\n", "\ncatalog\n", "How the orchestrator runs.", "design/system/orchestration.md"} {
		if strings.Contains(out, gone) {
			t.Errorf("text has %q:\n%s", gone, out)
		}
	}
	if !strings.Contains(out, fmt.Sprintf("size: %d bytes, %d lines, this header included", len(out), strings.Count(out, "\n"))) {
		t.Errorf("size is not the pack's (%d bytes):\n%s", len(out), out)
	}
	if strings.Count(out, "\n"+story("S-0001")+"\n") != 1 {
		t.Errorf("the story two reasons load printed other than once:\n%s", out)
	}

	out, errOut, code = runIn(t, root, append(both, "--part", "1")...)
	if code != 0 {
		t.Fatalf("--part 1: %s", errOut)
	}
	var v reasonPack
	if err := json.Unmarshal([]byte(out), &v); err != nil || v.Part != 1 || v.Parts != 1 || v.Role != "orchestrate" || strings.Join(v.Reasons, " ") != "thread:TH-0001 review:S-0001" || len(v.Items) < 3 {
		t.Errorf("--part 1: %v %+v %s", err, v, out)
	}
	if _, errOut, code = runIn(t, root, append(both, "--part", "2")...); code == 0 || !strings.Contains(errOut, "part 2 of a pack in 1 part") {
		t.Errorf("--part 2: %d %s", code, errOut)
	}
	if out, _, _ = runIn(t, root, append(both, "--budget", "1")...); !strings.Contains(out, "; budget 1 bytes\nover budget: the role's conventions alone exceed it; what the reasons load is still loaded") || !strings.Contains(out, "Which way round should the slice go?") {
		t.Errorf("--budget 1:\n%s", out)
	}

	out, errOut, code = runIn(t, root, "prime", "--role", "orchestrate")
	if code != 0 {
		t.Fatalf("--role orchestrate: %s", errOut)
	}
	for _, want := range []string{"for the whole project", "- Analysis rule.", "design/system/orchestration.md"} {
		if !strings.Contains(out, want) {
			t.Errorf("the general pack lacks %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "primed for") || strings.Contains(out, "Which way round should the slice go?") {
		t.Errorf("the general pack is primed for reasons:\n%s", out)
	}
}

// ADR-0217 §1: --reason is refused with another role, without a role, and
// with --story or --epic, each saying to give --role orchestrate alone; an
// unknown reason, and one without the ID it needs, are refused naming every
// form a reason takes.
func TestPrimeRefusesReasonsOutsideTheOrchestrator(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := primeReasonsProject(t)
	forms := "thread:TH-nnnn, review:S-nnnn, planned:E-nnnn or planned:S-nnnn, draft:S-nnnn, release, or ready"
	for _, c := range []struct {
		args []string
		want string
	}{
		{[]string{"prime", "--role", "plan", "--epic", "E-0001", "--reason", "release"}, "--reason primes the orchestrator for its run's reasons, and --role plan is not the orchestrator; give --role orchestrate and no --story or --epic"},
		{[]string{"prime", "--role", "analyze", "--reason", "ready"}, "--role analyze is not the orchestrator; give --role orchestrate"},
		{[]string{"prime", "--role", "verify", "--story", "S-0001", "--reason", "ready"}, "--role verify is not the orchestrator"},
		{[]string{"prime", "--reason", "ready"}, "no --role is not the orchestrator; give --role orchestrate"},
		{[]string{"prime", "--story", "S-0001", "--reason", "review:S-0001"}, "no --role is not the orchestrator"},
		{[]string{"prime", "--role", "orchestrate", "--story", "S-0001", "--reason", "review:S-0001"}, "not one item; give --role orchestrate and no --story or --epic"},
		{[]string{"prime", "--role", "orchestrate", "--epic", "E-0001", "--reason", "planned:E-0001"}, "not one item; give --role orchestrate and no --story or --epic"},
		{[]string{"prime", "--role", "orchestrate", "--reason", "ready", "--reason", "merge:S-0001"}, "flai prime --reason merge:S-0001: no such kind of reason"},
		{[]string{"prime", "--role", "orchestrate", "--reason", "review", "--json"}, "review needs the ID of a story (S-nnnn); a reason is " + forms},
		{[]string{"prime", "--role", "orchestrate", "--reason", "thread:TH-0099"}, "TH-0099"},
	} {
		out, errOut, code := runIn(t, root, c.args...)
		if code == 0 || !strings.Contains(errOut, c.want) || out != "" {
			t.Errorf("%v: code %d, out %q, %s; want %q", c.args, code, out, errOut, c.want)
		}
	}
	if _, errOut, code := runIn(t, root, "prime", "--role", "orchestrate", "--reason", "merge:S-0001"); code == 0 || !strings.Contains(errOut, "; a reason is "+forms) {
		t.Errorf("an unknown reason does not name the forms: code %d, %s", code, errOut)
	}
}
