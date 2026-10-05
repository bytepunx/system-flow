---
id: T-0848
type: task
nature: feature
title: The design, the scripts index, and the operator guide describe the tiers as flai-test.sh runs them
status: done
parent: S-0267
owner: alex
created: 2026-10-05T00:18:49Z
updated: 2026-10-05T00:20:39Z
transitions:
  - to: ready
    at: 2026-10-05T00:19:03Z
    by: agent-S-0267
  - to: in-progress
    at: 2026-10-05T00:19:53Z
    by: agent-S-0267
  - to: done
    at: 2026-10-05T00:20:39Z
    by: agent-S-0267
stream: S-0267
tags: []
touches: [design/system/devex.md, scripts/README.md, docs/operators/index.md, scripts/close-out.sh]
after: [T-0846]
usage:
  source: log
  seconds: 46
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 23
      output: 5130
      cache_read: 913674
      cache_write: 24438
      cost: 0.4809
---
# T-0848 The design, the scripts index, and the operator guide describe the tiers as flai-test.sh runs them

## Work

Describe the tiers as they run after T-0846: `design/system/devex.md`'s test tiers row, `scripts/README.md`'s rows for `test.sh`, `flai-test.sh`, and the new `flaiover-unit.sh`, and where `docs/operators/index.md` names `scripts/flai-test.sh` as a check. It waits for T-0846, because it describes what that task makes.

## Done when

- the three documents say that `flai-test.sh`, close-out, and `make flai-test` run the Go tests once, the full run, and that `make test` remains the short run

## Notes
