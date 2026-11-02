package cmd

import (
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/preview"
)

// Every acceptance lands while the main checkout has uncommitted files
// (S-0402, ADR-0189, ADR-0248): by hand, from the dashboard, --by
// orchestrator, and the merge queue's landing, flai accept --batch, alike.
// It commits the shared set, the files under wip, design, and the shared
// paths, and leaves the rest as they are, out of the lint fix over the main
// checkout's uncommitted files and out of the acceptance commit, which holds
// the acceptance and the shared set. A file of the rest the story's change
// would overwrite still stops its merge: the preview names it as the blocker
// dirty_main for every acceptance but the queue's landing, which holds the
// story back itself. --yes, by hand from the shell, includes the rest.

// leftUncommitted is what acceptance res leaves uncommitted in the main
// checkout: the rest, the files outside the shared set uncommitted before it
// began, as its preview read them; none when --yes includes them.
func leftUncommitted(res *preview.Acceptance) []string {
	if res == nil || res.IncludesUncommitted {
		return nil
	}
	return res.Uncommitted
}

// outsideFolder is the paths of paths not under folder, given with its
// trailing slash.
func outsideFolder(paths []string, folder string) []string {
	var out []string
	for _, p := range paths {
		if !strings.HasPrefix(p, folder) {
			out = append(out, p)
		}
	}
	return out
}

// addArgs is git add's arguments to stage everything the acceptance changed
// in the main checkout but the paths left.
func addArgs(left []string) []string {
	args := []string{"add", "-A"}
	if len(left) == 0 {
		return args
	}
	args = append(args, "--", ".")
	for _, l := range left {
		args = append(args, ":(exclude,literal)"+strings.TrimSuffix(l, "/"))
	}
	return args
}
