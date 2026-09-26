package serve

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// Removed is the notification flai serve sends a dashboard just before it
// drops a project that is no longer registered (S-0118), so that the
// dashboard's switcher drops it too instead of showing it as not connected.
const Removed = "removed"

// Unavailable says why a registered project cannot be served (S-0118): its
// folder is gone, it has no system-flow.yaml or one that does not load or has
// no key, or the agent credential cannot be read. Empty when it can be.
func (e Entry) Unavailable() string {
	info, err := os.Stat(e.Root)
	switch {
	case errors.Is(err, os.ErrNotExist):
		return "the folder is gone"
	case err != nil:
		return "the folder cannot be read: " + err.Error()
	case !info.IsDir():
		return "it is not a folder"
	}
	m, err := manifest.Load(filepath.Join(e.Root, manifest.File))
	switch {
	case errors.Is(err, os.ErrNotExist):
		return "it has no " + manifest.File
	case err != nil:
		return manifest.File + " does not load: " + err.Error()
	case m.Key == "":
		return manifest.File + " has no key"
	}
	if key, err := os.ReadFile(e.KeyFile); err != nil || strings.TrimSpace(string(key)) == "" {
		return "the agent credential " + e.KeyFile + " cannot be read; flai dashboard makes it"
	}
	return ""
}

// keyTaken is the entry other than the one at root that has key, if any: two
// roots with one key would take turns replacing each other's connection in
// the dashboard (S-0118).
func keyTaken(all []Entry, key, root string) error {
	for _, x := range all {
		if x.Key == key && x.Root != root {
			return fmt.Errorf("the key %s is already served for %s; give this project another key in its %s, or flai serve project remove %s", key, x.Root, manifest.File, x.Root)
		}
	}
	return nil
}
