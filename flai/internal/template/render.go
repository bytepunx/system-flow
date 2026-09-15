package template

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"
	"unicode"
)

// Options controls a render.
type Options struct {
	Vars   map[string]any    // variable values by name
	Layout map[string]string // chosen folder name per layout key; defaults from the manifest when absent
	Source Source            // recorded in system-flow.yaml via .template.*
	Force  bool              // overwrite existing files
	Now    time.Time         // render time; zero means time.Now()
}

// Result lists what a render did.
type Result struct {
	Written []string
	Skipped []string // existed already and Force was false
}

// Funcs are the functions available inside template files and defaults.
var Funcs = template.FuncMap{
	"initials": Initials,
	"slug":     Slug,
	"upper":    strings.ToUpper,
	"lower":    strings.ToLower,
	"quote":    strconv.Quote, // YAML-safe double-quoted scalar
}

// Data builds the template data map: variables, layout, template source, and
// timestamps. It is exported so commands can preview it.
func Data(m Manifest, opt Options) map[string]any {
	now := opt.Now
	if now.IsZero() {
		now = time.Now()
	}
	now = now.UTC()
	layout := map[string]any{}
	for k, def := range m.Layout {
		layout[k] = def
		if v, ok := opt.Layout[k]; ok && v != "" {
			layout[k] = v
		}
	}
	data := map[string]any{}
	for k, v := range opt.Vars {
		data[k] = v
	}
	data["layout"] = layout
	data["today"] = now.Format("2006-01-02")
	data["now"] = now.Format(time.RFC3339)
	repo := opt.Source.Repo
	if opt.Source.Local {
		repo = opt.Source.Dir // absolute, so the manifest is meaningful from anywhere
	}
	data["template"] = map[string]any{
		"repo":    repo,
		"ref":     opt.Source.Ref,
		"version": m.Version,
	}
	return data
}

// EvalDefault renders a variable's default, which may reference earlier
// variables and the template functions.
func EvalDefault(v Variable, data map[string]any) (string, error) {
	if !strings.Contains(v.Default, "{{") {
		return v.Default, nil
	}
	return renderString("default:"+v.Name, v.Default, data)
}

// Render writes the template's root/ tree into dest.
func Render(m Manifest, tplDir, dest string, opt Options) (Result, error) {
	data := Data(m, opt)
	root := filepath.Join(tplDir, RootDir)
	var res Result
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		relSlash := filepath.ToSlash(rel)
		if ignored(m.Render.Ignore, relSlash) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		target := filepath.Join(dest, filepath.FromSlash(rewriteLayout(m, data, relSlash)))
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if strings.HasSuffix(target, m.Render.Suffix) {
			target = strings.TrimSuffix(target, m.Render.Suffix)
			out, err := renderString(relSlash, string(content), data)
			if err != nil {
				return err
			}
			content = []byte(out)
		}
		display, _ := filepath.Rel(dest, target)
		if _, err := os.Stat(target); err == nil && !opt.Force {
			res.Skipped = append(res.Skipped, display)
			return nil
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(target, content, info.Mode().Perm()); err != nil {
			return err
		}
		res.Written = append(res.Written, display)
		return nil
	})
	return res, err
}

// rewriteLayout swaps the first path segment when it is a layout default
// whose chosen name differs.
func rewriteLayout(m Manifest, data map[string]any, rel string) string {
	first, rest, hasRest := strings.Cut(rel, "/")
	layout, _ := data["layout"].(map[string]any)
	for key, def := range m.Layout {
		if first != def {
			continue
		}
		chosen, _ := layout[key].(string)
		if chosen == "" || chosen == def {
			return rel
		}
		if hasRest {
			return chosen + "/" + rest
		}
		return chosen
	}
	return rel
}

func ignored(patterns []string, rel string) bool {
	for _, p := range patterns {
		if p == rel {
			return true
		}
		if ok, _ := filepath.Match(p, rel); ok {
			return true
		}
	}
	return false
}

// RenderText renders one template string with the standard functions.
func RenderText(name, text string, data map[string]any) (string, error) {
	return renderString(name, text, data)
}

func renderString(name, text string, data map[string]any) (string, error) {
	t, err := template.New(name).Funcs(Funcs).Option("missingkey=error").Parse(text)
	if err != nil {
		return "", fmt.Errorf("template %s: %w", name, err)
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("template %s: %w", name, err)
	}
	return buf.String(), nil
}

var nonAlnum = regexp.MustCompile(`[^a-z0-9]+`)

// Slug lowercases and replaces runs of non-alphanumerics with dashes.
func Slug(s string) string {
	s = nonAlnum.ReplaceAllString(strings.ToLower(s), "-")
	return strings.Trim(s, "-")
}

// Initials returns the first letter of each word, lowercased: "system-flow" is "sf".
func Initials(s string) string {
	var b strings.Builder
	inWord := false
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			if !inWord {
				b.WriteRune(unicode.ToLower(r))
				inWord = true
			}
		} else {
			inWord = false
		}
	}
	return b.String()
}
