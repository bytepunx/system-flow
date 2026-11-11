---
id: T-3076
type: task
nature: improvement
title: The design and the user guide say which moves of main drop a verified batch and that flai stats counts drops
status: done
parent: S-0675
owner: alex
created: 2026-11-10T10:36:56Z
updated: 2026-11-10T20:16:12Z
transitions:
  - to: ready
    at: 2026-11-10T20:14:17Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T20:14:17Z
    by: system-flow
  - to: done
    at: 2026-11-10T20:16:12Z
    by: agent-S-0675
stream: S-0675
tags: [docs, merge-queue]
touches: [design/system/workflow.md, design/system/metrics.md, design/system/flai-cli.md, docs/users/flai.md]
after: [T-3074, T-3075]
usage:
  source: log
  seconds: 115
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 38
      output: 10131
      cache_read: 3890290
      cache_write: 94161
      cost: 1.5727
---
# T-3076 The design and the user guide say which moves of main drop a verified batch and that flai stats counts drops

## Work

Describe the behaviour T-3074 and T-3075 built, linking T-3073's ADR. It waits for both because it describes what they built.

- `design/system/workflow.md` § What waits on an acceptance (S-0296): rewrite the "Main moving" bullet. Main's commits that change only `wip` paths no story of the batch changes keep the batch, and it lands on the new main. Any other move drops it, and the drop is logged.
- `design/system/metrics.md` § Merge queue (S-0404): add `dropped` to the `queue.weeks[]` row and a row of its own for `queue.weeks[].dropped`. Add the `dropped` line to the event log's list.
- `design/system/flai-cli.md` § The merge queue pass: say what the pass checks before it drops a verified batch.
- `docs/users/flai.md`: in § The merge queue, say when a batch is dropped. Where `flai stats` lists the `queue.weeks[]` fields, add `dropped`.
- `template/CHANGELOG.md`: one line under `## Unreleased`, if this repository records flai behaviour there.

## Done when

- Each document says what the code now does, links the ADR where it states the rule, and contradicts neither T-3074's tests nor T-3075's.
- The markdown lint and `flai check --strict` pass on the changed files.

## Notes
