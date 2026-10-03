---
id: S-0242
type: story
nature: improvement
title: A publish refused because the remote moved drops its tags, takes the remote by rebase or merge, and works conflicts through tasks and threads
status: done
owner: arobson
created: 2026-10-02T23:48:34Z
updated: 2026-10-03T02:36:30Z
transitions:
  - to: ready
    at: 2026-10-02T23:55:31Z
    by: alex
  - to: in-progress
    at: 2026-10-03T01:22:36Z
    by: agent-S-0242
  - to: review
    at: 2026-10-03T02:34:57Z
    by: agent-S-0242
  - to: done
    at: 2026-10-03T02:36:30Z
    by: alex
tags: [flai, template]
touches: [template/root/design/conventions/work-management.md, design/conventions/work-management.md, flai/cmd/release.go, docs/users/flai-reference.md, flai/cmd/release_moved_test.go, flai/internal/release/remote.go, flai/internal/release/remote_test.go, design/system/pushing-from-the-board.md, docs/users/flai.md, design/system/flai-cli.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 4366
  models:
    - model: claude-opus-5-5
      input: 306
      output: 71907
      cache_read: 16035883
      cache_write: 293515
      cost: 6.5877
    - model: claude-sonnet-5-5
      input: 42
      output: 9540
      cache_read: 698069
      cache_write: 87675
      cost: 0.4543
---
# S-0242 A publish refused because the remote moved drops its tags, takes the remote by rebase or merge, and works conflicts through tasks and threads

## Goal

ADR-0067 (accepted 2026-10-02) lets a clone whose remote moved take the remote's branch by rebase or merge, with conflicts coordinated on threads; the work-management baseline still says "merge the remote branch (never rebase what was accepted)". The designer decided on TH-0070: when the remote has changed, the release tags this clone made are invalid, so they are deleted, the clone rebases onto the remote or merges it, and the release is tagged again after verification; merge conflicts are worked through tasks on the stories in question, in threads with the operator.

## Acceptance criteria
- [x] The work-management baseline in `template/root/design/conventions/` and its copy in `design/conventions/` say: when a publish is refused because the remote moved, delete the release tags it made, rebase onto the remote branch or merge it, verify, and publish again, which tags again
- [x] The baseline says that merge conflicts are worked through tasks on the stories whose changes conflict, each discussed with the operator on a thread
- [x] `flai release --pending`, refused because a push found the remote moved, says which local tags it made and how to delete them, or deletes them itself, and the design says which
- [x] `design/system/pushing-from-the-board.md` and `docs/users/flai.md` describe the recovery

## Tasks
- T-0727 flai release --pending deletes the tags it made when the remote moved under its push
- T-0728 The work-management baseline says how to recover a publish the remote moved under
- T-0729 The design and the user guide describe recovering a publish the remote moved under

## Notes

Follow-up of S-0195 (TH-0070). On TH-0073 the designer chose threads over tasks for the second criterion: accepted stories take no new tasks, and a story opened for the conflicts would overlap and hold others, so conflicts are worked on a thread on each conflicting story. Wait for S-0196, which edits the same template conventions. S-0195 makes the publish refuse before tagging when the remote branch has commits the clone lacks, so tags are left behind only when the remote moves between that check and the push.
