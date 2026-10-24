package guard

import (
	"slices"
	"strings"
	"testing"
)

// itsCurator is a guard in the curator's session as flai serve starts it,
// with the commands that start a strategic agent and the digest's among
// flai's (S-0431, ADR-0196).
var itsCurator = Guard{Commands: append(slices.Clone(g.Commands), "analyze", "curate", "digest", "plan"), Role: RoleCurate, Agent: "curator", Session: "20261023T000000"}

// writesOnlyEntries is what each of the curator's refusals says.
const writesOnlyEntries = "the curator writes digest entries only with flai digest new and flai digest edit"

// S-0431, ADR-0196: the curator reads as a sub-agent does, logs its
// activities, opens and replies to threads, waits on events, and lists,
// harvests, writes, and edits digest entries, however flai is run.
func TestTheCuratorReadsAndWritesDigestEntries(t *testing.T) {
	allowed := []Event{{ToolName: "Read"}, {ToolName: "Grep"}, {ToolName: "Glob"}, {ToolName: "Agent"}}
	for _, tool := range append(slices.Clone(MCPReads), "inbox", "activity_log", "thread_open", "thread_reply", "wait_for_events") {
		allowed = append(allowed, itemOf(tool, "", ""))
	}
	for _, c := range []string{
		"flai digest harvest --since 2026-10-01 --json",
		"scripts/flai.sh digest harvest --kind answer,cause",
		"bin/flai --config c.json digest list --status active --json",
		"flai digest show D-0001",
		`flai digest new "Releases go out with flai push --pending" --kind answer --source TH-0001 --topic git`,
		"scripts/flai.sh digest new Fact --kind cause --source I-0001 --body-file -",
		"bin/flai digest edit D-0001 --asked-again TH-0002,TH-0003",
		"flai digest edit D-0002 --source TH-0004 --body 'merged'",
		"flai doc search digest && flai thread show TH-0001 && flai board --json",
		"flai digest new Fact --kind answer --source TH-0001 --by curator",
		"FLAI_AGENT=curator flai digest edit D-0001 --topic git",
		"FLAI_SESSION=20261023T000000 flai digest list",
		"flai digest new --help",
		"git log --oneline -5 && git status",
		"ls design/digest",
	} {
		allowed = append(allowed, bash("", c))
	}
	for _, e := range allowed {
		if why := itsCurator.Check(e); why != "" {
			t.Errorf("%s %q refused: %s", e.ToolName, e.ToolInput.Command, why)
		}
	}
}

// S-0431, ADR-0196: the curator edits no file, unlike the analyzer, which
// edits its report: Edit, MultiEdit, Write, and NotebookEdit are refused it
// on every file, a digest entry's and an analysis report's among them.
func TestTheCuratorEditsNoFile(t *testing.T) {
	for _, tool := range []string{"Edit", "MultiEdit", "Write", "NotebookEdit"} {
		for _, file := range []string{"design/digest/D-0001-releases.md", "design/analysis/2026-10-06-all.md", "/tmp/notes.md", ""} {
			why := itsCurator.Check(fileEdit(tool, file))
			if want := "the curator cannot use " + tool + " on "; !strings.HasPrefix(why, want) || !strings.Contains(why, writesOnlyEntries) {
				t.Errorf("%s %q: %q", tool, file, why)
			}
		}
	}
}

// S-0431, ADR-0196: of flai's MCP tools the curator calls only the reads and
// its own; every work item write, the issue tools, and the tools that start
// a strategic agent or publish are refused it, each saying it writes digest
// entries only with flai digest new and edit.
func TestTheCuratorCallsNoOtherFlaiToolThatWrites(t *testing.T) {
	events := []Event{itemNew("story", true), itemOf("item_edit", "S-0001", ""), itemOf("item_move", "S-0001", "backlog")}
	for _, tool := range []string{"issue_new", "issue_bump", "issue_story", "issue_close", "curate", "analyze", "plan", "release_publish", "thread_resolve", "criteria_tick", "adr_new", "message_get", "wait_for_work", "story_start", "agent_start", "queue_add", "task_done"} {
		events = append(events, itemOf(tool, "S-0001", ""))
	}
	for _, e := range events {
		tool := strings.TrimPrefix(e.ToolName, MCPPrefix)
		if why := itsCurator.Check(e); !strings.HasPrefix(why, "the curator cannot call "+tool+": "+writesOnlyEntries) {
			t.Errorf("%s: %q", tool, why)
		}
	}
}

