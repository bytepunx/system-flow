---
id: T-0774
type: task
nature: feature
title: flai check validates the activity documents' front matter and stops calling them orphan narratives
status: done
parent: S-0206
owner: alex
created: 2026-10-03T18:37:41Z
updated: 2026-10-03T18:59:43Z
transitions:
  - to: ready
    at: 2026-10-03T18:38:46Z
    by: agent-S-0206
  - to: in-progress
    at: 2026-10-03T18:46:12Z
    by: agent-S-0206
  - to: done
    at: 2026-10-03T18:59:43Z
    by: agent-S-0206
stream: S-0206
tags: []
touches: [flai/internal/check, design/system/flai-cli.md]
after: [T-0772]
usage:
  source: log
  seconds: 811
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 80
      output: 392
      cache_read: 3538042
      cache_write: 83970
      cost: 1.4709
---
# T-0774 flai check validates the activity documents' front matter and stops calling them orphan narratives

## Work

`flai check` reads `wip/agents/planner.md`, `orchestrator.md`, and `analyzer.md` when present with T-0772's strict parser: an error when the front matter does not parse, `kind` does not match the file, a total is negative, or `last_run` is not a timestamp; a warning when the totals disagree with the log. The orphan-narrative check passes over them. Say so in `design/system/flai-cli.md`'s check section.

Waits for T-0772: it validates with its parser. Runs beside T-0773 and T-0775: no path in common.

## Done when

- Tests cover a valid document passing, each error and the warning, and no orphan finding for the three files

## Notes
