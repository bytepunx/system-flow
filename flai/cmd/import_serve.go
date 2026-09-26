package cmd

import (
	"fmt"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// A project imported on the command line is served as one imported from the
// board is (S-0120): flai import registers it with flai serve when a flai host
// runs for the config, and otherwise says what serves it.

// importServed is what flai import did to have the dashboard show the project.
type importServed struct {
	// Served is true when the project is registered with flai serve.
	Served bool `json:"served"`
	// Dashboard is the address the dashboard shows it at.
	Dashboard string `json:"dashboard,omitempty"`
	// Reason says why it is not served, or not yet reached.
	Reason string `json:"reason,omitempty"`
	// Next is the one command that serves it, run in Root.
	Next string `json:"next,omitempty"`
	Root string `json:"root"`
}

// serveImported registers the imported project with flai serve, with the
// entry flai dashboard writes, when a flai host runs for the config. With
// none it registers nothing: flai dashboard in the project registers it and
// starts the dashboard and the host.
func (a *app) serveImported(repo *workitem.Repo) importServed {
	out := importServed{Root: repo.Root, Next: "flai dashboard"}
	if repo.MainRoot != "" {
		out.Root = repo.MainRoot
	}
	if _, alive := a.hostDir().ReadStatus(time.Now()); !alive {
		out.Reason = "flai host is not running"
		return out
	}
	if repo.Manifest.Key == "" {
		out.Reason, out.Next = "its system-flow.yaml has no key", "flai check"
		return out
	}
	s, err := a.dashboardSettings(repo, "", "", 0, "")
	if err != nil {
		out.Reason = err.Error()
		return out
	}
	dir := a.serveDir()
	if _, err := a.registerProject(repo, s); err != nil {
		out.Reason = "not registered: " + err.Error()
		return out
	}
	out.Served, out.Dashboard, out.Next = true, s.url(), ""
	if _, alive := dir.ReadStatus(time.Now()); !alive {
		out.Reason, out.Next = "flai host runs but flai serve is stopped", "flai serve start"
	}
	return out
}

func (a *app) printImportServed(s importServed) {
	switch {
	case s.Served && s.Next == "":
		fmt.Fprintf(a.out, "  dashboard: served by the host flai at %s; pick it in the project switcher\n", s.Dashboard)
	case s.Served:
		fmt.Fprintf(a.out, "  dashboard: registered for %s, but %s; %s serves it\n", s.Dashboard, s.Reason, s.Next)
	default:
		fmt.Fprintf(a.out, "  dashboard: not served yet (%s); %s in %s serves it\n", s.Reason, s.Next, s.Root)
	}
}