// S-0431, ADR-0196: of flai's commands the curator runs only the reads and
// digest harvest, new, and edit; digest seed, every other command that
// writes, the commands that start a strategic agent among them, and git's
// writes are refused it, each saying it writes digest entries only with flai
// digest new and edit.
func TestTheCuratorRunsNoOtherCommandThatWrites(t *testing.T) {
	for _, c := range []string{
		"flai digest seed",
		"scripts/flai.sh digest seed --since 2026-10-01 --dry-run",
		"flai digest harvest --json && bin/flai digest seed",
		"flai story new --epic E-0001 --draft 'T'",
		"flai task new --story S-0001 'T'",
		"flai edit S-0001 --touches flai/cmd",
		"flai move S-0001 backlog",
		"flai issue new 'x' --class defect",
		"flai issue bump I-0001",
		"flai thread new --on S-0001 'q' 'text'",
		"flai curate --since 2026-10-01",
		"flai analyze --focus risk",
		"flai plan E-0001",
		"flai accept S-0001",
		"flai release --pending",
		"flai adr new 'x'",
		"git add design/digest",
		"git commit -m digest",
		"git push",
		"bash -c 'flai digest list && git commit -am x'",
	} {
		why := itsCurator.Check(bash("", c))
		if !strings.HasPrefix(why, "the curator cannot run ") || !strings.Contains(why, writesOnlyEntries) {
			t.Errorf("%q: %q", c, why)
		}
	}
	if why := itsCurator.Check(bash("", "flai digest seed")); !strings.Contains(why, seedsUnjudged) {
		t.Errorf("digest seed says why: %q", why)
	}
}

// S-0431, ADR-0143, ADR-0196: the curator writes only as itself, curator,
// in its own session: a call that names another author, or sets or unsets
// FLAI_AGENT or FLAI_SESSION to another's, is refused, as the analyzer's is,
// and its refusal names the curator's own name.
func TestTheCuratorWritesOnlyAsItself(t *testing.T) {
	for _, c := range []string{
		"flai digest new Fact --kind answer --source TH-0001 --by alex",
		"flai --by alex digest edit D-0001 --topic git",
		"scripts/flai.sh digest edit D-0001 --asked-again TH-0002 --by=analyzer",
		"FLAI_AGENT=alex flai digest new Fact --kind answer --source TH-0001",
		"export FLAI_AGENT=orchestrator; flai digest edit D-0001 --topic git",
		"env -u FLAI_AGENT flai digest new Fact --kind answer --source TH-0001",
		"FLAI_SESSION=other flai digest edit D-0001 --topic git",
		"unset FLAI_SESSION",
		"bash -c 'flai digest new Fact --kind answer --source TH-0001 --by alex'",
		"flai order --by wsjf --placed-by alex",
	} {
		why := itsCurator.Check(bash("", c))
		if !strings.HasPrefix(why, "the curator cannot run ") {
			t.Errorf("%q: %q", c, why)
		}
	}
	why := itsCurator.Check(bash("", "flai digest new Fact --kind answer --source TH-0001 --by alex"))
	if !strings.Contains(why, "ADR-0143") || !strings.Contains(why, "a served agent writes only as itself, curator") {
		t.Errorf("--by another: %q", why)
	}
	// with no FLAI_AGENT, the name flai serve gives the run is its own
	unnamed := itsCurator
	unnamed.Agent = ""
	if why := unnamed.Check(bash("", "flai digest new Fact --kind answer --source TH-0001 --by curator")); why != "" {
		t.Errorf("--by curator refused: %s", why)
	}
	if why := unnamed.Check(bash("", "flai digest new Fact --kind answer --source TH-0001 --by analyzer")); !strings.Contains(why, "a served agent writes only as itself, curator") {
		t.Errorf("--by analyzer: %q", why)
	}
}

// S-0431, ADR-0196: the curator's sub-agents are held as every sub-agent
// is: they read, and write no digest entry.
func TestTheCuratorsSubAgentIsHeldAsASubAgent(t *testing.T) {
	for _, tool := range []string{"inbox", "thread_open", "activity_log", "curate"} {
		e := itemOf(tool, "", "")
		e.AgentID, e.AgentType = "a1", "explorer"
		if why := itsCurator.Check(e); !strings.Contains(why, "a sub-agent (explorer) cannot call "+tool) {
			t.Errorf("%s: %q", tool, why)
		}
	}
	if why := itsCurator.Check(bash("explorer", "flai digest list && git log -5")); why != "" {
		t.Errorf("a sub-agent's reads refused: %s", why)
	}
	for _, c := range []string{"flai digest new Fact --kind answer --source TH-0001", "flai digest harvest"} {
		if why := itsCurator.Check(bash("explorer", c)); !strings.Contains(why, "a sub-agent (explorer) cannot run") {
			t.Errorf("a sub-agent's %q: %q", c, why)
		}
	}
}
