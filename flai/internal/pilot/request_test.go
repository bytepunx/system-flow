package pilot

import (
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/harness"
)

// The values in the fake environment, which no refusal may print.
const (
	tokenValue = "s3cr3t-token-value"
	urlValue   = "http://host.docker.internal:4242/mcp"
	hostValue  = "127.0.0.1:4242"
)

// env is a fake environment holding the MCP server's URL, token, and Host.
func env(name string) string {
	return map[string]string{
		"FLAI_MCP_URL":   urlValue,
		"FLAI_MCP_TOKEN": tokenValue,
		"FLAI_MCP_HOST":  hostValue,
	}[name]
}

// isDirOnly finds only dir.
func isDirOnly(dir string) func(string) bool {
	return func(p string) bool { return p == dir }
}

// valid is a start request pilot runs, with cache as its cache.
func valid(cache string) harness.Request {
	return harness.Request{
		Story:       "S-0507",
		Name:        "agent-S-0507",
		Session:     "S-0507-1",
		Inner:       harness.ClaudeCode,
		MCPURLEnv:   "FLAI_MCP_URL",
		MCPTokenEnv: "FLAI_MCP_TOKEN",
		MCPHostEnv:  "FLAI_MCP_HOST",
		Cache:       cache,
	}
}

// oneLineWithoutSecrets fails t when err is nil, spans lines, or prints a
// value of the fake environment.
func oneLineWithoutSecrets(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("refused nothing")
	}
	msg := err.Error()
	if strings.Contains(msg, "\n") {
		t.Errorf("the refusal spans lines: %q", msg)
	}
	for _, v := range []string{tokenValue, urlValue, hostValue} {
		if strings.Contains(msg, v) {
			t.Errorf("the refusal prints %q: %s", v, msg)
		}
	}
}

func TestAValidStartRequestIsReadAndPassesTheCheck(t *testing.T) {
	cache := t.TempDir()
	in := " \n" + `{"Story":"S-0507","Name":"agent-S-0507","Session":"S-0507-1","inner":"claude-code",` +
		`"mcp_url_env":"FLAI_MCP_URL","mcp_token_env":"FLAI_MCP_TOKEN","mcp_host_env":"FLAI_MCP_HOST","cache":` +
		`"` + filepath.ToSlash(cache) + `"}` + "\n"
	req, err := Read(strings.NewReader(in))
	if err != nil {
		t.Fatal(err)
	}
	if want := valid(cache); req.Story != want.Story || req.Name != want.Name || req.Session != want.Session ||
		req.Inner != want.Inner || req.MCPURLEnv != want.MCPURLEnv || req.MCPTokenEnv != want.MCPTokenEnv ||
		req.MCPHostEnv != want.MCPHostEnv || req.Cache != want.Cache || req.Resume {
		t.Fatalf("read %+v", req)
	}
	if err := Check(req, env, isDirOnly(cache)); err != nil {
		t.Fatal(err)
	}
}

func TestReadRefusesWhatIsNotOneJSONObject(t *testing.T) {
	for _, c := range []struct{ name, in, want string }{
		{"empty", "", "is empty"},
		{"blank", " \n\t", "is empty"},
		{"malformed", `{"Story": "S-0507",}`, "not valid JSON at byte"},
		{"truncated", `{"Story": "S-0507"`, "ends before the object does"},
		{"two objects", `{"Story":"S-0507"} {"Story":"S-0508"}`, "more follows the object"},
		{"trailing garbage", `{"Story":"S-0507"}` + tokenValue, "more follows the object"},
		{"unknown field", `{"Story":"S-0507","inner":"claude-code","mcp_token":"FLAI_MCP_TOKEN"}`, `its field "mcp_token" is not one`},
		{"resume", `{"Story":"S-0507","Resume":true}`, `its field "Resume" is not one`},
		{"array", `[{"Story":"S-0507"}]`, "it is a JSON array, not an object"},
		{"null", `null`, "it is a JSON null, not an object"},
		{"string", `"` + tokenValue + `"`, "it is a JSON string, not an object"},
		{"wrong type", `{"Story":42,"cache":"` + urlValue + `"}`, `its field "Story" holds a JSON number where it takes a string`},
	} {
		t.Run(c.name, func(t *testing.T) {
			_, err := Read(strings.NewReader(c.in))
			oneLineWithoutSecrets(t, err)
			if err != nil && !strings.Contains(err.Error(), c.want) {
				t.Errorf("refused with %q, want it to say %q", err, c.want)
			}
		})
	}
}

