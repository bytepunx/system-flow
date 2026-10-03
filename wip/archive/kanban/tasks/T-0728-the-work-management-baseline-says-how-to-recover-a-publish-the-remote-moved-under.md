---
id: T-0728
type: task
nature: improvement
title: The work-management baseline says how to recover a publish the remote moved under
status: done
parent: S-0242
owner: arobson
created: 2026-10-03T01:24:33Z
updated: 2026-10-03T01:26:02Z
transitions:
  - to: ready
    at: 2026-10-03T01:24:42Z
    by: agent-S-0242
  - to: in-progress
    at: 2026-10-03T01:24:42Z
    by: agent-S-0242
  - to: done
    at: 2026-10-03T01:26:02Z
    by: agent-S-0242
stream: S-0242
tags: []
touches: [template/root/design/conventions/work-management.md, design/conventions/work-management.md]
usage:
  source: log
  seconds: 80
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 30
      output: 6973
      cache_read: 1555034
      cache_write: 28463
      cost: 0.6388
---
# T-0728 The work-management baseline says how to recover a publish the remote moved under

## Work

Replace "merge the remote branch (never rebase what was accepted)" in the baseline and its copy with: delete the release tags the refused publish made (flai does), rebase onto the remote branch or merge it, verify, and publish again, which tags again; merge conflicts are worked through tasks on the stories whose changes conflict, each discussed with the operator on a thread. Waits for nothing: first layer, no path in common with the code task.

## Done when

- Both files carry the same baseline text, the copy matching the template's
- The baseline names the tag deletion, rebase or merge, verification, publishing again, and the conflict tasks and threads

## Notes
