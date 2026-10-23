---
id: S-0431
type: story
nature: feature
title: The curator, a strategic agent behind the curate host action, takes in what the archive and the threads gained since its last run as digest entries written through flai alone
status: backlog
parent: E-0025
owner: alex
created: 2026-10-16T05:59:25Z
updated: 2026-10-23T17:43:05Z
transitions: []
tags: [cli, template]
topics: [cli, planning, analysis, conventions, template]
touches: [".claude/agents/curator.md", ".claude/settings.json", design/adrs, design/adrs/README.md, design/conventions/README.md, design/conventions/strategic-agents.md, design/system/agent-adapters.md, design/system/agent-narrative.md, design/system/flai-cli.md, design/system/project-manifest.md, design/system/strategic-agents.md, docs/operators/index.md, docs/operators/settings.md, docs/users/flai-reference.md, docs/users/flai.md, flai/cmd/curate.go, flai/cmd/curate_test.go, flai/cmd/guard.go, flai/cmd/mcp.go, flai/cmd/mcp_http.go, flai/cmd/root.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/internal/context/strategic.go, flai/internal/context/strategic_test.go, flai/internal/conventions/conventions.go, flai/internal/conventions/conventions_test.go, flai/internal/guard/authors.go, flai/internal/guard/call.go, flai/internal/guard/curate_test.go, flai/internal/guard/guard.go, flai/internal/guard/shared.go, flai/internal/harness/adapters.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/mcpserver/activity.go, flai/internal/mcpserver/curate.go, flai/internal/mcpserver/curate_test.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/folder_test.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go, flai/internal/serve/activity.go, flai/internal/serve/agents.go, flai/internal/serve/claudecheck.go, flai/internal/serve/curate.go, flai/internal/serve/curate_test.go, flai/internal/usage/usage.go, flai/internal/workitem/activity.go, flai/internal/workitem/activity_test.go, flai/internal/workitem/usage_test.go, template/CHANGELOG.md, template/root/.claude/agents/curator.md, template/root/.claude/settings.json, template/root/design/conventions/README.md, template/root/design/conventions/strategic-agents.md, template/root/wip/agents/README.md, template/template.yaml, wip/agents/README.md]
after: [S-0430]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 820
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 51
          output: 850
          cache_read: 5629905
          cache_write: 12112
          cost: 1.2411
cost_of_delay:
  value: 147.54
  by: planner-S-0431
  at: 2026-10-23T17:06:20Z
forecast:
  duration: 1h45m
  delivery: 2026-10-24T09:37:00Z
  basis: "Its own forecast of 1h45m; 22nd in the pull order with an in-progress limit of 5, behind S-0429, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0414, S-0415, S-0416, S-0417, S-0418, S-0419, S-0420, S-0421, S-0422, S-0423, S-0424, S-0425 and S-0430."
  by: flai
  at: 2026-10-23T17:43:05Z
finalized:
  by: orchestrator
  at: 2026-10-23T17:06:47Z
---
# S-0431 The curator, a strategic agent behind the curate host action, takes in what the archive and the threads gained since its last run as digest entries written through flai alone

## Goal

An agent keeps the digest, so that the decisions in narratives, which only a model can judge, are taken in and the seed's raw entries are merged. The curator (TH-0476 question 3) is a fourth strategic agent on the analyzer's pattern: `flai curate`, the host API's `curate.run`, and the MCP tool `curate` start it behind the host action `curate`, off by default, one run per project at a time. It runs in the main checkout as `curator` with `FLAI_ROLE=curate` and its own prompt and definition, primes with `flai prime --role curate`, reads `flai digest harvest --since` its last run, and writes or merges entries with `flai digest new` and `flai digest edit`, recording in an entry's `asked_again` each thread that asked what the entry already answered. `flai guard` holds it to the reads, the digest commands, and its thread and activity tools. Each run is an entry in `wip/agents/curator.md`.

## Acceptance criteria

