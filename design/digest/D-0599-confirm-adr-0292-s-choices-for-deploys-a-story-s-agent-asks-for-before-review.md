---
id: D-0599
title: "Confirm ADR-0292's choices for deploys a story's agent asks for before review: Recommended: confirm all eight as written"
kind: answer
sources: [TH-1103, S-0659]
topics: [kluster, deployment, flai, cli, go, code, all]
paths: [design/system/kluster.md, design/adrs/README.md, design/adrs/0254-flai-runs-one-kluster-per-project-from-the-config-its-manifest-s-kluster-block.md, design/adrs/0292-a-story-s-agent-asks-flai-serve-to-deploy-check-and-tear-down-its-own-story-s.md]
date: 2026-11-13
status: absorbed
absorbed_by: "design/adrs/0292-a-story-s-agent-asks-flai-serve-to-deploy-check-and-tear-down-its-own-story-s.md#consequences"
asked_again: []
cited: []
---

# D-0599 Confirm ADR-0292's choices for deploys a story's agent asks for before review: Recommended: confirm all eight as written

The operator, alex, answered TH-1103 at 2026-11-13T07:25:04Z.

What was asked:

```text
ADR-0292 is written, proposed, and refines ADR-0254 (§4, §5, §5a, §7, §8). `design/system/kluster.md` describes it as decided and marks every part not built with S-0660 to S-0663. The cluster is shared infrastructure on your machine, so these choices are yours to confirm before I accept the ADR.
```

The answer:

```text
**Recommended: confirm all eight as written.** I read ADR-0292 on `story/S-0659` against ADR-0254. Each choice keeps what ADR-0254 §7 protects: flai serve runs every job, the `kluster` host action gates it, and an agent reaches only its own namespace.

**Why each holds:**

- **1 and 2, the limit and the refusal.** Review deploys are the acceptance gate's evidence, so they must never wait on the limit. A namespace frees only at an acceptance or a cancellation, so a queued request could sit for days; a refusal naming the holders lets the agent go on.
- **3, one FIFO queue.** ADR-0254 §5 already runs every job one at a time behind one lock. A priority lane for review is a later change if the Deployments page shows acceptances waiting.
- **4, failures before review are the agent's.** It follows the close-out's way (ADR-0229): the agent fixes its own change. §5a's send-back still applies at review.
- **5, the scoped kubeconfig.** The guard cannot read inside scripts the agent runs, so only Kubernetes RBAC on the namespace actually holds the boundary. The token is for the local cluster only and lasts 24h.
- **6, keeping the namespace through a send-back.** The agent keeps the environment it is fixing the change in.
- **7, "at the head".** Commits under `design`, `docs`, and `wip` do not change what is deployed, so the close-out's issue commit does not force a redeploy.
- **8, the close-out's `kluster` step.** It stops only a change the cluster could not run at the head. Not every story needs the running system.

**One caveat for 7:** if a project's deploy ever builds from `docs` or `design`, a change there would count as "at the head" without being deployed. No project's deploy does that today, so I would not change it now.

**Why this needs you:** choice 5 creates credentials, and every choice changes how agents use shared infrastructure on your machine. That is yours alone, so flai serve does not confirm this after `confirm_after`.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/adrs/0254-flai-runs-one-kluster-per-project-from-the-config-its-manifest-s-kluster-block.md § 7. The `kluster` host action
```

S-0659 accepted ADR-0292 on that answer, with no choice changed. From its Decisions (`wip/archive/agents/S-0659.md`):

```text
The operator confirmed all eight of ADR-0292's choices on TH-1103 (2026-11-13). The orchestrator's caveat on "at the head" (a deploy that builds from design or docs) is recorded as a consequence in the ADR; no choice changed.
```

Sources: TH-1103; S-0659 (`wip/archive/agents/S-0659.md`, Decisions).
