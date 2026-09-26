package serve

import (
	"fmt"
	"path/filepath"
	"sort"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Found is a system-flow project flai serve found below a folder rather than
// in its registry: below the folder it was started in (S-0102), or a git
// repository with a system-flow.yaml below a folder named for import
// (S-0117). It is served as registered projects are, for as long as it is
// there, without being written to the registry.
type Found struct {
	Entry
	// Imported is true when it was found below a folder named for import.
	Imported bool `json:"imported,omitempty"`
	// Reason says why it is not served; empty when it is.
	Reason string `json:"reason,omitempty"`
}

// FindBelow lists the system-flow projects below folder, and the git
// repositories with a system-flow.yaml below importRoots, each once, by root,
// with the key and name its manifest gives or why the manifest does not load.
// Either may be empty.
func FindBelow(folder string, importRoots []string) []Found {
	var out []Found
	seen := map[string]bool{}
	add := func(root string, imported bool) {
		if seen[root] {
			return
		}
		seen[root] = true
		f := Found{Entry: Entry{Root: root, Name: filepath.Base(root)}, Imported: imported}
		m, err := manifest.Load(filepath.Join(root, manifest.File))
		if err != nil {
			f.Reason = fmt.Sprintf("its %s does not load: %s", manifest.File, err)
		} else {
			f.Key = m.Key
			if m.Name != "" {
				f.Name = m.Name
			}
		}
		out = append(out, f)
	}
	if folder != "" {
		for _, root := range workitem.FindProjects(folder) {
			add(root, false)
		}
	}
	for _, root := range repositories(importRoots) {
		if exists(filepath.Join(root, manifest.File)) {
			add(root, true)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Root < out[j].Root })
	return out
}

// Place decides which of the found projects are served, and for which
// dashboard: a registered one is left to the registry and left out; one
// whose manifest does not load or has no key, or whose key a registered or
// earlier project has, is not served, and says why; the rest are served for
// the first dashboard, by address, of dashboards, and with none known they
// are not served either.
func Place(found []Found, registered []Entry, dashboards map[string]Entry) []Found {
	roots, keys := map[string]bool{}, map[string]string{}
	for _, e := range registered {
		roots[e.Root], keys[e.Key] = true, e.Root
	}
	var d *Entry
	if len(dashboards) > 0 {
		urls := make([]string, 0, len(dashboards))
		for u := range dashboards {
			urls = append(urls, u)
		}
		sort.Strings(urls)
		first := dashboards[urls[0]]
		d = &first
	}
	var out []Found
	for _, f := range found {
		switch {
		case roots[f.Root]:
			continue
		case f.Reason != "":
		case f.Key == "":
			f.Reason = fmt.Sprintf("its %s has no key (flai check says how to add one)", manifest.File)
		case keys[f.Key] != "":
			f.Reason = fmt.Sprintf("its key %s is served already, for %s", f.Key, keys[f.Key])
		case d == nil:
			f.Reason = "no dashboard is known to serve it for (flai dashboard starts one)"
		default:
			f.URL, f.KeyFile = d.URL, d.KeyFile
			keys[f.Key] = f.Root
		}
		out = append(out, f)
	}
	return out
}
