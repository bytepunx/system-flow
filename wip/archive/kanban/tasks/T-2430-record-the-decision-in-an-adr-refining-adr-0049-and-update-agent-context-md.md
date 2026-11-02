---
id: T-2430
type: task
nature: feature
title: Record the decision in an ADR refining ADR-0049 and update agent-context.md, flai-cli.md, and the prime docs
status: done
parent: S-0563
owner: alex
created: 2026-11-02T04:27:41Z
updated: 2026-11-02T04:39:03Z
transitions:
  - to: ready
    at: 2026-11-02T04:27:59Z
    by: agent-S-0563
  - to: in-progress
    at: 2026-11-02T04:35:02Z
    by: agent-S-0563
  - to: done
    at: 2026-11-02T04:39:03Z
    by: agent-S-0563
stream: S-0563
tags: []
touches: [design/system/agent-context.md, design/system/flai-cli.md, design/adrs/README.md, docs/users/flai.md, docs/users/flai-reference.md, flai/cmd/prime.go, flai/internal/context/context.go]
after: [T-2429]
usage:
  source: log
  seconds: 241
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 21
      output: 5403
      cache_read: 1691162
      cache_write: 35298
      cost: 0.7288
  strategic:
    - kind: orchestrator
      seconds: 3
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 0
          output: 8
          cache_read: 22601
          cache_write: 730
          cost: 0.0054
---
# T-2430 Record the decision in an ADR refining ADR-0049 and update agent-context.md, flai-cli.md, and the prime docs

## Work

Record TH-0805's decision with `flai adr new --commit`, refining ADR-0049, and update `design/system/agent-context.md` and `design/system/flai-cli.md` where they describe the conventions case. It waits for T-2429, whose header text and behaviour the docs describe.

## Done when

The ADR is committed, and the design text matches what T-2429 built.

## Notes
