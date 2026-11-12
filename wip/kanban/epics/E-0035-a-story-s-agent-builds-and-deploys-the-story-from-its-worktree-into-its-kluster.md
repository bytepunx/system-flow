---
id: E-0035
type: epic
nature: feature
title: A story's agent builds and deploys the story from its worktree into its kluster namespace and checks it there, so criteria that need the running system are ticked before review
status: backlog
owner: alex
created: 2026-11-10T08:09:06Z
updated: 2026-11-10T09:35:46Z
transitions: []
tags: [cli, dashboard]
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: planner
      seconds: 387
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 114
          output: 37560
          cache_read: 3826452
          cache_write: 292869
          cost: 3.6018
    - kind: orchestrator
      seconds: 1000
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 168
          output: 1186
          cache_read: 6659792
          cache_write: 752811
          cost: 2.011
        - model: claude-sonnet-5-5
          input: 24
          output: 2543
          cache_read: 232771
          cache_write: 110824
          cost: 0.2584
cost_of_delay:
  inputs:
    revenue_per_week: 1250
    by: alex
    at: 2026-11-10T09:35:46Z
---
# E-0035 A story's agent builds and deploys the story from its worktree into its kluster namespace and checks it there, so criteria that need the running system are ticked before review

## Outcome

In a project with a kluster and a deploy, a story's agent builds and deploys the story from its worktree while the story is in progress, runs the project's checks against the story's namespace, and ticks the acceptance criteria that need the whole system running before it moves the story to review. A story reaches review with such a criterion unticked only when the agent could not verify it, with the reason in its Notes.

Today the cluster is out of a story's reach until review ([kluster.md](../../../design/system/kluster.md), [ADR-0254](../../../design/adrs/0254-flai-runs-one-kluster-per-project-from-the-config-its-manifest-s-kluster-block.md)):

| Today | What it costs |
|-------|---------------|
| `flai serve` deploys a story only once it enters review, and tears its namespace down when it goes back to in-progress | The agent has nothing running to verify against while it works |
| The close-out and `flai verify` do not use the cluster | A criterion that needs the running system cannot be verified before review |
| `flai guard` refuses every agent `flai kluster deploy`, and gives no other way to ask | The agent moves to review with the criterion unticked, and the operator finds out there whether it holds |
| A failed deploy or checks run in review opens a thread to the operator | The operator carries a failure the story's agent could have found and fixed |

What changes:

- A story's agent asks flai serve for a deploy of its story at its branch's head. flai serve runs the same job as in review: build in the worktree, load, deploy into the story's namespace, `s-nnnn`. The agent never runs a kluster command that writes.
- The agent runs the project's checks against the namespace, and its own commands with the namespace's variables, and reads each result and log.
- The namespace lives from the first deploy until the story is accepted or cancelled. The move to review deploys again only at a new head.
- The move to review, the close-out, and the review page say what was deployed and checked, at which commit.
- A failure before review is the agent's to fix. It reaches the operator only when the agent says the cause is outside the story.

Done when a story in a project with a kluster, such as authstar, is deployed and checked from its worktree by its own agent, and reaches review with every criterion ticked and the deploy and checks recorded at its head, with nobody but the agent acting.

## Stories
- S-0659 An ADR settles how a story in progress is deployed into its kluster namespace at its agent's request: who asks, what the agent may run, when the namespace goes, and how many live at once
- S-0660 A story's agent asks flai serve to deploy its story from its worktree while it is in progress, and the namespace stays until the story is accepted or cancelled
- S-0661 A story's agent runs the project's checks against its story's kluster namespace before review, and its own commands with the namespace's variables
- S-0662 The move to review, the close-out, and the review page say what was deployed and checked in the story's namespace and at which commit, and agents are told to verify there first
- S-0663 The Deployments page shows the namespace of each story in progress, and the operator deploys, checks, and tears down from it

## Notes

- The first story writes the ADR that amends ADR-0254 §4, §5, §7, and §8. No other story starts before it is accepted.
- The cluster is shared infrastructure on the operator's machine, which is why ADR-0254 §7 keeps agents from it. The ADR decides what replaces that: the agent asks, flai serve runs the job, as `flai verify --job` hands a close-out to flai serve, and the `kluster` host action still gates everything.
- More namespaces live at once: one per story in progress that asked, beside those in review. The ADR sets a limit and says what a request over it gets.
- Jobs still run one at a time per project. A deploy asked by an agent waits behind the queue, so the ADR says how the agent waits and for how long.
- ADR-0254 §5a's send-back of a failed deploy in review, the `kluster` blocker on acceptance, and `flai accept --without-kluster` are still not built. This epic does not build them, but a deploy and checks run that passed at the head before review is the evidence that blocker would read.
- Out of scope: a cluster per story, deploying uncommitted work, and letting an agent bring the cluster down, promote, or act on another story's namespace.