- [ ] An ADR, refining ADR-0099 and ADR-0075, records the curator: its host action, how it is started, one run per project, its agent (`digest.agent` merged over the project's `agent`), its environment, its pack (role `curate`, topics `digest` and `planning`), what the guard lets it do, where it logs, and how its run is settled (`worked` when it wrote or merged an entry or found nothing to take in, `failed` otherwise) and charged (to no item, its kind's project strategic total).
- [ ] `flai curate [--since <date>] [--json]`, `curate.run`, and the MCP tool `curate` start it, refusing as `flai analyze` does: the action off, a run under way, no harness and no command, and an agent `flai serve` started asking through MCP; tests cover each refusal and a start.
- [ ] `harness.Prompt` gives it a prompt of its own and the project's `.claude/agents/curator.md` (template and here) its tools: it reads the harvest since its last run (the last entry of `wip/agents/curator.md`, or everything on its first run), writes an entry for each fact not yet in the digest, merges one into the entry that already holds it, records `asked_again`, never writes a secret's value, edits no file, and ends with a one-line summary naming the entries it wrote and merged.
- [ ] With `FLAI_ROLE=curate`, `flai guard` passes the reads, `flai digest new`, `edit`, `harvest`, `list`, and `show`, and the MCP tools `inbox`, `activity_log`, `thread_open`, `thread_reply`, and `wait_for_events`; it refuses `Edit`, `Write`, and `NotebookEdit` on every file, every work item write, and every other command that writes; `.claude/settings.json` (template and here) runs the guard for `curate` as for `analyze`.
- [ ] The run is recorded in `serve/agents.json` under `curator`, its activity logged in `wip/agents/curator.md` with the entries it wrote as the entry's items, and `flai stats` reports the curator's strategic use as it does the analyzer's.
- [ ] A test starts the curator through a stub harness command on a fixture project with archived stories and threads, and sees its run recorded, its activity logged, and the entries the stub wrote checked. The first real run on this repository follows the story's release, when the operator turns `curate` on.
- [ ] The convention `strategic-agents.md` gains "As the curator" (template and here); `design/system/strategic-agents.md`, `design/system/project-manifest.md`, `docs/operators/index.md`, `docs/operators/settings.md`, `docs/users/flai.md`, and `template/CHANGELOG.md` describe the role, its action, and its settings.

## Tasks
- T-1823 An ADR records the curator: its host action, how it starts and runs, its pack, what the guard lets it do, where it logs, and how its run is settled and charged
- T-1824 The curate role is known to the conventions, the strategic pack, the activity documents, and the manifest's digest.agent
- T-1825 harness.Prompt gives the curator its prompt and the claude-code adapter starts it as the project's curator definition with FLAI_ROLE=curate
- T-1826 flai guard holds the curator to the reads, the digest commands, and its thread and activity tools, and refuses it every file edit and work item write
- T-1827 flai serve starts, records, and settles the curator's run, one per project, and logs its activity in wip/agents/curator.md
- T-1828 The curator's definition lists its tools and rules, and settings.json runs flai guard on its file edits, in the template and here
- T-1829 flai curate, the MCP tool curate, the host API's curate.run, and the curate host action start the curator
- T-1830 The strategic-agents convention gains As the curator, and the design and the user and operator docs describe the role, its action, and digest.agent

## Notes

- Drafted by the planner for E-0025. Built on the analyzer's machinery, as the epic's Notes allow; no dependency on E-0024. When E-0024 lands, its role files (S-0414) and guard rules (S-0416) take the curator in as a fourth baseline role.
- Writing `.claude/agents/curator.md` and `.claude/settings.json` in this repository needs the operator's permission through `permission_prompt`; plan for the wait.
- Triggers after acceptance and on a schedule, absorption, and compaction are the next story's. This story starts the curator only on demand.
- A run under the story's own build cannot happen before release, since `flai serve` runs the installed flai; criterion 6 tests the run with a stub instead.

### Planning

Revisited on 2026-10-23 by planner-S-0431 against the analyzer's files on today's main, mapped by the explorer; the plan of 2026-10-16 stands, with wider touches and the figures below.

