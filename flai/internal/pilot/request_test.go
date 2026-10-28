package pilot

import (
	"path/filepath"
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
