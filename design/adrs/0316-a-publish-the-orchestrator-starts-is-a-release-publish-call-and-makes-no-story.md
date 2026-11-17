---
id: ADR-0316
title: "A publish the orchestrator starts is a release_publish call and makes no story, and the baseline publish formula is retired"
status: proposed
date: 2026-11-17
supersedes: []
superseded_by: []
refines: [ADR-0094, ADR-0240, ADR-0273]
---

# ADR-0316 A publish the orchestrator starts is a release_publish call and makes no story, and the baseline publish formula is retired

## Context

ADR-0240 made each publish the orchestrator starts a story: it runs the baseline `publish` formula, a story's agent works its three steps (fetch, evaluate, publish), the story is closed out, reviewed, accepted, and archived, and `release_publish` refuses while the project defines the formula. ADR-0273 passes the orchestrator's judgement to that formula as an input.

The reason given was that a publish as a tool call was not worked, measured, or archived like other work. In this repository, from 2026-11-07 to 2026-11-17, that made 51 stories titled "Publish what is accepted and not yet released", costing 131 USD and 2.5 hours of agent time. Each takes an in-progress slot and then a review slot, is rebased and closed out like a story that changes code though its branch changes nothing, and waits for an acceptance of its own before the board is clear of it. On 2026-11-17 the operator asked that a publish no longer introduce a story.

The route without a story is still built: `release_publish` (ADR-0094) runs the publish of what is pending on the host, under the `publish` permission and the `push` host action, and a project with no `publish` formula uses it today.

## Decision

A publish is not a story. The orchestrator publishes with `release_publish`, and the baseline `publish` formula is retired.

1. **One route.** Under `orchestration.permissions.publish`, when `release_evaluate` is met the orchestrator calls `release_publish` with the figure that was met, or its reasoning under `judgement`, as the reason. It makes no story, task, branch, or worktree.
2. **What the formula's steps did, the tool does.** `release_publish` fetches the remote's branch and tags, evaluates the policy again, and publishes only when it is met and something is pending. When another publish already released what was pending it says which release covers each component and changes nothing. A moved remote or the `push` gate is refused with the same words as today.
3. **The record.** The release's own record, its tags, and the orchestrator's `activity_log` entry (the policy, its figures or the reasoning, the versions and tags, the items bundled) are the record of a publish. Its time and cost are the orchestrator's, measured as its other actions are.
4. **A refusal** is logged and raised on a thread to the operator on the most recently accepted story of the batch, as the role file already says for `release_publish`.
5. **The formula is retired.** The template no longer ships `design/workflows/formulas/publish.md`. `flai upgrade` removes a project's copy that it has not changed, and leaves one with project additions, reporting it. `release_publish` no longer refuses for a formula of that name, and flai refuses the orchestrator every formula. The gate in the pending publish on an agent flai serve started, the `push` host action, stays (ADR-0240's third point), as do a step's `in: main` and `stop:` (ADR-0271) for other formulas.
6. **The dashboard's Publish and the operator's own shell** are unchanged.

## Consequences

- No publish story enters the board, the WIP limits, the review queue, or the archive. The 51 archived ones stay as they are.
- A publish is no longer a work item with usage of its own; the flow metrics lose those short stories, which will move throughput and cycle-time figures at the date this lands.
- The `--judged` flag of the release's evaluation (ADR-0273) is kept only if `release_publish` uses it; the formula's `judgement` input goes.
- ADR-0312's `verify_at.publish` (S-0811) runs its suites inside the publish, whoever starts it, so it does not depend on the formula; S-0811's criterion about the publish formula's step is rewritten to speak of `release_publish`.

## Alternatives considered

- **Keep the formula and have flai serve accept its story unseen.** Rejected: the story is still made, worked by an agent, closed out, and archived; only the review wait goes.
- **Have flai serve publish as a job of its own, with no orchestrator.** Not taken now: the decision to publish under `judgement` is the orchestrator's, and the tool already exists. A later story may move the publish to a job as S-0804 did for rebase and resolve.
- **Delete the formula in this repository only.** Rejected as the whole answer: every project made from the template would keep making the stories. It is the stopgap until this lands.
