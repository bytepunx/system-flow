---
id: S-0046
type: story
nature: feature
title: Moving a story to done is acceptance
status: done
parent: E-0006
owner: alex
created: 2026-09-18T16:33:03Z
updated: 2026-09-18T16:57:43Z
transitions:
  - to: ready
    at: 2026-09-18T16:38:48Z
    by: alex
  - to: in-progress
    at: 2026-09-18T16:38:48Z
    by: alex
  - to: review
    at: 2026-09-18T16:50:45Z
    by: alex
  - to: done
    at: 2026-09-18T16:57:43Z
    by: alex
tags: [cli, dashboard]
touches: [flai/cmd/move.go, flai/cmd/accept.go, flaiover/src/routes/board, flaiover/src/routes/items]
---

# S-0046 Moving a story to done is acceptance

## Goal
There is one way for a story to become done, and it is acceptance. Dragging a card to done on the board, pressing a move button on the item page, `flai move <story> done`, and `flai accept <story>` all run the same flow: rebase and merge the story branch, archive, commit, tag, push, release. A story can no longer end up done but unmerged and unreleased.

## Acceptance criteria
- [x] `flai move <story> done` from review runs the acceptance flow (the same code as `flai accept`, with `--by`, `--deliver`, `--no-release`, `--no-push` passed through); from any other state it is refused with the rule text as today; tasks and epics keep their current move semantics, with epics still accepted by `flai accept`
- [x] `flai accept` completes an item that is already `done` but not archived (the state an older flai or a hand edit can leave behind): it skips the transition and carries on from the merge; `flai check` flags `done` stories that are unarchived, or whose `story/<id>` branch still exists, as `story.unaccepted`
- [x] The dashboard shows what acceptance will do before it happens: dropping a card on done (or pressing the move on the item page) opens a confirmation with the release plan from `flai release --dry-run` and the branch to be merged; confirming runs the flow and shows the tags; cancelling leaves the card in review
- [x] Failures before the first change (no committer identity, git unusable, dirty worktree, rebase conflicts) are shown verbatim from flai and leave the story in review; a push that cannot happen after the acceptance commit is reported as accepted locally and not pushed, never as a half-done story
- [x] S-0041's accept button uses the same endpoint; `design/conventions/work-management.md` project addition, workflow.md transitions table, docs/users, and the flaiover guide say that done means accepted
- [x] Tests: command tests for move-to-done with and without a story branch, the repair path, and the check rule; a dashboard test for the confirmation flow with a fake flai; I-0015 closed

## Tasks
- T-0137 flai: one acceptance flow shared by accept and move-to-done; accept completes done-but-unarchived stories; accept --dry-run --json preview
- T-0138 flai check: story.unaccepted for done stories that are unarchived or still have a story branch
- T-0139 flaiover: acceptance preview endpoint and confirmation dialog on the board drop and the item page move; failures leave the card in review
- T-0140 Tests on both sides; conventions, workflow, user and dashboard docs; close I-0015

## Notes
- Operator's request on 2026-09-18 after moving S-0045 to done from the board expecting acceptance (I-0015). S-0045 was landed by hand with accept's own steps.
- Agents still never move a story to done: the convention stands, and this story makes the operator's move sufficient on its own.
- Overlaps S-0041 (review and acceptance from the dashboard): this story owns the semantics and the confirmation; S-0041 keeps the review page with the branch diff.
- Found while building it: the dashboard container has git but no committer identity and no ssh. `flai dashboard` now passes the host's git identity into the container, acceptance preflights the identity before changing anything, and a push that cannot happen is reported rather than failing the acceptance.
