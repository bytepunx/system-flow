---
id: TH-0496
title: "S-0393 criterion 4: I am building the force-with-lease exception for story branches as recommended on TH-0463 question 4; confirm or redirect"
anchor:
  path: wip/kanban/stories/S-0393-under-pull-request-acceptance-a-story-s-move-to-review-has-flai-serve-push-its-branch-and-open-or-update-one-pull-request-from-the-project-s-pull-request-template.md
  item: S-0393
status: resolved
participants: [agent-S-0393, orchestrator, alex]
created: 2026-10-17T07:44:35Z
updated: 2026-10-17T16:03:13Z
---

# TH-0496 S-0393 criterion 4: I am building the force-with-lease exception for story branches as recommended on TH-0463 question 4; confirm or redirect

On wip/kanban/stories/S-0393-under-pull-request-acceptance-a-story-s-move-to-review-has-flai-serve-push-its-branch-and-open-or-update-one-pull-request-from-the-project-s-pull-request-template.md.

## Entries

### 2026-10-17T07:44:35Z agent-S-0393
TH-0463 question 4 is still unanswered. You finalized S-0393 with criterion 4 as the planner wrote it, so I am going ahead with that recommendation.

**What I am building:**

- After a sync rebased `story/S-nnnn`, flai serve pushes it with `--force-with-lease=story/S-nnnn:<head flai last pushed>`.
- The lease means the push never overwrites a commit someone else pushed.
- flai never forces `main` or any other branch.
- The ADR (T-1611) and both copies of `git.md` (T-1614) name this as the only exception to never force pushing.

**Confirm:** reply yes, and nothing changes.

**Alternative:** under pull-request acceptance, sync merges main into story branches instead of rebasing them. That needs no force push, but it changes the trial merge, the conflict log, and the close-out's sync check. I would rather leave that to a story of its own and take the force-push parts out of this one.

I will finish the rest of the story while this is open. The ADR records whatever you answer.

### 2026-10-17T08:56:58Z orchestrator (recommendation)
**Recommendation:** confirm the force-with-lease exception, and have S-0393 name it in `safety.md` as well as `git.md` before it is accepted.

## Why confirm

- **The exception is narrow.** Only `flai serve` pushes, never an agent. Only `story/S-nnnn` is forced, never main or any other branch.
- **The lease protects other people's work.** It names the head flai last pushed, so it never overwrites a commit someone else pushed.
- **The alternative costs more.** Merging main into story branches under pull-request mode changes the trial merge, the conflict log and the close-out's sync check.

## Why `safety.md` too

S-0393 adds the exception to "Never force push" in both copies of `git.md`, and ADR-0166 records it as accepted. `safety.md` § Rules still says force pushes always need confirmation, and that confirmation is per action and per context.

An automatic force-with-lease by flai serve, once a minute, on every rebased story in review, is a standing confirmation. Unless `safety.md` names the exception too, the two conventions contradict each other.

Both copies of `safety.md` need it: `design/conventions/safety.md` and `template/root/design/conventions/safety.md`. The line would read something like: "force pushes always do, except the force-with-lease push of a story's own branch that `flai serve` makes under pull-request acceptance (git.md, ADR-0166)."

## Why you

This is a change to a baseline safety rule, and it is not recorded yet. TH-0463 question 4 asks the same thing.

## Where S-0393 stands

S-0393 is in review, and this thread blocks its acceptance. Its last verify, at 2238cd4, is one commit behind the head, 6131ebd. Once you answer, I will verify it again, have the verifier check `safety.md` if you take that part, and accept it.

Source: design/conventions/safety.md § Rules

### 2026-10-17T16:03:13Z alex
Resolved: S-0393 was accepted
