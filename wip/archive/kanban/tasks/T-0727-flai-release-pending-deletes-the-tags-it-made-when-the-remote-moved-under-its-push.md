---
id: T-0727
type: task
nature: improvement
title: flai release --pending deletes the tags it made when the remote moved under its push
status: done
parent: S-0242
owner: arobson
created: 2026-10-03T01:24:33Z
updated: 2026-10-03T01:34:31Z
transitions:
  - to: ready
    at: 2026-10-03T01:24:42Z
    by: agent-S-0242
  - to: in-progress
    at: 2026-10-03T01:24:42Z
    by: agent-S-0242
  - to: done
    at: 2026-10-03T01:34:31Z
    by: agent-S-0242
stream: S-0242
tags: []
touches: [flai/cmd/release.go, flai/internal/release/remote.go, flai/internal/release/remote_test.go, flai/cmd/release_moved_test.go, docs/users/flai-reference.md]
usage:
  source: log
  seconds: 589
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 91
      output: 21434
      cache_read: 4780033
      cache_write: 87492
      cost: 1.9637
---
# T-0727 flai release --pending deletes the tags it made when the remote moved under its push

## Work

When the push of `flai release --pending` fails, ask the remote again, past the kept answer, whether its branch moved. When it did, delete the release tags of the batches that did not reach the remote, name the ones an earlier tags-only batch already pushed, and say how to take the remote by rebase or merge, verify, and publish again (exit 3). Push each batch with `--atomic`, so a refused branch takes its batch's tags with it. Any other push failure keeps its tags, as before, so a rerun finishes the push. Update the command's long help and regenerate the reference. Waits for nothing: first layer.

## Done when

- A test moves the remote between the check and the push, and the publish exits 3, deletes the tag it made, names it, and says to rebase or merge and publish again; after a merge, a rerun tags again and pushes
- The resume-after-failed-push test still passes
- `docs/users/flai-reference.md` matches the help

## Notes
