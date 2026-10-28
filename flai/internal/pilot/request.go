// Package pilot is the entrypoint of flai's agent container images (S-0507,
// ADR-0212). It reads one start request on standard input, a harness.Request
// that names the inner harness to run, the mounted cache, and the variables
// that hold the host's HTTP MCP server's URL and token, checks it, and starts
// the inner harness through its adapter, relaying what it writes. The request
// names variables, never their values, and pilot never prints a value.
package pilot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/harness"
)

// ExitUsage is the exit code of a refused start request.
const ExitUsage = 2

// container is the inner harness pilot refuses whatever the adapters hold:
// pilot runs inside a container, so it cannot start another (S-0507).
const container = "container"

// unknownField is how encoding/json begins the error for a field the
// decoded struct does not have, when unknown fields are disallowed; the
// package has no typed error for it.
const unknownField = "json: unknown field "

// Read decodes the start request: exactly one JSON object, whose fields are
// all harness.Request's. It refuses empty input, anything after the object,
// and a field Request does not have, so that a misspelt field is not lost.
// Its error is one line and never quotes a value from the input.
func Read(r io.Reader) (harness.Request, error) {
	var raw json.RawMessage
	dec := json.NewDecoder(r)
	if err := dec.Decode(&raw); err != nil {
		if errors.Is(err, io.EOF) {
			return harness.Request{}, errors.New("the start request on standard input is empty; write one JSON object, a harness.Request with inner, cache, mcp_url_env, and mcp_token_env, to pilot's standard input")
		}
		return harness.Request{}, refused(why(err))
	}
	if _, err := dec.Token(); !errors.Is(err, io.EOF) {
		return harness.Request{}, refused("more follows the object")
	}
	if raw[0] != '{' {
		return harness.Request{}, refused("it is a JSON " + kind(raw[0]) + ", not an object")
	}
	var req harness.Request
	strict := json.NewDecoder(bytes.NewReader(raw))
	strict.DisallowUnknownFields()
	if err := strict.Decode(&req); err != nil {
		return harness.Request{}, refused(why(err))
	}
	return req, nil
}

// refused is Read's error for a request that is not one JSON object, for
// the reason given.
func refused(reason string) error {
	return fmt.Errorf("the start request on standard input is not one JSON object: %s; write one JSON object, a harness.Request, to pilot's standard input", reason)
}

// why says what is wrong with the JSON in err's terms, without the
// input's values: encoding/json's own messages quote a number or a
// character of the input.
func why(err error) string {
	var syntax *json.SyntaxError
	var typ *json.UnmarshalTypeError
	switch {
	case errors.Is(err, io.ErrUnexpectedEOF):
		return "it ends before the object does"
	case errors.As(err, &syntax):
		return fmt.Sprintf("it is not valid JSON at byte %d", syntax.Offset)
	case errors.As(err, &typ):
		// Value is "number 5" for a number: its first word is the kind.
		value, _, _ := strings.Cut(typ.Value, " ")
		return fmt.Sprintf("its field %s holds a JSON %s where it takes a %s", strconv.Quote(typ.Field), value, typ.Type)
	case strings.HasPrefix(err.Error(), unknownField):
		name, uerr := strconv.Unquote(strings.TrimPrefix(err.Error(), unknownField))
		if uerr != nil {
			name = strings.TrimPrefix(err.Error(), unknownField)
		}
		return fmt.Sprintf("its field %s is not one harness.Request has; check the field's spelling", strconv.Quote(name))
	default:
		return "it cannot be decoded"
	}
}

// kind names the JSON value that begins with c.
func kind(c byte) string {
	switch c {
	case '[':
		return "array"
	case '"':
		return "string"
	case 't', 'f':
		return "boolean"
	case 'n':
		return "null"
	default:
		return "number"
	}
}

// Check refuses a start request pilot cannot run: one that names no inner
// harness, or one pilot cannot start; no cache, or one that is not an
// absolute path to a folder isDir finds; or not the variables of the host's
// HTTP MCP server's URL and token, which pilot always reaches, or a variable
// it names that getenv finds unset or empty. It reports the first refusal,
// in one line, naming the field and the variable but never a variable's
// value.
func Check(req harness.Request, getenv func(string) string, isDir func(string) bool) error {
	names := innerNames()
	switch {
	case req.Inner == "":
		return fmt.Errorf("the start request names no inner harness; set inner to the harness pilot runs, one of %s", strings.Join(names, ", "))
	case req.Inner == container:
		return fmt.Errorf("the start request's inner is %s, and pilot runs inside a container, so it cannot start another; set inner to one of %s", container, strings.Join(names, ", "))
	case !slices.Contains(names, req.Inner):
		return fmt.Errorf("the start request's inner %s is not a harness pilot can run; set inner to one of %s", strconv.Quote(req.Inner), strings.Join(names, ", "))
	case req.Cache == "":
		return errors.New("the start request names no cache; set cache to the absolute path of the folder mounted to keep the harness's configuration and sessions")
	case !filepath.IsAbs(req.Cache):
		return fmt.Errorf("the start request's cache %s is not an absolute path; set cache to the mounted folder's absolute path", strconv.Quote(req.Cache))
	case !isDir(req.Cache):
		return fmt.Errorf("the start request's cache %s is not a folder in the container; mount the cache folder there, or set cache to where it is mounted", strconv.Quote(req.Cache))
	case req.MCPURLEnv == "":
		return errors.New("the start request names no mcp_url_env; pilot reaches the host's HTTP MCP server, so set mcp_url_env to the name of the variable that holds its URL")
	case req.MCPTokenEnv == "":
		return errors.New("the start request names no mcp_token_env; pilot reaches the host's HTTP MCP server, so set mcp_token_env to the name of the variable that holds the agent's bearer token")
	}
	for _, v := range []struct{ field, name string }{{"mcp_url_env", req.MCPURLEnv}, {"mcp_token_env", req.MCPTokenEnv}, {"mcp_host_env", req.MCPHostEnv}} {
		if v.name != "" && getenv(v.name) == "" {
			return fmt.Errorf("the start request's %s names the variable %s, which is unset or empty in pilot's environment; pass it to the container, or name the variable that holds it", v.field, strconv.Quote(v.name))
		}
	}
	return nil
}

// innerNames are the inner harnesses pilot can run, sorted: every adapter
// but a container.
func innerNames() []string {
	return slices.DeleteFunc(harness.Names(), func(n string) bool { return n == container })
}
