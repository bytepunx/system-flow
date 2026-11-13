package manifest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ADR-0261: an agent that names no harness falls back to flai's built-in
// default, claude-code with no model, under its own words; one that names a
// harness is itself.
func TestOrBuiltinFallsBackToClaudeCodeOnlyWhenNoHarnessIsNamed(t *testing.T) {
	var none *Agent
	if got := none.OrBuiltin(); !got.Same(&Agent{Harness: "claude-code"}) {
		t.Errorf("nil: %+v", got)
	}
	if b := Builtin(); b.Harness != BuiltinHarness || b.Model != "" || len(b.Config) != 0 || len(b.Roles) != 0 {
		t.Errorf("built-in: %+v", b)
	}

	own := &Agent{Model: "claude-sonnet-5", Provider: "litellm", Config: map[string]string{"effort": "high"},
		Roles: map[string]Role{"explore": {Model: "haiku"}}, Image: &AgentImage{Layer: "images/agent/Dockerfile"}}
	got := own.OrBuiltin()
	if got.Harness != "claude-code" || got.Model != "claude-sonnet-5" || got.Provider != "litellm" || got.Config["effort"] != "high" || got.Roles["explore"].Model != "haiku" {
		t.Errorf("only a model: %+v", got)
	}
	if got.Image == nil || got.Image.Layer != "images/agent/Dockerfile" {
		t.Errorf("the image was not kept: %+v", got.Image)
	}
	if own.Harness != "" {
		t.Errorf("the agent was changed: %+v", own)
	}

	for _, h := range []string{"command", "none", "claude-code"} {
		named := &Agent{Harness: h, Model: "m"}
		if got := named.OrBuiltin(); got != named {
			t.Errorf("harness %s: %+v", h, got)
		}
	}
}

func TestAgentWithMergesAndOverrides(t *testing.T) {
	def := &Agent{Harness: "claude-code", Model: "claude-opus-5-5", Config: map[string]string{"effort": "high", "max_turns": "50"}}
	got := def.With(&Agent{Model: "claude-sonnet-5", Config: map[string]string{"effort": "low"}})
	if got.Harness != "claude-code" || got.Model != "claude-sonnet-5" || got.Config["effort"] != "low" || got.Config["max_turns"] != "50" {
		t.Errorf("merged: %+v", got)
	}
	if def.Config["effort"] != "high" {
		t.Error("the default was changed by a merge")
	}
	var none *Agent
	if none.With(nil) != nil || none.With(&Agent{}) != nil {
		t.Error("nothing with nothing is nil")
	}
	if got := none.With(&Agent{Harness: "codex"}); got == nil || got.Harness != "codex" {
		t.Errorf("an override with no default: %+v", got)
	}
}

func TestAgentValidate(t *testing.T) {
	good := &Agent{Harness: "claude-code", Model: "anthropic/claude-opus-5-5@2026", Config: map[string]string{"max_turns": "50", "permission.mode": "plan"}}
	if err := good.Validate(); err != nil {
		t.Errorf("good: %v", err)
	}
	for _, bad := range []*Agent{
		{Harness: "Claude Code"},
		{Model: "has space"},
		{Config: map[string]string{"Bad Key": "x"}},
		{Config: map[string]string{"k": "two\nlines"}},
	} {
		if err := bad.Validate(); err == nil {
			t.Errorf("accepted %+v", bad)
		}
	}
}

