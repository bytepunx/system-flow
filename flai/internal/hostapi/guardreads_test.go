package hostapi

import (
	"slices"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/guard"
)

// hostAPIWrites are the methods of flai hostapi's table that are not on
// guard.HostAPIReads: each writes something, a work item, a thread, a
// document, a file, a commit, a setting, a mark of what was seen, or kept
// state, starts or stops a process, reaches the network to change it, or is
// left to the story's agent for a reason given beside it (S-0612, I-0307).
// ImportMethods' import.preview and import.run are not here: flai serve adds
// them, and flai hostapi's table, MethodsFor's, has neither.
var hostAPIWrites = []string{
	"accept.run",               // accepts a story: merges, archives, commits
	"adr.accept",               // write ADRs
	"adr.new",                  // write ADRs
	"adr.template",             // prints a body, but flai adr new is not a guard read
	"agent.commit",             // commits a story's worktree
	"agent.restart",            // starts or stops an agent
	"agent.start",              // starts or stops an agent
	"agent.stop",               // starts or stops an agent
	"analyze.run",              // start or stop a strategic agent
	"board.archive",            // archives items
	"board.limit",              // sets a WIP limit
	"checks.cancel",            // start or stop a check run
	"checks.run",               // start or stop a check run
	"checks.status",            // read a check run, but flai checks is not a guard read
	"checks.tail",              // read a check run, but flai checks is not a guard read
	"components.start",         // start or stop a dev server
	"components.status",        // marks each dev server used
	"components.stop",          // start or stop a dev server
	"curate.run",               // start or stop a strategic agent
	"dashboard.check",          // pulls a docker image
	"dashboard.restart",        // change the dashboard's container
	"dashboard.status",         // host-level, not a guard read
	"dashboard.stop",           // change the dashboard's container
	"dashboard.upgrade",        // change the dashboard's container
	"dashboard.versions",       // host-level, not a guard read
	"doc.save",                 // write documents and files
	"file.get",                 // reads, but flai file show is not a guard read
	"file.save",                // write documents and files
	"host.check",               // host-level, not a guard read
	"host.process",             // change the host's processes or its flai
	"host.status",              // host-level, not a guard read
	"host.upgrade",             // change the host's processes or its flai
	"host.versions",            // host-level, not a guard read
	"images.ensure",            // asks flai serve for an image job
	"images.get",               // runs docker and the project's flai.command, not checked as reads
	"images.log",               // reads flai serve's state, which no agent's flai hostapi has
	"issue.bump",               // write issues and stories
	"issue.close",              // write issues and stories
	"issue.new",                // write issues and stories
	"issue.story",              // write issues and stories
	"item.block",               // write work items
	"item.criteria",            // write work items
	"item.edit",                // write work items
	"item.finalize",            // write work items
	"item.move",                // write work items
	"item.new",                 // write work items and the board
	"item.order",               // write work items and the board
	"item.template",            // prints a body, but flai story new is not a guard read
	"item.unblock",             // write work items and the board
	"kluster.checks",           // asks flai serve for a checks run against a story's namespace
	"kluster.deploy",           // change the cluster
	"kluster.deployments",      // reach the cluster with kubectl, not checked as reads
	"kluster.down",             // change the cluster
	"kluster.forward.start",    // change the cluster or start a port-forward
	"kluster.forward.stop",     // change the cluster or start a port-forward
	"kluster.log",              // reach the cluster with kubectl, not checked as reads
	"kluster.promote",          // change the cluster
	"kluster.teardown",         // change the cluster
	"kluster.up",               // change the cluster or start a port-forward
	"kluster.workloads",        // reach the cluster with kubectl, not checked as reads
	"log.read",                 // reads flai serve's logs for the dashboard, kept from agents as images.log is
	"orchestrate.start",        // start or stop a strategic agent
	"orchestrate.stop",         // start or stop a strategic agent
	"plan.run",                 // start an agent
	"publish.run",              // tags and pushes
	"queue.add",                // write the merge queue's state
	"queue.hold",               // write the merge queue's state
	"queue.remove",             // write the merge queue's state
	"role.run",                 // start an agent
	"settings.action",          // write host settings, or check them as a write's step
	"settings.agent",           // write host settings, or check them as a write's step
	"settings.check",           // write host settings, or check them as a write's step
	"settings.checks_timeout",  // sets the checks' time limit
	"settings.dashboard_hosts", // sets the hosts the dashboard accepts
	"settings.dashboard_token", // rotates the dashboard token
	"settings.default_agent",   // sets the default agent
	"settings.harness",         // write host settings, or check them as a write's step
	"settings.import",          // write host settings or the manifest
	"settings.manifest",        // write host settings or the manifest
	"settings.mcp_token",       // rotates the project's MCP token
	"settings.provider",        // write host settings or the manifest
	"settings.serve",           // write host settings or the manifest
	"settings.shared",          // write the manifest or host settings
	"settings.shared_check",    // reads, as flai shared check does, but T-2702 left it off guard.HostAPIReads
	"settings.subproject",      // write the manifest or host settings
	"settings.unserve",         // write the manifest or host settings
	"story.start",              // move items, open worktrees, commit
	"stream.answer",            // write narratives, threads, or a branch
	"stream.log",               // write narratives, threads, or a branch
	"stream.restore",           // write narratives, threads, or a branch
	"stream.state",             // write narratives, threads, or a branch
	"task.done",                // move items, open worktrees, commit
	"test.run",                 // run tiers; verify records its verdict and may commit
	"thread.confirm",           // write threads
	"thread.new",               // write threads
	"thread.reply",             // write threads
	"thread.resolve",           // write threads
	"verify.run",               // run tiers; verify records its verdict and may commit
}

// S-0612, I-0307: flai guard lets every agent run the host API methods on
// guard.HostAPIReads and refuses the rest. The guard cannot import this
// package, which imports it, so it keeps that list apart from the method
// table, and this test keeps the two in step: every method of the table is on
// exactly one of guard.HostAPIReads and hostAPIWrites, and neither names a
// method the table lacks.
func TestTheGuardsHostAPIReadsClassEveryMethod(t *testing.T) {
	table := MethodsFor("test", nil, Host{Enabled: func(string, string) bool { return true }})
	for _, list := range []struct {
		name  string
		names []string
	}{{"guard.HostAPIReads", guard.HostAPIReads}, {"hostAPIWrites", hostAPIWrites}} {
		for _, name := range list.names {
			if _, ok := table[name]; !ok {
				t.Errorf("%s names %s, which flai hostapi's table does not have: take it off the list", list.name, name)
			}
		}
	}
	for _, name := range guard.HostAPIReads {
		if slices.Contains(hostAPIWrites, name) {
			t.Errorf("%s is on both guard.HostAPIReads and hostAPIWrites: a method that writes anything is not a read, so take it off guard.HostAPIReads", name)
		}
	}
	for name := range table {
		if !slices.Contains(guard.HostAPIReads, name) && !slices.Contains(hostAPIWrites, name) {
			t.Errorf("the host API's method %s is on neither guard.HostAPIReads nor hostAPIWrites: add it to guard.HostAPIReads in flai/internal/guard/guard.go if it writes nothing (no file, commit, mark of what was seen, setting, or process), so that every agent may run it with flai hostapi, else to hostAPIWrites here; the guard cannot import this package, so this test keeps its list in step with the table", name)
		}
	}
	if !slices.IsSorted(hostAPIWrites) || len(slices.Compact(slices.Clone(hostAPIWrites))) != len(hostAPIWrites) {
		t.Errorf("hostAPIWrites is not sorted, or names a method twice")
	}
}
