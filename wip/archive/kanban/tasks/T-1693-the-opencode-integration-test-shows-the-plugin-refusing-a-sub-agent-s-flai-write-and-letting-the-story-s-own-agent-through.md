---
id: T-1693
type: task
nature: feature
title: The OpenCode integration test shows the plugin refusing a sub-agent's flai write and letting the story's own agent through
status: done
parent: S-0407
owner: alex
created: 2026-10-16T05:15:47Z
updated: 2026-11-07T21:21:17Z
transitions:
  - to: ready
    at: 2026-11-07T21:16:12Z
    by: agent-S-0407
  - to: in-progress
    at: 2026-11-07T21:16:15Z
    by: agent-S-0407
  - to: done
    at: 2026-11-07T21:21:17Z
    by: agent-S-0407
stream: S-0407
tags: [cli]
touches: [flai/tests/integration/opencode_test.go]
after: [T-1691, T-1692]
usage:
  source: log
  seconds: 302
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 45
      output: 16473
      cache_read: 6782755
      cache_write: 84395
      cost: 2.2962
---
# T-1693 The OpenCode integration test shows the plugin refusing a sub-agent's flai write and letting the story's own agent through

## Work

- Extend `flai/tests/integration/opencode_test.go`: with the plugin in the test project, the agent calls a `flai` write itself and through a sub-agent; assert the sub-agent's is refused with the guard's reason and the agent's own passes. Where `Guard` is false, assert instead that the plugin still refuses what `Decide` refuses for the agent's own calls, such as another author's `--by`.
- Keep the skip on a missing or unknown `opencode`.
- Waits for T-1691 (the plugin) and T-1692 (the capability the test asserts on).

## Done when

The test passes on the known `opencode` version and is skipped, saying so, without it.

## Notes
