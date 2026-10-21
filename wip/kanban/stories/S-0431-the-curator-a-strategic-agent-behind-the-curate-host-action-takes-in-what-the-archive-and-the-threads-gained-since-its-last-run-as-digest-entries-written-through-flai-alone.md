---
id: S-0431
type: story
nature: feature
title: The curator, a strategic agent behind the curate host action, takes in what the archive and the threads gained since its last run as digest entries written through flai alone
status: backlog
parent: E-0025
owner: alex
created: 2026-10-16T05:59:25Z
updated: 2026-10-21T11:12:19Z
transitions: []
tags: [cli, template]
topics: [cli, planning, analysis, conventions, template]
touches: [".claude/agents/curator.md", ".claude/settings.json", design/adrs, design/adrs/README.md, design/conventions/strategic-agents.md, design/system/flai-cli.md, design/system/project-manifest.md, design/system/strategic-agents.md, docs/operators/index.md, docs/operators/settings.md, docs/users/flai-reference.md, docs/users/flai.md, flai/cmd/curate.go, flai/cmd/curate_test.go, flai/cmd/guard.go, flai/cmd/root.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/internal/context/strategic.go, flai/internal/context/strategic_test.go, flai/internal/conventions/conventions.go, flai/internal/conventions/conventions_test.go, flai/internal/guard/curate_test.go, flai/internal/guard/guard.go, flai/internal/harness/adapters.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/mcpserver/curate.go, flai/internal/mcpserver/folder.go, flai/internal/serve/activity.go, flai/internal/serve/agents.go, flai/internal/serve/curate.go, flai/internal/serve/curate_test.go, flai/internal/workitem/activity.go, flai/internal/workitem/activity_test.go, template/CHANGELOG.md, template/root/.claude/agents/curator.md, template/root/.claude/settings.json, template/root/design/conventions/strategic-agents.md, template/template.yaml]
after: [S-0430]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
forecast:
  duration: 1h23m
  delivery: 2026-10-22T02:34:00Z
  basis: "Its own forecast of 1h23m; 25th in the pull order with an in-progress limit of 5, behind S-0460, S-0461, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0414, S-0415, S-0416, S-0417, S-0418, S-0419, S-0420, S-0421, S-0422, S-0423, S-0424, S-0425, S-0428, S-0429 and S-0430."
  by: flai
  at: 2026-10-21T11:12:19Z
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

- **Touches.** The story declared none. Every touch comes from its eight tasks, file by file. Each of the curator's files mirrors one of the analyzer's that the explorer mapped.
  - **From layout:** `flai/internal/serve/curate.go` (as `analyze.go`), `agents.go`, and `activity.go`; `flai/internal/harness/harness.go` and `adapters.go`; `flai/internal/guard/guard.go` with a new `curate_test.go`, and `flai/cmd/guard.go`; `flai/cmd/curate.go` (as `analyze.go`) and `root.go`; `flai/internal/mcpserver/curate.go` and `folder.go`; `flai/internal/hostapi/writes.go`; `flai/cmd/serve_actions.go`; `flai/internal/conventions/conventions.go`; `flai/internal/context/strategic.go`; `flai/internal/workitem/activity.go`; `flai/internal/manifest/manifest.go` and `settings.go`; the tests of each; `.claude/agents/curator.md` and `.claude/settings.json` with their template copies.
  - **From design:** the ADR and its index; for criterion 7, the convention in both copies, `strategic-agents.md`, `project-manifest.md`, `flai-cli.md`, `docs/operators/index.md`, `docs/operators/settings.md`, `docs/users/flai.md`, `template/CHANGELOG.md`, and `template/template.yaml`.
  - **Co-change, added:** `flai/internal/harness/harness_test.go` (92% of the commits that changed the analyzer's files), `template/CHANGELOG.md` (33%), `template/template.yaml` (18%), `design/system/flai-cli.md` (20%), `flai/internal/serve/agents.go` (18%), `adapters.go` (15%), `flai/cmd/serve_actions.go` (13%), and the convention `strategic-agents.md` (12%).
  - **Co-change, not added:** `work-management.md` and `delegation.md` in both copies (18% and 15%), since nothing changes in how a story's agent works or delegates; `flai/internal/serve/agents_test.go` (12%), since `curate_test.go` covers the curator's record.
  - **Folder touch kept:** `design/adrs`, because the ADR's file name is known only once `flai adr new` makes it. It is in `claims.shared`, so it holds no story.
- **Forecast.** Kept at flai's 1h23m. Most of the 45 touches mirror the analyzer's code, which is faster than new code; the waits on `permission_prompt` for the `.claude/` writes are not agent time. Delivery is flai's play-out moved by S-0429's raise.
- **Cost of delay.** Not set. Neither the story nor E-0025 has inputs, so `flai cod` refuses. The epic's input is asked on TH-0476, question 1, and the value follows its answer.
- **Overlap.** `guard.go`, `harness.go`, `adapters.go`, `serve/agents.go`, `conventions.go`, and `strategic.go` are also claimed by E-0024's S-0414 to S-0419, and `hostapi/writes.go` and `serve_actions.go` by E-0021's and E-0022's drafts, so the holds run them one after the other.
