---
id: D-0488
title: "S-0565 needs four host steps before flai serve can start its proof stories in a container: Recommended: do steps 1 to 3 yourself on the host, in this order, then reply \"done\" here"
kind: answer
sources: [TH-0836]
topics: [harness, docker, all]
paths: [design/experiments/S-0565-a-story-worked-end-to-end-by-flai-serve-on-harness-container.md, design/experiments/README.md, design/issues, design/issues/summary.md]
date: 2026-11-07
status: active
absorbed_by: ""
asked_again: []
cited: [TH-0826, TH-0836]
---

# D-0488 S-0565 needs four host steps before flai serve can start its proof stories in a container: Recommended: do steps 1 to 3 yourself on the host, in this order, then reply "done" here

The operator, alex, answered TH-0836 at 2026-11-07T04:19:06Z.

What was asked:

```text
S-0565 is in progress, but its proof runs cannot start yet. Four of T-2449's five host preconditions do not hold, and each is yours: `flai guard` refuses me host configuration writes (ADR-0220), and upgrading the host restarts `flai serve`.
```

The answer:

````text
**Recommended:** do steps 1 to 3 yourself on the host, in this order, then reply "done" here:

```bash
flai host upgrade                                   # 1. to 1.84.0, which has S-0520's mounts (ADR-0253)
make agent-images                                   # 2. from the main checkout once it is at 1.84.0
flai serve agent harness container --login ~/.claude   # 3. the container login
```

**Why only you:** each is a host change: an upgrade, Docker images, and a credential path. flai guard refuses both me and S-0565's agent these. The runbook says to rebuild the images after a change to flai the agent should have. Without 1.84.0 the container has none of the main checkout's `wip/` and `system-flow.yaml`.

**Step 4, S-0558 to ready:** I have not moved it. `flai promote --candidates` lists it at rank 1 (12m, throughput), and it now carries `harness: container`. But if it goes to ready before steps 1 to 3, flai serve starts it on 1.83.0 with stale images and no login, and the proof is used up. Once you reply here that 1 to 3 are done, move S-0558 to ready, and only S-0558. Or tell me to, and I will move it under `promote_to_ready` on my next run. S-0556 waits until S-0558 is in review, as S-0565's agent asks.

**About TH-0826 and TH-0828:** S-0558 and S-0556 now carry `harness: container`, so their concern, an ordinary agent working the proof outside the container, applies only until step 4. Resolve them once S-0558 is in ready.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: docs/operators/runbooks/agent-container.md § Build
````

The operator, alex, answered TH-0836 again at 2026-11-07T06:09:25Z.

The answer:

```text
**Noted, agent-S-0565:** I hold S-0556 in backlog. I will not promote it until alex says so here, even though `flai promote --candidates` lists it at rank 1 (40m, throughput). My recommendation of 04:41:56Z on S-0558 is out of date: S-0558 reached review and I accepted it at 6e4f917cdcd4.

**Recommended, alex:** run criterion 2 on S-0556's container, as agent-S-0565 proposes.

1. Add the line to S-0556's Notes yourself, for example: "Before the close-out, ask the operator on a thread on S-0556 to confirm the change to the guard's read list, and end your session until it is answered." A finalized story's words are yours, so I do not write it.
2. When you are at the dashboard, write "move S-0556" here, or move it to ready yourself. In-progress is at 5 of 5, so flai serve starts it only once a slot frees.
3. Once its agent has moved S-0556 to in-progress, press Stop on its activity page, wait about 30 seconds, then press Retry. Say here what the card showed each time.
4. Answer its question on S-0556's thread. That is the restart after an answered thread.

**Why only you:** ADR-0058 gives the stop and the Retry to the operator. A stopped run gets no new agent until it is retried. An automatic confirmation would let S-0556 start with nobody at the dashboard, and the second proof would be used up.

**Alternative:** leave the Notes as they are. agent-S-0565 then records the restart half of criterion 2 as not met, and the stop half still runs on S-0556.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/adrs/0058-the-agent-host-action-lets-the-dashboard-stop-a-story-s-agent-and-a-stop.md § Decision
```

Source: TH-0836.
