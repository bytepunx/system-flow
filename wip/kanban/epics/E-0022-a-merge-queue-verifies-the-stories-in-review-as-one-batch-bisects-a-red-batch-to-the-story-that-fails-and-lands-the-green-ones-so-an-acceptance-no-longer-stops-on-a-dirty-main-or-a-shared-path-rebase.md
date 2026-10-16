---
id: E-0022
type: epic
nature: improvement
title: A merge queue verifies the stories in review as one batch, bisects a red batch to the story that fails, and lands the green ones, so an acceptance no longer stops on a dirty main or a shared-path rebase
status: backlog
owner: alex
created: 2026-10-16T04:27:57Z
updated: 2026-10-16T04:27:57Z
transitions: []
tags: [cli, dashboard]
topics: [cli, orchestration, analysis]
touches: [flai/internal/serve, flai/internal/workitem, flai/internal/hostapi, flai/cmd, flaiover/src, design/system/workflow.md, design/system/agent-coordination.md, design/system/metrics.md, design/conventions/strategic-agents.md, docs/operators/index.md, docs/users/flai.md, template]
---
# E-0022 A merge queue verifies the stories in review as one batch, bisects a red batch to the story that fails, and lands the green ones, so an acceptance no longer stops on a dirty main or a shared-path rebase

## Outcome

Acceptance is one story at a time, and it stops: on uncommitted files in the main checkout that no story claims, on a rebase that conflicts inside a shared path, and on gates the orchestrator raised itself. In the week from 2026-10-12 three stories carried 259 hours of review wait for these reasons while the done count moved by one story ([I-0150](../../../design/issues/I-0150-a-story-in-review-waits-days-for-the-operator-when-the-orchestrator-s-acceptance-stops-on-a-dirty-main-a-rebase-on-a-shared-path-or-a-gate-it-raised-itself.md), [2026-10-15-bottlenecks.md](../../../design/analysis/2026-10-15-bottlenecks.md)). A close-out that main outran runs the integration and smoke tiers twice, past the agent's tool limit ([I-0148](../../../design/issues/I-0148-a-close-out-whose-run-main-outruns-takes-two-rounds-of-the-integration-and-smoke-tiers-about-24-minutes-past-the-bash-tool-s-10-minute-limit-so-the-agent-s-call-is-killed-mid-round.md), [ADR-0156](../../../design/adrs/0156-a-close-out-that-main-outran-syncs-and-verifies-again-at-most-three-rounds-and.md)). Gas Town's Refinery batches pending merges, runs the gates once over the stack, bisects a red batch, and lands the green ones, so no agent pushes to main and no story waits on another's conflict. This epic brings that queue to flai, keeping the acceptance of each story as it is: merge, done, archive, commit, by the operator or the orchestrator under its permission.

- [ ] Stories in review whose last `flai verify` passed at their branch head, with every criterion ticked and no thread open, form a queue in the order they entered review; `flai board`, the MCP `board`, and the dashboard show the queue, the batch under verification and its step, and each story's outcome.
- [ ] The queue runs on the host, in `flai serve`, behind a host action and under `orchestration.permissions.accept_reviews`, as the orchestrator's acceptance does ([ADR-0093](../../../design/adrs/0093-with-accept-reviews-on-the-orchestrator-accepts-a-story-in-review-through-fl.md)): the operator can hold it, and `flai accept` by hand still works and takes a story out of the queue.
- [ ] A batch is main plus the queued branches rebased in order, built in a worktree of its own, verified once with the tiers `flai verify` runs, so that a batch of several stories costs one run of the integration and smoke tiers, not one per story, and no agent's tool call waits on it.
- [ ] A green batch lands each story in order through the acceptance of today, so that each story's transitions, acceptance commit, archive, and evidence read as they do now and metrics are unchanged per story; a red batch is bisected to the story that fails, the rest land, and the failing story goes back to in-progress with the failing step in its notes and a thread to its agent, as a send-back does.
- [ ] A conflict inside the manifest's shared paths, such as `design/issues/summary.md` or `design/adrs/README.md`, is settled by the queue's regenerate-and-merge steps ([ADR-0126](../../../design/adrs/0126-a-rebase-of-a-story-branch-merges-an-issue-file-both-sides-changed-by-its.md), [ADR-0157](../../../design/adrs/0157-flai-stream-sync-and-flai-accept-merge-design-adrs-readme-md-s-index-row-by-row.md)) and never sends a story back; a conflict outside them is the story's to resolve and is reported to its agent by message, as a trial-merge conflict is ([ADR-0121](../../../design/adrs/0121-flai-tells-two-stories-of-a-trial-merge-conflict-or-a-grown-overlap-with-a.md)).
- [ ] Uncommitted files in the main checkout that no queued story claims stop no landing: an ADR decides how a batch lands on main while the checkout is dirty, and the landing tells the operator what their checkout is behind, instead of waiting a day for them to commit or stash.
- [ ] With `accept_reviews` on, the orchestrator enqueues rather than accepting one by one, never bisects or lands by hand, and opens no gate thread of its own on a story the queue will verify; its activity document records each batch it enqueued with the policy figure.
- [ ] `metrics.md` defines, by an ADR, the queue wait per story, the batch size, and the bisections per week, and `flai stats` and the dashboard's charts report them, so the analyzer can see whether the queue removed the waits I-0150 counts.
- [ ] `docs/operators/index.md`, `docs/users/flai.md`, `design/system/workflow.md` § What waits on an acceptance, and the convention `strategic-agents.md` describe the queue, and the template carries its settings.

## Stories

To be drafted by the planner. Order: the queue's state and its view first, read-only, so the operator sees what would batch; then the batch worktree and the single verification; then landing green and bisecting red; then the shared-path and dirty-checkout rules, each with its ADR; then the orchestrator's enqueueing; then metrics and charts; documentation and the template with each story.

## Notes

- From the comparison of system-flow with Gas Town, Symphony, and Paperclip on 2026-10-16, and from the analyzer's bottleneck report of 2026-10-15: this epic is the remediation I-0150 and I-0148 wait on.
- The planner reads first: `design/system/workflow.md` § What waits on an acceptance and § Transitions and who makes them; `design/system/agent-coordination.md`; ADR-0093, ADR-0096, ADR-0121, ADR-0126, ADR-0155, ADR-0156, ADR-0157; `flai/internal/serve` and the acceptance in `flai/internal/workitem`.
- Open to the designer: whether the queue is a fourth strategic process of flai's own, with no model, or a step the orchestrator asks for; what the batch's upper size is; whether a story whose branch undoes main's change ([ADR-0155](../../../design/adrs/0155-a-rebased-story-branch-that-undoes-main-s-change-to-a-path-is-recorded-at-sync.md)) may be in a batch at all, since only the operator keeps a removal.
- Cost of delay inputs are the operator's; I-0150 estimates 96 hours of calendar time lost per cycle.
