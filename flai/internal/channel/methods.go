package channel

import (
	"context"
	"encoding/json"
	"path/filepath"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// ProjectInfo is what project.info answers: the project as its manifest
// names it, and the flai that serves it.
type ProjectInfo struct {
	Name        string            `json:"name"`
	Key         string            `json:"key"`
	Description string            `json:"description,omitempty"`
	Owner       string            `json:"owner,omitempty"`
	Layout      map[string]string `json:"layout"`
	Flai        string            `json:"flai"`
}

// Methods is the table flai serve offers. It is the whole of what a
// dashboard can ask this flai to do (ADR-0029); S-0072 has one method.
func Methods(version string) map[string]Method {
	return map[string]Method{
		"project.info": func(_ context.Context, p Project, _ json.RawMessage) (any, *Error) {
			m, err := manifest.Load(filepath.Join(p.Root, manifest.File))
			if err != nil {
				return nil, &Error{Code: CodeInternal, Message: err.Error()}
			}
			return ProjectInfo{Name: m.Name, Key: m.Key, Description: m.Description, Owner: m.Owner, Layout: m.Layout, Flai: version}, nil
		},
	}
}