- **Touches.** 64 paths, each a file but `design/adrs`, the union of the eight tasks' touches. Every one of the 45 the story declared is kept.
  - **Declared (45).** The first plan's, each a file that mirrors one of the analyzer's: `flai/internal/serve/curate.go` and `curate_test.go` as `analyze.go` and its test, `flai/cmd/curate.go` as `analyze.go`, `flai/internal/mcpserver/curate.go` as `analyze.go`, `flai/internal/guard/curate_test.go` beside the analyzer's cases in `guard_test.go`, `.claude/agents/curator.md` as `analyzer.md`, and the registrations in `agents.go`, `activity.go`, `harness.go`, `adapters.go`, `guard.go`, `cmd/guard.go`, `root.go`, `folder.go`, `hostapi/writes.go`, `serve_actions.go`, `conventions.go`, `context/strategic.go`, `workitem/activity.go`, `manifest.go`, `settings.go`, their tests, the ADR index, the convention and docs criterion 7 names, `template/CHANGELOG.md`, and `template/template.yaml`.
  - **Layout, added (19).** What the explorer found the analyzer also lives in: `flai/internal/usage/usage.go`, whose `StrategicKinds` is hard-coded and pinned equal to `ActivityKinds` by `flai/internal/workitem/usage_test.go`, so T-1824 fails its own tests without both; `flai/internal/mcpserver/activity.go`, the kind wording of `activity_log`; `flai/internal/guard/call.go`, whose `Decide` switches on the role with no default, so an unknown role passes every call, and `guard/authors.go` and `guard/shared.go`, which give the agent the name the guard holds it to (ADR-0143); `flai/internal/serve/claudecheck.go`, which gates the Claude Code check on the host actions and lists the strategic agents; `flai/internal/mcpserver/server.go`, `flai/cmd/mcp.go`, and `flai/cmd/mcp_http.go`, which wire the tool as they wire `Analyses`, with `mcpserver/folder_test.go` and `server_test.go`, which pin the tool names, and a new `mcpserver/curate_test.go` as `analyze_test.go`; `docs/operators/settings.md` on T-1829 as well as T-1830, since `make flai-reference` regenerates its flag table and `reference_test.go` pins it; `design/system/agent-narrative.md` and `design/system/agent-adapters.md`, which list the activity documents and the strategic agents; `design/conventions/README.md`, `wip/agents/README.md`, and their template copies, which name the three agents.
  - **Co-change, not added.** `flai touches suggest` lists `design/system/flaiover-dashboard.md` (12%), `docs/users/flaiover.md` (10%), `design/system/workflow.md` (8%), `design/issues/summary.md` (7%), `design/conventions/work-management.md` (6%), `flaiover/src/lib/server/agent.ts` (5%), `flai/internal/guard/guard_test.go` (3%), and `flai/internal/serve/agents_test.go` (3%): the story has no dashboard page (the settings page lists whatever host actions flai reports), the workflow and work management do not change, the summary is the close-out's, the guard's cases go in `curate_test.go`, and `serve/curate_test.go` covers the run record. `flai/internal/mcpserver/server.go` (7%) and `server_test.go` (5%) are added from layout above.
  - **Folder touch kept.** `design/adrs`: the ADR's file name is known only once `flai adr new` makes it. It is in `claims.shared`, so it holds no story.
- **Forecast.** 1h45m, delivery 2026-10-24T09:54:00Z. flai gives 2h8m (median 107 s per unit of size over 85 done feature stories on claude-opus-5-5 in the large band, size 71 from 7 criteria and 64 touches), up from 1h34m at the declared 45 touches. Brought down because the 19 added paths are one-line registrations, kind-list pins, and index rows the per-unit median overstates, and S-0223, which built the analyzer whole with nine tasks, took 47 min of agent time against a 2h forecast; kept above 1h34m because the eight tasks run through six layers of sub-agents and the two `.claude/` writes wait on `permission_prompt`. Delivery is flai's play-out with the shorter duration: 23rd in the pull order behind S-0428, S-0429, S-0430, and the E-0024 and OpenCode drafts.
- **Cost of delay.** 147.54 USD a week, flai's figure: the story's share of E-0025's 600.00 USD a week (from the epic's `time_lost_per_cycle` of 4h), 1h45m of the 7h7m forecast over the epic's eight open stories without inputs. No input of the story's own; none asked for, since the epic's covers it.
- **Findings for the tasks.** Three things the first plan's `## Work` does not name, proposed on the plan thread: `serve/agents.go`'s `orchestrates()` takes any run with no story, no item, and no focus for the orchestrator's, so a curator run would overwrite `Orchestrator` in `serve/agents.json` unless the run carries a mark of its own (T-1827); `guard/call.go`'s `Decide` needs a `curate` case or the role passes everything (T-1826); `usage.StrategicKinds` and its pin (T-1824). `flai digest new`, `edit`, `harvest`, `list`, and `show` do not exist on main yet: S-0428 (in progress) and S-0429 build them, and the story waits on S-0430, which waits on S-0429.
- **Overlap.** `guard.go`, `harness.go`, `adapters.go`, `serve/agents.go`, `conventions.go`, and `context/strategic.go` are also claimed by E-0024's S-0414 to S-0419, and `hostapi/writes.go` and `serve_actions.go` by E-0021's and E-0022's drafts, so the holds run them one after the other.
