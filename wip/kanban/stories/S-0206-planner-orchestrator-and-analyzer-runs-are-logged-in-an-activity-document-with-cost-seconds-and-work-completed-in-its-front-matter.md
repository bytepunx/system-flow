---
id: S-0206
type: story
nature: feature
title: Planner, orchestrator, and analyzer runs are logged in an activity document, with cost, seconds, and work completed in its front matter
status: ready
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:13Z
updated: 2026-10-03T05:34:34Z
transitions:
  - to: ready
    at: 2026-10-03T05:34:34Z
    by: alex
tags: [flai, dashboard]
touches: [flai/internal/serve, flai/internal/usage, flai/internal/workitem, flai/internal/check, design/system/agent-narrative.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0206 Planner, orchestrator, and analyzer runs are logged in an activity document, with cost, seconds, and work completed in its front matter

## Goal

Story agents have narratives; the planner, orchestrator, and analyzer have no story and so no record. Each needs a log document of what it did, how long it took, and what it cost, and totals in front matter the metrics can read. Decided by the designer on 2026-10-02.

## Acceptance criteria
- [ ] Each agent kind has one activity document per project, `wip/agents/planner.md`, `orchestrator.md`, `analyzer.md`, with front matter `kind`, `accrued_cost`, `accrued_seconds`, `tasks_completed`, `last_run`, and `## Log` entries flai writes per activity: the timestamp, a one-line summary the agent gives, the items touched, the wall-clock duration, and the estimated cost of that activity (from the run's usage so far, apportioned as task usage is, ADR-0051)
- [ ] `flai serve` writes the entry when an activity ends, measured from the run's stream-json log as story runs are, and updates the front matter totals; a run that spans activities (the orchestrator) reports each through an MCP tool `activity_log` that the agent calls with its summary
- [ ] The documents are listed in `wip/agents/index.md` under a heading of their own and lint clean; `flai check` validates their front matter
- [ ] `flai stats` reads the front matter and log for the strategic-agent metrics
- [ ] `design/system/agent-narrative.md` and the user guide describe them; tests cover an entry written, totals accrued, and a resumed run

## Tasks

## Notes
