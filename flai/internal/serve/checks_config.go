package serve

import (
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// Substitute replaces {story} with story and {root} with root in each
// argument; nothing else is interpreted. Shared by the agent command
// (S-0079) and checks commands (S-0082), the two places an operator's
// argument list is run as it stands, never through a shell.
func Substitute(args []string, story, root string) []string {
	out := make([]string, len(args))
	r := strings.NewReplacer("{story}", story, "{root}", root)
	for i, a := range args {
		out[i] = r.Replace(a)
	}
	return out
}

// ChecksConfig is the operator's say about running checks for a story in
// review (S-0082): the named commands, run in order in the story's
// worktree, and how long one run of all of them together may take.
type ChecksConfig struct {
	Enabled  bool
	Commands []manifest.NamedCommand
	Timeout  time.Duration
}

// DefaultChecksTimeout is used when the operator has not set one.
const DefaultChecksTimeout = 15 * time.Minute

// ResolveChecks is the operator's real choice of where to name checks: the
// host's own configuration, when it names any, else the manifest's. Never
// both, never a merge of the two.
func ResolveChecks(hostCommands, manifestCommands []manifest.NamedCommand) []manifest.NamedCommand {
	if len(hostCommands) > 0 {
		return hostCommands
	}
	return manifestCommands
}
