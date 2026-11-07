package harness

import (
	"bytes"
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/bytepunx/system-flow/flai/internal/mcphttp"
	"github.com/bytepunx/system-flow/flai/internal/protected"
)

// openCodeMCPEntry is a local MCP server as OpenCode's configuration names
// one: the program and its arguments, on, with its per-call timeout in
// milliseconds.
type openCodeMCPEntry struct {
	Type    string   `json:"type"`
	Command []string `json:"command"`
	Enabled bool     `json:"enabled"`
	Timeout int64    `json:"timeout"`
}

// openCodeConfigContent is OPENCODE_CONFIG_CONTENT for r's run (ADR-0260
// §1): flai's MCP server alone, mcp.flai, run on stdio as flai mcp --agent
// <name>, with mcphttp.ClientTimeout as its timeout, so that the calls that
// hold, wait_for_events and permission_prompt, run past OpenCode's default.
// It names no provider, no model, and no key.
func openCodeConfigContent(r Request) (string, error) {
	command := []string{r.Flai, "mcp"}
	if r.Name != "" {
		command = append(command, "--agent", r.Name)
	}
	entry := openCodeMCPEntry{Type: "local", Command: command, Enabled: true, Timeout: mcphttp.ClientTimeout.Milliseconds()}
	data, err := json.Marshal(map[string]map[string]openCodeMCPEntry{"mcp": {"flai": entry}})
	if err != nil {
		return "", fmt.Errorf("the OpenCode configuration for flai's MCP server: %w", err)
	}
	return string(data), nil
}

// openCodeProjectFiles are the project's configuration files at its root,
// in the order OpenCode merges them, each over the one before: it finds
// opencode.jsonc then opencode.json in a folder and reverses the list
// (config/paths.ts at v1.19.5).
var openCodeProjectFiles = []string{"opencode.json", "opencode.jsonc"}

// openCodeKept is the project's permission that a permission object from
// flai would replace whole, as OpenCode merges: the whole permission given
// as one action, or its edit given as one action. At most one is set.
type openCodeKept struct {
	whole, edit string
}

// openCodeProjectPermission reads what of the project's permission flai must
// keep from its opencode.json and opencode.jsonc in fsys, its root, merged in
// OpenCode's order (ADR-0260 §3). A later permission given as one action
// replaces an earlier one whole; a later object replaces an earlier action
// and keeps an earlier edit it does not give. An absent or empty file gives
// nothing; one that is not JSONC refuses the start, naming its path under
// root, since OpenCode would refuse it too.
func openCodeProjectPermission(root string, fsys fs.FS) (openCodeKept, error) {
	var kept openCodeKept
	for _, name := range openCodeProjectFiles {
		data, err := fs.ReadFile(fsys, name)
		if errors.Is(err, fs.ErrNotExist) || (err == nil && len(data) == 0) {
			continue
		}
		path := filepath.Join(root, name)
		if err != nil {
			return openCodeKept{}, fmt.Errorf("opencode cannot start: flai could not read the project's %s for its permission: %w; make it readable, or remove it", path, err)
		}
		var config map[string]json.RawMessage
		if err := json.Unmarshal(stripJSONC(data), &config); err != nil {
			return openCodeKept{}, fmt.Errorf("opencode cannot start: the project's %s is not a JSONC object, which OpenCode would refuse too: %w; fix it", path, err)
		}
		raw, ok := config["permission"]
		if !ok {
			continue
		}
		if action, ok := jsonString(raw); ok {
			kept = openCodeKept{whole: action}
			continue
		}
		var keys map[string]json.RawMessage
		if json.Unmarshal(raw, &keys) != nil {
			continue
		}
		kept.whole = ""
		if edit, ok := keys["edit"]; ok {
			kept.edit, _ = jsonString(edit)
		}
	}
	return kept, nil
}

// jsonString is raw's string, and whether raw is a JSON string.
func jsonString(raw json.RawMessage) (string, bool) {
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return "", false
	}
	return s, true
}

// openCodeDenies are the edit patterns that deny writing each path of l, in
// l's order (ADR-0260 §2): a folder as <name>/* and */<name>/*, a file as
// <name> and */<name>, a root folder as <path>/*, and a root file as
// itself. OpenCode matches a pattern against the path from the worktree's
// root, and its * crosses /.
func openCodeDenies(l protected.List) []string {
	var out []string
	for _, f := range l.Folders {
		out = append(out, f+"/*", "*/"+f+"/*")
	}
	for _, n := range l.Files {
		out = append(out, n, "*/"+n)
	}
	for _, f := range l.RootFolders {
		out = append(out, f+"/*")
	}
	return append(out, l.RootFiles...)
}

