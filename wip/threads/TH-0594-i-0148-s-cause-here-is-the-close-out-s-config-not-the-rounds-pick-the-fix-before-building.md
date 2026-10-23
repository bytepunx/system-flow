---
id: TH-0594
title: "I-0148's cause here is the close-out's config, not the rounds: pick the fix before building"
anchor:
  path: wip/kanban/stories/S-0466-a-close-out-whose-run-main-outruns-takes-two-rounds-of-the-integration-and-smoke-tiers-about-24-minutes-past-the-bash-tool-s-10-minute-limit-so-the-agent-s-call-is-killed-mid-round.md
  item: S-0466
status: resolved
participants: [agent-S-0466, orchestrator, alex]
created: 2026-10-22T21:00:30Z
updated: 2026-10-23T03:33:49Z
---

# TH-0594 I-0148's cause here is the close-out's config, not the rounds: pick the fix before building

On wip/kanban/stories/S-0466-a-close-out-whose-run-main-outruns-takes-two-rounds-of-the-integration-and-smoke-tiers-about-24-minutes-past-the-bash-tool-s-10-minute-limit-so-the-agent-s-call-is-killed-mid-round.md.

## Entries

### 2026-10-22T21:00:30Z agent-S-0466
**Finding.** ADR-0174 already meant to close I-0148: while the merge queue can batch, flai verify leaves integration and smoke to the batch, and its Consequences say "I-0148 closes once this runs with the queue on". The queue is on for this project in flai serve's configuration (~/.flai/config.json: merge-queue, accept_reviews on, acceptance local). But scripts/close-out.sh runs the tree's flai through scripts/flai.sh, which reads .flai-cache/config.json, whose host_actions hold only push. So queueBatches (flai/cmd/verify.go:392) answers false in every close-out here, and integration and smoke run anyway.

Evidence from .flai-cache/verify, all after S-0405 landed ADR-0174 on 2026-10-19:

| Story | Run | integration and smoke |
|---|---|---|
| S-0471 | 14m27s | ran |
| S-0473 | 12m35s | ran |
| S-0479 | 10m32s | ran |
| S-0483 | 5m04s | queued (an MCP verify, which is the installed flai reading ~/.flai/config.json) |

So one round alone already exceeds the 600000 ms call, and the close-out and the MCP verify disagree about the same story.

**Options.**

1. **(Recommended) Make the close-out see the queue, and say when it does not.** (a) You turn merge-queue on for this project in the dev config: `scripts/flai.sh serve enable merge-queue` (I make no change under ~/.flai, and I won't change .flai-cache/config.json without your yes). (b) flai verify adds one line when it runs a tier marked queue because the queue cannot batch, naming which condition is off and the configuration file it read, so this mismatch shows up in the close-out's output instead of as a 14-minute call. (c) devex.md and docs/contributors say the dev config needs the host action too. Small: about an hour. I-0148 closes on a close-out that leaves both tiers queued.
2. **The planner's plan.** A Go `flai close-out` that runs the rounds detached under .flai-cache/close-out and is followed by calls that each return within the tool's limit (T-2030 ADR, T-2031, T-2033, plus the convention, template, harness, and verifier changes its touches name). This covers projects whose tiers are long with the queue off. It is several hours, and on its own it does not stop this repo running 12-minute tiers the queue would run again anyway.
3. Both: 1 now in this story, and 2 as a backlog story from I-0148's second half (queue off), if you want that case covered.

Which one? With 1 or 3 I cancel T-2030, T-2031 and T-2033 and write the tasks for 1. With 2 I build the planner's tasks as drafted.

### 2026-10-22T21:07:14Z orchestrator (recommendation)
Recommended for the operator: option 3. Do option 1 in S-0466 now, and record the queue-off case for later.

- **The agent's evidence settles the cause here.** ADR-0174 already leaves integration and smoke to the merge queue, but the tree's `scripts/flai.sh` reads `.flai-cache/config.json`, which lacks the `merge-queue` host action. So every close-out in this repo runs 10 to 14 minutes of tiers the queue then runs again. The MCP verify of S-0483, which reads `~/.flai/config.json`, queued them and took 5 minutes.
- **Option 1 is small and closes I-0148 for this repo.** You run `scripts/flai.sh serve enable merge-queue`. flai verify names the condition that is off when it runs a `queue` tier. The dev docs say the dev config needs the host action.
- **The detached `flai close-out` still has value**, for projects whose tiers run long with the queue off, which is the template's default. But it is hours of work for a case this repo no longer hits. Keep it as a backlog story from I-0148, valued on its own.

This replaces my recommendation on TH-0596, which I made before reading this thread.

Source: design/adrs/0174-while-the-merge-queue-can-batch-flai-verify-and-the-close-out-leave-the-test.md

### 2026-10-23T03:21:35Z alex
Confirmed the recommendation of 2026-10-22T21:07:14Z orchestrator.

Source: design/adrs/0174-while-the-merge-queue-can-batch-flai-verify-and-the-close-out-leave-the-test.md

### 2026-10-23T03:33:49Z agent-S-0466
Resolved: Option 3 confirmed. S-0466 built option 1: T-2061 adds the line flai verify prints when it runs a queue tier because the queue cannot batch, and T-2062 documents that the dev configuration needs the host action. T-2030, T-2031, T-2033, T-2046, T-2050 and T-2054 are cancelled. The detached close-out is left for a backlog story from I-0148.
