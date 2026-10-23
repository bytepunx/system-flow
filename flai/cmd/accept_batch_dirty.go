package cmd

import (
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/preview"
)

// The merge queue's landing, flai accept --batch, lands while the main
// checkout has uncommitted files outside wip (S-0402), and so does the
// orchestrator's acceptance (ADR-0189): neither refuses them, as an
// acceptance by hand does unless --yes includes them, and each leaves them
// as they are, out of the lint fix over the main checkout's uncommitted files
// and out of the acceptance commit, which holds only the acceptance. A file
// the story's change would overwrite still stops its merge: the orchestrator's
// preview names it as the blocker dirty_main, and git's own check stops the
// queue's.

// leftUncommitted is what acceptance res leaves uncommitted in the main
// checkout: under --batch or by the orchestrator, the files outside wip
// uncommitted before it began, as its preview read them; none by hand.
func leftUncommitted(res *preview.Acceptance) []string {
	if res == nil || (res.Batch == "" && !res.Orchestrator) {
		return nil
	}
	return res.Uncommitted
}

// withoutLeft is paths without each path left names, or that lies below a
// folder left names, as git status names an untracked folder.
func withoutLeft(paths, left []string) []string {
	if len(left) == 0 {
		return paths
	}
	var out []string
	for _, p := range paths {
		if !isLeft(p, left) {
			out = append(out, p)
		}
	}
	return out
}

// isLeft reports whether p is, or lies below, a path left names.
func isLeft(p string, left []string) bool {
	for _, l := range left {
		l = strings.TrimSuffix(l, "/")
		if p == l || strings.HasPrefix(p, l+"/") {
			return true
		}
	}
	return false
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
