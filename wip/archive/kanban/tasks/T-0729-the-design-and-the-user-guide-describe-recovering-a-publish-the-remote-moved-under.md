---
id: T-0729
type: task
nature: improvement
title: The design and the user guide describe recovering a publish the remote moved under
status: done
parent: S-0242
owner: arobson
created: 2026-10-03T01:24:39Z
updated: 2026-10-03T01:40:32Z
transitions:
  - to: ready
    at: 2026-10-03T01:24:42Z
    by: agent-S-0242
  - to: in-progress
    at: 2026-10-03T01:34:31Z
    by: agent-S-0242
  - to: done
    at: 2026-10-03T01:40:32Z
    by: agent-S-0242
stream: S-0242
tags: []
touches: [design/system/pushing-from-the-board.md, docs/users/flai.md, design/system/flai-cli.md]
after: [T-0727]
usage:
  source: log
  seconds: 361
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 30
      output: 7101
      cache_read: 1583482
      cache_write: 28983
      cost: 0.6505
---
# T-0729 The design and the user guide describe recovering a publish the remote moved under

## Work

Describe the recovery in `design/system/pushing-from-the-board.md`, `docs/users/flai.md`, and the `flai release --pending` row of `design/system/flai-cli.md`: flai deletes the release tags a publish made when its push finds the remote moved (and says which, and why it keeps ones already pushed), then the clone takes the remote by rebase or merge, conflicts are worked through tasks on the stories whose changes conflict with threads with the operator, the result is verified, and publishing again tags again. Waits for T-0727, whose behaviour it describes.

## Done when

- The design says flai deletes the tags itself, which ones, and what it leaves
- The user guide gives the recovery steps

## Notes