// openCodePermission is OPENCODE_PERMISSION for a run under deny_protected
// (ADR-0260 §2 and §3): an edit object with a deny rule for every path
// OpenCode's protected list holds, after the project's own action as the
// rule for *, when its opencode.json or opencode.jsonc in fsys, its root,
// gives its permission or its edit as one action, which flai's object would
// otherwise replace. A whole permission given as one action is kept as the
// top-level * too, so that the other tools keep it.
func openCodePermission(root string, fsys fs.FS) (string, error) {
	kept, err := openCodeProjectPermission(root, fsys)
	if err != nil {
		return "", err
	}
	var edit, permission orderedObject
	if action := cmp.Or(kept.whole, kept.edit); action != "" {
		edit = append(edit, member{"*", action})
	}
	for _, p := range openCodeDenies(protected.For(protected.OpenCode)) {
		edit = append(edit, member{p, "deny"})
	}
	if kept.whole != "" {
		permission = append(permission, member{"*", kept.whole})
	}
	data, err := json.Marshal(append(permission, member{"edit", edit}))
	if err != nil {
		return "", fmt.Errorf("the OpenCode permission that denies the protected paths: %w", err)
	}
	return string(data), nil
}

// orderedObject is a JSON object whose keys keep the order they were added
// in, since OpenCode reads a permission's rules in the order written and the
// last that matches wins.
type orderedObject []member

// member is one key of an orderedObject and its value.
type member struct {
	key   string
	value any
}

// MarshalJSON writes o's members in order.
func (o orderedObject) MarshalJSON() ([]byte, error) {
	var b bytes.Buffer
	b.WriteByte('{')
	for i, m := range o {
		if i > 0 {
			b.WriteByte(',')
		}
		key, err := json.Marshal(m.key)
		if err != nil {
			return nil, err
		}
		value, err := json.Marshal(m.value)
		if err != nil {
			return nil, err
		}
		b.Write(key)
		b.WriteByte(':')
		b.Write(value)
	}
	b.WriteByte('}')
	return b.Bytes(), nil
}

// stripJSONC makes data, JSONC as OpenCode reads it, into JSON: each comment
// outside a string, // to the end of its line or /* to */, becomes a space,
// and a comma with nothing but space before a closing } or ] goes. A /* with
// no */ is left as it stands, for the JSON parser to refuse, as OpenCode's
// does.
func stripJSONC(data []byte) []byte {
	out := make([]byte, 0, len(data))
	inString, escaped := false, false
	for i := 0; i < len(data); i++ {
		c := data[i]
		if inString {
			out = append(out, c)
			switch {
			case escaped:
				escaped = false
			case c == '\\':
				escaped = true
			case c == '"':
				inString = false
			}
			continue
		}
		next := byte(0)
		if i+1 < len(data) {
			next = data[i+1]
		}
		switch {
		case c == '"':
			inString = true
			out = append(out, c)
		case c == '/' && next == '/':
			end := bytes.IndexByte(data[i:], '\n')
			if end < 0 {
				end = len(data) - i
			}
			out = append(out, ' ')
			i += end - 1
		case c == '/' && next == '*':
			end := bytes.Index(data[i+2:], []byte("*/"))
			if end < 0 {
				return dropTrailingCommas(append(out, data[i:]...))
			}
			out = append(out, ' ')
			i += 2 + end + 1
		default:
			out = append(out, c)
		}
	}
	return dropTrailingCommas(out)
}

// dropTrailingCommas is data, JSON with no comments, less each comma
// outside a string that only space parts from a closing } or ].
func dropTrailingCommas(data []byte) []byte {
	out := make([]byte, 0, len(data))
	inString, escaped := false, false
	for i, c := range data {
		switch {
		case inString && escaped:
			escaped = false
		case inString && c == '\\':
			escaped = true
		case c == '"':
			inString = !inString
		case !inString && c == ',':
			rest := bytes.TrimLeft(data[i+1:], " \t\r\n")
			if len(rest) > 0 && (rest[0] == '}' || rest[0] == ']') {
				continue
			}
		}
		out = append(out, c)
	}
	return out
}