// S-0189: an agent's roles merge role by role and key by key, compare, check,
// print, and read back as they were written.
func TestAgentRoles(t *testing.T) {
	def := &Agent{Harness: "claude-code", Model: "claude-opus-5-5", Roles: map[string]Role{
		"explore": {Model: "haiku"},
		"verify":  {Model: "sonnet", Config: map[string]string{"effort": "medium"}},
	}}
	got := def.With(&Agent{Roles: map[string]Role{"verify": {Config: map[string]string{"effort": "high"}}, "plan": {Model: "opus"}}})
	if got.Roles["explore"].Model != "haiku" || got.Roles["verify"].Model != "sonnet" || got.Roles["verify"].Config["effort"] != "high" || got.Roles["plan"].Model != "opus" {
		t.Errorf("merged: %+v", got.Roles)
	}
	if def.Roles["verify"].Config["effort"] != "medium" {
		t.Error("the default's role was changed by a merge")
	}
	if def.Same(got) || !got.Same(def.With(got)) {
		t.Error("Same does not compare roles")
	}
	if (&Agent{Roles: map[string]Role{"verify": {Model: "sonnet"}}}).IsZero() {
		t.Error("an agent with only roles is not empty")
	}
	if s := got.String(); s != "claude-code, claude-opus-5-5; explore: haiku; plan: opus; verify: sonnet, effort=high" {
		t.Errorf("String: %s", s)
	}
	for _, bad := range []*Agent{
		{Roles: map[string]Role{"Verify": {Model: "sonnet"}}},
		{Roles: map[string]Role{"verify": {}}},
		{Roles: map[string]Role{"verify": {Model: "has space"}}},
		{Roles: map[string]Role{"verify": {Config: map[string]string{"Bad": "x"}}}},
	} {
		if err := bad.Validate(); err == nil {
			t.Errorf("accepted %+v", bad.Roles)
		}
	}
	path := filepath.Join(t.TempDir(), File)
	if err := os.WriteFile(path, []byte("version: 1\nname: Harbour\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := WriteAgent(path, got); err != nil {
		t.Fatal(err)
	}
	m, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if !m.Agent.Same(got) {
		data, _ := os.ReadFile(path)
		t.Errorf("read back %+v from:\n%s", m.Agent, data)
	}
}

// S-0349: an agent's provider is checked as a harness is, carried by a merge
// and a copy, compared, printed, and read back as it was written.
func TestAgentProvider(t *testing.T) {
	def := &Agent{Harness: "claude-code", Model: "claude-opus-5-5", Provider: "litellm"}
	if err := def.Validate(); err != nil {
		t.Errorf("good: %v", err)
	}
	for _, bad := range []string{"Lite LLM", "LiteLLM", "-litellm"} {
		err := (&Agent{Provider: bad}).Validate()
		if err == nil || !strings.Contains(err.Error(), `agent provider "`+bad+`"`) {
			t.Errorf("%q: %v", bad, err)
		}
	}
	if got := def.With(&Agent{Model: "claude-sonnet-5"}); got.Provider != "litellm" {
		t.Errorf("a merge that leaves provider out dropped it: %+v", got)
	}
	if got := def.With(&Agent{Provider: "openrouter"}); got.Provider != "openrouter" || got.Model != "claude-opus-5-5" {
		t.Errorf("a merge that sets provider: %+v", got)
	}
	if def.Provider != "litellm" {
		t.Error("the default was changed by a merge")
	}
	if got := (*Agent)(nil).With(&Agent{Provider: "litellm"}); got == nil || got.Provider != "litellm" {
		t.Errorf("only a provider: %+v", got)
	}
	if (&Agent{Provider: "litellm"}).IsZero() {
		t.Error("an agent with only a provider is not empty")
	}
	if c := def.Clone(); c.Provider != "litellm" || !c.Same(def) {
		t.Errorf("clone: %+v", c)
	}
	if def.Same(&Agent{Harness: "claude-code", Model: "claude-opus-5-5"}) {
		t.Error("Same does not compare providers")
	}
	if s := (&Agent{Harness: "claude-code", Model: "claude-opus-5-5", Provider: "litellm", Config: map[string]string{"effort": "high"}}).String(); s != "claude-code, claude-opus-5-5, provider litellm, effort=high" {
		t.Errorf("String: %s", s)
	}
	block := AgentBlock(def, "")
	if !strings.Contains(block, "  provider: litellm\n") {
		t.Errorf("block:\n%s", block)
	}
	path := filepath.Join(t.TempDir(), File)
	if err := os.WriteFile(path, []byte("version: 1\nname: Harbour\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := WriteAgent(path, def); err != nil {
		t.Fatal(err)
	}
	m, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if !m.Agent.Same(def) {
		data, _ := os.ReadFile(path)
		t.Errorf("read back %+v from:\n%s", m.Agent, data)
	}
}

func TestParseRoles(t *testing.T) {
	got, err := ParseRoles([]string{"verify=claude-code"}, []string{"explore=haiku", "verify = sonnet"}, []string{"verify.effort=medium"})
	if err != nil {
		t.Fatal(err)
	}
	if got["explore"].Model != "haiku" || got["verify"].Harness != "claude-code" || got["verify"].Model != "sonnet" || got["verify"].Config["effort"] != "medium" {
		t.Errorf("parsed: %+v", got)
	}
	if got, err := ParseRoles(nil, nil, nil); got != nil || err != nil {
		t.Errorf("nothing: %v %v", got, err)
	}
	for _, bad := range [][]string{{"", "verify", ""}, {"", "", "verify=x"}, {"", "", ".effort=x"}} {
		if _, err := ParseRoles(split(bad[0]), split(bad[1]), split(bad[2])); err == nil {
			t.Errorf("accepted %q", bad)
		}
	}
}

func split(s string) []string {
	if s == "" {
		return nil
	}
	return []string{s}
}

// flai owns the agent block of system-flow.yaml and nothing else of it.
func TestWriteAgentKeepsEverythingElse(t *testing.T) {
	path := filepath.Join(t.TempDir(), File)
	orig := "version: 1\nname: Harbour\n# a comment the operator wrote\nlayout:\n  design: design\n  docs: docs\n  wip: wip\ndashboard:\n  port: 4242\n"
	if err := os.WriteFile(path, []byte(orig), 0o644); err != nil {
		t.Fatal(err)
	}
	a := &Agent{Harness: "claude-code", Model: "claude-opus-5-5", Config: map[string]string{"effort": "high", "note": "yes"}}
	if err := WriteAgent(path, a); err != nil {
		t.Fatal(err)
	}
	m, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if m.Agent.Harness != "claude-code" || m.Agent.Model != "claude-opus-5-5" || m.Agent.Config["effort"] != "high" || m.Agent.Config["note"] != "yes" || m.Dashboard.Port != 4242 {
		t.Errorf("read back: %+v %+v", m.Agent, m.Dashboard)
	}
	data, _ := os.ReadFile(path)
	if !strings.HasPrefix(string(data), orig) || !strings.Contains(string(data), `note: "yes"`) {
		t.Errorf("written:\n%s", data)
	}
	// replaced, not appended twice, and removed when cleared
	if err := WriteAgent(path, &Agent{Model: "claude-sonnet-5"}); err != nil {
		t.Fatal(err)
	}
	data, _ = os.ReadFile(path)
	if strings.Count(string(data), "agent:") != 1 || strings.Contains(string(data), "claude-code") {
		t.Errorf("replaced:\n%s", data)
	}
	if err := WriteAgent(path, nil); err != nil {
		t.Fatal(err)
	}
	if data, _ = os.ReadFile(path); string(data) != orig {
		t.Errorf("cleared:\n%s", data)
	}
	if err := WriteAgent(path, &Agent{Harness: "Not Valid"}); err == nil {
		t.Error("an invalid agent was written")
	}
}

// S-0510: the project's agent names an image layer, which counts as set,
// copies and compares, and is never carried by a merge into a story's agent
// or a strategic role's.
func TestAgentImage(t *testing.T) {
	def := &Agent{Harness: "claude-code", Model: "claude-opus-5-5", Image: &AgentImage{Layer: "images/agent/Dockerfile"}}
	if (&Agent{Image: &AgentImage{Layer: "Dockerfile"}}).IsZero() {
		t.Error("an agent with only an image is not empty")
	}
	if c := def.Clone(); !c.Same(def) || c.Image == def.Image {
		t.Errorf("clone shares or loses the image: %+v", c)
	}
	if def.Same(&Agent{Harness: "claude-code", Model: "claude-opus-5-5"}) || def.Same(&Agent{Harness: "claude-code", Model: "claude-opus-5-5", Image: &AgentImage{Layer: "other/Dockerfile"}}) {
		t.Error("Same does not compare images")
	}
	// S-0745: the args are compared, and copied rather than shared.
	withArgs := func(args map[string]string) *Agent {
		return &Agent{Image: &AgentImage{Layer: "Dockerfile", Args: args}}
	}
	for _, c := range []struct {
		name string
		a, b map[string]string
		same bool
	}{
		{"the same args", map[string]string{"GO_VERSION": "1.26.0"}, map[string]string{"GO_VERSION": "1.26.0"}, true},
		{"no args either side", nil, map[string]string{}, true},
		{"another value", map[string]string{"GO_VERSION": "1.26.0"}, map[string]string{"GO_VERSION": "1.25.0"}, false},
		{"another name", map[string]string{"GO_VERSION": "1.26.0"}, map[string]string{"GO": "1.26.0"}, false},
		{"args on one side", map[string]string{"GO_VERSION": "1.26.0"}, nil, false},
	} {
		if got := withArgs(c.a).Same(withArgs(c.b)); got != c.same {
			t.Errorf("%s: Same = %v, want %v", c.name, got, c.same)
		}
	}
	args := withArgs(map[string]string{"GO_VERSION": "1.26.0"})
	for name, c := range map[string]*Agent{"clone": args.Clone(), "OrBuiltin": args.OrBuiltin()} {
		c.Image.Args["GO_VERSION"] = "0"
		if args.Image.Args["GO_VERSION"] != "1.26.0" {
			t.Errorf("%s shares the args", name)
		}
		args.Image.Args["GO_VERSION"] = "1.26.0"
	}
	for _, c := range []struct {
		name string
		got  *Agent
	}{
		{"a story with nothing of its own", def.With(nil)},
		{"a story with a model of its own", def.With(&Agent{Model: "claude-sonnet-5"})},
		{"an override that names an image", (&Agent{Model: "claude-sonnet-5"}).With(def)},
	} {
		if c.got == nil || c.got.Image != nil || c.got.Harness != "claude-code" {
			t.Errorf("%s: %+v", c.name, c.got)
		}
	}
	if got := (&Agent{Image: def.Image}).With(nil); got != nil {
		t.Errorf("an agent of only an image merged into a story: %+v", got)
	}
}

// S-0510: agent.image.layer is a clean path inside the repository on one
// line, an image block sets it, and only the project's agent names one.
func TestAgentImageProblems(t *testing.T) {
	for _, c := range []struct {
		what  string
		layer *AgentImage
		want  string
	}{
		{"agent", &AgentImage{Layer: "flai/images/system-flow/Dockerfile"}, ""},
		{"agent", &AgentImage{Layer: "Dockerfile"}, ""},
		{"agent", &AgentImage{}, "agent image sets nothing; name a Dockerfile with layer"},
		{"agent", &AgentImage{Layer: "/srv/Dockerfile"}, `agent image layer "/srv/Dockerfile" is an absolute path`},
		{"agent", &AgentImage{Layer: "C:/images/Dockerfile"}, `agent image layer "C:/images/Dockerfile" is an absolute path`},
		{"agent", &AgentImage{Layer: "../other/Dockerfile"}, `agent image layer "../other/Dockerfile" leaves the repository`},
		{"agent", &AgentImage{Layer: ".."}, `agent image layer ".." leaves the repository`},
		{"agent", &AgentImage{Layer: "images//Dockerfile"}, `agent image layer "images//Dockerfile" is not a clean path`},
		{"agent", &AgentImage{Layer: "./Dockerfile"}, `agent image layer "./Dockerfile" is not a clean path`},
		{"agent", &AgentImage{Layer: `images\Dockerfile`}, `agent image layer "images\\Dockerfile" is not a clean path; write it with forward slashes`},
		{"agent", &AgentImage{Layer: "images/\nDockerfile"}, "agent image layer spans lines"},
		{"planning.agent", &AgentImage{Layer: "Dockerfile"}, "planning.agent image is the project's alone; name the layer under agent.image"},
		// S-0745, ADR-0295: args beside the layer, ARG names other than BASE,
		// with literal values on one line.
		{"agent", &AgentImage{Layer: "Dockerfile", Args: map[string]string{"GO_VERSION": "1.26.0", "_pnpm2": "10.34.5", "EMPTY": ""}}, ""},
		{"agent", &AgentImage{Layer: "Dockerfile", Args: map[string]string{"GO-VERSION": "1.26.0"}}, `agent image args name "GO-VERSION" is not a Dockerfile ARG name such as GO_VERSION`},
		{"agent", &AgentImage{Layer: "Dockerfile", Args: map[string]string{"2GO": "1"}}, `agent image args name "2GO" is not a Dockerfile ARG name`},
		{"agent", &AgentImage{Layer: "Dockerfile", Args: map[string]string{"BASE": "debian"}}, "agent image args set BASE, which flai sets to the harness image it builds the layer over; leave BASE out of agent.image.args"},
		{"agent", &AgentImage{Layer: "Dockerfile", Args: map[string]string{"GO_VERSION": "1.26.0\nRUN x"}}, "agent image args GO_VERSION spans lines; write its value in agent.image.args as one literal line"},
		{"agent", &AgentImage{Args: map[string]string{"GO_VERSION": "1.26.0"}}, "agent image args name no layer to build; name its Dockerfile with agent.image.layer"},
		{"agent", &AgentImage{Layer: "Dockerfile", Args: map[string]string{"Z-1": "", "BASE": "x", "A-1": ""}}, `agent image args set BASE, which flai sets to the harness image it builds the layer over; leave BASE out of agent.image.args; agent image args name "Z-1"`},
		{"agent", &AgentImage{Layer: "/srv/Dockerfile", Args: map[string]string{"BASE": "x"}}, `agent image layer "/srv/Dockerfile" is an absolute path`},
	} {
		got := strings.Join((&Agent{Image: c.layer}).problems(c.what), "; ")
		if c.want == "" && got != "" || !strings.Contains(got, c.want) {
			t.Errorf("%s %q %v: got %q, want %q", c.what, c.layer.Layer, c.layer.Args, got, c.want)
		}
	}
	// Each problem is reported by name, in order, the layer's first.
	got := (&Agent{Image: &AgentImage{Layer: "/srv/Dockerfile", Args: map[string]string{"Z-1": "", "BASE": "x", "A-1": ""}}}).problems("agent")
	if len(got) != 4 || !strings.Contains(got[0], "absolute path") || !strings.Contains(got[1], `"A-1"`) || !strings.Contains(got[2], "BASE") || !strings.Contains(got[3], `"Z-1"`) {
		t.Errorf("problems in order: %q", got)
	}
}

// S-0510: the image is written under the agent block, reads back as written,
// and stays when flai agent set or clear rewrites the block without one.
func TestWriteAgentKeepsTheImage(t *testing.T) {
	path := filepath.Join(t.TempDir(), File)
	orig := "version: 1\nname: Harbour\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"
	if err := os.WriteFile(path, []byte(orig), 0o644); err != nil {
		t.Fatal(err)
	}
	image := &AgentImage{Layer: "images/agent/Dockerfile"}
	argsImage := &AgentImage{Layer: "images/agent/Dockerfile", Args: map[string]string{"GO_VERSION": "1.26.0", "PNPM_VERSION": "10.34.5", "NODE": "true"}}
	for _, c := range []struct {
		name  string
		write *Agent
		want  *Agent
	}{
		{"only an image", &Agent{Image: image}, &Agent{Image: image}},
		{"an agent set without one", &Agent{Model: "claude-opus-5-5"}, &Agent{Model: "claude-opus-5-5", Image: image}},
		{"the agent cleared", nil, &Agent{Image: image}},
		{"another image", &Agent{Harness: "claude-code", Image: &AgentImage{Layer: "Dockerfile"}}, &Agent{Harness: "claude-code", Image: &AgentImage{Layer: "Dockerfile"}}},
		{"an image with args", &Agent{Image: argsImage}, &Agent{Image: argsImage}},
		{"an agent set beside args", &Agent{Model: "claude-opus-5-5"}, &Agent{Model: "claude-opus-5-5", Image: argsImage}},
	} {
		if err := WriteAgent(path, c.write); err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}
		m, err := Load(path)
		data, _ := os.ReadFile(path)
		if err != nil || !m.Agent.Same(c.want) || strings.Count(string(data), "image:") != 1 {
			t.Errorf("%s: read back %+v, %v, from:\n%s", c.name, m.Agent, err, data)
		}
	}
	for _, bad := range []*AgentImage{{Layer: "/abs/Dockerfile"}, {Layer: "Dockerfile", Args: map[string]string{"BASE": "x"}}} {
		if err := WriteAgent(path, &Agent{Image: bad}); err == nil {
			t.Errorf("an invalid image was written: %+v", bad)
		}
	}
}
