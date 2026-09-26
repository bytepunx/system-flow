package cmd

import (
	"fmt"

	"github.com/bytepunx/system-flow/flai/internal/serve"
)

// The system-flow projects below the folders named for import that are not
// registered (S-0117): flai serve serves them, and flai serve status says so,
// and why any it does not serve is not.

// belowImportRoots is what flai serve does with the projects below the
// folders named for import: from its state when it runs, and otherwise worked
// out here, where each it would serve is not served because it does not run.
func (a *app) belowImportRoots(st serveStatus) (served []serve.Entry, unserved []serve.Found) {
	if st.Running {
		return st.Status.ImportProjects, st.Status.Unserved
	}
	placed := serve.Place(serve.FindBelow("", a.importRoots()), st.Projects, a.serveDir().KnownDashboards(st.Projects))
	for _, p := range placed {
		if p.Reason == "" {
			p.Reason = "flai serve is not running"
		}
		unserved = append(unserved, p)
	}
	return nil, unserved
}

func (a *app) printBelowImportRoots(st serveStatus) {
	served, unserved := a.belowImportRoots(st)
	if len(served) > 0 {
		fmt.Fprintln(a.out, "served from the folders named for import (flai serve import list):")
		for _, p := range served {
			line := "not connected"
			if c, ok := st.Status.Connections[p.Root]; ok && c.Connected {
				line = "connected since " + c.Since
			}
			fmt.Fprintf(a.out, "  %s  %s  %s\n    %s\n", p.Key, p.URL, line, p.Root)
		}
	}
	if len(unserved) > 0 {
		fmt.Fprintln(a.out, "not served:")
		for _, p := range unserved {
			fmt.Fprintf(a.out, "  %s\n    %s\n", p.Root, p.Reason)
		}
	}
}
