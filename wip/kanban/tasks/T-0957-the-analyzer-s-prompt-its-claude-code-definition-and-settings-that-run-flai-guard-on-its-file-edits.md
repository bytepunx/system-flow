---
id: T-0957
type: task
nature: feature
title: The analyzer's prompt, its claude-code definition, and settings that run flai guard on its file edits
status: in-progress
parent: S-0223
owner: alex
created: 2026-10-05T05:46:27Z
updated: 2026-10-06T20:08:35Z
transitions:
  - to: ready
    at: 2026-10-06T20:08:35Z
    by: agent-S-0223
  - to: in-progress
    at: 2026-10-06T20:08:35Z
    by: agent-S-0223
stream: S-0223
tags: [flai]
touches: [flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/harness/adapters.go, ".claude/agents/analyzer.md", ".claude/settings.json", template/root/.claude/agents/analyzer.md, template/root/.claude/settings.json]
usage:
  source: log
  seconds: 1149
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 85
      output: 435
      cache_read: 3196578
      cache_write: 164474
      cost: 1.4643
---
# T-0957 The analyzer's prompt, its claude-code definition, and settings that run flai guard on its file edits

## Work

Give the harness an analyzer request, as it has the planner's (`flai/internal/harness/harness.go`, `planPrompt`): a `Request` with `Role` `analyze` and a focus (`bottlenecks`, `intent`, `risk`, or none for all three) gets `analyzePrompt`, which tells the agent to:

- prime with the MCP tool `prime` and role `analyze`, then call `inbox`
- read the metrics with `flai stats --json`, the design with `doc_search` and `doc_get`, and the issues, and hand wide search of the code to the explorer
- write one report, `design/analysis/<date>-<focus>.md`, with front matter `title`, `updated`, `status`, `focus`, and the window, and one section per finding with its evidence (metric figures, file paths, design sections quoted), severity, and estimated impact (time lost per cycle, or revenue or penalty where the design states them), and add it to `design/analysis/README.md`
- edit nothing else, author no stories, and end with a one-line summary that names the report

The adapter runs it with `FLAI_ROLE=analyze` in its environment, as it sets `plan` for the planner (`adapters.go`).

Add the claude-code definition `.claude/agents/analyzer.md` and its template copy `template/root/.claude/agents/analyzer.md`, with the read tools, Edit and Write, Bash, Agent, and the MCP reads, `activity_log`, `thread_open`, `thread_reply`, and `wait_for_events`, and no item writes. Widen the `Edit|Write|NotebookEdit` hook in `.claude/settings.json` and its template copy so that it runs `flai guard` for `FLAI_ROLE` `analyze` as well as `plan`. A story's agent cannot write `.claude/` (I-0069): send the operator the whole file to paste, and test it once it is in place.

The template's changelog is the last task's, so that one task writes it.

This task waits for nothing. It runs with the manifest task, whose paths it does not share.

## Done when

- a test pins `analyzePrompt` for each focus and for none, and that it names the report path, the README, and `--role analyze`
- a test pins `FLAI_ROLE=analyze` in the analyzer's environment
- `.claude/agents/analyzer.md` and its template copy match, and the settings hook runs the guard for `analyze`
- `go test ./internal/harness/` passes

## Notes
