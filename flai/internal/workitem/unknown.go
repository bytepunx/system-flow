package workitem

import (
	"fmt"
	"log/slog"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"sync"
)

// Field is a top-level front-matter key this flai does not know, kept as it
// was written so that a write gives it back unchanged (S-0181). A newer flai
// wrote it; an older one reads past it with a warning, and flai check
// reports it.
type Field struct {
	Name string
	Raw  string // the key's line and every line of its value, each ending in a newline
}

// topKey matches a line that starts a top-level mapping key.
var topKey = regexp.MustCompile(`^([A-Za-z0-9_][A-Za-z0-9_.-]*)\s*:`)

// UnknownFields returns the top-level keys of front matter fm that v's
// struct does not name in its yaml tags, in the order written, each with its
// lines. A value's lines run to the next top-level key, so indented blocks
// and sequences written at column zero stay with their key.
func UnknownFields(fm string, v any) []Field {
	known := yamlKeys(reflect.TypeOf(v))
	var out []Field
	var cur *Field
	for _, line := range strings.SplitAfter(fm, "\n") {
		if line == "" {
			continue
		}
		if m := topKey.FindStringSubmatch(line); m != nil {
			cur = nil
			if !known[m[1]] {
				out = append(out, Field{Name: m[1]})
				cur = &out[len(out)-1]
			}
		}
		if cur != nil {
			if !strings.HasSuffix(line, "\n") {
				line += "\n"
			}
			cur.Raw += line
		}
	}
	return out
}

// FieldNames lists the fields' names.
func FieldNames(fields []Field) []string {
	names := make([]string, len(fields))
	for i, f := range fields {
		names[i] = f.Name
	}
	return names
}

// UnknownFieldErrors is what flai check reports for each unknown field.
func UnknownFieldErrors(fields []Field) []string {
	var errs []string
	for _, f := range fields {
		errs = append(errs, fmt.Sprintf("field %q is not one this flai knows: a newer flai wrote it, or it is misspelled", f.Name))
	}
	return errs
}

// warned is what WarnUnknown has said, by path and field names, so that a
// long-running flai that lists often says it once.
var warned sync.Map

// WarnUnknown logs that a listing read past the unknown fields of the file at
// path, which this flai keeps but does not act on, once per process for each
// file and set of fields.
func WarnUnknown(path string, fields []Field) {
	if len(fields) == 0 {
		return
	}
	if _, seen := warned.LoadOrStore(path+"\x00"+strings.Join(FieldNames(fields), ","), true); seen {
		return
	}
	slog.Default().Warn("front matter has fields this flai does not know; it reads past them and keeps them on write. Upgrade flai to act on them",
		"component", "workitem", "path", path, "fields", strings.Join(FieldNames(fields), ", "))
}

// WriteFields appends the kept fields to front matter being written.
func WriteFields(b *strings.Builder, fields []Field) {
	for _, f := range fields {
		b.WriteString(f.Raw)
	}
}

func yamlKeys(t reflect.Type) map[string]bool {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	keys := map[string]bool{}
	for i := range t.NumField() {
		tag := t.Field(i).Tag.Get("yaml")
		name, _, _ := strings.Cut(tag, ",")
		if name == "" || name == "-" {
			continue
		}
		keys[name] = true
	}
	return keys
}

// KnownFields lists the front-matter keys v's struct reads, sorted.
func KnownFields(v any) []string {
	keys := yamlKeys(reflect.TypeOf(v))
	out := make([]string, 0, len(keys))
	for k := range keys {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// fieldTypes names, for each item field that only some types of item carry,
// the types that do. Validate refuses the field on any other type, so a flai
// that carries one on a new type must raise flai.minimum: an older flai
// refuses the item. front-matter-fields.txt lists it for that reason
// (TestFieldsFileIsTheCode, release.RaiseMinimum).
var fieldTypes = map[string][]string{
	"after":         {Story, Task}, // a task's since S-0176
	"agent":         {Story},
	"cost_of_delay": {Epic, Story},
	"draft":         {Story},
	"finalized":     {Story},
	"forecast":      {Story},
	"topics":        {Epic, Story},
}

// Carries reports whether an item of type typ may carry the field: true for
// every field that is not restricted to some types.
func Carries(typ, field string) bool {
	types, ok := fieldTypes[field]
	return !ok || contains(types, typ)
}

// FieldTypes renders fieldTypes as one line, field=type,type, sorted by
// field: what front-matter-fields.txt lists after "item.types:".
func FieldTypes() string {
	fields := make([]string, 0, len(fieldTypes))
	for f := range fieldTypes {
		fields = append(fields, f)
	}
	sort.Strings(fields)
	for i, f := range fields {
		fields[i] = f + "=" + strings.Join(fieldTypes[f], ",")
	}
	return strings.Join(fields, " ")
}
