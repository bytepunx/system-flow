---
id: T-3507
type: task
nature: remediation
title: The design and the user guide say the guard forgets a foreground sub-agent once its result is in the session transcript
status: done
parent: S-0819
owner: alex
created: 2026-11-17T03:35:00Z
updated: 2026-11-17T07:29:25Z
transitions:
  - to: ready
    at: 2026-11-17T07:17:18Z
    by: system-flow
  - to: in-progress
    at: 2026-11-17T07:17:19Z
    by: system-flow
  - to: done
    at: 2026-11-17T07:29:25Z
    by: agent-S-0819
stream: S-0819
tags: [guard, docs]
touches: [design/system/agent-context.md, design/system/flai-cli.md, docs/users/flai.md]
after: [T-3504]
usage:
  source: log
  seconds: 726
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 27
      output: 6731
      cache_read: 2343318
      cache_write: 51414
      cost: 0.9542
---
# T-3507 The design and the user guide say the guard forgets a foreground sub-agent once its result is in the session transcript

## Work

Describe T-3504's rule where ADR-0092's is described; it waits for T-3504 because it states that task's decision, and runs beside T-3505, whose paths it shares none of.

- `design/system/agent-context.md` § Waiting for a sub-agent: the record bullet gains the prune, and the sentence that a session's stale entry "costs nothing" is replaced by what now happens to a foreground sub-agent that ends with no `SubagentStop`, linking the new ADR.
- `design/system/flai-cli.md`: the `flai guard` row of § Commands says the record is pruned from the transcript at a wait. The same row says the lock is "taken over after 5 s", which the code does not do (`flai/internal/filelock`); correct it to what the code does.
- `docs/users/flai.md` § Sub-agents: say that a foreground sub-agent that ended on an error no longer holds the wait.

## Done when

- The three documents state the rule as the ADR does, with `updated` bumped, and the markdown lint passes on them.

## Notes
