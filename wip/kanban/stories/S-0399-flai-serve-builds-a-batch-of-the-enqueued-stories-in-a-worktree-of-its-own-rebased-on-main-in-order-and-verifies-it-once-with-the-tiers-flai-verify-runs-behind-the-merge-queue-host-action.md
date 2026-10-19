---
id: S-0399
type: story
nature: feature
title: flai serve builds a batch of the enqueued stories in a worktree of its own, rebased on main in order, and verifies it once with the tiers flai verify runs, behind the merge-queue host action
status: backlog
parent: E-0022
owner: alex
created: 2026-10-16T04:57:53Z
updated: 2026-10-19T04:52:50Z
transitions: []
tags: [cli, dashboard]
topics: [orchestration]
touches: [flai/internal/mergequeue/state.go, flai/internal/mergequeue/state_test.go, flai/internal/mergequeue/batch.go, flai/internal/mergequeue/batch_test.go, flai/internal/verify/story.go, flai/internal/verify/story_test.go, flai/internal/serve/mergequeue.go, flai/internal/serve/mergequeue_test.go, flai/internal/serve/serve.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/config/config.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/internal/manifest/manifest.go, flai/internal/manifest/settings.go, flai/internal/manifest/manifest_test.go, flai/cmd/queue.go, flai/cmd/queue_test.go, flai/cmd/accept.go, flai/cmd/accept_queue_test.go, flaiover/src/lib/components/Review.svelte, flaiover/src/lib/components/Review.svelte.test.ts, flaiover/src/lib/components/MergeQueue.svelte, flaiover/src/lib/components/MergeQueue.svelte.test.ts, flaiover/src/routes/board/+page.svelte, "flaiover/src/routes/api/items/[id]/queue/+server.ts", flaiover/src/routes/api/queue/+server.ts, template/root/system-flow.yaml.tmpl, template/template.yaml, template/CHANGELOG.md, docs/operators/index.md, docs/operators/settings.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/users/flai.md, docs/users/flai-reference.md, docs/users/flaiover.md]
after: [S-0398]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
forecast:
  duration: 1h30m
  delivery: 2026-10-18T21:58:00Z
  basis: "Its own forecast of 1h30m; 3rd in the pull order with an in-progress limit of 5, behind S-0241 and S-0398."
  by: flai
  at: 2026-10-18T19:29:13Z
finalized:
  by: alex
  at: 2026-10-19T04:52:50Z
---
# S-0399 flai serve builds a batch of the enqueued stories in a worktree of its own, rebased on main in order, and verifies it once with the tiers flai verify runs, behind the merge-queue host action

## Goal

The stories enqueued for acceptance are verified together, once, by `flai serve` on the host, not one by one in an agent's tool call. A batch is main with the enqueued stories' branches rebased onto it in order, in a worktree of its own, verified with the tiers `flai verify` runs, so that a batch of several stories costs one run of the integration and smoke tiers. This story builds and verifies the batch; landing it is the next story's.

## Acceptance criteria

- [ ] `flai queue add S-nnnn` and `flai queue remove S-nnnn`, and the host writes behind the review page's queue button, enqueue and dequeue a story in review, recording who and when in the queue's state. `add` refuses a story that is not eligible, naming each reason, and refuses under pull-request acceptance, saying the queue runs under local acceptance only. `flai accept` by hand takes the story out of the queue, and drops a batch that holds it.
- [ ] With a new host action `merge-queue` on, `orchestration.permissions.accept_reviews` on, and `acceptance.mode` local, `flai serve` takes up to `merge_queue.batch_max` (4 unless set) eligible enqueued stories, in queue order, and builds the batch in `.flai-cache/worktrees/` on a branch of its own: main, then each story's branch rebased onto the one before. A story whose rebase stops is left out of the batch, with the stop and its paths as its outcome; the rest go on. Under pull-request acceptance it starts no batch and logs why once.
- [ ] The batch is verified once with the tiers `flai verify` selects for the paths the batch changes against main, in `flai serve`'s process. Its step and result are kept in the queue's state, and `flai queue`, `flai board`, and the dashboard show them as they change. A story's own `flai verify` result is not changed by it.
- [ ] `flai queue hold` and `flai queue resume`, and the host write behind a hold switch on the board, hold the queue: no batch starts while it is held, and one under verification finishes and waits.
- [ ] `merge_queue.batch_max` is validated by the manifest's checks and set with `flai manifest set`; `docs/operators/settings.md` indexes it and the host action, `docs/operators/index.md` says how to turn the queue on and hold it, and that it runs under local acceptance only, the template's manifest carries the key commented out, and `template/CHANGELOG.md` names it. `design/system/flai-cli.md` and `docs/users/flai.md` describe the commands.

## Tasks
- T-1640 The manifest takes merge_queue.batch_max and the host takes the merge-queue host action, off by default
- T-1641 flai queue add, remove, hold, and resume write the queue's state, and flai accept by hand takes its story out of the queue
- T-1657 The mergequeue package builds a batch worktree from main and the enqueued branches rebased in order, and verifies it once with the tiers flai verify selects
- T-1658 Host writes enqueue, dequeue, and hold the queue, and the dashboard's review page and board offer them and show the batch's step
- T-1668 flai serve runs the queue for each project with the merge-queue host action and accept_reviews on, one batch at a time, and honours a hold
- T-1676 The operators' guide, the settings index, the template, and the CLI and dashboard docs describe turning the queue on, enqueueing, holding, and batch_max

## Notes

- Drafted by the planner for E-0022, after S-0398, whose state file and view this story writes to.
- Landing a green batch and bisecting a red one are the next story's. Until it lands, a verified batch waits in the state, and stories are still accepted by hand.
- A rebase that stops on a shared path is settled by a later story of E-0022; here it leaves the story out of the batch.
- Local acceptance only, as S-0398's ADR records: under pull-request acceptance a batch's rebase would leave every story stopped on ADR-0169's `forge_stale`.

### Planning

- **Touches.** None were declared; every touch is the planner's.
  - Layout, from the explorer's map:
    - host actions: `flai/internal/hostapi/writes.go` `Actions`, `flai/internal/config/config.go`, `flai/cmd/serve_actions.go`
    - the loop: `flai/internal/serve/serve.go` starts the orchestrator; the queue's loop is a new `serve/mergequeue.go` beside it
    - verify: `flai/internal/verify/story.go` gains an entry point for a batch
    - the manifest key: `flai/internal/manifest/manifest.go` and `settings.go`
    - acceptance: `flai/cmd/accept.go` takes the story out of the queue.
  - Design: the epic's criteria 2 and 3, and TH-0465 questions 2 and 3.
  - Documentation: the operator rules in `documentation.md` need the settings index and the operators' guide. The template's manifest is `template/root/system-flow.yaml.tmpl`.
  - Co-change: `flai touches suggest` added `docs/users/flai-reference.md`, `docs/users/flaiover.md`, `design/system/flaiover-dashboard.md`, and `template/template.yaml`, the template's version.
  - No folder touch.
  - E-0021's S-0395 has since changed `flai/cmd/accept.go` and `flai/internal/hostapi/writes.go`; S-0397 changes them again. The touches stand. Whoever goes second rebases onto the other.
- **Forecast.** Raised from flai's 1h9m to 1h30m. A long-running loop in `flai serve` that makes worktrees and rebases needs tests on real repositories, and flai's size measure does not weigh that. The local-only check adds a test case, not a figure.
- **Cost of delay.** Waits on TH-0465 question 1.
