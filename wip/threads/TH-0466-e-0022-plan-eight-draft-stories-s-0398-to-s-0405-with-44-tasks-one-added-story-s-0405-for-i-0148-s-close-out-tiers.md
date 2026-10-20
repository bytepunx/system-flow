---
id: TH-0466
title: "E-0022 plan: eight draft stories, S-0398 to S-0405, with 44 tasks; one added story (S-0405) for I-0148's close-out tiers"
anchor:
  path: wip/kanban/epics/E-0022-a-merge-queue-verifies-the-stories-in-review-as-one-batch-bisects-a-red-batch-to-the-story-that-fails-and-lands-the-green-ones-so-an-acceptance-no-longer-stops-on-a-dirty-main-or-a-shared-path-rebase.md
  item: E-0022
status: resolved
participants: [planner-E-0022, alex]
created: 2026-10-16T05:06:37Z
updated: 2026-10-20T19:00:27Z
---

# TH-0466 E-0022 plan: eight draft stories, S-0398 to S-0405, with 44 tasks; one added story (S-0405) for I-0148's close-out tiers

On wip/kanban/epics/E-0022-a-merge-queue-verifies-the-stories-in-review-as-one-batch-bisects-a-red-batch-to-the-story-that-fails-and-lands-the-green-ones-so-an-acceptance-no-longer-stops-on-a-dirty-main-or-a-shared-path-rebase.md.

## Entries

### 2026-10-16T05:06:37Z planner-E-0022
I drafted eight stories for E-0022, all in the backlog, with their tasks. Cost of delay values wait on TH-0465 question 1. The design answers there shape the ADRs, not the plan's order.

## Stories and order

| Story | What | After |
|-------|------|-------|
| S-0398 | `flai queue`: eligible stories and the others' reasons; the queue on the board, MCP, and dashboard; the queue's ADR | none |
| S-0399 | enqueue, hold, `merge-queue` host action, `batch_max`; `flai serve` builds and verifies one batch | S-0398 |
| S-0400 | land a green batch through today's acceptance; bisect a red one; send the failing story back with a thread | S-0399 |
| S-0401 | shared-path stops settled in the batch; other conflicts told by message | S-0399 |
| S-0402 | landing on a dirty main checkout, and telling the operator what it is behind | S-0400 |
| S-0403 | the orchestrator enqueues with evidence and stops gate threads | S-0400 |
| S-0404 | queue wait, batch size, and bisections in `metrics.md`, `flai stats`, and a chart | S-0400 |
| S-0405 | close-out leaves the tiers marked `queue` to the batch (added, I-0148) | S-0399 |

S-0401 and S-0405 may start once S-0399 is accepted. S-0400 and S-0401 both change `flai/internal/serve/mergequeue.go` and `design/system/workflow.md`, so the holds run them one after the other. Each story carries its own documentation and template changes.

## Tasks and layers

- **S-0398:** T-1638, T-1639; then T-1655, T-1656; then T-1667; then T-1675.
- **S-0399:** T-1640, T-1641; then T-1657, T-1658; then T-1668; then T-1676.
- **S-0400:** T-1642, T-1643, T-1644; then T-1659, T-1660; then T-1669; then T-1677.
- **S-0401:** T-1645, T-1646; then T-1661; then T-1670; then T-1678.
- **S-0402:** T-1647, T-1648; then T-1662; then T-1671; then T-1679.
- **S-0403:** T-1649, T-1650; then T-1663; then T-1672.
- **S-0404:** T-1651, T-1652; then T-1664; then T-1673; then T-1680.
- **S-0405:** T-1653, T-1654; then T-1665, T-1666; then T-1674; then T-1681.

Each story that decides something writes its ADR in its first layer and its documentation last.

## Assumptions

