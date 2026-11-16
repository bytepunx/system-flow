---
id: S-0804
type: story
nature: improvement
title: Rebase and resolve rebases the story's branch and runs its close-out as a job, with its progress shown on the review page, and sends the story back to its agent only when the rebase stops on conflicts or a step fails
status: backlog
owner: alex
created: 2026-11-16T23:12:36Z
updated: 2026-11-16T23:12:36Z
transitions: []
tags: [flai, flaiover]
topics: [acceptance, verify]
touches: [flai/cmd/accept_resolve.go, flai/cmd/accept.go, flai/cmd/verify.go, flai/internal/closeout, flai/internal/serve/verifyjob.go, flai/internal/serve/verifyjob_container.go, flai/internal/hostapi, flai/internal/verify, flaiover/src/lib/components/Review.svelte, flaiover/src/routes/api/items, flaiover/src/lib/review.ts, design/adrs, design/system/flai-cli.md, design/system/flaiover-dashboard.md, design/system/workflow.md, docs, scripts/close-out.sh, template]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    time_lost_per_cycle: 6h
    by: alex
    at: 2026-11-16T23:12:36Z
---
# S-0804 Rebase and resolve rebases the story's branch and runs its close-out as a job, with its progress shown on the review page, and sends the story back to its agent only when the rebase stops on conflicts or a step fails

## Goal

The operator asked for this on 2026-11-16, to raise throughput. Today `flai accept <S-nnnn> --resolve` and the review page's Rebase and resolve button send a review story straight back to in-progress and start its agent ([ADR-0302](../../../design/adrs/0302-the-operator-sends-a-review-story-that-cannot-be-accepted-back-to-its-agent-to.md)), though most such stories need only a rebase onto main and their tiers run again, which no agent has to do. On 2026-11-15 five stories were sent back so in one day.

Rebase and resolve first does the work itself, with no agent:

1. it rebases the story's branch onto the main branch, as `flai stream sync` does;
2. it runs the story's close-out at the rebased head, with the output captured;
3. only when the rebase stops on conflicts, or a step of the close-out fails, does it send the story back to in-progress and open the send-back thread for its agent, naming the conflicting paths or the step that failed and where its output is.

A story whose rebase is clean and whose close-out passes stays in review, ready to accept. While the job runs, the review page shows what it is doing and its output in a collapsible stream, as it shows an agent's activity, so that the operator can see it has not hung.

## Acceptance criteria

- [ ] `flai accept <S-nnnn> --resolve`, and the host method `accept.resolve`, start a job in flai serve for a story in review and answer at once with the job; the story stays in review and no agent is started while the job runs.
- [ ] The job rebases the story's branch onto the main branch as `flai stream sync` does, then runs the story's close-out at the rebased head, with all of its output kept in the job's log.
- [ ] A rebase that stops on conflicts is aborted, and the story is sent back to in-progress with the send-back thread naming the conflicting paths, which flai serve starts its agent on.
- [ ] A close-out that stops at a step or fails sends the story back to in-progress with the send-back thread naming the step, its findings, and the job's log, which flai serve starts its agent on.
- [ ] A clean rebase and a passing close-out leave the story in review, with nothing sent back and no agent started, and the acceptance preview no longer lists what the rebase and the close-out settled.
- [ ] The job's record says which step it is on, the rebase, each step and tier of the close-out as it starts and ends, and how it ended, and flai gives it to the dashboard while the job runs.
- [ ] The review page shows a running Rebase and resolve job in a collapsible activity stream with its current step, its elapsed time, and its output as it is written, follows it to its end without a reload, and then shows whether the story stayed in review or was sent back and why.
- [ ] A second request while a job runs for the story answers with the job under way and starts nothing; a job flai serve's restart interrupts is recorded as interrupted and sends nothing back.
- [ ] An ADR records the change and what of ADR-0302 it replaces, and the design, the user guide, and `flai accept`'s help describe it.

## Tasks

## Notes

What exists to build on, read on 2026-11-16:

- The close-out job: `flai verify <S-nnnn> --job` asks flai serve through request files under `.flai-cache/verify/jobs` (`flai/internal/closeout`), and `flai/internal/serve/verifyjob.go` runs `scripts/close-out.sh` on the host, or `verifyjob_container.go` its steps for a story whose agent runs in a container, with the output in `<S-nnnn>.log` and the record in `<S-nnnn>.job.json`. It announces a `verified` event at its end, and a job a restart interrupts is swept as `interrupted`.
- The close-out syncs only when its verify stops at sync ([I-0383](../../../design/issues/I-0383-a-close-out-asked-for-a-story-whose-branch-is-behind-main-stops-at-sync-without.md)); this job syncs first.
- The job's step is derived from its log by `closeout.Running`, and names only `flai verify` while every tier runs, because flai verify prints its steps once it has ended. `verify.TierOptions.Started` and `Ended` exist and only the merge queue's batch uses them.
- The send-back: `sendBackBlocked` and `sendBackThread` in `flai/cmd`, and `sendBackUnsettled`, which aborts a stopped rebase (ADR-0248). flai serve starts the agent on a thread by `flai` whose title begins `Acceptance sent back:` and resolves it once the story is back in review (`flai/internal/serve/agents_restart.go`).
- The dashboard: `Review.svelte` posts to `/api/items/[id]/resolve`, which calls `accept.resolve`. `LogFeed.svelte` is a collapsible live log that reads through `log.read`, which may read a `.log` under `.flai-cache/verify/jobs`. The verification section shows a running close-out as one sentence, from `verify.status`, and does not follow it.
- Related: [I-0384](../../../design/issues/I-0384-after-a-sync-moves-a-story-s-head-the-review-page-flai-queue-and-the-acceptance.md), the stale report shown while a close-out runs, and [I-0381](../../../design/issues/I-0381-the-acceptance-preview-lists-verify-behind-main-while-the-merge-queue-batches.md), the preview's blocker while the queue batches, each with its own story.
