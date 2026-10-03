---
id: T-0772
type: task
nature: feature
title: "Activity documents for the planner, orchestrator, and analyzer: read, append an entry with accrued totals, and list in the agents index"
status: done
parent: S-0206
owner: alex
created: 2026-10-03T18:37:26Z
updated: 2026-10-03T18:46:11Z
transitions:
  - to: ready
    at: 2026-10-03T18:38:46Z
    by: agent-S-0206
  - to: in-progress
    at: 2026-10-03T18:38:47Z
    by: agent-S-0206
  - to: done
    at: 2026-10-03T18:46:11Z
    by: agent-S-0206
stream: S-0206
tags: []
touches: [flai/internal/workitem/activity.go, flai/internal/workitem/activity_test.go, flai/internal/workitem/narrative.go, flai/internal/workitem/narrative_test.go, design/adrs, design/system/agent-narrative.md, template/root/wip/agents/README.md, wip/agents/README.md]
usage:
  source: log
  seconds: 444
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 84
      output: 505
      cache_read: 3826676
      cache_write: 107151
      cost: 1.5975
---
# T-0772 Activity documents for the planner, orchestrator, and analyzer: read, append an entry with accrued totals, and list in the agents index

## Work

Model the activity document in `flai/internal/workitem`: the kinds `planner`, `orchestrator`, `analyzer`, the path `wip/agents/<kind>.md`, front matter `kind`, `accrued_cost`, `accrued_seconds`, `tasks_completed`, `last_run` parsed strictly, and `## Log` entries each with its timestamp, one-line summary, items touched, wall-clock seconds, and estimated cost. Appending an entry creates the document when it is missing, adds the entry last, and accrues the totals; the result lints clean. `WriteIndex` lists the documents that exist under a heading of their own. Record the decision in an ADR and describe the documents in `design/system/agent-narrative.md` and the template's `wip/agents/README.md`.

Waits for nothing: every other task builds on this model.

## Done when

- Tests cover an entry appended to a missing and to an existing document, the totals accrued, the log parsed back, a strict parse refusing an unknown key, and the index listing the documents under their heading
- The written document passes flai's markdown lint
- The ADR, the design, and the template README say what the documents are

## Notes
