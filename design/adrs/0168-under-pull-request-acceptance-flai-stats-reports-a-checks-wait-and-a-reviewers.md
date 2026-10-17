---
id: ADR-0168
title: "Under pull-request acceptance flai stats reports a checks wait and a reviewers wait of their own, read from the kept forge events and falling inside the story's review intervals"
status: accepted
date: 2026-10-17
supersedes: []
superseded_by: []
refines: [ADR-0081, ADR-0114, ADR-0167]
---

# ADR-0168 Under pull-request acceptance flai stats reports a checks wait and a reviewers wait of their own, read from the kept forge events and falling inside the story's review intervals

## Context

[ADR-0081](0081-flai-stats-reports-forecast-error-cost-of-delay-waiting-holds-touches-drift-and.md) has `flai stats` report the time a story's agent waits: on its threads while in progress, and in review. [ADR-0114](0114-flai-stats-lists-the-longest-waits-of-the-window-each-with-its-story-its-thread.md) lists the longest of those waits. Under pull-request acceptance, a story's time in review is spent waiting on two parties the review wait does not tell apart: the forge's CI, running the required checks, and the reviewers.

[ADR-0167](0167-flai-serve-reads-each-pull-request-s-required-checks-and-reviews-into-kept.md) has flai serve keep each story's pull request under `.flai-cache/forge/status/S-nnnn.json`, with an event for each change it sees: a new head commit, a required check that appeared or changed state, a review that is new or changed. Each event carries the time flai serve saw it. flai serve reads a pull request only while its story is in review.

The forge contract reports each review's state. It does not report how many approvals the main branch's rule requires.

## Decision

`flai stats` reports two waits of their own under pull-request acceptance, the checks wait and the reviewers wait, read from the events ADR-0167 keeps. Both fall inside the story's `review` intervals. `design/system/metrics.md` § Waiting defines them.

- **The checks wait** runs from each push of the story's head to the time every required check on that head has finished. A push is an event of kind `head`. The checks have finished at the first read, at or after the push, after which no required check is `pending`. A head whose checks had all finished when it was first read waits 0. A head with no required check waits 0.
- **The reviewers wait** runs from then to the first review event that approves or requests changes. A wait that starts with an approval or a requested change standing from an earlier review ends as it starts.
- **A push ends both.** A new head ends the open checks wait or reviewers wait of the head before it, at the push, and starts a checks wait of its own.
- **Only review time counts.** Each wait counts only its part inside the story's `review` intervals. A wait still open counts up to now, or to the end of the review interval it is in.
- **The first approval ends the reviewers wait.** The forge contract does not carry the number of approvals the main branch's rule requires, so flai takes the first approval as the one that meets it. A later story can read the rule's count through the forge contract and change this definition with an ADR of its own.
- **Fields.**
  - `items[].wait_checks_seconds` and `items[].wait_reviewers_seconds`: the sums of the item's waits of each kind. Present, 0 or more, when the item has kept forge state. Absent when it has none.
  - `checks` and `reviewers` in each of `waiting.weeks[]`, each with `total_seconds` and `mean_seconds` as `review` has them.
  - `waiting.longest[]` gains the kinds `checks` and `reviewers`. `awaited` is the required check that finished last, or the reviewer whose review ended the wait.
  - `flai stats` prints a line for each under the review line.
- **Local mode reads nothing.** A project in `local` acceptance mode loads no kept forge state, even if a file remains from an earlier mode. Its items carry neither field, its weeks report 0, and its longest waits gain no kind. Its other figures do not change.

## Consequences

- A slow CI or a slow reviewer shows where the operator and the analyzer already look for waiting: `flai stats`, its JSON, and the dashboard's agent-waiting chart and table.
- The review wait is unchanged. The checks and reviewers waits are parts of it, not taken out of it, so `review` stays comparable across projects in either mode.
- The events carry the time flai serve saw each change, once a minute, not the forge's own time. Each wait is accurate to about the forge pass's interval.
- The kept state is under `.flai-cache` and not in git. A clone without it, or a host that never ran the forge pass, reports no checks or reviewers waits.
- A branch rule that requires two approvals is reported as waiting on reviewers only until the first.

## Alternatives considered

- **Use the forge's own times, the checks' completion and the reviews' submission.** Rejected for the checks: the kept state holds only each check's latest `since`, so earlier heads' times are lost. The events are the record ADR-0167 keeps for this.
- **Take the waits out of the review wait.** Rejected. It would change the meaning of `wait_review_seconds` and every week's `review` figures for projects under pull-request acceptance.
- **Start the reviewers wait only when every check passed.** Rejected. A failed check sends the story back under ADR-0167, which ends its review interval, so the wait it would count is already near 0. The plain definition is easier to state and to check.
- **Read the required approvals count now.** Deferred. The forge contract has no call for it, and adding one widens this story into the forge adapter.
