---
id: S-0206
type: story
nature: feature
title: Planner, orchestrator, and analyzer runs are logged in an activity document, with cost, seconds, and work completed in its front matter
status: done
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:13Z
updated: 2026-10-03T20:33:24Z
transitions:
  - to: ready
    at: 2026-10-03T05:34:34Z
    by: alex
  - to: in-progress
    at: 2026-10-03T18:34:00Z
    by: agent-S-0206
  - to: review
    at: 2026-10-03T19:49:24Z
    by: agent-S-0206
  - to: done
    at: 2026-10-03T20:33:24Z
    by: alex
tags: [flai, dashboard]
touches: [flai/internal/serve, flai/internal/usage, flai/internal/workitem, flai/internal/check, design/system/agent-narrative.md, design/adrs, template/root/wip/agents/README.md, wip/agents/README.md, flai/internal/metrics, flai/cmd/stats.go, design/system/metrics.md, flai/internal/mcpserver, flai/cmd/mcp.go, flai/cmd/mcp_http.go, flai/cmd/activity.go, flai/cmd/activity_test.go, flai/internal/guard, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, design/issues, flai/cmd/check_stats_test.go, flai/cmd/hostapi_reads_test.go, flai/internal/hostapi/reads.go, design/conventions/strategic-agents.md, template/CHANGELOG.md, template/root/design/conventions/strategic-agents.md, template/template.yaml]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 4577
  models:
    - model: claude-haiku-4-5-20251001
      input: 98
      output: 7027
      cache_read: 430141
      cache_write: 60012
      cost: 0.1533
    - model: claude-opus-5-5
      input: 600
      output: 189970
      cache_read: 40485215
      cache_write: 720222
      cost: 16.3111
    - model: claude-sonnet-5
      input: 120
      output: 33899
      cache_read: 3235294
      cache_write: 206186
      cost: 1.5018
---
# S-0206 Planner, orchestrator, and analyzer runs are logged in an activity document, with cost, seconds, and work completed in its front matter

## Goal

Story agents have narratives; the planner, orchestrator, and analyzer have no story and so no record. Each needs a log document of what it did, how long it took, and what it cost, and totals in front matter the metrics can read. Decided by the designer on 2026-10-02.

## Acceptance criteria
- [x] Each agent kind has one activity document per project, `wip/agents/planner.md`, `orchestrator.md`, `analyzer.md`, with front matter `kind`, `accrued_cost`, `accrued_seconds`, `tasks_completed`, `last_run`, and `## Log` entries flai writes per activity: the timestamp, a one-line summary the agent gives, the items touched, the wall-clock duration, and the estimated cost of that activity (from the run's usage so far, apportioned as task usage is, ADR-0051)
- [x] `flai serve` writes the entry when an activity ends, measured from the run's stream-json log as story runs are, and updates the front matter totals; a run that spans activities (the orchestrator) reports each through an MCP tool `activity_log` that the agent calls with its summary
- [x] The documents are listed in `wip/agents/index.md` under a heading of their own and lint clean; `flai check` validates their front matter
- [x] `flai stats` reads the front matter and log for the strategic-agent metrics
- [x] `design/system/agent-narrative.md` and the user guide describe them; tests cover an entry written, totals accrued, and a resumed run

## Tasks
- T-0772 Activity documents for the planner, orchestrator, and analyzer: read, append an entry with accrued totals, and list in the agents index
- T-0773 flai serve measures an activity from its run's stream-json log and writes its entry, at an activity's end and at a strategic run's end
- T-0774 flai check validates the activity documents' front matter and stops calling them orphan narratives
- T-0775 flai stats reports each strategic agent's accrued cost, seconds, activities, and log from its activity document
- T-0776 The MCP tool activity_log reports a strategic agent's activity with its summary and items, and the user guide describes the documents

## Notes

- Nothing starts a strategic run yet. `serve.LogRunEnd` is the hook S-0208 (planner) and S-0218 (orchestrator) call when a run they start ends, and their logs are named `<key>-<kind>-<start>.log` for `serve.ActivityLogs`. Tests on fixture logs cover an activity's end through `activity_log` and a run's end (`flai/internal/serve/activity_test.go`, `flai/cmd/activity_test.go`).
- `flai stats --json` and the dashboard's `stats.get` report `strategic`: each document's totals and its entries in the window (ADR-0079, `metrics.md`). The per-day series and the comparison with delivery are S-0205's.
