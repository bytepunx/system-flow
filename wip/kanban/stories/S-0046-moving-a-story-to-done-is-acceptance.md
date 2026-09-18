---
id: S-0046
type: story
nature: feature
title: Moving a story to done is acceptance
status: backlog
parent: E-0006
owner: alex
created: 2026-09-18T16:33:03Z
updated: 2026-09-18T16:33:03Z
transitions: []
tags: [cli, dashboard]
touches: [flai/cmd/move.go, flai/cmd/accept.go, flaiover/src/routes/board, flaiover/src/routes/items]
---

# S-0046 Moving a story to done is acceptance

## Goal
There is one way for a story to become done, and it is acceptance. Dragging a card to done on the board, pressing a move button on the item page, `flai move <story> done`, and `flai accept <story>` all run the same flow: rebase and merge the story branch, archive, commit, tag, push, release. A story can no longer end up done but unmerged and unreleased.

## Acceptance criteria
- [ ] `flai move <story> done` from review runs the acceptance flow (the same code as `flai accept`, with `--by`, `--deliver`, `--no-release`, `--no-push` passed through); from any other state it is refused with the rule text as today; tasks and epics keep their current move semantics, with epics still accepted by `flai accept`
- [ ] `flai accept` completes an item that is already `done` but not archived (the state an older flai or a hand edit can leave behind): it skips the transition and carries on from the merge; `flai check` flags `done` stories that are unarchived, or whose `story/<id>` branch still exists, as `story.unaccepted`
- [ ] The dashboard shows what acceptance will do before it happens: dropping a card on done (or pressing the move on the item page) opens a confirmation with the release plan from `flai release --dry-run` and the branch to be merged; confirming runs the flow and shows the tags; cancelling leaves the card in review
- [ ] Failures (dirty worktree, rebase conflicts, a failed push) are shown verbatim from flai and leave the story in review, never half done
- [ ] S-0041's accept button uses the same endpoint; `design/conventions/work-management.md` project addition, workflow.md transitions table, docs/users, and the flaiover guide say that done means accepted
- [ ] Tests: command tests for move-to-done with and without a story branch, the repair path, and the check rule; a dashboard test for the confirmation flow with a fake flai; I-0015 closed

## Tasks

## Notes
- Operator's request on 2026-09-18 after moving S-0045 to done from the board expecting acceptance (I-0015). S-0045 was landed by hand with accept's own steps.
- Agents still never move a story to done: the convention stands, and this story makes the operator's move sufficient on its own.
- Overlaps S-0041 (review and acceptance from the dashboard): this story owns the semantics and the confirmation; S-0041 keeps the review page with the branch diff.
