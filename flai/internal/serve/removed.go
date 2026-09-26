package serve

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
)

// The project folders the operator removed from the dashboard (S-0123):
// flai serve does not serve one of them from below a folder named for
// import or the folder it was started in, which it otherwise would for as
// long as the project is there. flai serve project add takes one off.

func (d Dir) removed() string { return filepath.Join(string(d), "removed.json") }

// RemovedRoots reads the list of removed project folders; none is an empty list.
func (d Dir) RemovedRoots() ([]string, error) {
	data, err := os.ReadFile(d.removed())
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var out []string
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("%s: %w", d.removed(), err)
	}
	return out, nil
}

// RemoveRoot adds a project folder to the list; there already is not an error.
func (d Dir) RemoveRoot(root string) error {
	all, err := d.RemovedRoots()
	if err != nil {
		return err
	}
	if slices.Contains(all, root) {
		return nil
	}
	all = append(all, root)
	slices.Sort(all)
	return d.write(d.removed(), all)
}

// RestoreRoot takes a project folder off the list, and says whether it was on it.
func (d Dir) RestoreRoot(root string) (bool, error) {
	all, err := d.RemovedRoots()
	if err != nil || !slices.Contains(all, root) {
		return false, err
	}
	return true, d.write(d.removed(), slices.DeleteFunc(all, func(r string) bool { return r == root }))
}

// RemovedSet is the list as a set; a list that does not read is empty, so
// that a broken file serves what it would without one rather than nothing.
func (d Dir) RemovedSet() (map[string]bool, error) {
	all, err := d.RemovedRoots()
	set := map[string]bool{}
	for _, r := range all {
		set[r] = true
	}
	return set, err
}
