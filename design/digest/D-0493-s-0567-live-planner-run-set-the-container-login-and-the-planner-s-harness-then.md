---
id: D-0493
title: "S-0567 live planner run: set the container login and the planner's harness, then plan S-0559 from the dashboard: Recommended: do agent-S-0567's three steps in its order, then set `planning.agent`…"
kind: answer
sources: [TH-0843]
topics: [harness, docker, planning, all]
paths: [wip/agents/planner.md, design/system/strategic-agents.md, design/issues, design/issues/summary.md]
date: 2026-11-07
status: active
absorbed_by: ""
asked_again: []
cited: [TH-0843]
---

# D-0493 S-0567 live planner run: set the container login and the planner's harness, then plan S-0559 from the dashboard: Recommended: do agent-S-0567's three steps in its order, then set `planning.agent`…

The operator, alex, answered TH-0843 at 2026-11-07T05:38:10Z.

What was asked:

```text
The preconditions an agent can check are in place. Three steps are yours, because the guard refuses an agent host-configuration writes (ADR-0220) and the manifest is yours to set.
```

The answer:

```text
**Recommended:** do agent-S-0567's three steps in its order, then set `planning.agent` back to the host harness as soon as the planner's run on S-0559 ends.

1. `flai serve agent harness container --login <folder>`. Use the folder that holds the `.credentials.json` your login uses. ADR-0247 mounts that folder read-only at `/run/claude-login`, and only `.credentials.json` is read. This is one host setting, and S-0565 needs it too. TH-0836 step 3 names `~/.claude` and this thread names `~/.claude-rmh`. Pick the one that holds the credentials, and it serves both proofs.
2. `flai manifest set planning.agent='{"harness":"container","model":"claude-opus-5-5"}'`.
3. Plan S-0559 from the dashboard.

**Why only you:** each step is a credential path, host configuration, or a manifest key, and the guard refuses all of them to agents (ADR-0220).

**What I do meanwhile:**
- I start no planner with `plan` while `planning.agent` names `container`, so no planner of mine lands in the container beside the proof run.
- I skip S-0559 in `flai plan --candidates` while this thread is open, so the proof run is not used up on the host harness.

For the operator alone: this is a host change.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/adrs/0247-a-strategic-agent-s-container-mounts-the-main-checkout-read-write-at-its-host.md
```

Source: TH-0843.
