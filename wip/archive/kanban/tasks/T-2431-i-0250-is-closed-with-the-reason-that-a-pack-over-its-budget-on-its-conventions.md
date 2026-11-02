---
id: T-2431
type: task
nature: remediation
title: I-0250 is closed with the reason that a pack over its budget on its conventions now keeps its briefs
status: done
parent: S-0563
owner: alex
created: 2026-11-02T04:29:06Z
updated: 2026-11-02T04:39:25Z
transitions:
  - to: ready
    at: 2026-11-02T04:39:12Z
    by: agent-S-0563
  - to: in-progress
    at: 2026-11-02T04:39:13Z
    by: agent-S-0563
  - to: done
    at: 2026-11-02T04:39:25Z
    by: agent-S-0563
stream: S-0563
tags: [issues]
touches: [design/issues/I-0250-the-orchestrator-s-conventions-sit-within-about-1-5kb-of-its-80kb-pack-budget.md, design/issues/summary.md]
after: [T-2429, T-2430]
usage:
  source: log
  seconds: 12
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 6
      output: 1418
      cache_read: 443688
      cache_write: 9261
      cost: 0.1912
  strategic:
    - kind: orchestrator
      seconds: 3
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 0
          output: 8
          cache_read: 22600
          cache_write: 730
          cost: 0.0054
---
# T-2431 I-0250 is closed with the reason that a pack over its budget on its conventions now keeps its briefs

## Work

Close I-0250 from the story's worktree with `flai issue close I-0250 --reason "<reason>"`. The reason names what fixed it: T-2429's change to `addDesign` and the ADR T-2430 records. A pack whose conventions and open-issue table pass its budget now keeps `exceeded: conventions` and still loads what is named and every brief. So no number of open issues or few lines of convention drops the orchestrator's design briefs. `flai issue close` rewrites the issue file and `design/issues/summary.md`.

It waits for T-2429 and T-2430: the reason names the change T-2429 built and the ADR T-2430 records, and criterion 2 asks for what fixed it.

## Done when

- I-0250's status is `closed` and its reason names the `addDesign` change and the ADR.
- `design/issues/summary.md` no longer lists I-0250.
- `flai check --strict` is clean.

## Notes
