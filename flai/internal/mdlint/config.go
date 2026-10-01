package mdlint

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
)

// ConfigFiles are the markdownlint configuration files flai reads from a
// project's root, in markdownlint-cli2's order of precedence. The cli2 files
// carry the rules under their config key.
var ConfigFiles = []string{
	".markdownlint-cli2.jsonc", ".markdownlint-cli2.yaml",
	".markdownlint.jsonc", ".markdownlint.json", ".markdownlint.yaml", ".markdownlint.yml",
}

// Config is a project's markdownlint configuration as flai applies it: which
// of the rules flai implements are on, and their options.
type Config struct {
	// Path is the file it was read from, relative to the project root.
	Path    string
	enabled map[string]bool
	options map[string]map[string]any
}

// Load reads the project's markdownlint configuration from root. A project
// without one gets nil and no error: flai lints nothing there.
func Load(root string) (*Config, error) {
	for _, name := range ConfigFiles {
		data, err := os.ReadFile(filepath.Join(root, name))
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, err
		}
		c, err := Parse(data, strings.HasPrefix(name, ".markdownlint-cli2"), strings.Contains(name, ".json"))
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		c.Path = name
		return c, nil
	}
	return nil, nil
}

// Parse reads markdownlint configuration in YAML, JSON, or JSONC. Keys are
// applied in order, as markdownlint applies them: default first, then each
// rule, alias, or tag. cli2 says the rules are under the config key; json
// that the file is JSON or JSONC, whose comments are taken out first.
func Parse(data []byte, cli2, json bool) (*Config, error) {
	if json {
		data = stripComments(data)
	}
	var top yaml.MapSlice
	if err := yaml.UnmarshalWithOptions(data, &top, yaml.UseOrderedMap()); err != nil {
		return nil, err
	}
	if cli2 {
		var rules yaml.MapSlice
		for _, it := range top {
			if k, _ := it.Key.(string); k == "config" {
				rules, _ = it.Value.(yaml.MapSlice)
			}
		}
		top = rules
	}
	c := &Config{enabled: map[string]bool{}, options: map[string]map[string]any{}}
	def := true
	for _, it := range top {
		if k, _ := it.Key.(string); strings.EqualFold(k, "default") {
			if b, ok := it.Value.(bool); ok {
				def = b
			}
		}
	}
	for _, r := range rules {
		c.enabled[r.id] = def
	}
	for _, it := range top {
		key, _ := it.Key.(string)
		if strings.EqualFold(key, "default") || strings.EqualFold(key, "extends") || strings.HasPrefix(key, "$") {
			continue
		}
		on, opts := true, map[string]any(nil)
		switch v := it.Value.(type) {
		case bool:
			on = v
		case yaml.MapSlice:
			opts = map[string]any{}
			for _, o := range v {
				if k, ok := o.Key.(string); ok {
					opts[k] = o.Value
				}
			}
		case nil:
			on = false
		}
		for _, id := range resolve(key) {
			c.enabled[id] = on
			if opts != nil {
				c.options[id] = opts
			}
		}
	}
	return c, nil
}

// Enabled reports whether a rule, by ID, is on.
func (c *Config) Enabled(id string) bool { return c != nil && c.enabled[id] }

func (c *Config) intOpt(id, name string, def int) int {
	switch v := c.options[id][name].(type) {
	case uint64:
		return int(v)
	case int64:
		return int(v)
	case int:
		return v
	case float64:
		return int(v)
	}
	return def
}

func (c *Config) boolOpt(id, name string, def bool) bool {
	if v, ok := c.options[id][name].(bool); ok {
		return v
	}
	return def
}

// strOpt returns the option, present reporting whether the file set it at
// all: an empty string is not the same as no value for front_matter_title.
func (c *Config) strOpt(id, name, def string) (string, bool) {
	v, ok := c.options[id][name]
	if !ok || v == nil {
		return def, false
	}
	return fmt.Sprint(v), true
}

func (c *Config) listOpt(id, name string) []string {
	l, _ := c.options[id][name].([]any)
	var out []string
	for _, v := range l {
		out = append(out, fmt.Sprint(v))
	}
	return out
}

// resolve turns a configuration key, a rule ID, alias, or tag, into the IDs
// of the rules flai implements that it names.
func resolve(key string) []string {
	k := strings.ToLower(key)
	var out []string
	for _, r := range rules {
		if strings.ToLower(r.id) == k || contains(r.aliases, k) || contains(r.tags, k) {
			out = append(out, r.id)
		}
	}
	return out
}

func contains(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}

// stripComments removes JSONC comments outside strings so the YAML parser,
// which reads JSON, reads the file.
func stripComments(data []byte) []byte {
	s := string(data)
	if !strings.Contains(s, "//") && !strings.Contains(s, "/*") {
		return data
	}
	var b strings.Builder
	inStr, esc := false, false
	for i := 0; i < len(s); i++ {
		ch := s[i]
		switch {
		case inStr:
			b.WriteByte(ch)
			switch {
			case esc:
				esc = false
			case ch == '\\':
				esc = true
			case ch == '"':
				inStr = false
			}
		case ch == '"':
			inStr = true
			b.WriteByte(ch)
		case ch == '/' && i+1 < len(s) && s[i+1] == '/':
			for i < len(s) && s[i] != '\n' {
				i++
			}
			if i < len(s) {
				b.WriteByte('\n')
			}
		case ch == '/' && i+1 < len(s) && s[i+1] == '*':
			end := strings.Index(s[i+2:], "*/")
			if end < 0 {
				return []byte(b.String())
			}
			b.WriteString(strings.Repeat("\n", strings.Count(s[i:i+2+end+2], "\n")))
			i += 2 + end + 1
		default:
			b.WriteByte(ch)
		}
	}
	return []byte(b.String())
}
