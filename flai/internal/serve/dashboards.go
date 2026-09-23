package serve

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Dashboard is a dashboard flai dashboard started outside any project
// (S-0101): nothing is registered for it, but flai serve offers it the
// repositories to import (S-0098), which otherwise go only to the dashboards
// served projects dial.
type Dashboard struct {
	URL     string `json:"url"`
	KeyFile string `json:"key_file"`
}

func (d Dir) dashboards() string { return filepath.Join(string(d), "dashboards.json") }

// Dashboards reads the dashboards recorded; none is an empty list.
func (d Dir) Dashboards() ([]Dashboard, error) {
	data, err := os.ReadFile(d.dashboards())
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var out []Dashboard
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("%s: %w", d.dashboards(), err)
	}
	return out, nil
}

// AddDashboard records a dashboard, or replaces the one at its address.
func (d Dir) AddDashboard(db Dashboard) error {
	if db.URL == "" || db.KeyFile == "" {
		return errors.New("a dashboard needs an address and a credential file")
	}
	all, err := d.Dashboards()
	if err != nil {
		return err
	}
	out := []Dashboard{db}
	for _, x := range all {
		if x.URL != db.URL {
			out = append(out, x)
		}
	}
	return d.write(d.dashboards(), out)
}

// ForgetDashboard removes the dashboard at an address; absent is not an error.
func (d Dir) ForgetDashboard(url string) error {
	all, err := d.Dashboards()
	if err != nil || len(all) == 0 {
		return err
	}
	out := []Dashboard{}
	for _, x := range all {
		if x.URL != url {
			out = append(out, x)
		}
	}
	return d.write(d.dashboards(), out)
}
