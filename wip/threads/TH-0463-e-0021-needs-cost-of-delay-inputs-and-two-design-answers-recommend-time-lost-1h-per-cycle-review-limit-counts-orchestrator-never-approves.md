---
id: TH-0463
title: "E-0021 needs cost of delay inputs and two design answers: recommend time lost 1h per cycle, review limit counts, orchestrator never approves"
anchor:
  path: wip/kanban/epics/E-0021-a-project-may-accept-a-story-through-a-pull-request-and-its-ci-instead-of-a-local-merge-with-the-checks-and-review-comments-read-back-to-the-story-s-agent.md
  item: E-0021
status: open
participants: [planner-E-0021]
created: 2026-10-16T04:29:21Z
updated: 2026-10-17T00:46:06Z
---

# TH-0463 E-0021 needs cost of delay inputs and two design answers: recommend time lost 1h per cycle, review limit counts, orchestrator never approves

On wip/kanban/epics/E-0021-a-project-may-accept-a-story-through-a-pull-request-and-its-ci-instead-of-a-local-merge-with-the-checks-and-review-comments-read-back-to-the-story-s-agent.md.

## Entries

### 2026-10-16T04:29:21Z planner-E-0021
I am drafting E-0021's stories now. Three inputs are yours. Each has my recommendation first.

## 1. Cost of delay inputs (blocking the stories' values)

`flai cod E-0021` refuses: the epic has no inputs, and its stories have none either.

- **Recommended:** `time_lost_per_cycle: 1h` on E-0021. At `planning.hour_rate` 150, that is 150 USD a week, shared among the stories by duration. This repository accepts locally, so the delay costs adoption, not work lost here. One hour a week is a modest stand-in for that.
- Alternatives: `30m` (75 USD a week), or `penalty_per_week` with your own figure for adoption.

## 2. Does a story in review count toward review's WIP limit while it waits on CI?

- **Recommended:** yes. It stays in `review` and counts, as today. The new checks-and-reviewers wait in `metrics.md` makes a slow CI visible, and the limit keeps pressure on it. No new column and no change to ADR-0073's hold.
- Alternative: a wait on the forge does not count toward the limit. That needs a new state or a flag on the review lane, plus an ADR.

## 3. May the orchestrator approve a pull request's review when `accept_reviews` is on?

- **Recommended:** no. It merges only a pull request that a person approved, wherever the project requires a review. flai opens the request with the operator's token, so the operator is its author on GitHub, and GitHub does not let an author approve their own request. Approval by the orchestrator would also turn the gate into one that the gated party opens itself.
- Alternative: the orchestrator approves through a second token of its own, set on the host.

Meanwhile I plan everything that needs no answer. Once these are answered, I write the stories' values.

### 2026-10-16T04:34:35Z planner-E-0021
A fourth input, found while drafting the second story.

## 4. May flai serve force-push a story's own branch to the remote after a sync rebased it?

`flai stream sync` rebases `story/S-nnnn` onto main. Once the branch is on the remote under pull-request acceptance, the next push after a sync is refused unless it is forced. `git.md` says never force push, and `safety.md` says a force push always needs confirmation.

- **Recommended:** yes, for story branches only. flai serve updates `story/S-nnnn` on the remote with `--force-with-lease`, naming the head it last pushed, so that it never overwrites a commit someone else pushed. It never forces `main` or any other branch. An ADR records this as the one exception, and `git.md` says so in both the baseline and the template.
- Alternative: under pull-request acceptance, `flai stream sync` merges main into the story branch instead of rebasing it. No force push is needed, but the history differs between the two modes, and the trial merge, the conflict log, and the close-out's sync check all change.

### 2026-10-17T00:46:06Z planner-E-0021
Still open after 20 hours. Your answers unblock three things:

- **Question 1:** the cost of delay values of S-0392 to S-0397. All six have none, so the cost of delay and WSJF orderings rank them last.
- **Question 4:** S-0393's criterion 4, and its tasks T-1611 and T-1612.
- **Questions 2 and 3:** the ADRs of S-0395 and S-0397.

S-0392, which you finalized, needs no answer to start.
