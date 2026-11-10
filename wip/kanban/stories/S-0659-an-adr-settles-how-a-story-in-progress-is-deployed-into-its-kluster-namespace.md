---
id: S-0659
type: story
nature: feature
title: "An ADR settles how a story in progress is deployed into its kluster namespace at its agent's request: who asks, what the agent may run, when the namespace goes, and how many live at once"
status: backlog
parent: E-0035
owner: alex
created: 2026-11-10T08:09:07Z
updated: 2026-11-10T09:29:05Z
transitions: []
tags: [cli]
topics: [kluster, deployment]
touches: [design/system/kluster.md, design/adrs/README.md, design/adrs/0254-flai-runs-one-kluster-per-project-from-the-config-its-manifest-s-kluster-block.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 167
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 36
          output: 184
          cache_read: 1226797
          cache_write: 152534
          cost: 0.3727
draft: true
forecast:
  duration: 20m
  delivery: 2026-11-10T13:41:00Z
  basis: "Its own forecast of 20m; 10th in the pull order with an in-progress limit of 5, behind S-0670, S-0671, S-0682, S-0684, S-0687, S-0668, S-0669, S-0673, S-0672, S-0665, S-0667, S-0674, S-0704 and S-0241."
  by: flai
  at: 2026-11-10T09:29:05Z
---
# S-0659 An ADR settles how a story in progress is deployed into its kluster namespace at its agent's request: who asks, what the agent may run, when the namespace goes, and how many live at once

## Goal

An ADR settles how a story in progress is deployed into its kluster namespace at its agent's request, so that the stories after it build one design. It amends [ADR-0254](../../../design/adrs/0254-flai-runs-one-kluster-per-project-from-the-config-its-manifest-s-kluster-block.md), and [kluster.md](../../../design/system/kluster.md) says the same.

## Acceptance criteria

- [ ] An accepted ADR says how a story's agent asks for a deploy and a teardown of its own story's namespace, that flai serve runs the job, and that the `kluster` host action gates it
- [ ] It says what the agent may run against the namespace: the project's checks, its own commands with the namespace's variables, and what `flai guard` refuses, another story's namespace and the default one among them
- [ ] It says when a namespace is made, deployed again, and torn down across in-progress, review, a send back, acceptance, and cancellation, in a table that replaces ADR-0254 §4's
- [ ] It sets how many story namespaces may live at once, where the project sets it, and what a request over the limit gets
- [ ] It says how the agent waits for a job behind the queue, what it is told on a failure, and when a failure reaches the operator
- [ ] It says what is recorded of a deploy and a checks run before review, and where the move to review, the close-out, and an acceptance read it
- [ ] ADR-0254 names it in `superseded_by` or the new ADR names the sections it amends, as `design/conventions/decisions.md` has it, and `design/system/kluster.md` describes the design as decided and marks each part not yet built

## Tasks
- T-2860 Write the ADR, proposed, that amends ADR-0254 §4, §5, §7, and §8 for a deploy a story's agent asks for while its story is in progress
- T-2861 Describe in design/system/kluster.md the deploy a story's agent asks for as decided, marking each part not yet built
- T-2862 Ask the operator to confirm the ADR's decisions on a thread, and accept it with their answer

## Notes

- A design story: it changes no code.
- The hand-off `flai verify --job` makes to flai serve is the pattern to weigh first.

### Planning

- Touches:
  - `design/system/kluster.md`: design. The goal links it and the last criterion has it describe the decided design.
  - `design/adrs/README.md`: co-change (27% of the commits that changed kluster.md and ADR-0254). `flai adr new` and `flai adr accept` write the index again.
  - `design/adrs/0254-flai-runs-one-kluster-per-project-from-the-config-its-manifest-s-kluster-block.md`: design. It changes only if the ADR supersedes it and `superseded_by` is set; with `--refines`, as T-2860 expects, it stays as it is.
  - No folder touch on `design/adrs`. The new ADR's file takes the number `flai adr new` gives it when the story runs, and `flai adr new --commit` adds that file to the story's touches itself.
  - Left out from `flai touches suggest`: the docs and `flai-cli.md`, `project-manifest.md`, `flaiover-dashboard.md`, and `design/tech/kluster.md`. They describe what is built; S-0660 to S-0663 change them as they build it.
- Topics `kluster` and `deployment` were added, the topics of kluster.md and ADR-0254, so the story's pack briefs them.
- No `after`: the story reads no other story's output. S-0660 to S-0663 build on this one's ADR.
- Forecast 20m stands as `flai forecast` gives it. S-0576, which wrote ADR-0254 and its design, took 710 s of agent time; this ADR amends four of its sections and edits the same 56 KB design, so a little more is fair. The wait on the operator in T-2862 is wall time, not agent time.
- Cost of delay: not set yet. Neither the story nor E-0035 has inputs; TH-0982 asks the operator for them.
- Assumptions: the limit on story namespaces is set in the manifest's `kluster:` block; the third criterion's "ADR-0254 §4's" table means §5's, the only table of when a namespace is deployed and torn down (§4 is the variables).
