---
id: TH-1103
title: Confirm ADR-0292's choices for deploys a story's agent asks for before review
anchor:
  path: wip/kanban/stories/S-0659-an-adr-settles-how-a-story-in-progress-is-deployed-into-its-kluster-namespace.md
  item: S-0659
status: resolved
participants: [agent-S-0659, orchestrator, alex]
created: 2026-11-12T16:08:15Z
updated: 2026-11-13T07:26:16Z
---

# TH-1103 Confirm ADR-0292's choices for deploys a story's agent asks for before review

On wip/kanban/stories/S-0659-an-adr-settles-how-a-story-in-progress-is-deployed-into-its-kluster-namespace.md.

## Entries

### 2026-11-12T16:08:15Z agent-S-0659
ADR-0292 is written, proposed, and refines ADR-0254 (§4, §5, §5a, §7, §8). `design/system/kluster.md` describes it as decided and marks every part not built with S-0660 to S-0663. The cluster is shared infrastructure on your machine, so these choices are yours to confirm before I accept the ADR.

**Recommended: confirm all eight as written.** Reply `confirm`, or name the numbers to change and how.

| # | Choice | Recommendation | Main alternative |
|---|--------|----------------|------------------|
| 1 | How many stories in progress may hold a namespace | `kluster.agent_namespaces`, default **2**; 0 turns deploys before review off. Review namespaces do not count, so a review deploy is never refused | Count review too: a story entering review could then not be deployed |
| 2 | A deploy request over the limit | **Refused at once**, naming who holds the namespaces. Not queued, nothing evicted | Queue it: a namespace frees only at another acceptance, so the agent could wait days |
| 3 | How the agent waits | Job id at once, then a `kluster` event via `wait_for_events`, as the close-out's `verified`. **One FIFO queue** shared with review deploys and promotions; no wait limit beyond the jobs ahead times `kluster.timeout` | Review deploys first, so acceptances never wait behind agents |
| 4 | When a failure reaches you | A failure **before review is the agent's**: no thread, no send-back. It reaches you only when the agent says on a thread that the cause is outside the story, or at review (§5a as S-0694 built it) | flai serve opens its failure thread for agent deploys too |
| 5 | How the agent is held to its own namespace | A **kubeconfig scoped to the namespace** (ServiceAccount with the `edit` role there, 24h token) for the agent's own commands, plus guard refusals of `up`, `down`, `promote`, other namespaces, and flai's cluster kubeconfig | Guard alone: it cannot see inside scripts the agent runs |
| 6 | When a namespace goes | Kept from its first deploy **until acceptance, cancellation, or a move back to ready or backlog**, including when anyone sends a story from review back to in-progress (today that tears it down) | Tear down on any leave from review, as ADR-0254 §5 had it |
| 7 | What "at the head" means for the evidence | The record's commit, or a later head whose commits change **only `design`, `docs`, and `wip`**, so the close-out's issue commit does not force a redeploy | The exact head only |
| 8 | The close-out | A new `kluster` step reads the record and runs nothing; it **fails only on a deploy or checks run that failed at the head**. Not deploying is no failure | Report only, never fail |

What each answer leads to: `confirm` accepts ADR-0292 as it stands. A change to any number: I edit the ADR and `kluster.md` to match, then accept it.

Source: design/adrs/0292-a-story-s-agent-asks-flai-serve-to-deploy-check-and-tear-down-its-own-story-s.md § Decision

### 2026-11-12T16:09:58Z orchestrator (recommendation)
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

### 2026-11-13T07:25:04Z alex
Confirmed the recommendation of 2026-11-12T16:09:58Z orchestrator.

Source: design/adrs/0254-flai-runs-one-kluster-per-project-from-the-config-its-manifest-s-kluster-block.md § 7. The `kluster` host action

### 2026-11-13T07:26:16Z agent-S-0659
Resolved: The operator confirmed all eight choices; ADR-0292 is accepted as written, with the confirmation and the caveat on "at the head" among its consequences (commit 7abf2f1).
