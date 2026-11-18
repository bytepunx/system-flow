---
id: T-3458
type: task
nature: feature
title: The design and the guides say the batch runs the suites verify_at.batch lists, and the baseline runs every one
status: done
parent: S-0809
owner: alex
created: 2026-11-16T23:49:38Z
updated: 2026-11-18T09:08:01Z
transitions:
  - to: ready
    at: 2026-11-18T09:03:18Z
    by: agent-S-0809
  - to: in-progress
    at: 2026-11-18T09:03:18Z
    by: agent-S-0809
  - to: done
    at: 2026-11-18T09:08:01Z
    by: agent-S-0809
stream: S-0809
tags: [flai]
touches: [design/system/flai-cli.md, design/system/workflow.md, design/system/project-manifest.md, docs/users/flai.md, docs/operators/index.md]
after: [T-3453]
usage:
  source: log
  seconds: 283
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 36
      output: 9128
      cache_read: 2509994
      cache_write: 86065
      cost: 1.2232
---
# T-3458 The design and the guides say the batch runs the suites verify_at.batch lists, and the baseline runs every one

## Work

ADR-0312 §6 and §8. Say that the batch and its bisection run the tiers of the suites `verify_at.batch` lists that the batch's paths select, every one under `all`, the baseline's value, which is what the batch ran before; that the suites are read from the main checkout's workflow file; and that `flai queue` and `state.json` name them. Write it where the batch's tiers are described today: `design/system/flai-cli.md` § The merge queue pass, `design/system/workflow.md` where it says the batch runs every tier its paths select, the `queue` paragraph of `design/system/project-manifest.md`, the merge queue sections of `docs/users/flai.md`, and the batch's verification and `state.json` paragraphs of `docs/operators/index.md`. Link ADR-0312 rather than restate it. Waits for T-3453, whose behaviour and record it describes; it runs beside T-3456, since the two share no file.

## Done when

- [ ] Each of the five documents says which tiers the batch runs under `verify_at.batch`, with a link to ADR-0312, and none still says the batch runs every tier its paths select without naming the suites.
- [ ] `updated` is bumped on each design document changed.
- [ ] `flai test` on the changed paths passes, the markdown lint among its tiers.

## Notes