func TestReadNeverQuotesAValueOfTheInput(t *testing.T) {
	for _, in := range []string{
		`{"Story":` + "123456789" + `}`,
		`{"Story":"x"` + "9876" + `}`,
	} {
		_, err := Read(strings.NewReader(in))
		oneLineWithoutSecrets(t, err)
		for _, v := range []string{"123456789", "9876"} {
			if err != nil && strings.Contains(err.Error(), v) {
				t.Errorf("Read(%s) quotes %s: %s", in, v, err)
			}
		}
	}
}

func TestCheckRefusesAStartRequestPilotCannotRun(t *testing.T) {
	cache := t.TempDir()
	for _, c := range []struct {
		name string
		edit func(*harness.Request)
		env  func(string) string
		want string
	}{
		{"no inner", func(r *harness.Request) { r.Inner = "" }, env, "names no inner harness; set inner to the harness pilot runs, one of claude-code, command"},
		{"unknown inner", func(r *harness.Request) { r.Inner = "codex" }, env, `inner "codex" is not a harness pilot can run; set inner to one of claude-code, command`},
		{"container inner", func(r *harness.Request) { r.Inner = "container" }, env, "pilot runs inside a container, so it cannot start another"},
		{"no cache", func(r *harness.Request) { r.Cache = "" }, env, "names no cache"},
		{"relative cache", func(r *harness.Request) { r.Cache = "cache" }, env, `cache "cache" is not an absolute path`},
		{"missing cache", func(r *harness.Request) { r.Cache = filepath.Join(cache, "missing") }, env, "is not a folder in the container"},
		{"no url variable", func(r *harness.Request) { r.MCPURLEnv = "" }, env, "names no mcp_url_env"},
		{"no token variable", func(r *harness.Request) { r.MCPTokenEnv = "" }, env, "names no mcp_token_env"},
		{"url variable unset", func(*harness.Request) {}, without("FLAI_MCP_URL"), `mcp_url_env names the variable "FLAI_MCP_URL", which is unset or empty`},
		{"token variable unset", func(*harness.Request) {}, without("FLAI_MCP_TOKEN"), `mcp_token_env names the variable "FLAI_MCP_TOKEN", which is unset or empty`},
		{"host variable unset", func(*harness.Request) {}, without("FLAI_MCP_HOST"), `mcp_host_env names the variable "FLAI_MCP_HOST", which is unset or empty`},
		{"token variable misnamed", func(r *harness.Request) { r.MCPTokenEnv = "FLAI_MCP_TOKN" }, env, `mcp_token_env names the variable "FLAI_MCP_TOKN"`},
	} {
		t.Run(c.name, func(t *testing.T) {
			req := valid(cache)
			c.edit(&req)
			err := Check(req, c.env, isDirOnly(cache))
			oneLineWithoutSecrets(t, err)
			if err != nil && !strings.Contains(err.Error(), c.want) {
				t.Errorf("refused with %q, want it to say %q", err, c.want)
			}
		})
	}
}

