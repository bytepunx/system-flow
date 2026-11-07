---
id: T-2575
type: task
nature: remediation
title: Say in the operators guide and flai-cli.md that a story already in ready gets its agent once the agent action is turned on
status: done
parent: S-0589
owner: alex
created: 2026-11-07T10:06:34Z
updated: 2026-11-07T10:23:08Z
transitions:
  - to: ready
    at: 2026-11-07T10:22:16Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T10:22:19Z
    by: system-flow
  - to: done
    at: 2026-11-07T10:23:08Z
    by: agent-S-0589
stream: S-0589
tags: [docs, serve]
touches: [docs/operators/index.md, design/system/flai-cli.md]
after: [T-2573]
usage:
  source: log
  seconds: 49
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 4
      output: 602
      cache_read: 479913
      cache_write: 6815
      cost: 0.1618
  strategic:
    - kind: orchestrator
      seconds: 3
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 1
          output: 8
          cache_read: 28935
          cache_write: 760
          cost: 0.0074
---
# T-2575 Say in the operators guide and flai-cli.md that a story already in ready gets its agent once the agent action is turned on

## Work

Bring the documentation in step with what T-2573 found, so that no one reads the start rule as I-0300 did. It waits for T-2573, whose test pins the behaviour the documents describe.

- `docs/operators/index.md` § "Starting an agent when a story becomes ready", the **When one starts** item: "A story is started once each time it enters ready" reads as if a story already in ready when the `agent` action is turned on, or when `flai serve` starts, waits for **Start agent**. Say that such a story gets its agent at the next look, within a minute, as ADR-0041 says, and that one whose agent names no harness, with no command set, waits with that reason on its card and in `flai serve`'s log, not in the journal, which records only starts and failures.
- `design/system/flai-cli.md`, the `flai serve agent` row: where it says a story is started "whether it entered ready before or after `flai serve` started", add "or before the `agent` action was turned on", with S-0589.

## Done when

- Both documents say it, the markdown lint passes on them, and `updated` is bumped where the front matter has it.

## Notes
