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

Source: TH-0836.