// I-0224: a probe is read, and passes the check with the story it runs in.
func TestAProbeIsReadAndPassesTheCheck(t *testing.T) {
	cache := t.TempDir()
	in := `{"Story":"S-0507","Name":"probe-S-0507","probe":true,"inner":"claude-code",` +
		`"mcp_url_env":"FLAI_MCP_URL","mcp_token_env":"FLAI_MCP_TOKEN","cache":"` + filepath.ToSlash(cache) + `"}`
	req, err := Read(strings.NewReader(in))
	if err != nil {
		t.Fatal(err)
	}
	if !req.Probe || req.Name != "probe-S-0507" || req.Story != "S-0507" {
		t.Fatalf("read %+v, want a probe of S-0507 named probe-S-0507", req)
	}
	if err := Check(req, env, isDirOnly(cache)); err != nil {
		t.Fatal(err)
	}
}

// I-0224: a probe is a trial of the image, never a strategic agent's run.
func TestCheckRefusesAProbeOfARoleOrAnItem(t *testing.T) {
	cache := t.TempDir()
	for _, c := range []struct {
		name string
		edit func(*harness.Request)
		want string
	}{
		{"role", func(r *harness.Request) { r.Story, r.Role = "", "plan" }, `is a probe and names the Role "plan"; a probe is a trial of the agent image, never a strategic agent's run, so drop Role`},
		{"item", func(r *harness.Request) { r.Item = "E-0029" }, `is a probe and names the Item "E-0029"; a probe is a trial of the agent image, never a strategic agent's run, so drop Item`},
	} {
		t.Run(c.name, func(t *testing.T) {
			req := valid(cache)
			req.Probe = true
			c.edit(&req)
			err := Check(req, env, isDirOnly(cache))
			oneLineWithoutSecrets(t, err)
			if err != nil && !strings.Contains(err.Error(), c.want) {
				t.Errorf("refused with %q, want it to say %q", err, c.want)
			}
		})
	}
}

// S-0705: a defined role's start request, its name in role_name, passes the
// check.
func TestCheckTakesADefinedRolesStartRequest(t *testing.T) {
	cache := t.TempDir()
	req := valid(cache)
	req.Story, req.Role, req.RoleName, req.Item, req.Name = "", "triage", "triager", "S-0042", "triager-S-0042"
	if err := Check(req, env, isDirOnly(cache)); err != nil {
		t.Fatal(err)
	}
}

// S-0705: role_name names a role the project defines, so Check refuses it
// with no Role, with a Story, on a probe, and as anything but a role's
// kebab-case name, which pilot joins to the role files' folder.
func TestCheckRefusesARoleNamePilotCannotRun(t *testing.T) {
	cache := t.TempDir()
	for _, c := range []struct {
		name string
		edit func(*harness.Request)
		want string
	}{
		{"no role", func(r *harness.Request) { r.Role = "" }, `names the role_name "triager" and no Role; a role the project defines runs in its file's verb, so set Role to it, or drop role_name`},
		{"a story", func(r *harness.Request) { r.Story = "S-0507" }, `names the role_name "triager" and the Story "S-0507"; a role the project defines never works a story, so drop Story`},
		{"a probe", func(r *harness.Request) { r.Probe, r.Role, r.Item = true, "", "" }, `is a probe and names the role_name "triager"; a probe is a trial of the agent image, never a strategic agent's run, so drop role_name`},
		{"a path", func(r *harness.Request) { r.RoleName = "../../etc/passwd" }, `role_name "../../etc/passwd" is not a role's name, which is kebab-case and names its role file, so holds no path; set role_name to the role file's name without .md`},
		{"a folder", func(r *harness.Request) { r.RoleName = "roles/triager" }, `role_name "roles/triager" is not a role's name`},
		{"dot dot", func(r *harness.Request) { r.RoleName = ".." }, `role_name ".." is not a role's name`},
		{"dot", func(r *harness.Request) { r.RoleName = "." }, `role_name "." is not a role's name`},
		{"not kebab-case", func(r *harness.Request) { r.RoleName = "Triager" }, `role_name "Triager" is not a role's name`},
	} {
		t.Run(c.name, func(t *testing.T) {
			req := valid(cache)
			req.Story, req.Role, req.RoleName, req.Item, req.Name = "", "triage", "triager", "S-0042", "triager-S-0042"
			c.edit(&req)
			err := Check(req, env, isDirOnly(cache))
			oneLineWithoutSecrets(t, err)
			if err != nil && !strings.Contains(err.Error(), c.want) {
				t.Errorf("refused with %q, want it to say %q", err, c.want)
			}
		})
	}
}

