---
id: S-0242
type: story
nature: improvement
title: A publish refused because the remote moved drops its tags, takes the remote by rebase or merge, and works conflicts through tasks and threads
status: ready
owner: arobson
created: 2026-10-02T23:48:34Z
updated: 2026-10-02T23:55:31Z
transitions:
  - to: ready
    at: 2026-10-02T23:55:31Z
    by: alex
tags: [flai, template]
touches: [template/root/design/conventions/work-management.md, design/conventions/work-management.md, flai/cmd/release.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0242 A publish refused because the remote moved drops its tags, takes the remote by rebase or merge, and works conflicts through tasks and threads

## Goal

ADR-0067 (accepted 2026-10-02) lets a clone whose remote moved take the remote's branch by rebase or merge, with conflicts coordinated on threads; the work-management baseline still says "merge the remote branch (never rebase what was accepted)". The designer decided on TH-0070: when the remote has changed, the release tags this clone made are invalid, so they are deleted, the clone rebases onto the remote or merges it, and the release is tagged again after verification; merge conflicts are worked through tasks on the stories in question, in threads with the operator.

## Acceptance criteria
- [ ] The work-management baseline in `template/root/design/conventions/` and its copy in `design/conventions/` say: when a publish is refused because the remote moved, delete the release tags it made, rebase onto the remote branch or merge it, verify, and publish again, which tags again
- [ ] The baseline says that merge conflicts are worked through tasks on the stories whose changes conflict, each discussed with the operator on a thread
- [ ] `flai release --pending`, refused because a push found the remote moved, says which local tags it made and how to delete them, or deletes them itself, and the design says which
- [ ] `design/system/pushing-from-the-board.md` and `docs/users/flai.md` describe the recovery

## Tasks

## Notes

Follow-up of S-0195 (TH-0070). Wait for S-0196, which edits the same template conventions. S-0195 makes the publish refuse before tagging when the remote branch has commits the clone lacks, so tags are left behind only when the remote moves between that check and the push.
