---
id: TH-0836
title: S-0565 needs four host steps before flai serve can start its proof stories in a container
anchor:
  path: wip/kanban/stories/S-0565-one-story-of-this-repository-is-worked-end-to-end-by-flai-serve-on-harness.md
  item: S-0565
status: open
participants: [agent-S-0565, orchestrator]
created: 2026-11-07T03:44:23Z
updated: 2026-11-07T03:46:17Z
---

# TH-0836 S-0565 needs four host steps before flai serve can start its proof stories in a container

On wip/kanban/stories/S-0565-one-story-of-this-repository-is-worked-end-to-end-by-flai-serve-on-harness.md.

## Entries

### 2026-11-07T03:44:23Z agent-S-0565
S-0565 is in progress, but its proof runs cannot start yet. Four of T-2449's five host preconditions do not hold, and each is yours: `flai guard` refuses me host configuration writes (ADR-0220), and upgrading the host restarts `flai serve`.

## What I found

| Precondition | Found | Holds |
|--------------|-------|-------|
| Installed flai holds S-0508, S-0510, S-0517, S-0520 | `flai serve` runs 1.83.0. S-0520's mounts (ADR-0253) are only in 1.84.0, published but not installed | No |
| Images built | All three were built 2026-11-02 ~18:52Z. That is before S-0510's final layer commits (2026-11-03T01:19Z) and before 1.84.0 | No |
| Container harness login | `flai serve agent harness` says `login: not set` | No |
| Container reaches the host MCP listener (I-0254) | `wget` from the default bridge to `host.docker.internal:4243/mcp` answers 401, so the route works (Docker Desktop) | Yes |
| Proof stories startable | S-0558 and S-0556 are in backlog | No |

Done already: S-0558 and S-0556 now carry `harness: container`. The results document records the table with each command, on `story/S-0565`.

## What I ask, in this order

1. **Upgrade the host to flai 1.84.0:** `flai host upgrade`.
2. **Rebuild the images** from the main checkout at 1.84.0: `make agent-images`. That builds bedrock, claude-code, and `flai-agent-sf-claude-code:local` from the current layer.
3. **Set the login:** `flai serve agent harness container --login ~/.claude`.
4. **Move S-0558 to ready**, and only S-0558 for now. S-0556 goes to ready once S-0558 is in review (T-2452 waits on T-2451).

Reply here once 1–4 are done. While S-0558's agent works in its container, I will ask you on this thread to stop it once from the dashboard (criterion 2).

**Recommendation:** do all four as listed. To keep the experiment smaller, you could skip step 1 for criteria 1 and 2 only. I don't recommend that: without S-0520's `wip/` and `system-flow.yaml` mounts, `flai` inside the container lacks the main checkout's state, so a failure would not show which layer broke.

When the answer comes, I check each precondition again, record it, and close T-2449.

### 2026-11-07T03:46:17Z orchestrator (recommendation)
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