// ADR-0289, I-0265: the operator's outside folders in the start request are
// read and kept, pass the check, and reach the inner claude's arguments as
// --add-dir, in their order, through Main as the container runs it, with the
// empty Host pilot starts the inner harness with.
func TestTheOutsideFoldersReachTheInnerHarnessAsAddDir(t *testing.T) {
	dir := fake(t)
	cache := t.TempDir()
	outside := []string{filepath.Join(t.TempDir(), "kluster"), filepath.Join(t.TempDir(), "scratch")}
	in := `{"Story":"S-0507","Name":"agent-S-0507","inner":"claude-code","mcp_url_env":"FLAI_MCP_URL","mcp_token_env":"FLAI_MCP_TOKEN",` +
		`"cache":"` + filepath.ToSlash(cache) + `","outside":["` + filepath.ToSlash(outside[0]) + `","` + filepath.ToSlash(outside[1]) + `"]}`
	req, err := Read(strings.NewReader(in))
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(req.Outside, outside) {
		t.Fatalf("read outside %q, want %q", req.Outside, outside)
	}
	if err := Check(req, env, isDirOnly(cache)); err != nil {
		t.Fatal(err)
	}
	r := run(t, t.Context(), req, environ(dir), nil)
	if r.code != 0 {
		t.Fatalf("pilot exited %d; stderr: %s", r.code, r.stderr)
	}
	args := argv(t, dir)
	var dirs []string
	for i, a := range args {
		if a == "--add-dir" && i+1 < len(args) {
			dirs = append(dirs, args[i+1])
		}
	}
	if !slices.Equal(dirs, outside) {
		t.Errorf("the inner claude's --add-dir are %q, want the outside folders %q; argv %q", dirs, outside, args)
	}
}

// ADR-0289: an outside folder that is not an absolute path is refused, since
// the container mounts each at its own host path.
func TestCheckRefusesAnOutsideFolderThatIsNotAbsolute(t *testing.T) {
	cache := t.TempDir()
	for _, dir := range []string{"kluster", "../kluster", "--dangerously-skip-permissions", ""} {
		req := valid(cache)
		req.Outside = []string{cache, dir}
		err := Check(req, env, isDirOnly(cache))
		oneLineWithoutSecrets(t, err)
		if want := "outside folder " + strconv.Quote(dir) + " is not an absolute path"; err != nil && !strings.Contains(err.Error(), want) {
			t.Errorf("%q: refused with %q, want it to say %q", dir, err, want)
		}
	}
}

func TestCheckTakesARequestThatNamesNoHostVariable(t *testing.T) {
	cache := t.TempDir()
	req := valid(cache)
	req.MCPHostEnv = ""
	if err := Check(req, without("FLAI_MCP_HOST"), isDirOnly(cache)); err != nil {
		t.Fatal(err)
	}
}

func TestCheckReportsTheFirstRefusal(t *testing.T) {
	// no inner, no cache, and no variables: the inner is reported
	err := Check(harness.Request{}, func(string) string { return "" }, func(string) bool { return false })
	oneLineWithoutSecrets(t, err)
	if err != nil && !strings.Contains(err.Error(), "no inner harness") {
		t.Errorf("reported %q first", err)
	}
	// a cache that is missing is reported before the variables
	req := valid("/nowhere")
	req.MCPURLEnv = ""
	err = Check(req, env, func(string) bool { return false })
	if err == nil || !strings.Contains(err.Error(), "not a folder") {
		t.Errorf("reported %v first", err)
	}
}

// without is the fake environment with name unset.
func without(name string) func(string) string {
	return func(n string) string {
		if n == name {
			return ""
		}
		return env(n)
	}
}
