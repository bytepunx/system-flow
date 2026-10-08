---
id: T-1351
type: task
nature: improvement
title: The user guide, the dashboard guide, and the CLI design say when flai verify resumes and how to force a full run
status: done
parent: S-0341
owner: alex
created: 2026-10-08T08:05:45Z
updated: 2026-10-08T08:59:29Z
transitions:
  - to: ready
    at: 2026-10-08T08:52:10Z
    by: agent-S-0341
  - to: in-progress
    at: 2026-10-08T08:52:11Z
    by: agent-S-0341
  - to: done
    at: 2026-10-08T08:59:29Z
    by: agent-S-0341
stream: S-0341
tags: [cli, docs]
touches: [docs/users/flai.md, design/system/flai-cli.md, docs/users/flaiover.md]
after: [T-1344, T-1346]
usage:
  source: log
  seconds: 438
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 36
      output: 11206
      cache_read: 2508273
      cache_write: 63063
      cost: 1.147
---
# T-1351 The user guide, the dashboard guide, and the CLI design say when flai verify resumes and how to force a full run

## Work

Criterion 5, and the dashboard guide for criterion 2. It waits for T-1344 and T-1346, so it describes the text, the flag, and the page as built. It shares no path with T-1347, so the two run together.

- `docs/users/flai.md` § Verify a story before review: when a run resumes (same head, same main commit, same tiers, manifest `tests`, and commands), that the four cheap steps always run, the `reused` state and its `from` time in the text and `reused_from` in `--json`, `--fresh` in the example block, and that a commit, a sync, or a manifest edit makes the next run full. Say that the close-out resumes too and its last line names the reused tiers.
- `design/system/flai-cli.md`: the `flai verify` row of the commands table gains `[--fresh]`, the resume condition, the fingerprint in the record, `reused_from`, and that `flai test` never resumes.
- `docs/users/flaiover.md` § Reviewing a story: the step states gain `reused`, shown with the time of the run it comes from.
- Bump each file's `updated`.

## Done when

- The three documents say what the built code does, read against T-1343's and T-1344's code.
- `flai test docs/users/flai.md design/system/flai-cli.md docs/users/flaiover.md` passes.

## Notes

None.
