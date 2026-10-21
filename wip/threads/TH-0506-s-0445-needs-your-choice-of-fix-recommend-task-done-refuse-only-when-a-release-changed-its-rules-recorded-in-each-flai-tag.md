---
id: TH-0506
title: "S-0445 needs your choice of fix: recommend task done refuse only when a release changed its rules, recorded in each flai tag"
anchor:
  path: wip/kanban/stories/S-0445-the-host-s-installed-flai-is-again-behind-the-project-s-releases-after-i-0149-closed-so-the-mcp-task-done-refuses-every-close.md
  item: S-0445
status: answered
participants: [planner-S-0445, orchestrator, alex]
created: 2026-10-18T02:30:52Z
updated: 2026-10-20T21:12:38Z
---

# TH-0506 S-0445 needs your choice of fix: recommend task done refuse only when a release changed its rules, recorded in each flai tag

On wip/kanban/stories/S-0445-the-host-s-installed-flai-is-again-behind-the-project-s-releases-after-i-0149-closed-so-the-mcp-task-done-refuses-every-close.md.

## Entries

### 2026-10-18T02:30:52Z planner-S-0445
**Recommended: B.** S-0445 needs one decision before I can draft its tasks. Which fix should it build?

What the instances show: all six instances of I-0157 have the same cause.

- Every feature story publishes a minor flai release. Ten minor releases were tagged between 2026-10-16T22:45Z and 2026-10-18T02:19Z, from 1.53.0 to 1.62.0.
- `auto-upgrade` is off for every project on this host. `flai serve` asked about each release on a thread, as ADR-0158 and ADR-0161 say. That is 14 threads, TH-0469 to TH-0505, all still open. The host was upgraded once, from 1.52.0 to 1.57.0.
- ADR-0142's minor threshold therefore refuses `task_done` nearly all the time. Each agent fell back to the tree's `scripts/flai.sh task done`, at 3m per occurrence.
- ADR-0142 only bites in a project whose history holds `flai/vX.Y.Z` tags. In practice that is this repository.

The options:

- **B (recommended): refuse only when a release changed task done's rules.**
  - flai carries a task-done rules number, `buildinfo.TaskDoneRules`, and `flai version --json` prints it.
  - `flai release` writes `task-done-rules: N` into each `flai/vX.Y.Z` tag's annotation. The tags are already annotated (`release.go` `Tag`, `TagPending`).
  - `taskdone.refuseOutdated` and ADR-0160's caught-up check compare the rules numbers instead of the minor versions. A tag without the line keeps today's minor rule.
  - A new ADR refines ADR-0142. You chose its threshold on TH-0401, so this needs your word.
  - B works with `auto-upgrade` off. It refuses only on the rare release that changes how a task closes. It also covers the minutes between a tag and its publication, which auto-upgrade cannot.
  - Forecast: about 4h of agent time, in six tasks.
- **A: you turn on auto-upgrade for system-flow** (`flai serve enable auto-upgrade` in the main checkout). Nothing to build: ADR-0158 built it for this. But each release still refuses closes for 10 to 15 minutes, until it is published and installed. The story then shrinks to one task: confirm and close I-0157.
- **C: keep the refusal, and make the threads easier to act on.** `flai serve` would resolve an upgrade thread once a newer release supersedes it, so you would see one thread instead of 14. The cause would still recur until you upgrade.

My assumptions for B:

- Bumping `TaskDoneRules` becomes part of any change to task done's commit, move, log, or touches steps. The design and the ADR will say so.
- C is worth a story of its own whatever you choose. I would draft it only if you say so.

Reply with A, B, or C, or your own fix. I will draft the tasks and the forecast for it.

### 2026-10-18T02:31:12Z orchestrator (recommendation)
**Recommendation:** A now, and B as S-0445.

- **A now.** Run `flai serve enable auto-upgrade` in the main checkout. It costs nothing to build, and it stops the refusals within one publish cycle today.
- **B as S-0445.** Make `task done` refuse only when a release changed its rules. That ends the refusals for good, including the minutes between a tag and its publication, which A cannot cover.
- **C is not needed with A on.** With auto-upgrade on, the upgrade threads stop piling up.

## A note on my part in this

I have published after every acceptance under the `judgement` policy: ten minor flai releases in about 28 hours. That cadence is what keeps the host behind ADR-0142's minor threshold.

If you prefer, I can batch releases instead, for example once per epic or once a day, with the reasons in the log. That would cut the churn without any change to flai. Say so here, and I will hold publishing to that rhythm from now on.

## Why you

ADR-0142's threshold is yours (TH-0401), and B refines it with a new ADR. Turning on a host action is also yours. Neither is mine to settle. The planner drafts S-0445's tasks once you answer.

Source: design/adrs/0142-flai-task-done-refuses-to-close-a-task-while-the-running-flai-s-major-or-minor.md § Decision

### 2026-10-20T20:43:30Z planner-S-0445
You finalized S-0445 while this thread was still open, and an agent has pulled it. I have planned it as B, the fix both recommendations here name: T-1920 to T-1926. The plan is in its own thread on S-0445. T-1920's ADR comes first and still needs your acceptance. Answer A or C here instead, and that task is rewritten and the code tasks are dropped.

### 2026-10-20T21:12:38Z alex
Confirmed the recommendation of 2026-10-18T02:31:12Z orchestrator.

Source: design/adrs/0142-flai-task-done-refuses-to-close-a-task-while-the-running-flai-s-major-or-minor.md § Decision
