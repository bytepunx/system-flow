---
id: ADR-0169
title: "Under pull-request acceptance, every acceptance, the operator's and the orchestrator's alike, waits for the pull request's required checks to pass, for no requested change to stand, and for the approvals the base branch requires, read from the state flai serve keeps at the branch's head, and the orchestrator never approves a review"
status: accepted
date: 2026-10-17
supersedes: []
superseded_by: []
refines: [ADR-0093, ADR-0165, ADR-0167]
topics: [release, orchestration]
---

# ADR-0169 Under pull-request acceptance, every acceptance, the operator's and the orchestrator's alike, waits for the pull request's required checks to pass, for no requested change to stand, and for the approvals the base branch requires, read from the state flai serve keeps at the branch's head, and the orchestrator never approves a review

## Context

[ADR-0167](0167-flai-serve-reads-each-pull-request-s-required-checks-and-reviews-into-kept.md) has flai serve read each pull request's required checks and reviews into kept state under `.flai-cache`, in `forge/status/S-nnnn.json`. On a failing required check or a requested change, flai serve sends the story back to in-progress.

Nothing stops an acceptance while CI or a review has not passed. `flai accept` merges the branch locally whatever the pull request says. [ADR-0093](0093-with-accept-reviews-on-the-orchestrator-accepts-a-story-in-review-through-flai.md) sets the conditions for the orchestrator's acceptance, and none of them is the forge's.

[ADR-0165](0165-a-project-chooses-local-or-pull-request-acceptance-in-its-manifest-flai-reaches.md)'s forge contract does not say how many approvals the base branch requires. It left that to S-0395.

## Decision

Under pull-request acceptance, every acceptance, the operator's and the orchestrator's alike, waits for the pull request's required checks to pass, for no requested change to stand, and for the approvals the base branch requires. The gate reads the state flai serve keeps at the branch's head. The orchestrator never approves a review.

- **The blockers.** Under `acceptance.mode: pull-request`, the preview of `flai accept S-nnnn`, and its `--json`, lists a blocker for each of these, by name. The real run refuses on them before anything is merged, as it does on any blocker. Each has a code. The preview lists them in a new `forge_blockers` list, and each is in `blockers` as well.

  | Code | When |
  |------|------|
  | `pull_request_missing` | flai serve has kept no pull request for the story: none is open, or flai serve has not read it yet. |
  | `forge_stale` | The kept state was read at a commit other than the story branch's head. flai serve has not yet read the head, so a stale pass never lets a story through. |
  | `check_failing` | A required check failed. One blocker per check, with its log's URL. |
  | `check_pending` | A required check has not finished. One blocker per check. |
  | `changes_requested` | A reviewer's latest decisive review requests changes. One blocker per reviewer, with the review's URL. |
  | `approvals_missing` | Fewer reviewers' latest decisive reviews approve than the base branch requires. |

- **A reviewer's decision.** It is their latest review that approves, requests changes, or was dismissed, as the forge counts it. flai serve's send back counts reviews the same way.
- **Required approvals.** The forge contract gains `RequiredApprovals(base)`: the most approving reviews the base branch's rulesets or classic protection require. flai serve keeps it as `required_approvals`, with the checks and reviews. A rule out of the token's reach counts as none, as it does for required checks.
- **Kept state only.** The gate reads only the kept state and needs no token. The operator's shell, the dashboard, and the orchestrator judge the same state. The dashboard accepts through the host API's `accept.run`, which runs `flai accept`.
- **Whoever accepts.** The operator's acceptance waits for these blockers too. An open thread is different: only the orchestrator's acceptance waits for one (ADR-0093). The orchestrator's acceptance, `flai accept --by orchestrator`, obeys these blockers beside ADR-0093's conditions.
- **The orchestrator never approves.** flai opens the pull request with the operator's token, so the forge counts the operator as its author. GitHub does not let an author approve their own request. An approval by the orchestrator would also let the gated party open its own gate. The planner recommended this on TH-0463 question 3. S-0395's criterion 3, which the operator finalized, says so. TH-0501 asks the operator to confirm it.
- **Under `local`.** The preview is as before.
- **The merge.** The merge itself still happens locally, until S-0397 merges through the remote.

## Consequences

- Under pull-request acceptance, a story in review cannot be accepted until CI and its reviewers pass it.
- A stale read blocks until flai serve's next look.
- The dashboard's acceptance confirmation lists the same blockers, from the preview's `--json`.
- A project whose token cannot read branch rules sees 0 required approvals and no required checks. Its gate is only as strong as what the token can read.

## Alternatives considered

- **The orchestrator approves with a second token of its own on the host.** Rejected for now. It would be a story of its own.
- **The gate asks the forge live at acceptance.** Rejected. Every shell, and the orchestrator, would need the token, which [ADR-0166](0166-under-pull-request-acceptance-flai-serve-pushes-a-story-s-branch-and-opens-or.md) keeps out of agents' environments.
- **Only the orchestrator's acceptance waits for the forge.** Rejected. CI and the team's reviewers are the gate whoever accepts.
