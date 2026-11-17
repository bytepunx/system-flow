---
id: T-3538
type: task
nature: remediation
title: docs/users/flai.md says flai release --pending waits out a held index lock as acceptance does
status: done
parent: S-0833
owner: alex
created: 2026-11-17T08:30:02Z
updated: 2026-11-17T09:33:02Z
transitions:
  - to: ready
    at: 2026-11-17T09:32:48Z
    by: agent-S-0833
  - to: in-progress
    at: 2026-11-17T09:32:48Z
    by: agent-S-0833
  - to: done
    at: 2026-11-17T09:33:02Z
    by: agent-S-0833
stream: S-0833
tags: [docs, release]
touches: [docs/users/flai.md]
after: [T-3535]
usage:
  source: log
  seconds: 14
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 3
      output: 606
      cache_read: 245197
      cache_write: 8072
      cost: 0.1215
---
# T-3538 docs/users/flai.md says flai release --pending waits out a held index lock as acceptance does

## Work

`docs/users/flai.md` § Accept and release says acceptance waits up to about nine seconds for another git process's index lock and never removes the lock file (S-0307). It says nothing of the publish. Add one short paragraph where the section describes `flai release --pending` saying that its publish commit waits the same way, what it leaves when the lock is still held (nothing tagged or pushed, the release's files staged), and that running it again once the other process is done publishes. Link the acceptance paragraph rather than repeat it.

This task waits for T-3535, because the paragraph describes what T-3535's change does and the message its test shows.

## Done when

- `docs/users/flai.md` says that `flai release --pending` waits out a held index lock and what to do when it still fails.
- The markdown lint passes on `docs/users/flai.md` (`flai test docs/users/flai.md`).

## Notes