1. TH-0465's recommended answers hold. A different answer changes the ADRs of S-0398, S-0399, and S-0405, not the order.
2. The queue's state and its event log live under `.flai-cache` in the main checkout, host-local, as the verify result and the conflict log do. Nothing new goes into front matter.
3. A story enters the queue only when someone enqueues it, the orchestrator with evidence or the operator. Eligibility alone does not enqueue.
4. The queue lands each story through `flai accept`'s own flow, so the transitions, archive, evidence, and per-story metrics stay as they are (epic criterion 4). The verified commit is the story's rebased head in a passing batch (S-0400's ADR).
5. `flai serve` becomes a second maker of review to in-progress, for the queue's send-back only (S-0400's ADR).
6. A shared-path stop that the existing steps cannot settle, such as `design/system/flai-cli.md` in I-0150, needs a new merge. S-0401's ADR decides it; I suggest section by section, by heading, as ADR-0157 merges rows.
7. How a batch lands on a dirty checkout is S-0402's ADR. I suggest landing through the batch's worktree and moving main with git's own overwrite check.
8. `flai accept` by hand is unchanged and takes a story out of the queue.

## Proposals

- **Added, drafted:** S-0405, from TH-0465 question 5. The queue alone leaves I-0148 open, since every close-out still runs integration and smoke. Drop it if you answer no there.
- **Overlap with E-0021, not changed:** E-0021's drafts S-0395 and S-0397 also change `flai/cmd/accept.go`, `flai/internal/preview/accept.go`, and `flai/internal/hostapi/writes.go`. Whichever epic goes second rebases onto the other's acceptance. I recommend running E-0022 first. It remediates the waits I-0150 measures, and E-0021's pull-request mode can then gate the queue's landing rather than a single acceptance.
- **Not drafted:** a queue for releases, publishing after each landed batch. Publishing stays the orchestrator's `publish` permission, which the epic leaves alone.

## Figures

| Story | Forecast | flai's | Delivery |
|-------|----------|--------|----------|
| S-0398 | 47m | 47m | 2026-10-16T08:09Z |
| S-0399 | 1h30m | 1h9m | 09:47Z |
| S-0400 | 1h15m | 40m | 11:09Z |
| S-0401 | 50m | 28m | 10:43Z |
| S-0402 | 45m | 34m | 12:02Z |
| S-0403 | 42m | 42m | 12:00Z |
| S-0404 | 37m | 37m | 11:54Z |
| S-0405 | 41m | 41m | 11:23Z |

Each story's `### Planning` notes give the reason for each raise. The deliveries play out the pull order behind E-0021's stories.

### 2026-10-17T17:57:10Z planner-E-0022
E-0021 went first, so I'm updating the overlap note. S-0395's T-1625 makes `flai accept` refuse on forge blockers before it merges: under pull-request acceptance, a required check not passed or a review outstanding stops it. S-0400's landing goes through `preview.AcceptWith`, so the queue inherits that refusal.

I propose two changes once S-0395 is accepted:

1. **S-0398:** a story with forge blockers is not eligible for a batch. The queue shows those blockers by name, as it shows ADR-0155's undone paths.
2. **S-0400:** a landing refused on a forge blocker after a green batch leaves the story in review with the reason. It is neither bisected nor sent back.

I haven't changed the drafts yet, because S-0395 is still in progress and its blockers may change. I'll make both changes once it is accepted, unless you say otherwise.

### 2026-10-18T02:20:32Z planner-E-0022
S-0395 is accepted. Its ADR-0169 changed what I proposed on 2026-10-17, so I made a smaller change instead.

ADR-0169 blocks every acceptance under pull-request acceptance on `forge_stale` when the kept forge state was read at a commit other than the branch's head. A batch rebases each story onto main, so a queue landing in that mode would always stop there. The forge's CI already verifies each story in that mode, and S-0397 merges through the remote. So I made the queue run under `acceptance.mode: local` only:

- **S-0398:** a new criterion says `flai queue`, the board, and the dashboard report the queue as off under pull-request acceptance, with the reason. The ADR's criterion now records this. T-1638 and T-1639 were updated to match.
- **S-0399:** `flai queue add` refuses under pull-request acceptance, and `flai serve` starts no batch there and logs why once. The operators' guide says so. T-1641 and T-1668 were updated to match.
- **S-0400:** no change. Under local acceptance no forge blocker arises, so the change I proposed for S-0400 is withdrawn.

Forecasts are unchanged, since each change is one check and a test case. `flai check --strict` passes. A queue that feeds the forge's own merge queue under pull-request acceptance would be a story of its own. I haven't drafted it; say if you want it.

### 2026-10-20T19:00:27Z alex
Resolved: S-0404 was accepted
