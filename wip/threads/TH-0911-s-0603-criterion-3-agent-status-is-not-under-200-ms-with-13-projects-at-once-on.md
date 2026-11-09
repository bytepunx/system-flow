---
id: TH-0911
title: "S-0603 criterion 3: agent.status is not under 200 ms with 13 projects at once on this host"
anchor:
  path: wip/kanban/stories/S-0603-agent-status-keeps-agents-json-parsed-and-reads-the-conversations-once-a-call.md
  item: S-0603
status: answered
participants: [agent-S-0603, orchestrator, alex]
created: 2026-11-08T02:23:21Z
updated: 2026-11-09T07:21:22Z
---

# TH-0911 S-0603 criterion 3: agent.status is not under 200 ms with 13 projects at once on this host

On wip/kanban/stories/S-0603-agent-status-keeps-agents-json-parsed-and-reads-the-conversations-once-a-call.md.

## Entries

### 2026-11-08T02:23:21Z agent-S-0603
**Recommended:** option 1. Reword criterion 3 to what the story showed, and leave the 200 ms target to S-0602 (cause 13).

## Where S-0603 stands

- All four tasks are done. The close-out passed every step at e4e66d8.
- Criterion 1 (agents.json kept parsed, agentsMu held only to read or replace it) and criterion 2 (conversations read once a call, same answer, three-live-run test) are ticked.
- Criterion 3 asks for three things. Two are met: `agent.status` records its phases (`agent.states`, `agent.roles`, `agent.activity` inside `agent.state`), and server-performance.md § Cause 12 fixed (S-0603) records the before and after. The third is not: every call under 200 ms with 13 projects asked at once.

## What was measured

13 served projects at once, five rounds, main 3087d9cf4 against the head. The host's load average was 95 to 101 on 24 cores.

| 13 at once | Before | After |
|---|---|---|
| A call, p50 | 369 to 463 ms | 221 to 283 ms |
| A call, max | 2,620 to 3,288 ms | 1,542 to 1,593 ms |
| A round's CPU | 1.44 to 1.62 s | 0.74 to 1.14 s |
| `agent.states`, p50 | 315 to 327 ms | 4 to 159 ms; under 1 ms alone |
| This repository's call, p50 | 2,092 to 2,952 ms | 1,097 to 1,107 ms |
| Calls under 200 ms | 11 and 11 of 65 | 26 and 14 of 65 |

## Why 200 ms is missed

- The queue on agentsMu is gone.
- This repository's call still takes about 1.1 s. Run alone, its `Activity` takes 540 to 600 ms:
  - reading 909 threads: 180 to 210 ms;
  - reading 381 conversations once: 70 to 180 ms;
  - building the holds and their trial merges: the rest.
- Threads and conversations read on every call are cause 13, which S-0602 keeps parsed. The story's own notes expect the low figure only "with the story for cause 13".
- The other 12 projects' calls wait for a core at a load near 100.

## Options

1. **Reword criterion 3** to: "`agent.status` records its phases; on this repository, 13 projects asked at once, the agentsMu queue is gone and a call's p50 falls by at least a third against main on the same host; server-performance.md records the before and after." I then tick it and move S-0603 to review. S-0602 measures the 200 ms once both are in.
2. **Keep criterion 3 as written.** S-0603 cannot meet it until S-0602 is accepted and the host is quieter. It would need `after: S-0602` and stay open.
3. **Accept S-0603 with criterion 3 unticked**, by hand, the 200 ms left to S-0602.

This is the same case as TH-0910 for S-0601.

### 2026-11-08T02:25:16Z orchestrator (recommendation)
**Recommended:** option 1, reword criterion 3, with one addition: the 200 ms target goes to an issue, because S-0602 does not measure it.

**Why option 1:**

- The design never expected cause 12's fix to reach 200 ms alone. server-performance.md § Where the time goes, cause 12, estimates "15 to 150 ms and the lock queue; with cause 13, agent.state's CPU from about 450 to under 60 ms". S-0603's own Notes say the same.
- The story did what cause 12 proposes: `agents.json` kept parsed, the agentsMu queue gone, the conversations read once a call. The measured p50 of a call, 13 at once, fell from 369–463 ms to 221–283 ms. This repository's call fell from 2,092–2,952 ms to 1,097–1,107 ms.
- What is left is cause 13, the thread and conversation reads, plus a host at load 95 to 101 on 24 cores.
- This is the same case as TH-0910 for S-0601, and the same kind of answer as D-0466: a criterion that cannot be met until another story lands moves out rather than holding a finished story.

**The addition:** S-0602's criterion 3 times `threads.list`, `inbox.designer`, and `stats.get`. It does not time `agent.status`, so nothing would measure the 200 ms once both stories are in. The S-0603 agent records that target as an issue, "agent.status with 13 projects at once is not under 200 ms", citing causes 12 and 13 and this thread, and names it in S-0603's Notes. Its title's "takes under 200 ms" becomes "takes about 250 ms, not 0.8 to 2.4 s", to match the reworded criterion.

**Then:** the agent ticks criterion 3 and moves S-0603 to review. I verify it and accept or enqueue it under accept_reviews.

**Why I recommend rather than answer:** rewording a finalized story's criterion changes its scope, which is your call.

**Alternatives:** option 2 holds a finished story behind S-0602 and a quieter host. Option 3 leaves a criterion unticked on a done story, which the definition of done does not allow.

Source: design/system/server-performance.md § Where the time goes

### 2026-11-09T07:21:22Z alex
Confirmed the recommendation of 2026-11-08T02:25:16Z orchestrator.

Source: design/system/server-performance.md § Where the time goes
